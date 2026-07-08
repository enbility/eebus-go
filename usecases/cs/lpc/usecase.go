package lpc

import (
	"errors"
	"sync"

	"github.com/enbility/eebus-go/api"
	features "github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/eebus-go/usecases/usecase"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type LPC struct {
	*usecase.UseCaseBase

	pendingMux    sync.Mutex
	pendingLimits map[model.MsgCounterType]*spineapi.Message

	pendingDeviceConfigMux sync.Mutex
	pendingDeviceConfigs   map[model.MsgCounterType]*spineapi.Message

	heartbeatDiag *features.DeviceDiagnosis

	heartbeatKeoWorkaround bool // required because KEO Stack uses multiple identical entities for the same functionality, and it is not clear which to use
}

var _ ucapi.CsLPCInterface = (*LPC)(nil)

// Add support for the Limitation of Power Consumption (LPC) use case
// as a Controllable System actor
//
// Note: if the Monitoring of Power Consumption (MPC) or Monitoring of Grid Connection Point (MGCP) will be supported, add them first
//
//	Parameters:
//	 - localEntity: The local entity which should support the use case
//	 - eventCB: The callback to be called when an event is triggered (optional, can be nil)
func NewLPC(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) *LPC {
	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeEnergyGuard}
	validEntityTypes := []model.EntityTypeType{
		model.EntityTypeTypeGridGuard,
		model.EntityTypeTypeCEM,
	}
	useCaseScenarios := []api.UseCaseScenario{
		{
			Scenario:  model.UseCaseScenarioSupportType(1),
			Mandatory: true,
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(2),
			Mandatory: true,
		},
		{
			Scenario:       model.UseCaseScenarioSupportType(3),
			Mandatory:      true,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
		},
		{
			Scenario:  model.UseCaseScenarioSupportType(4),
			Mandatory: true,
		},
	}

	usecase := usecase.NewUseCaseBase(
		localEntity,
		model.UseCaseActorTypeControllableSystem,
		model.UseCaseNameTypeLimitationOfPowerConsumption,
		"1.0.0",
		"release",
		useCaseScenarios,
		eventCB,
		UseCaseSupportUpdate,
		validActorTypes,
		validEntityTypes,
		false,
	)

	uc := &LPC{
		UseCaseBase:          usecase,
		pendingLimits:        make(map[model.MsgCounterType]*spineapi.Message),
		pendingDeviceConfigs: make(map[model.MsgCounterType]*spineapi.Message),
	}

	_ = localEntity.Device().Events().Subscribe(uc)

	return uc
}

