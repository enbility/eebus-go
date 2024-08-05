package mpc

import (
	"github.com/stretchr/testify/assert"
)

func (s *MuMPCSuite) Test_Power() {
	err := s.sut.SetPower(5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_PowerPerPhase() {
	err := s.sut.SetPowerPerPhase(5.0, 5.0, 5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_EnergyConsumed() {
	err := s.sut.SetEnergyConsumed(5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_EnergyProduced() {
	err := s.sut.SetEnergyProduced(5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_CurrentPerPhase() {
	err := s.sut.SetCurrentPerPhase(5.0, 5.0, 5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_VoltagePerPhase() {
	err := s.sut.SetVoltagePerPhase(5.0, 5.0, 5.0)
	assert.Nil(s.T(), err)
}

func (s *MuMPCSuite) Test_Frequency() {
	err := s.sut.SetFrequency(5.0)
	assert.Nil(s.T(), err)
}
