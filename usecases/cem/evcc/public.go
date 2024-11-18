package evcc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the current charge state of the EV
func (e *EVCC) ChargeState(entity spineapi.EntityRemoteInterface) (ucapi.EVChargeStateType, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return ucapi.EVChargeStateTypeUnplugged, nil
	}

	evDeviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, entity)
	if err != nil {
		return ucapi.EVChargeStateTypeUnplugged, err
	}

	diagnosisState, err := evDeviceDiagnosis.GetState()
	if err != nil {
		return ucapi.EVChargeStateTypeUnknown, err
	}

	operatingState := diagnosisState.OperatingState
	if operatingState == nil {
		return ucapi.EVChargeStateTypeUnknown, api.ErrDataNotAvailable
	}

	switch *operatingState {
	case model.DeviceDiagnosisOperatingStateTypeNormalOperation:
		return ucapi.EVChargeStateTypeActive, nil
	case model.DeviceDiagnosisOperatingStateTypeStandby:
		return ucapi.EVChargeStateTypePaused, nil
	case model.DeviceDiagnosisOperatingStateTypeFailure:
		return ucapi.EVChargeStateTypeError, nil
	case model.DeviceDiagnosisOperatingStateTypeFinished:
		return ucapi.EVChargeStateTypeFinished, nil
	}

	return ucapi.EVChargeStateTypeUnknown, nil
}

// return if an EV is connected
//
// this includes all required features and
// minimal data being available
func (e *EVCC) EVConnected(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil || entity.Device() == nil {
		return false
	}

	// getting current charge state should work
	if _, err := e.ChargeState(entity); err != nil {
		return false
	}

	remoteDevice := e.LocalEntity.Device().RemoteDeviceForSki(entity.Device().Ski())
	if remoteDevice == nil {
		return false
	}

	// check if the device still has an entity assigned with the provided entities address
	return remoteDevice.Entity(entity.Address().Entity) == entity
}

func (e *EVCC) deviceConfigurationValueForKeyName(
	entity spineapi.EntityRemoteInterface,
	keyname model.DeviceConfigurationKeyNameType,
	valueType model.DeviceConfigurationKeyValueTypeType) (*model.DeviceConfigurationKeyValueDataType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	evDeviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	// check if device configuration descriptions has an communication standard key name
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyname,
	}
	if _, err = evDeviceConfiguration.GetKeyValueDescriptionsForFilter(filter); err != nil {
		return nil, err
	}

	filter.ValueType = &valueType
	data, err := evDeviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data, nil
}

// return the current communication standard type used to communicate between EVSE and EV
//
// if an EV is connected via IEC61851, no ISO15118 specific data can be provided!
// sometimes the connection starts with IEC61851 before it switches
// to ISO15118, and sometimes it falls back again. so the error return is
// never absolut for the whole connection time, except if the use case
// is not supported
//
// the values are not constant and can change due to communication problems, bugs, and
// sometimes communication starts with IEC61851 before it switches to ISO
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EVCC) CommunicationStandard(entity spineapi.EntityRemoteInterface) (model.DeviceConfigurationKeyValueStringType, error) {
	unknown := EVCCCommunicationStandardUnknown

	if !e.IsCompatibleEntityType(entity) {
		return unknown, api.ErrNoCompatibleEntity
	}

	data, err := e.deviceConfigurationValueForKeyName(entity, model.DeviceConfigurationKeyNameTypeCommunicationsStandard, model.DeviceConfigurationKeyValueTypeTypeString)
	if err != nil || data == nil || data.Value == nil || data.Value.String == nil {
		return unknown, api.ErrDataNotAvailable
	}

	return *data.Value.String, nil
}

// return if the EV supports asymmetric charging
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
func (e *EVCC) AsymmetricChargingSupport(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	data, err := e.deviceConfigurationValueForKeyName(entity, model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported, model.DeviceConfigurationKeyValueTypeTypeBoolean)
	if err != nil || data == nil || data.Value == nil || data.Value.Boolean == nil {
		return false, api.ErrDataNotAvailable
	}

	return *data.Value.Boolean, nil
}

// return the identifications of the currently connected EV or nil if not available
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EVCC) Identifications(entity spineapi.EntityRemoteInterface) ([]ucapi.IdentificationItem, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	evIdentification, err := client.NewIdentification(e.LocalEntity, entity)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	identifications, err := evIdentification.GetDataForFilter(model.IdentificationDataType{})
	if err != nil {
		return nil, err
	}

	var ids []ucapi.IdentificationItem
	for _, identification := range identifications {
		item := ucapi.IdentificationItem{}

		typ := identification.IdentificationType
		if typ != nil {
			item.ValueType = *typ
		}

		value := identification.IdentificationValue
		if value != nil {
			item.Value = string(*value)
		}

		ids = append(ids, item)
	}

	return ids, nil
}

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *EVCC) ManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	ucapi.ManufacturerData,
	error,
) {
	if !e.IsCompatibleEntityType(entity) {
		return ucapi.ManufacturerData{}, api.ErrNoCompatibleEntity
	}

	return internal.ManufacturerData(e.LocalEntity, entity)
}

// return the minimum, maximum charging and, standby power of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EVCC) ChargingPowerLimits(entity spineapi.EntityRemoteInterface) (float64, float64, float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0.0, 0.0, 0.0, api.ErrNoCompatibleEntity
	}

	evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return 0.0, 0.0, 0.0, api.ErrDataNotAvailable
	}

	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	elParamDesc, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
	if err != nil || len(elParamDesc) == 0 || elParamDesc[0].ParameterId == nil {
		return 0.0, 0.0, 0.0, api.ErrDataNotAvailable
	}

	filter2 := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: elParamDesc[0].ParameterId,
	}
	dataSet, err := evElectricalConnection.GetPermittedValueSetForFilter(filter2)
	if err != nil || len(dataSet) == 0 ||
		dataSet[0].PermittedValueSet == nil ||
		len(dataSet[0].PermittedValueSet) != 1 ||
		dataSet[0].PermittedValueSet[0].Range == nil ||
		len(dataSet[0].PermittedValueSet[0].Range) != 1 {
		return 0.0, 0.0, 0.0, api.ErrDataNotAvailable
	}

	var minValue, maxValue, standByValue float64
	if dataSet[0].PermittedValueSet[0].Range[0].Min != nil {
		minValue = dataSet[0].PermittedValueSet[0].Range[0].Min.GetValue()
	}
	if dataSet[0].PermittedValueSet[0].Range[0].Max != nil {
		maxValue = dataSet[0].PermittedValueSet[0].Range[0].Max.GetValue()
	}
	if len(dataSet[0].PermittedValueSet[0].Value) > 0 {
		standByValue = dataSet[0].PermittedValueSet[0].Value[0].GetValue()
	}

	return minValue, maxValue, standByValue, nil
}

// is the EV in sleep mode
// returns operatingState, lastErrorCode, error
func (e *EVCC) IsInSleepMode(
	entity spineapi.EntityRemoteInterface,
) (bool, error) {
	if !e.IsCompatibleEntityType(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	evseDeviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, entity)
	if err != nil {
		return false, err
	}

	data, err := evseDeviceDiagnosis.GetState()
	if err != nil {
		return false, err
	}

	if data.OperatingState != nil &&
		*data.OperatingState == model.DeviceDiagnosisOperatingStateTypeStandby {
		return true, nil
	}

	return false, nil
}
