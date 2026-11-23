package mpc

import (
	"testing"
	"time"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MuMpcBcSuite struct {
	suite.Suite
	*MuMPCSuite
}

// Test suite testing an MPC MonitoredUnit that supports metering for 2 phases (phases AB)
func TestMuMpcAbSuite(t *testing.T) {
	suite.Run(t, new(MuMpcBcSuite))
}

func (s *MuMpcBcSuite) BeforeTest(suiteName, testName string) {
	s.MuMPCSuite = NewMuMPCSuite(
		&s.Suite,
		&MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeBc,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
		&MonitorEnergyConfig{
			ValueSourceProduction:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
		&MonitorVoltageConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeB:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeBc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
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

func (s *MuMpcBcSuite) Test_Power() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerTotal(5.0, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	power, err := s.sut.Power()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, power)
}

func (s *MuMpcBcSuite) Test_PowerPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerPhaseB(6.0, util.Ptr(time.Now()), nil),
		s.sut.UpdateDataPowerPhaseC(7.0, util.Ptr(time.Now()), util.Ptr(model.MeasurementValueStateTypeError)),
	)
	assert.Nil(s.T(), err)
	expectedPowerPerPhase := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeB: 6.0,
		model.ElectricalConnectionPhaseNameTypeC: 7.0,
	}

	powerPerPhases, err := s.sut.PowerPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedPowerPerPhase, powerPerPhases)

	err = s.sut.Update(
		s.sut.UpdateDataPowerPhaseA(5.0, util.Ptr(time.Now()), nil),
	)
	assert.NotNil(s.T(), err)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []float64{6.0, 7.0}, values)
}

func (s *MuMpcBcSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyConsumed(5.0, util.Ptr(time.Now()), nil, util.Ptr(time.Now()), util.Ptr(time.Now())),
	)
	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyConsumed)
}

func (s *MuMpcBcSuite) Test_EnergyProduced() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyProduced(5.0, nil, nil, nil, nil),
	)
	assert.Nil(s.T(), err)

	energyProduced, err := s.sut.EnergyProduced()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyProduced)
}

func (s *MuMpcBcSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataCurrentPhaseB(3.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseC(1.0, nil, nil),
	)
	assert.Nil(s.T(), err)
	expectedCurrentPerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeB: 3.0,
		model.ElectricalConnectionPhaseNameTypeC: 1.0,
	}

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedCurrentPerPhases, currentPerPhases)

	err = s.sut.Update(
		s.sut.UpdateDataCurrentPhaseA(5.0, nil, nil),
	)
	assert.NotNil(s.T(), err)
}

func (s *MuMpcBcSuite) Test_VoltagePerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataVoltagePhaseB(6.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseC(7.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseBToC(9.0, nil, nil),
	)
	assert.Nil(s.T(), err)
	expectedVoltagePerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeB:  6.0,
		model.ElectricalConnectionPhaseNameTypeC:  7.0,
		model.ElectricalConnectionPhaseNameTypeBc: 9.0,
	}

	voltagePerPhases, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedVoltagePerPhases, voltagePerPhases)

	err = s.sut.Update(
		s.sut.UpdateDataVoltagePhaseA(5.0, nil, nil),
	)
	assert.NotNil(s.T(), err)

	err = s.sut.Update(
		s.sut.UpdateDataVoltagePhaseAToB(5.0, nil, nil),
	)
	assert.NotNil(s.T(), err)

	err = s.sut.Update(
		s.sut.UpdateDataVoltagePhaseAToC(5.0, nil, nil),
	)
	assert.NotNil(s.T(), err)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, "", ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []float64{6.0, 7.0, 9.0}, values)
}

func (s *MuMpcBcSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.UpdateDataFrequency(5.0, nil, nil),
	)
	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, frequency)
}
