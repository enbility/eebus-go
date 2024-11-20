package mpc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type MuMpcAbcSuite struct {
	suite.Suite
	*MuMPCSuite
}

func TestMuMpcAbcSuite(t *testing.T) {
	suite.Run(t, new(MuMpcAbcSuite))
}

func (s *MuMpcAbcSuite) BeforeTest(suiteName, testName string) {
	s.MuMPCSuite = NewMuMPCSuite(
		&s.Suite,
		&MonitorPowerConfig{
			ConnectedPhases:   ConnectedPhasesABC,
			ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorVoltageConfig{
			SupportPhaseToPhase:  true,
			ValueSourcePhaseA:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseB:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseC:    util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseAToB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseBToC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseCToA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(1),
			}),
		},
	)
	s.MuMPCSuite.BeforeTest(suiteName, testName)
}

func (s *MuMpcAbcSuite) Test_Power() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerTotal(5.0, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	power, err := s.sut.Power()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, power)
}

func (s *MuMpcAbcSuite) Test_PowerPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerPhaseA(5.0, util.Ptr(time.Now()), nil),
		s.sut.UpdateDataPowerPhaseB(6.0, util.Ptr(time.Now()), nil),
		s.sut.UpdateDataPowerPhaseC(7.0, util.Ptr(time.Now()), util.Ptr(model.MeasurementValueStateTypeError)),
	)
	assert.Nil(s.T(), err)

	powerPerPhases, err := s.sut.PowerPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0}, powerPerPhases)
}

func (s *MuMpcAbcSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyConsumed(5.0, util.Ptr(time.Now()), nil, util.Ptr(time.Now()), util.Ptr(time.Now())),
	)
	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyConsumed)
}

func (s *MuMpcAbcSuite) Test_EnergyProduced() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyProduced(5.0, nil, nil, nil, nil),
	)
	assert.Nil(s.T(), err)

	energyProduced, err := s.sut.EnergyProduced()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyProduced)
}

func (s *MuMpcAbcSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataCurrentPhaseA(5.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseB(3.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseC(1.0, nil, nil),
	)
	assert.Nil(s.T(), err)

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 3.0, 1.0}, currentPerPhases)
}

func (s *MuMpcAbcSuite) Test_VoltagePerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataVoltagePhaseA(5.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseB(6.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseC(7.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseAToB(8.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseBToC(9.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseCToA(10.0, nil, nil),
	)
	assert.Nil(s.T(), err)

	voltagePerPhases, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0}, voltagePerPhases)
}

func (s *MuMpcAbcSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.UpdateDataFrequency(5.0, nil, nil),
	)
	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, frequency)
}
