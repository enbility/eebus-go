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

// Test suite testing an MPC MonitoredUnit that uses ValueConstraints on its measurements
func TestMuMpcConstraintSuite(t *testing.T) {
	suite.Run(t, new(MuMpcConstraintSuite))
}

func (s *MuMpcConstraintSuite) BeforeTest(suiteName, testName string) {
	s.MuMPCSuite = NewMuMPCSuite(
		&s.Suite,
		&MonitorPowerConfig{
			ConnectedPhases:  model.ElectricalConnectionPhaseNameTypeA,
			ValueSourceTotal: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			ValueConstraintsTotal: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(0.1),
			}),
			ValueConstraintsPerPhase: PhaseMeasurementConstraintsMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementConstraintsDataType{
					ValueRangeMin: model.NewScaledNumberType(0),
					ValueStepSize: model.NewScaledNumberType(0.1),
				}),
			},
		},
		&MonitorEnergyConfig{
			ValueSourceConsumption: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			ValueConstraintsConsumption: util.Ptr(model.MeasurementConstraintsDataType{
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueStepSize: model.NewScaledNumberType(100),
			}),
		},
		&MonitorCurrentConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			ValueConstraintsPerPhase: PhaseMeasurementConstraintsMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementConstraintsDataType{
					ValueRangeMin: model.NewScaledNumberType(0),
					ValueRangeMax: model.NewScaledNumberType(32),
					ValueStepSize: model.NewScaledNumberType(0.1),
				}),
			},
		},
		&MonitorVoltageConfig{
			ValueSourcePerPhase: PhaseMeasurementSourceMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
			},
			ValueConstraintsPerPhase: PhaseMeasurementConstraintsMap{
				model.ElectricalConnectionPhaseNameTypeA: util.Ptr(model.MeasurementConstraintsDataType{
					ValueStepSize: model.NewScaledNumberType(0.1),
				}),
			},
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
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.Power()
		assert.Error(s.T(), err)
	}

	// Testing updating the power total and read this update
	{
		err := s.sut.Update(
			s.sut.UpdateDataPowerTotal(5.7, util.Ptr(time.Now()), nil),
		)
		assert.Nil(s.T(), err)

		power, err := s.sut.Power()
		expectedPowerValue := 5.7
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), expectedPowerValue, power)

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
}

func (s *MuMpcConstraintSuite) Test_PowerPerPhase() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.PowerPerPhase()
		assert.Error(s.T(), err)
	}

	// Test when updating power power per phase and read this update
	{
		err := s.sut.Update(
			s.sut.UpdateDataPowerPhaseA(5.7, util.Ptr(time.Now()), nil),
		)
		expectedPowerPerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
			model.ElectricalConnectionPhaseNameTypeA: 5.7,
		}
		assert.Nil(s.T(), err)

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
		assert.Equal(s.T(), []float64{5.7}, values)
	}

	// Test updating the energy consumed when it is not supported
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, nil, s.currentConfig, s.voltageConfig, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		err = mpcInstance.Update(
			mpcInstance.UpdateDataEnergyConsumed(100, nil, nil, nil, nil),
		)
		assert.Error(s.T(), err)
	}
}

func (s *MuMpcConstraintSuite) Test_EnergyConsumed() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.EnergyConsumed()
		assert.Error(s.T(), err)
	}

	// Test when acEnergyConsumed has no measurement id
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, s.energyConfig, s.currentConfig, s.voltageConfig, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		mpcInstance.acEnergyConsumed = nil
		_, err = mpcInstance.EnergyConsumed()
		assert.Error(s.T(), err)
	}

	// Test when updating the energy consumed and get this update
	{
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
}

func (s *MuMpcConstraintSuite) Test_EnergyProduced() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.EnergyProduced()
		assert.Error(s.T(), err)
	}

	// Test when updating the energy produced and read this update
	{
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

	// Test updating the energy produced when it is not supported
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, nil, s.currentConfig, s.voltageConfig, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		err = mpcInstance.Update(
			mpcInstance.UpdateDataEnergyProduced(100, nil, nil, nil, nil),
		)
		assert.Error(s.T(), err)
	}
}

func (s *MuMpcConstraintSuite) Test_CurrentPerPhase() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.CurrentPerPhase()
		assert.Error(s.T(), err)
	}

	// Test when updating the current per phase and read this update
	{
		err := s.sut.Update(
			s.sut.UpdateDataCurrentPhaseA(0.1, nil, nil),
		)
		assert.Nil(s.T(), err)
		expectedCurrentPerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
			model.ElectricalConnectionPhaseNameTypeA: 0.1,
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
		assert.Equal(s.T(), []float64{0.1}, values)
	}

	// Test updating the current per phase when it is not supported
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, s.energyConfig, nil, s.voltageConfig, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		err = mpcInstance.Update(
			mpcInstance.UpdateDataCurrentPhaseA(10, nil, nil),
			mpcInstance.UpdateDataCurrentPhaseB(10, nil, nil),
			mpcInstance.UpdateDataCurrentPhaseC(10, nil, nil),
		)
		assert.Error(s.T(), err)
	}
}

func (s *MuMpcConstraintSuite) Test_VoltagePerPhase() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.VoltagePerPhase()
		assert.Error(s.T(), err)
	}

	// Test when updating the voltage per phase and read this update
	{
		err := s.sut.Update(
			s.sut.UpdateDataVoltagePhaseA(230, nil, nil),
		)
		assert.Nil(s.T(), err)

		expectedVoltagePerPhases := map[model.ElectricalConnectionPhaseNameType]float64{
			model.ElectricalConnectionPhaseNameTypeA: 230,
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
		assert.Equal(s.T(), []float64{230}, values)
	}

	// Test updating the voltage per phase when it is not supported
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, s.energyConfig, s.currentConfig, nil, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		err = mpcInstance.Update(
			mpcInstance.UpdateDataVoltagePhaseA(230, nil, nil),
			mpcInstance.UpdateDataVoltagePhaseB(230, nil, nil),
			mpcInstance.UpdateDataVoltagePhaseC(230, nil, nil),
			mpcInstance.UpdateDataVoltagePhaseAToB(0, nil, nil),
			mpcInstance.UpdateDataVoltagePhaseBToC(0, nil, nil),
			mpcInstance.UpdateDataVoltagePhaseAToC(0, nil, nil),
		)
		assert.Error(s.T(), err)
	}
}

func (s *MuMpcConstraintSuite) Test_Frequency() {
	// Test when getMeasurementId returns error
	{
		_, err := s.sut.Frequency()
		assert.Error(s.T(), err)
	}

	// Test when acFrequency has no measurement id
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, s.energyConfig, s.currentConfig, s.voltageConfig, s.frequencyConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		mpcInstance.acFrequency = nil
		_, err = mpcInstance.Frequency()
		assert.Error(s.T(), err)
	}

	// Test when updating the frequency and read this update
	{
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

	// Test covering updating frequency when it is not supported
	{
		mpcInstance, err := NewMPC(s.sut.LocalEntity, s.Event, s.powerConfig, s.energyConfig, s.currentConfig, s.voltageConfig, nil)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), mpcInstance)

		err = mpcInstance.AddFeatures()
		assert.Nil(s.T(), err)
		mpcInstance.AddUseCase()

		err = mpcInstance.Update(
			mpcInstance.UpdateDataFrequency(50, nil, nil),
		)
		assert.Error(s.T(), err)
	}
}
