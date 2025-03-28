package mpc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

func (e *MPC) SetValue(value float64, mid model.MeasurementIdType) error {
	measure, err := server.NewMeasurement(e.LocalEntity)

	if err != nil {
		return err
	}

	ids := []api.MeasurementDataForID{
		{
			Id: mid,
			Data: model.MeasurementDataType{
				Value:     model.NewScaledNumberType(value),
				ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
			},
		},
	}

	return measure.UpdateDataForIds(ids)
}

func (e *MPC) SetValues(values []float64, mids []model.MeasurementIdType) error {
	measure, err := server.NewMeasurement(e.LocalEntity)

	if err != nil {
		return err
	}

	ids := []api.MeasurementDataForID{
		{
			Id: mids[0],
			Data: model.MeasurementDataType{
				Value:     model.NewScaledNumberType(values[0]),
				ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
			},
		},
		{
			Id: mids[1],
			Data: model.MeasurementDataType{
				Value:     model.NewScaledNumberType(values[1]),
				ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
			},
		},
		{
			Id: mids[2],
			Data: model.MeasurementDataType{
				Value:     model.NewScaledNumberType(values[2]),
				ValueType: util.Ptr(model.MeasurementValueTypeTypeValue),
			},
		},
	}

	return measure.UpdateDataForIds(ids)
}

// Scenario 1

// set the current power
func (e *MPC) SetPower(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(0))
}

// set the PowerPerPhase
func (e *MPC) SetPowerPerPhase(values []float64) error {
	return e.SetValues(values, []model.MeasurementIdType{model.MeasurementIdType(1), model.MeasurementIdType(2), model.MeasurementIdType(3)})
}

// Scenario 2

// set the EnergyConsumed
func (e *MPC) SetEnergyConsumed(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(4))
}

// set the EnergyProduced
func (e *MPC) SetEnergyProduced(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(5))
}

// Scenario 3

// set the CurrentPerPhase
func (e *MPC) SetCurrentPerPhase(values []float64) error {
	return e.SetValues(values, []model.MeasurementIdType{model.MeasurementIdType(6), model.MeasurementIdType(7), model.MeasurementIdType(8)})
}

// Scenario 4

// set the VoltagePerPhase
func (e *MPC) SetVoltagePerPhase(values []float64) error {
	return e.SetValues(values, []model.MeasurementIdType{model.MeasurementIdType(9), model.MeasurementIdType(10), model.MeasurementIdType(11)})
}

// Scenario 5

// set the Frequency
func (e *MPC) SetFrequency(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(12))
}
