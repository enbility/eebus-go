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

type MuMpcConstraintSuite struct {
	suite.Suite
	*MuMPCSuite
}

func TestMuMpcConstraintSuite(t *testing.T) {
	suite.Run(t, new(MuMpcConstraintSuite))
}

func (s *MuMpcConstraintSuite) BeforeTest(suiteName, testName string) {
	s.MuMPCSuite = NewMuMPCSuite(
		&s.Suite,
		&MonitorPowerConfig{
			ConnectedPhases:   ConnectedPhasesA,
			ValueSourceTotal:  util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsTotal: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(0.1),
			}),
			ValueConstraintsPhaseA: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(0.1),
			}),
		},
		&MonitorEnergyConfig{
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsConsumption: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(100),
			}),
		},
		&MonitorCurrentConfig{
			ValueSourcePhaseA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsPhaseA: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(32),
				ValueStepSize: model.NewScaledNumberType(0.1),
			}),
		},
		&MonitorVoltageConfig{
			SupportPhaseToPhase: false,
			ValueSourcePhaseA:   util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsPhaseA: util.Ptr(model.MeasurementConstraintsDataType{
				ValueStepSize: model.NewScaledNumberType(0.1),
			}),
		},
		&MonitorFrequencyConfig{
			ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraints: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(40),
				ValueRangeMax: model.NewScaledNumberType(60),
				ValueStepSize: model.NewScaledNumberType(0.01),
			}),
		},
	)
	s.MuMPCSuite.BeforeTest(suiteName, testName)
}

func (s *MuMpcConstraintSuite) Test_Power() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerTotal(5.7, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	power, err := s.sut.Power()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.7, power)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, nil)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.7}, values)
}

func (s *MuMpcConstraintSuite) Test_PowerPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerPhaseA(5.7, util.Ptr(time.Now()), nil),
	)
	assert.Nil(s.T(), err)

	powerPerPhases, err := s.sut.PowerPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.7}, powerPerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.7}, values)
}

func (s *MuMpcConstraintSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyConsumed(570, util.Ptr(time.Now()), nil, util.Ptr(time.Now()), util.Ptr(time.Now())),
	)
	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 570.0, energyConsumed)

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
	assert.Equal(s.T(), 570.0, (*values[0].Value).GetValue())
}

func (s *MuMpcConstraintSuite) Test_EnergyProduced() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyProduced(5.0, nil, nil, nil, nil),
	)
	assert.NotNil(s.T(), err)

	_, err = s.sut.EnergyProduced()
	assert.NotNil(s.T(), err)

	// Check if the client filter works (it shouldn't)
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACEnergyProduced),
	}
	measurement, err := server.NewMeasurement(s.sut.LocalEntity)
	assert.Nil(s.T(), err)
	_, err = measurement.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
}

func (s *MuMpcConstraintSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataCurrentPhaseA(0.1, nil, nil),
	)
	assert.Nil(s.T(), err)

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{0.1}, currentPerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{0.1}, values)
}

func (s *MuMpcConstraintSuite) Test_VoltagePerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataVoltagePhaseA(230, nil, nil),
	)
	assert.Nil(s.T(), err)

	voltagePerPhases, err := s.sut.VoltagePerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230}, voltagePerPhases)

	// Check if the client filter works
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}
	values, err := s.measurementPhaseSpecificDataForFilter(filter, "", ucapi.PhaseNameMapping)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230}, values)
}

func (s *MuMpcConstraintSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.UpdateDataFrequency(50, nil, nil),
	)
	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, frequency)

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
	assert.Equal(s.T(), 50.0, (*values[0].Value).GetValue())
}
