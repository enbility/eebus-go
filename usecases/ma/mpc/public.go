package mpc

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

// return the momentary active power consumption or production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	values, err := internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, nil)
	if err != nil {
		return 0, err
	}
	if len(values) != 1 {
		return 0, api.ErrDataNotAvailable
	}

	return values[0], nil
}

// return the momentary active phase specific power consumption or production per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
}

// Scenario 2

// return the total consumption energy
//
//   - positive values are used for consumption
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
	}
	values, err := measurement.GetDataForFilter(filter)
	if err != nil || len(values) == 0 {
		return 0, api.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if values[0].ValueState != nil && *values[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	return value.GetValue(), nil
}

// return the total feed in energy
//
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
	}
	values, err := measurement.GetDataForFilter(filter)
	if err != nil || len(values) == 0 {
		return 0, api.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if values[0].ValueState != nil && *values[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	return value.GetValue(), nil
}

// Scenario 3

// return the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
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

// Scenario 4

// return the phase specific voltage details
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
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

// Scenario 5

// return frequency
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - and others
func (e *MPC) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	}
	data, err := measurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 || data[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	// therefore we return an error
	if data[0].ValueState != nil && *data[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	// take the first item
	value := data[0].Value

	return value.GetValue(), nil
}
