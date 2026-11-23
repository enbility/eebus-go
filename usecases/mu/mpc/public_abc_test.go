package mpc

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MuMpcAbcSuite struct {
	suite.Suite
	*MuMPCSuite
}

// Test suite testing an MPC MonitoredUnit that supports 3-phase metering (phases ABC)
func TestMuMpcAbcSuite(t *testing.T) {
	suite.Run(t, new(MuMpcAbcSuite))
}

func (s *MuMpcAbcSuite) BeforeTest(suiteName, testName string) {
	s.MuMPCSuite = NewMuMPCSuite(
		&s.Suite,
		&MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeAbc,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
		&MonitorEnergyConfig{
			ValueSourceProduction: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsProduction: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(0.001),
			}),
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		},
		&MonitorCurrentConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeB: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
		},
		&MonitorVoltageConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeB:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeC:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeAb: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeBc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
				model.ElectricalConnectionPhaseNameTypeAc: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
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

func (s *MuMpcAbcSuite) Test_Power() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerTotal(5.0, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	power, err := s.sut.Power()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, power)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, nil)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0}, values)
}

func (s *MuMpcAbcSuite) Test_PowerPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerPhaseA(5.0, util.Ptr(time.Now()), nil),
		s.sut.UpdateDataPowerPhaseB(6.0, util.Ptr(time.Now()), nil),
		s.sut.UpdateDataPowerPhaseC(7.0, util.Ptr(time.Now()), util.Ptr(model.MeasurementValueStateTypeError)),
	)
	assert.Nil(s.T(), err)
	expectedPowerPerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeA: 5.0,
		model.ElectricalConnectionPhaseNameTypeB: 6.0,
		model.ElectricalConnectionPhaseNameTypeC: 7.0,
	}

	powerPerPhases, err := s.sut.PowerPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedPowerPerPhases, powerPerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []float64{5.0, 6.0, 7.0}, values)
}

func (s *MuMpcAbcSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyConsumed(5.0, util.Ptr(time.Now()), nil, util.Ptr(time.Now()), util.Ptr(time.Now())),
	)
	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyConsumed)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyConsumed),
	}
	measurement, err := server.NewMeasurement(s.sut.LocalEntity)
	assert.Nil(s.T(), err)
	values, err := measurement.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(values))
	assert.Equal(s.T(), 5.0, (*values[0].Value).GetValue())
}

func (s *MuMpcAbcSuite) Test_EnergyProduced() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyProduced(5.0, nil, nil, nil, nil),
	)
	assert.Nil(s.T(), err)

	energyProduced, err := s.sut.EnergyProduced()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyProduced)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
	}
	measurement, err := server.NewMeasurement(s.sut.LocalEntity)
	assert.Nil(s.T(), err)
	values, err := measurement.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(values))
	assert.Equal(s.T(), 5.0, (*values[0].Value).GetValue())
}

func (s *MuMpcAbcSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataCurrentPhaseA(5.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseB(3.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseC(1.0, nil, nil),
	)
	assert.Nil(s.T(), err)
	expectedCurrentPerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeA: 5.0,
		model.ElectricalConnectionPhaseNameTypeB: 3.0,
		model.ElectricalConnectionPhaseNameTypeC: 1.0,
	}

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedCurrentPerPhases, currentPerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []float64{5.0, 3.0, 1.0}, values)
}

func (s *MuMpcAbcSuite) Test_VoltagePerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataVoltagePhaseA(5.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseB(6.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseC(7.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseAToB(8.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseBToC(9.0, nil, nil),
		s.sut.UpdateDataVoltagePhaseAToC(10.0, nil, nil),
	)
	assert.Nil(s.T(), err)
	expectedVoltagePerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
		model.ElectricalConnectionPhaseNameTypeA:  5.0,
		model.ElectricalConnectionPhaseNameTypeB:  6.0,
		model.ElectricalConnectionPhaseNameTypeC:  7.0,
		model.ElectricalConnectionPhaseNameTypeAb: 8.0,
		model.ElectricalConnectionPhaseNameTypeBc: 9.0,
		model.ElectricalConnectionPhaseNameTypeAc: 10.0,
	}

	voltagePerPhases, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedVoltagePerPhases, voltagePerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, "", ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0}, values)
}

func (s *MuMpcAbcSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.UpdateDataFrequency(5.0, nil, nil),
	)
	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, frequency)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	}
	measurements, err := server.NewMeasurement(s.sut.LocalEntity)
	assert.Nil(s.T(), err)
	values, err := measurements.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(values))
	assert.Equal(s.T(), 5.0, (*values[0].Value).GetValue())
}
