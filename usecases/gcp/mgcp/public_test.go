package mgcp

import (
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *GcpMpcgSuite) Test_PowerLimitationFactor() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerLimitationFactor(0.5),
	)

	assert.Nil(s.T(), err)

	powerLimitationFactor, err := s.sut.PowerLimitationFactor()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0.5, powerLimitationFactor)

	// Test client getter
	keyname := model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor

	deviceConfiguration, err := server.NewDeviceConfiguration(s.localEntity)
	assert.Nil(s.T(), err)

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyname,
	}

	_, err = deviceConfiguration.GetKeyValueDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)

	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	assert.Equal(s.T(), 0.5, data.Value.ScaledNumber.GetValue())
}

func (s *GcpMpcgSuite) Test_PowerTotal() {
	err := s.sut.Update(
		s.sut.UpdateDataPowerTotal(5.0, nil, nil),
	)

	assert.Nil(s.T(), err)

	totalPower, err := s.sut.PowerTotal()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, totalPower)

	// Test client getter
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}

	data, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, nil)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, data[0])
}

func (s *GcpMpcgSuite) Test_EnergyConsumed() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyConsumed(5.0, nil, nil, nil, nil),
	)

	assert.Nil(s.T(), err)

	energyConsumed, err := s.sut.EnergyConsumed()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5.0, energyConsumed)

	// Test client getter
	measurement, err := server.NewMeasurement(s.localEntity)
	assert.Nil(s.T(), err)

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridConsumption),
	}

	result, err := measurement.GetDataForFilter(filter)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 5.0, result[0].Value.GetValue())
}

func (s *GcpMpcgSuite) Test_EnergyFeedIn() {
	err := s.sut.Update(
		s.sut.UpdateDataEnergyFeedIn(6.0, nil, nil, nil, nil),
	)

	assert.Nil(s.T(), err)

	energyProduced, err := s.sut.EnergyFeedIn()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6.0, energyProduced)

	// Test client getter
	measurement, err := server.NewMeasurement(s.localEntity)
	assert.Nil(s.T(), err)

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeGridFeedIn),
	}

	result, err := measurement.GetDataForFilter(filter)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), 6.0, result[0].Value.GetValue())
}

func (s *GcpMpcgSuite) Test_CurrentPerPhase() {
	err := s.sut.Update(
		s.sut.UpdateDataCurrentPhaseA(5.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseB(6.0, nil, nil),
		s.sut.UpdateDataCurrentPhaseC(7.0, nil, nil),
	)

	assert.Nil(s.T(), err)

	currentPerPhases, err := s.sut.CurrentPerPhase()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0}, currentPerPhases)

	// Test client getter
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}

	data, err := s.measurementPhaseSpecificDataForFilter(filter, model.EnergyDirectionTypeConsume, nil)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0}, data)
}

func (s *GcpMpcgSuite) Test_VoltagePerPhase() {
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

	// Test client getter
	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeVoltage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACVoltage),
	}

	data, err := s.measurementPhaseSpecificDataForFilter(filter, "", nil)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{5.0, 6.0, 7.0, 8.0, 9.0, 10.0}, data)
}

func (s *GcpMpcgSuite) Test_Frequency() {
	err := s.sut.Update(
		s.sut.UpdateDataFrequency(50.0, nil, nil),
	)

	assert.Nil(s.T(), err)

	frequency, err := s.sut.Frequency()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, frequency)

	// Test client getter
	measurement, err := server.NewMeasurement(s.localEntity)
	assert.Nil(s.T(), err)

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeFrequency),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACFrequency),
	}

	data, err := measurement.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data[0].Value.GetValue())
}
