package lpc

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the current loadcontrol limit data
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *LPC) ConsumptionLimit() (limit ucapi.LoadLimit, resultErr error) {
	limit = ucapi.LoadLimit{
		Value:        0.0,
		IsChangeable: false,
		IsActive:     false,
		Duration:     0,
	}
	resultErr = api.ErrDataNotAvailable

	lc, limidId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return limit, err
	}

	value, err := lc.GetLimitDataForId(limidId)
	if err != nil || value == nil || value.LimitId == nil || value.Value == nil {
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

	return limit, nil
}

// set the current loadcontrol limit data
func (e *LPC) SetConsumptionLimit(limit ucapi.LoadLimit) (resultErr error) {
	loadControlf, limidId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return err
	}

	limitData := []api.LoadControlLimitDataForFilter{
		{
			Data: model.LoadControlLimitDataType{
				IsLimitChangeable: util.Ptr(limit.IsChangeable),
				IsLimitActive:     util.Ptr(limit.IsActive),
				Value:             model.NewScaledNumberType(limit.Value),
			},
			Filter: model.LoadControlLimitDescriptionDataType{
				LimitId: &limidId,
			},
		},
	}
	if limit.Duration > 0 {
		limitData[0].Data.TimePeriod = &model.TimePeriodType{
			EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(limit.Duration),
		}
	}

	deleteSelector := &model.LoadControlLimitListDataSelectorsType{
		LimitId: &limidId,
	}

	deleteTimePeriod := &model.LoadControlLimitDataElementsType{
		TimePeriod: util.Ptr(model.TimePeriodElementsType{}),
	}

	return loadControlf.UpdateLimitDataForFilters(limitData, deleteSelector, deleteTimePeriod)
}

// return the currently pending incoming consumption write limits
func (e *LPC) PendingConsumptionLimits() map[model.MsgCounterType]ucapi.LoadLimit {
	result := make(map[model.MsgCounterType]ucapi.LoadLimit)

	_, limitId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return result
	}

	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	for key, msg := range e.pendingLimits {
		data := msg.Cmd.LoadControlLimitListData

		// elements are only added to the map if all required fields exist
		// therefor no checks for these are needed here

		// find the item which contains the limit for this usecase
		for _, item := range data.LoadControlLimitData {
			if item.LimitId == nil ||
				limitId != *item.LimitId {
				continue
			}

			limit := ucapi.LoadLimit{}

			if item.TimePeriod != nil {
				if duration, err := item.TimePeriod.GetDuration(); err == nil {
					limit.Duration = duration
				}
			}

			if item.IsLimitActive != nil {
				limit.IsActive = *item.IsLimitActive
			}

			if item.Value != nil {
				limit.Value = item.Value.GetValue()
			}

			result[key] = limit
		}
	}

	return result
}

// accept or deny an incoming consumption write limit
//
// use PendingConsumptionLimits to get the list of currently pending requests
func (e *LPC) ApproveOrDenyConsumptionLimit(msgCounter model.MsgCounterType, approve bool, reason string) {
	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	msg, ok := e.pendingLimits[msgCounter]
	if !ok {
		// no pending limit for this msgCounter, this is a caller error
		return
	}

	e.approveOrDenyConsumptionLimit(msg, approve, reason)

	delete(e.pendingLimits, msgCounter)
}

// Scenario 2

// return Failsafe limit for the consumed active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *LPC) FailsafeConsumptionActivePowerLimit() (limit float64, isChangeable bool, resultErr error) {
	limit = 0
	isChangeable = false
	resultErr = api.ErrDataNotAvailable

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
	}
	keyData, err := dc.GetKeyValueDataForFilter(filter)
	if err != nil || keyData == nil || keyData.KeyId == nil || keyData.Value == nil || keyData.Value.ScaledNumber == nil {
		return
	}

	limit = keyData.Value.ScaledNumber.GetValue()
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set Failsafe limit for the consumed active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *LPC) SetFailsafeConsumptionActivePowerLimit(value float64, changeable bool) error {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit
	keyValue := model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(value),
	}

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return err
	}

	data := model.DeviceConfigurationKeyValueDataType{
		Value:             &keyValue,
		IsValueChangeable: &changeable,
	}
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(keyName),
	}
	return dc.UpdateKeyValueDataForFilter(data, nil, filter)
}

