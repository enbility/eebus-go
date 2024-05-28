package mpc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the momentary active power consumption or production
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *MPC) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
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
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *MPC) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, internal.PhaseNameMapping)
}

// Scenario 2

// return the total consumption energy
//
//   - positive values are used for consumption
func (e *MPC) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
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

	return value.GetValue(), nil
}

// return the total feed in energy
//
//   - negative values are used for production
func (e *MPC) EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
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

	return value.GetValue(), nil
}

// Scenario 3

// return the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *MPC) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, model.EnergyDirectionTypeConsume, internal.PhaseNameMapping)
}

// Scenario 4

// return the phase specific voltage details
func (e *MPC) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}
	return internal.MeasurementPhaseSpecificDataForFilter(e.LocalEntity, entity, filter, "", internal.PhaseNameMapping)
}

// Scenario 5

// return frequency
func (e *MPC) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !e.IsCompatibleEntity(entity) {
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

	// take the first item
	value := data[0].Value

	return value.GetValue(), nil
}