func (e *LPC) loadControlServerAndLimitId() (lc *server.LoadControl, limitid model.LoadControlLimitIdType, err error) {
	limitid = model.LoadControlLimitIdType(0)

	lc, err = server.NewLoadControl(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
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

func (e *LPC) approveOrDenyConsumptionLimit(msg *spineapi.Message, approve bool, reason string) {
	f := e.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	result := model.ErrorType{
		ErrorNumber: model.ErrorNumberType(0),
	}

	if !approve {
		result.ErrorNumber = model.ErrorNumberType(7)
		result.Description = util.Ptr(model.DescriptionType(reason))
	}
	f.ApproveOrDenyWrite(msg, result)
}

// callback invoked on incoming write messages to this
// loadcontrol server feature.
// the implementation only considers write messages for this use case and
// approves all others
func (e *LPC) loadControlWriteCB(msg *spineapi.Message) {
	if msg.RequestHeader == nil || msg.RequestHeader.MsgCounter == nil ||
		msg.Cmd.LoadControlLimitListData == nil {
		logging.Log().Debug("LPC loadControlWriteCB: invalid message")
		return
	}

	_, limitId, err := e.loadControlServerAndLimitId()
	if err != nil {
		logging.Log().Debug("LPC loadControlWriteCB: error getting limit id")
		return
	}

	data := msg.Cmd.LoadControlLimitListData

	// we assume there is always only one limit
	if data == nil || data.LoadControlLimitData == nil ||
		len(data.LoadControlLimitData) == 0 {
		logging.Log().Debug("LPC loadControlWriteCB: no data")
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
			e.EventCB(msg.DeviceRemote.Ski(), msg.DeviceRemote, msg.EntityRemote, LimitWriteApprovalRequired)
			return
		}
	}
	e.pendingMux.Unlock()

	// approve, because this is no request for this usecase
	go e.approveOrDenyConsumptionLimit(msg, true, "")
}

func (e *LPC) approveOrDenyDeviceConfiguration(msg *spineapi.Message, approve bool, reason string) {
	f := e.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

	result := model.ErrorType{
		ErrorNumber: model.ErrorNumberType(0),
	}

	if !approve {
		result.ErrorNumber = model.ErrorNumberType(7)
		result.Description = util.Ptr(model.DescriptionType(reason))
	}

	f.ApproveOrDenyWrite(msg, result)
}

// callback invoked on incoming write messages to this
// DeviceConfiguration server feature.
// the implementation only considers write messages for this use case and
// approves all others
func (e *LPC) deviceConfigurationWriteCB(msg *spineapi.Message) {
	configsToApprove := map[model.DeviceConfigurationKeyNameType]struct{}{
		model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit: {},
		model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum:             {},
	}
	approvalRequired, err := internal.ConfigurationWriteRequiresApproval(msg, e.LocalEntity, configsToApprove)
	if err != nil {
		logging.Log().Errorf("LPC deviceConfigurationWriteCB: %s", err.Error())
		return
	}

	if approvalRequired {
		e.pendingDeviceConfigMux.Lock()
		if _, exists := e.pendingDeviceConfigs[*msg.RequestHeader.MsgCounter]; !exists {
			e.pendingDeviceConfigs[*msg.RequestHeader.MsgCounter] = msg
			// Unlock before calling EventCB to avoid deadlock (EventCB will need to read pendingDeviceConfigs)
			e.pendingDeviceConfigMux.Unlock()
			e.EventCB(msg.DeviceRemote.Ski(), msg.DeviceRemote, msg.EntityRemote, ConfigurationWriteApprovalRequired)
			return
		}
		e.pendingDeviceConfigMux.Unlock()
		return
	}

	// If neither a failsafe duration nor a failsafe limit were set this message does not pertain to this use case so we accept
	go e.approveOrDenyDeviceConfiguration(msg, true, "")
}

func (e *LPC) AddFeatures() error {
	// client features
	if f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient); f == nil {
		return errors.New("feature not found: DeviceDiagnosis")
	}

	// server features
	f := e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	if f == nil {
		return errors.New("feature not found: LoadControl")
	}

	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	if err := f.AddWriteApprovalCallback(e.loadControlWriteCB); err != nil {
		return err
	}

	measurementId := internal.GetPowerTotalMeasurementId(e.LocalEntity)
	newLimitDesc := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		MeasurementId:  util.Ptr(measurementId),
		Unit:           util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	lc, err := server.NewLoadControl(e.LocalEntity)
	if err != nil {
		return err
	}

	limitId := lc.AddLimitDescription(newLimitDesc)

	newLimiData := []api.LoadControlLimitDataForID{
		{
			Data: model.LoadControlLimitDataType{
				Value:             model.NewScaledNumberType(0),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
			},
			Id: *limitId,
		},
	}
	if err := lc.UpdateLimitDataForIds(newLimiData); err != nil {
		return err
	}

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)
	_ = f.AddWriteApprovalCallback(e.deviceConfigurationWriteCB)

	dcs, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return err
	}

	dcs.AddKeyValueDescription(
		model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
			ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			Unit:      util.Ptr(model.UnitOfMeasurementTypeW),
		},
	)

	// only add if it doesn't exist yet
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
	}
	if data, err := dcs.GetKeyValueDescriptionsForFilter(filter); err == nil && len(data) == 0 {
		dcs.AddKeyValueDescription(
			model.DeviceConfigurationKeyValueDescriptionDataType{
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
		)
	}

	value := &model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(0),
	}
	if err := dcs.UpdateKeyValueDataForFilter(
		model.DeviceConfigurationKeyValueDataType{
			Value:             value,
			IsValueChangeable: util.Ptr(true),
		},
		nil,
		model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
		},
	); err != nil {
		return err
	}

	value = &model.DeviceConfigurationKeyValueValueType{
		Duration: model.NewDurationType(0),
	}
	if err := dcs.UpdateKeyValueDataForFilter(
		model.DeviceConfigurationKeyValueDataType{
			Value:             value,
			IsValueChangeable: util.Ptr(true),
		},
		nil,
		model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
		},
	); err != nil {
		return err
	}

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	f = e.LocalEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, false)

	ec, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		return err
	}
	electricalConnectionId := internal.GetElectricalConnectionId(e.LocalEntity)
	parameterId := internal.GetParameterIdForACPowerTotalMeasurement(e.LocalEntity, electricalConnectionId, measurementId)
	newCharData := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(electricalConnectionId),
		ParameterId:            util.Ptr(parameterId),
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(e.characteristicType()),
		Unit:                   util.Ptr(model.UnitOfMeasurementTypeW),
	}
	if _, err := ec.AddCharacteristic(newCharData); err != nil {
		return err
	}

	return nil
}
