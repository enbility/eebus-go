package mpc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
)

// ------------------------- Getters ------------------------- //

// Scenario 1

// get the momentary active power consumption or production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Power() (float64, error) {
	if e.acPowerTotal == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acPowerTotal)
}

// get the momentary active power consumption or production per phase
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) PowerPerPhase() ([]float64, error) {
	if e.acPower[0] == nil || e.acPower[1] == nil || e.acPower[2] == nil {
		return nil, api.ErrMissingData
	}

	phaseA, err := e.getMeasurementDataForId(e.acPower[0])
	if err != nil {
		return nil, err
	}

	phaseB, err := e.getMeasurementDataForId(e.acPower[1])
	if err != nil {
		return nil, err
	}

	phaseC, err := e.getMeasurementDataForId(e.acPower[2])
	if err != nil {
		return nil, err
	}

	return []float64{phaseA, phaseB, phaseC}, nil
}

// Scenario 2

// get the total feed in energy
//
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) EnergyConsumed() (float64, error) {
	if e.acEnergyConsumed == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acEnergyConsumed)
}

// get the total feed in energy
//
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) EnergyProduced() (float64, error) {
	if e.acEnergyProduced == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acEnergyProduced)
}

// Scenario 3

// get the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) CurrentPerPhase() ([]float64, error) {
	if e.acCurrent[0] == nil || e.acCurrent[1] == nil || e.acCurrent[2] == nil {
		return nil, api.ErrMissingData
	}

	phaseA, err := e.getMeasurementDataForId(e.acCurrent[0])
	if err != nil {
		return nil, err
	}

	phaseB, err := e.getMeasurementDataForId(e.acCurrent[1])
	if err != nil {
		return nil, err
	}

	phaseC, err := e.getMeasurementDataForId(e.acCurrent[2])
	if err != nil {
		return nil, err
	}

	return []float64{phaseA, phaseB, phaseC}, nil
}

// Scenario 4

// get the phase specific voltage details
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) VoltagePerPhase() ([]float64, error) {
	for _, id := range e.acVoltage {
		if id == nil {
			return nil, api.ErrMissingData
		}
	}

	phaseA, err := e.getMeasurementDataForId(e.acVoltage[0])
	if err != nil {
		return nil, err
	}

	phaseB, err := e.getMeasurementDataForId(e.acVoltage[1])
	if err != nil {
		return nil, err
	}

	phaseC, err := e.getMeasurementDataForId(e.acVoltage[2])
	if err != nil {
		return nil, err
	}

	phaseAToB, err := e.getMeasurementDataForId(e.acVoltage[3])
	if err != nil {
		return nil, err
	}

	phaseBToC, err := e.getMeasurementDataForId(e.acVoltage[4])
	if err != nil {
		return nil, err
	}

	phaseCToA, err := e.getMeasurementDataForId(e.acVoltage[5])
	if err != nil {
		return nil, err
	}

	return []float64{phaseA, phaseB, phaseC, phaseAToB, phaseBToC, phaseCToA}, nil
}

// Scenario 5

// get frequency
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Frequency() (float64, error) {
	if e.acFrequency == nil {
		return 0, api.ErrMissingData
	}

	return e.getMeasurementDataForId(e.acFrequency)
}

// ------------------------- Setters ------------------------- //

// Scenario 1

// use MPC.MeasuredAcPowerTotal in MPC.Update to set the momentary active power consumption or production
func (e *MPC) MeasuredAcPowerTotal(acPowerTotal float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acPowerTotal),
		Id:   *e.acPowerTotal,
	}
}

// use MPC.MeasuredAcPowerPhaseA in MPC.Update to set the momentary active power consumption or production per phase
func (e *MPC) MeasuredAcPowerPhaseA(acPowerPhaseA float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acPowerPhaseA),
		Id:   *e.acPower[0],
	}
}

