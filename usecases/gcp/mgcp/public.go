package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
)

// -------- Getters -------- //

// -------- Setters -------- //

// Scenario 1

// Use MGCP.ConfigurationPvFeedInLimitationFactor in MGCP.Update to set the current power limitation factor
func (m *MGCP) ConfigurationPvFeedInLimitationFactor(pvFeedInLimitationFactor float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeConfiguration,
		updateTypeConfiguration: model.DeviceConfigurationKeyValueDataType{
			KeyId: m.pvFeedInLimitationFactor,
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(pvFeedInLimitationFactor),
			},
		},
	}
}

// Scenario 2

// Use MGCP.MeasurementAcPowerTotal in MGCP.Update to set the current total power
func (m *MGCP) MeasurementAcPowerTotal(acPowerTotal float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(acPowerTotal),
			Id:   *m.acPowerTotal,
		},
	}
}

// Scenario 3

// Use MGCP.MeasurementAcEnergyFeedIn in MGCP.Update to set the total feed in energy
func (m *MGCP) MeasurementAcEnergyFeedIn(energy float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(energy),
			Id:   *m.gridFeedIn,
		},
	}
}

// Scenario 4

// Use MGCP.MeasurementAcEnergyConsumed in MGCP.Update to set the total feed in energy
func (m *MGCP) MeasurementAcEnergyConsumed(energy float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(energy),
			Id:   *m.gridConsumption,
		},
	}
}

// Scenario 5

// Use MGCP.MeasurementAcCurrentPhaseA in MGCP.Update to set the current of phase A
func (m *MGCP) MeasurementAcCurrentPhaseA(current float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(current),
			Id:   *m.acCurrent[0],
		},
	}
}

// Use MGCP.MeasurementAcCurrentPhaseB in MGCP.Update to set the current of phase B
func (m *MGCP) MeasurementAcCurrentPhaseB(current float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(current),
			Id:   *m.acCurrent[1],
		},
	}
}

// Use MGCP.MeasurementAcCurrentPhaseC in MGCP.Update to set the current of phase C
func (m *MGCP) MeasurementAcCurrentPhaseC(current float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(current),
			Id:   *m.acCurrent[2],
		},
	}
}

// Scenario 6

// Use MGCP.MeasurementAcVoltagePhaseA in MGCP.Update to set the voltage of phase A
func (m *MGCP) MeasurementAcVoltagePhaseA(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[0],
		},
	}
}

// Use MGCP.MeasurementAcVoltagePhaseB in MGCP.Update to set the voltage of phase B
func (m *MGCP) MeasurementAcVoltagePhaseB(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[1],
		},
	}
}

// Use MGCP.MeasurementAcVoltagePhaseC in MGCP.Update to set the voltage of phase C
func (m *MGCP) MeasurementAcVoltagePhaseC(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[2],
		},
	}
}

// Use MGCP.MeasurementAcVoltagePhaseAToB in MGCP.Update to set the voltage between phase A and B
func (m *MGCP) MeasurementAcVoltagePhaseAToB(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[3],
		},
	}
}

// Use MGCP.MeasurementAcVoltagePhaseBToC in MGCP.Update to set the voltage between phase B and C
func (m *MGCP) MeasurementAcVoltagePhaseBToC(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[4],
		},
	}
}

// Use MGCP.MeasurementAcVoltagePhaseCToA in MGCP.Update to set the voltage between phase C and A
func (m *MGCP) MeasurementAcVoltagePhaseCToA(voltage float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(voltage),
			Id:   *m.acVoltage[5],
		},
	}
}

// Scenario 7

// Use MGCP.MeasurementAcFrequency in MGCP.Update to set the frequency
func (m *MGCP) MeasurementAcFrequency(frequency float64) UpdateValueType {
	return UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Data: measuredValue(frequency),
			Id:   *m.acFrequency,
		},
	}
}

// Update the data

// use MPC.Update to update the data of the MGCP Usecase
// use it like this:
//
//	mgcp.Update(
//	  mgcp.MeasuredAcPowerTotal(1000),
//	  mgcp.MeasuredAcPowerPhaseA(500),
//	  ...
//	)
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (m *MGCP) Update(updateValueType ...UpdateValueType) []error {
	measurements := make([]api.MeasurementDataForID, 0)
	configurations := make([]model.DeviceConfigurationKeyValueDataType, 0)
	errors := make([]error, 0)

	for _, update := range updateValueType {
		if update.updateValueType == SupportedUpdateValueTypeMeasurement {
			measurements = append(measurements, update.updateTypeMeasurement)
		} else {
			configurations = append(configurations, update.updateTypeConfiguration)
		}
	}

	if len(measurements) > 0 {
		_measurements, err := server.NewMeasurement(m.LocalEntity)
		if err != nil {
			errors = append(errors, err)
		}

		err = _measurements.UpdateDataForIds(measurements)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(configurations) == 1 {
		_configurations, err := server.NewDeviceConfiguration(m.LocalEntity)
		if err != nil {
			errors = append(errors, err)
		}

		err = _configurations.UpdateKeyValueDataForKeyId(configurations[0], nil, *configurations[0].KeyId)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