// return minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
func (e *LPC) FailsafeDurationMinimum() (duration time.Duration, isChangeable bool, resultErr error) {
	duration = 0
	isChangeable = false
	resultErr = api.ErrDataNotAvailable

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
	}
	keyData, err := dc.GetKeyValueDataForFilter(filter)
	if err != nil || keyData == nil || keyData.KeyId == nil || keyData.Value == nil || keyData.Value.Duration == nil {
		return
	}

	durationValue, err := keyData.Value.Duration.GetTimeDuration()
	if err != nil {
		return
	}

	duration = durationValue
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
//
// parameters:
//   - duration: has to be >= 2h and <= 24h
//   - changeable: boolean if the client service can change this value
func (e *LPC) SetFailsafeDurationMinimum(duration time.Duration, changeable bool) error {
	if duration < time.Duration(time.Hour*2) || duration > time.Duration(time.Hour*24) {
		return errors.New("duration outside of allowed range")
	}
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum
	keyValue := model.DeviceConfigurationKeyValueValueType{
		Duration: model.NewDurationType(duration),
	}

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return err
	}

	data := model.DeviceConfigurationKeyValueDataType{
		Value:             &keyValue,
		IsValueChangeable: &changeable,
	}
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(keyName),
	}
	return dc.UpdateKeyValueDataForFilter(data, nil, filter)
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

func (e *LPC) IsHeartbeatWithinDuration() bool {
	if e.heartbeatDiag == nil {
		return false
	}

	return e.heartbeatDiag.IsHeartbeatWithinDuration(2 * time.Minute)
}

// Scenario 4

// return nominal maximum active (real) power the Controllable System is allowed to consume.
//
// If the local device type is an EnergyManagementSystem, the contractual consumption
// nominal max is returned, otherwise the power consumption nominal max is returned.
func (e *LPC) ConsumptionNominalMax() (value float64, resultErr error) {
	value = 0
	resultErr = api.ErrDataNotAvailable

	ec, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		resultErr = err
		return
	}

	filter := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(e.characteristicType()),
	}
	charData, err := ec.GetCharacteristicsForFilter(filter)
	if err != nil || len(charData) == 0 ||
		charData[0].CharacteristicId == nil ||
		charData[0].Value == nil {
		return
	}

	return charData[0].Value.GetValue(), nil
}

// set nominal maximum active (real) power the Controllable System is allowed to consume.
//
// If the local device type is an EnergyManagementSystem, the contractual consumption
// nominal max is set, otherwise the power consumption nominal max is set.
//
// parameters:
//   - value: nominal max power consumption in W
func (e *LPC) SetConsumptionNominalMax(value float64) error {
	ec, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		return err
	}

	electricalConnectionid := util.Ptr(model.ElectricalConnectionIdType(0))
	parameterId := util.Ptr(model.ElectricalConnectionParameterIdType(0))
	charList, err := ec.GetCharacteristicsForFilter(model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: electricalConnectionid,
		ParameterId:            parameterId,
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(e.characteristicType()),
	})
	if err != nil || len(charList) == 0 {
		return api.ErrDataNotAvailable
	}

	data := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: electricalConnectionid,
		ParameterId:            parameterId,
		CharacteristicId:       charList[0].CharacteristicId,
		Value:                  model.NewScaledNumberType(value),
	}
	return ec.UpdateCharacteristic(data, nil)
}

// returns the characteristictype depending on the local entities device devicetype
func (e *LPC) characteristicType() model.ElectricalConnectionCharacteristicTypeType {
	deviceType := e.LocalEntity.Device().DeviceType()

	// According to LPC V1.0 2.2, lines 400ff:
	// - a HEMS provides contractual consumption nominal max
	// - any other devices provides power consupmtion nominal max
	characteristic := model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax
	if deviceType == nil || *deviceType == model.DeviceTypeTypeEnergyManagementSystem {
		characteristic = model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax
	}

	return characteristic
}