// use MPC.MeasuredAcPowerPhaseB in MPC.Update to set the momentary active power consumption or production per phase
func (e *MPC) MeasuredAcPowerPhaseB(acPowerPhaseB float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acPowerPhaseB),
		Id:   *e.acPower[1],
	}
}

// use MPC.MeasuredAcPowerPhaseC in MPC.Update to set the momentary active power consumption or production per phase
func (e *MPC) MeasuredAcPowerPhaseC(acPowerPhaseC float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acPowerPhaseC),
		Id:   *e.acPower[2],
	}
}

// Scenario 2

// use MPC.MeasuredAcEnergyConsumed in MPC.Update to set the total feed in energy
func (e *MPC) MeasuredAcEnergyConsumed(energyConsumed float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(energyConsumed),
		Id:   *e.acEnergyConsumed,
	}
}

// use MPC.MeasuredAcEnergyProduced in MPC.Update to set the total feed in energy
func (e *MPC) MeasuredAcEnergyProduced(energyProduced float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(energyProduced),
		Id:   *e.acEnergyProduced,
	}
}

// Scenario 3

// use MPC.MeasuredAcCurrentPhaseA in MPC.Update to set the momentary phase specific current consumption or production
func (e *MPC) MeasuredAcCurrentPhaseA(acCurrentPhaseA float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acCurrentPhaseA),
		Id:   *e.acCurrent[0],
	}
}

// use MPC.MeasuredAcCurrentPhaseB in MPC.Update to set the momentary phase specific current consumption or production
func (e *MPC) MeasuredAcCurrentPhaseB(acCurrentPhaseB float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acCurrentPhaseB),
		Id:   *e.acCurrent[1],
	}
}

// use MPC.MeasuredAcCurrentPhaseC in MPC.Update to set the momentary phase specific current consumption or production
func (e *MPC) MeasuredAcCurrentPhaseC(acCurrentPhaseC float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(acCurrentPhaseC),
		Id:   *e.acCurrent[2],
	}
}

// Scenario 4

// use MPC.MeasuredAcVoltagePhaseA in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseA(voltagePhaseA float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseA),
		Id:   *e.acVoltage[0],
	}
}

// use MPC.MeasuredAcVoltagePhaseB in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseB(voltagePhaseB float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseB),
		Id:   *e.acVoltage[1],
	}
}

// use MPC.MeasuredAcVoltagePhaseC in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseC(voltagePhaseC float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseC),
		Id:   *e.acVoltage[2],
	}
}

// use MPC.MeasuredAcVoltagePhaseAToB in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseAToB(voltagePhaseAToB float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseAToB),
		Id:   *e.acVoltage[3],
	}
}

// use MPC.MeasuredAcVoltagePhaseBToC in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseBToC(voltagePhaseBToC float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseBToC),
		Id:   *e.acVoltage[4],
	}
}

// use MPC.MeasuredAcVoltagePhaseCToA in MPC.Update to set the phase specific voltage details
func (e *MPC) MeasuredAcVoltagePhaseCToA(voltagePhaseCToA float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(voltagePhaseCToA),
		Id:   *e.acVoltage[5],
	}
}

// Scenario 5

// use MPC.MeasuredAcFrequency in MPC.Update to set the frequency
func (e *MPC) MeasuredAcFrequency(frequency float64) api.MeasurementDataForID {
	return api.MeasurementDataForID{
		Data: measuredValue(frequency),
		Id:   *e.acFrequency,
	}
}

// Update the measurement data

// use MPC.Update to update the measurement data
// use it like this:
//
//	mpc.Update(
//	  mpc.MeasuredAcPowerTotal(1000),
//	  mpc.MeasuredAcPowerPhaseA(500),
//	  ...
//	)
//
// possible errors:
//   - ErrMissingData if the id is not available
//   - and others
func (e *MPC) Update(measurementDataForIds ...api.MeasurementDataForID) error {
	measurements, err := server.NewMeasurement(e.LocalEntity)
	if err != nil {
		return err
	}

	return measurements.UpdateDataForIds(measurementDataForIds)
}
