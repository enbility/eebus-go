package mpc

import "github.com/stretchr/testify/assert"

func (s *MuMPCSuite) Test_Power() {
	power, err := s.sut.Power()
	assert.Equal(s.T(), 0.0, power)
	assert.Nil(s.T(), err)

	err = s.sut.SetPower(17.0)
	assert.Nil(s.T(), err)

	power, err = s.sut.Power()
	assert.Equal(s.T(), 17.0, power)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_PowerPerPhase() {
	phases, err := s.sut.PowerPerPhase()

	assert.Equal(s.T(), 3, len(phases))
	assert.Equal(s.T(), 0.0, phases[0])
	assert.Equal(s.T(), 0.0, phases[1])
	assert.Equal(s.T(), 0.0, phases[2])
	assert.Nil(s.T(), err)

	err = s.sut.SetPowerPerPhase(17.0, 18.0, 19.0)
	assert.Nil(s.T(), err)

	phases, err = s.sut.PowerPerPhase()
	assert.Equal(s.T(), 3, len(phases))
	assert.Equal(s.T(), 17.0, phases[0])
	assert.Equal(s.T(), 18.0, phases[1])
	assert.Equal(s.T(), 19.0, phases[2])
}

func (s *MuMPCSuite) Test_EnergyConsumed() {
	energy, err := s.sut.EnergyConsumed()
	assert.Equal(s.T(), 0.0, energy)
	assert.Nil(s.T(), err)

	err = s.sut.SetEnergyConsumed(17.0)
	assert.Nil(s.T(), err)

	energy, err = s.sut.EnergyConsumed()
	assert.Equal(s.T(), 17.0, energy)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_EnergyProduced() {
	energy, err := s.sut.EnergyProduced()
	assert.Equal(s.T(), 0.0, energy)
	assert.Nil(s.T(), err)

	err = s.sut.SetEnergyProduced(17.0)
	assert.Nil(s.T(), err)

	energy, err = s.sut.EnergyProduced()
	assert.Equal(s.T(), 17.0, energy)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_CurrentPerPhase() {
	current, err := s.sut.CurrentPerPhase()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(current))
	assert.Equal(s.T(), 0.0, current[0])
	assert.Equal(s.T(), 0.0, current[1])
	assert.Equal(s.T(), 0.0, current[2])

	err = s.sut.SetCurrentPerPhase(17.0, 18.0, 19.0)
	assert.Nil(s.T(), err)

	current, err = s.sut.CurrentPerPhase()
	assert.Equal(s.T(), 3, len(current))
	assert.Equal(s.T(), 17.0, current[0])
	assert.Equal(s.T(), 18.0, current[1])
	assert.Equal(s.T(), 19.0, current[2])
}

func (s *MuMPCSuite) Test_VoltagePerPhase() {
	voltages, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 3, len(voltages))
	assert.Equal(s.T(), 0.0, voltages[0])
	assert.Equal(s.T(), 0.0, voltages[1])
	assert.Equal(s.T(), 0.0, voltages[2])

	err = s.sut.SetVoltagePerPhase(1.0, 2.0, 3.0)
	assert.Nil(s.T(), err)

	voltages, err = s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 3, len(voltages))
	assert.Equal(s.T(), 1.0, voltages[0])
	assert.Equal(s.T(), 2.0, voltages[1])
	assert.Equal(s.T(), 3.0, voltages[2])
}

func (s *MuMPCSuite) Test_Frequency() {
	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0.0, frequency)

	err = s.sut.SetFrequency(50.0)
	assert.Nil(s.T(), err)

	frequency, err = s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, frequency)
}
