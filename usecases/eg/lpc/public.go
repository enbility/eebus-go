package lpc

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	internal "github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the current loadcontrol limit data
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *LPC) ConsumptionLimit(entity spineapi.EntityRemoteInterface) (
	limit ucapi.LoadLimit, resultErr error) {
	limit = ucapi.LoadLimit{
		Value:        0.0,
		IsChangeable: false,
		IsActive:     false,
	}

	resultErr = api.ErrNoCompatibleEntity
	if !e.IsCompatibleEntityType(entity) {
		return
	}

	resultErr = api.ErrDataNotAvailable
	loadControl, err := client.NewLoadControl(e.LocalEntity, entity)
	if err != nil || loadControl == nil {
		return
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}
	limitDescriptions, err := loadControl.GetLimitDescriptionsForFilter(filter)
	if err != nil || len(limitDescriptions) != 1 {
		return
	}

	value, err := loadControl.GetLimitDataForId(*limitDescriptions[0].LimitId)
	if err != nil || value.Value == nil {
		return
	}

	limit.Value = value.Value.GetValue()
	limit.IsChangeable = (value.IsLimitChangeable != nil && *value.IsLimitChangeable)
	limit.IsActive = (value.IsLimitActive != nil && *value.IsLimitActive)
	if value.TimePeriod != nil && value.TimePeriod.EndTime != nil {
		if duration, err := value.TimePeriod.GetDuration(); err == nil {
			limit.Duration = duration
		}
	}

	resultErr = nil

	return
}

// send new LoadControlLimits
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - limit: load limit data
//   - resultCB: callback function for handling the result response
func (e *LPC) WriteConsumptionLimit(
	entity spineapi.EntityRemoteInterface,
	limit ucapi.LoadLimit,
	resultCB func(result model.ResultDataType),
) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	return internal.WriteLoadControlLimit(e.LocalEntity, entity, filter, limit, resultCB)
}

// Scenario 2

// return Failsafe limit for the consumed active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *LPC) FailsafeConsumptionActivePowerLimit(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, api.ErrDataNotAvailable
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   &keyname,
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
	}
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil || data == nil || data.Value == nil || data.Value.ScaledNumber == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data.Value.ScaledNumber.GetValue(), nil
}

// send new Failsafe Consumption Active Power Limit
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - value: the new limit in W
func (e *LPC) WriteFailsafeConsumptionActivePowerLimit(entity spineapi.EntityRemoteInterface, value float64) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil || deviceConfiguration == nil {
		return nil, api.ErrDataNotAvailable
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyname,
	}
	data, err := deviceConfiguration.GetKeyValueDescriptionsForFilter(filter)
	if err != nil || data == nil || len(data) != 1 {
		return nil, api.ErrDataNotAvailable
	}

	keyData := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: data[0].KeyId,
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(value),
			},
		},
	}

	msgCounter, err := deviceConfiguration.WriteKeyValues(keyData)

	return msgCounter, err
}

// return minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
func (e *LPC) FailsafeDurationMinimum(entity spineapi.EntityRemoteInterface) (time.Duration, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, api.ErrDataNotAvailable
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   &keyname,
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
	}
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil || data == nil || data.Value == nil || data.Value.Duration == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data.Value.Duration.GetTimeDuration()
}

// send new Failsafe Duration Minimum
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - duration: the duration, between 2h and 24h
func (e *LPC) WriteFailsafeDurationMinimum(entity spineapi.EntityRemoteInterface, duration time.Duration) (*model.MsgCounterType, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	if duration < time.Duration(time.Hour*2) || duration > time.Duration(time.Hour*24) {
		return nil, errors.New("duration outside of allowed range")
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil || deviceConfiguration == nil {
		return nil, api.ErrDataNotAvailable
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyname,
	}
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil || data == nil {
		return nil, api.ErrDataNotAvailable
	}

	keyData := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: data.KeyId,
			Value: &model.DeviceConfigurationKeyValueValueType{
				Duration: model.NewDurationType(duration),
			},
		},
	}

	msgCounter, err := deviceConfiguration.WriteKeyValues(keyData)

	return msgCounter, err
}

// Scenario 3

// start sending heartbeat from the local entity supporting this usecase
//
// the heartbeat is started by default when a non 0 timeout is set in the service configuration
func (e *LPC) StartHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		_ = hm.StartHeartbeat()
	}
}

// stop sending heartbeat from the local CEM entity
func (e *LPC) StopHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		hm.StopHeartbeat()
	}
}

// check wether there was a heartbeat received within the last 2 minutes
//
// returns true, if the last heartbeat is within 2 minutes, otherwise false
func (e *LPC) IsHeartbeatWithinDuration(entity spineapi.EntityRemoteInterface) bool {
	lf, err := client.NewDeviceDiagnosis(e.LocalEntity, entity)
	if err != nil {
		return false
	}

	return lf.IsHeartbeatWithinDuration(2 * time.Minute)
}

// Scenario 4

// return nominal maximum active (real) power the Controllable System is
// able to consume according to the contract (EMS), device label or data sheet.
func (e *LPC) ConsumptionNominalMax(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil || electricalConnection == nil {
		return 0, err
	}

	filter := model.ElectricalConnectionCharacteristicDataType{
		CharacteristicContext: util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:    util.Ptr(e.characteristicType(entity)),
	}
	data, err := electricalConnection.GetCharacteristicsForFilter(filter)
	if err != nil {
		return 0, err
	} else if len(data) == 0 || data[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data[0].Value.GetValue(), nil
}

// returns the characteristictype depending on the remote entities device devicetype
func (e *LPC) characteristicType(entity spineapi.EntityRemoteInterface) model.ElectricalConnectionCharacteristicTypeType {
	// According to LPC V1.0 2.2, lines 400ff:
	// - a HEMS provides contractual consumption nominal max
	// - any other devices provides power consupmtion nominal max
	characteristic := model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax

	if entity == nil || entity.Device() == nil {
		return characteristic
	}

	deviceType := entity.Device().DeviceType()
	if deviceType == nil || *deviceType == model.DeviceTypeTypeEnergyManagementSystem {
		characteristic = model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax
	}

	return characteristic
}
