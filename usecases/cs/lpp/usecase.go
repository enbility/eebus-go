package lpp

import (
	"sync"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	features "github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/usecase"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
)

type CsLPP struct {
	*usecase.UseCaseBase

	pendingMux    sync.Mutex
	pendingLimits map[model.MsgCounterType]*spineapi.Message

	heartbeatDiag *features.DeviceDiagnosis

	heartbeatKeoWorkaround bool // required because KEO Stack uses multiple identical entities for the same functionality, and it is not clear which to use
}

var _ ucapi.CsLPPInterface = (*CsLPP)(nil)

func NewCsLPP(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *CsLPP {
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeGridGuard,
		model.EntityTypeTypeCEM, // KEO uses this entity type for an SMGW whysoever
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeControllableSystem,
		model.UseCaseNameTypeLimitationOfPowerProduction,
		"1.0.0",
		"release",
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4},
		eventCB,
		validEntityTypes,
	)

	uc := &CsLPP{
		UseCaseBase:   usecase,
		pendingLimits: make(map[model.MsgCounterType]*spineapi.Message),
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (e *CsLPP) loadControlServerAndLimitId() (lc *server.LoadControl, limitid model.LoadControlLimitIdType, err error) {
	limitid = model.LoadControlLimitIdType(0)

	lc, err = server.NewLoadControl(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}
	descriptions, err := lc.GetLimitDescriptionsForFilter(filter)
	if err != nil || len(descriptions) != 1 || descriptions[0].LimitId == nil {
		return
	}
	description := descriptions[0]

	if description.LimitId == nil {
		return
	}

	return lc, *description.LimitId, nil
}

// callback invoked on incoming write messages to this
// loadcontrol server feature.
// the implementation only considers write messages for this use case and
// approves all others
func (e *CsLPP) loadControlWriteCB(msg *spineapi.Message) {
	if msg.RequestHeader == nil || msg.RequestHeader.MsgCounter == nil ||
		msg.Cmd.LoadControlLimitListData == nil {
		return
	}

	_, limitId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return
	}

	data := msg.Cmd.LoadControlLimitListData

	// we assume there is always only one limit
	if data == nil || data.LoadControlLimitData == nil ||
		len(data.LoadControlLimitData) == 0 {
		return
	}

	e.pendingMux.Lock()

	// check if there is a matching limitId in the data
	for _, item := range data.LoadControlLimitData {
		if item.LimitId == nil ||
			limitId != *item.LimitId {
			continue
		}

		if _, ok := e.pendingLimits[*msg.RequestHeader.MsgCounter]; !ok {
			e.pendingLimits[*msg.RequestHeader.MsgCounter] = msg
			e.pendingMux.Unlock()
			e.EventCB(msg.DeviceRemote.Ski(), msg.DeviceRemote, msg.EntityRemote, WriteApprovalRequired)
			return
		}
	}

	e.pendingMux.Unlock()

	// approve, because this is no request for this usecase
	go e.ApproveOrDenyProductionLimit(*msg.RequestHeader.MsgCounter, true, "")
}

func (e *CsLPP) AddFeatures() {
	// client features
	_ = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)

	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	_ = f.AddWriteApprovalCallback(e.loadControlWriteCB)

	newLimitDesc := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
		MeasurementId:  util.Ptr(model.MeasurementIdType(0)), // This is a fake Measurement ID, as there is no Electrical Connection server defined, it can't provide any meaningful. But KEO requires this to be set :(
		Unit:           util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}
	if lc, err := server.NewLoadControl(e.LocalEntity); err == nil {
		limitId := lc.AddLimitDescription(newLimitDesc)

		newLimiData := model.LoadControlLimitDataType{
			Value:             model.NewScaledNumberType(0),
			IsLimitChangeable: util.Ptr(true),
			IsLimitActive:     util.Ptr(false),
		}
		_ = lc.UpdateLimitDataForId(newLimiData, nil, *limitId)
	}

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)

	if dcs, err := server.NewDeviceConfiguration(e.LocalEntity); err == nil {
		dcs.AddKeyValueDescription(
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
				Unit:      util.Ptr(model.UnitOfMeasurementTypeW),
			},
		)
		dcs.AddKeyValueDescription(
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
		)

		value := &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(0),
		}
		_ = dcs.UpdateKeyValueDataForFilter(
			model.DeviceConfigurationKeyValueDataType{
				Value:             value,
				IsValueChangeable: util.Ptr(false),
			},
			nil,
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
		)

		value = &model.DeviceConfigurationKeyValueValueType{
			Duration: model.NewDurationType(0),
		}
		_ = dcs.UpdateKeyValueDataForFilter(
			model.DeviceConfigurationKeyValueDataType{
				Value:             value,
				IsValueChangeable: util.Ptr(false),
			},
			nil,
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		)
	}

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, false)

	if ec, err := server.NewElectricalConnection(e.LocalEntity); err == nil {
		// ElectricalConnectionId and ParameterId should be identical to the ones used
		// in a MPC Server role implementation, which is not done here (yet)
		newCharData := model.ElectricalConnectionCharacteristicDataType{
			ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
			CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
			CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualProductionNominalMax),
			Unit:                   util.Ptr(model.UnitOfMeasurementTypeW),
		}
		_, _ = ec.AddCharacteristic(newCharData)
	}
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *CsLPP) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEnergyGuard,
		e.UseCaseName,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeDeviceDiagnosis,
		},
	) {
		return false, nil
	}

	if _, err := client.NewDeviceDiagnosis(e.LocalEntity, entity); err != nil {
		return false, api.ErrFunctionNotSupported
	}

	return true, nil
}
