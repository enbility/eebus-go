package features_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	measurement *features.Measurement
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*MeasurementSuite)(nil)

func (s *MeasurementSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeMeasurement,
				functions: []model.FunctionType{
					model.FunctionTypeMeasurementDescriptionListData,
					model.FunctionTypeMeasurementConstraintsListData,
					model.FunctionTypeMeasurementListData,
				},
			},
			{
				featureType: model.FeatureTypeTypeElectricalConnection,
				functions: []model.FunctionType{
					model.FunctionTypeElectricalConnectionDescriptionListData,
					model.FunctionTypeElectricalConnectionParameterDescriptionListData,
					model.FunctionTypeElectricalConnectionPermittedValueSetListData,
				},
			},
		},
	)

	var err error
	s.measurement, err = features.NewMeasurement(model.RoleTypeServer, model.RoleTypeClient, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.measurement)
}

func (s *MeasurementSuite) Test_RequestDescriptions() {
	err := s.measurement.RequestDescriptions()
	assert.Nil(s.T(), err)
}

func (s *MeasurementSuite) Test_RequestConstraints() {
	err := s.measurement.RequestConstraints()
	assert.Nil(s.T(), err)
}

func (s *MeasurementSuite) Test_RequestValues() {
	counter, err := s.measurement.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *MeasurementSuite) Test_GetValuesForTypeCommodityScope() {
	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent

	data, err := s.measurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.measurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.measurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurement = model.MeasurementTypeTypeArea
	data, err = s.measurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetDescriptionsForScope() {
	data, err := s.measurement.GetDescriptionsForScope(model.ScopeTypeTypeACCurrent)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.measurement.GetDescriptionsForScope(model.ScopeTypeTypeACCurrent)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetConstraints() {
	data, err := s.measurement.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	data, err = s.measurement.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetValues() {
	data, err := s.measurement.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	s.addDescription()

	data, err = s.measurement.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.measurement.GetValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *MeasurementSuite) addDescription() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				Unit:            util.Ptr(model.UnitOfMeasurementTypeA),
			},
			{
				MeasurementId:   util.Ptr(model.MeasurementIdType(1)),
				MeasurementType: util.Ptr(model.MeasurementTypeTypePercentage),
				CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementDescriptionListData, fData, nil, nil)
}

func (s *MeasurementSuite) addConstraints() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []model.MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ValueRangeMin: model.NewScaledNumberType(2),
				ValueRangeMax: model.NewScaledNumberType(16),
				ValueStepSize: model.NewScaledNumberType(0.1),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueRangeMin: model.NewScaledNumberType(0),
				ValueRangeMax: model.NewScaledNumberType(100),
				ValueStepSize: model.NewScaledNumberType(0.1),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementConstraintsListData, fData, nil, nil)
}

func (s *MeasurementSuite) addData() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))

	t := time.Now()
	fData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(9),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(t),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(9),
				Timestamp:     model.NewAbsoluteOrRelativeTimeTypeFromTime(t),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeMeasurementListData, fData, nil, nil)
}
