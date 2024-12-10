package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	internal "github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the current power limitation factor
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) PowerLimitationFactor(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, err
	}

	// check if device configuration description has curtailment limit factor key name
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyname,
	}
	_, err = deviceConfiguration.GetKeyValueDescriptionsForFilter(filter)
	if err != nil {
		return 0, err
	}

	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil || data == nil || data.Value == nil || data.Value.ScaledNumber == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data.Value.ScaledNumber.GetValue(), nil
}

// Scenario 2

// return the momentary power consumption or production at the grid connection point
//
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	data, err := internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, nil)
	if err != nil {
		return 0, err
	}

	if len(data) != 1 {
		return 0, api.ErrDataNotAvailable
	}

	return data[0], nil
}

// Scenario 3

// return the total feed in energy at the grid connection point
//
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) EnergyFeedIn(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil || measurement == nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
	}
	result, err := measurement.GetDataForFilter(filter)
	if err != nil || len(result) == 0 || result[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if result[0].ValueState != nil && *result[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	return result[0].Value.GetValue(), nil
}

// Scenario 4

// return the total consumption energy at the grid connection point
//
//   - positive values are used for consumption
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil || measurement == nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
	}
	result, err := measurement.GetDataForFilter(filter)
	if err != nil || len(result) == 0 || result[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if result[0].ValueState != nil && *result[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	return result[0].Value.GetValue(), nil
}

// Scenario 5

// return the momentary current consumption or production at the grid connection point
//
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
}

// Scenario 6

// return the voltage phase details at the grid connection point
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, "", ucapi.PhaseNameMapping)
}

// Scenario 7

// return frequency at the grid connection point
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MGCP) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil || measurement == nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	}
	result, err := measurement.GetDataForFilter(filter)
	if err != nil || len(result) == 0 || result[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if result[0].ValueState != nil && *result[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	return result[0].Value.GetValue(), nil
}
