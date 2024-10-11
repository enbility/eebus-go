package mpc

import (
	"github.com/stretchr/testify/assert"
)

func (s *MuMPCSuite) Test_Power() {
	err := s.sut.Update(
		s.sut.MeasuredAcPowerTotal(5.0),
	)
	assert.Nil(s.T(), err)

	power, err := s.sut.Power()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, power)
}

func (s *MuMPCSuite) Test_PowerPerPhase() {
	err := s.sut.Update(
		s.sut.MeasuredAcPowerPhaseA(5.0),
		s.sut.MeasuredAcPowerPhaseB(6.0),
		s.sut.MeasuredAcPowerPhaseC(7.0),
	)
	assert.Nil(s.T(), err)

	powerPerPhases, err := s.sut.PowerPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0}, powerPerPhases)
}

func (s *MuMPCSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.MeasuredAcEnergyConsumed(5.0),
	)
	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyConsumed)
}

func (s *MuMPCSuite) Test_EnergyProduced() {
	err := s.sut.Update(
		s.sut.MeasuredAcEnergyProduced(5.0),
	)
	assert.Nil(s.T(), err)

	energyProduced, err := s.sut.EnergyProduced()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyProduced)
}

func (s *MuMPCSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.MeasuredAcCurrentPhaseA(5.0),
		s.sut.MeasuredAcCurrentPhaseB(5.0),
		s.sut.MeasuredAcCurrentPhaseC(5.0),
	)
	assert.Nil(s.T(), err)

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 5.0, 5.0}, currentPerPhases)
}

func (s *MuMPCSuite) Test_VoltagePerPhase() {
	err := s.sut.Update(
		s.sut.MeasuredAcVoltagePhaseA(5.0),
		s.sut.MeasuredAcVoltagePhaseB(6.0),
		s.sut.MeasuredAcVoltagePhaseC(7.0),
		s.sut.MeasuredAcVoltagePhaseAToB(8.0),
		s.sut.MeasuredAcVoltagePhaseBToC(9.0),
		s.sut.MeasuredAcVoltagePhaseCToA(10.0),
	)
	assert.Nil(s.T(), err)

	voltagePerPhases, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0}, voltagePerPhases)
}

func (s *MuMPCSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.MeasuredAcFrequency(5.0),
	)
	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, frequency)
}
