package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

func (e *MGCP) SetValue(value float64, mid model.MeasurementIdType) error {
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

func (e *MGCP) SetValues(values []float64, mids []model.MeasurementIdType) error {
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

// Scenario 2

// set the current power limitation factor
func (e *MGCP) SetPowerLimitationFactor(value float64) error {
	dcs, err := server.NewDeviceConfiguration(e.LocalEntity)

	if err != nil {
		return err
	}

	val := &model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(value),
	}
	_ = dcs.UpdateKeyValueDataForFilter(
		model.DeviceConfigurationKeyValueDataType{
			Value:             val,
			IsValueChangeable: util.Ptr(true),
		},
		nil,
		model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
		},
	)

	return nil
}

// set the current power
func (e *MGCP) SetPower(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(0))
}

// Scenario 3

// set the EnergyFeedIn
func (e *MGCP) SetEnergyFeedIn(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(1))
}

// Scenario 4

// set the EnergyConsumed
func (e *MGCP) SetEnergyConsumed(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(2))
}

// // Scenario 5

// set the CurrentPerPhase
func (e *MGCP) SetCurrentPerPhase(values []float64) error {
	return e.SetValues(values, []model.MeasurementIdType{model.MeasurementIdType(3), model.MeasurementIdType(4), model.MeasurementIdType(5)})
}

// // Scenario 6

// set the VoltagePerPhase
func (e *MGCP) SetVoltagePerPhase(values []float64) error {
	return e.SetValues(values, []model.MeasurementIdType{model.MeasurementIdType(6), model.MeasurementIdType(7), model.MeasurementIdType(8)})
}

// Scenario 7

// set the Frequency
func (e *MGCP) SetFrequency(value float64) error {
	return e.SetValue(value, model.MeasurementIdType(9))
}
