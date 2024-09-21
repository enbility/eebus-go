package internal_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.MeasurementCommon
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
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

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalMeasurement(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteMeasurement(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *MeasurementSuite) Test_MeasurementCheckPayloadDataForScope() {
	scopeType := model.ScopeTypeTypeACPower
	filter := model.MeasurementDescriptionDataType{
		ScopeType: &scopeType,
	}
	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	exists = s.localSut.CheckEventPayloadDataForFilter(scopeType, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(scopeType, filter)
	assert.False(s.T(), exists)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				ScopeType: util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	fErr := s.localFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)

	data = &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				Value: model.NewScaledNumberType(80),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.True(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(data, filter)
	assert.True(s.T(), exists)
}

func (s *MeasurementSuite) Test_GetValueForMeasurementId() {
	measurement := model.MeasurementIdType(0)

	data, err := s.localSut.GetDataForId(measurement)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForId(measurement)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDataForId(measurement)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForId(measurement)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetDataForId(measurement)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 9.0, data.Value.GetValue())
	data, err = s.remoteSut.GetDataForId(measurement)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 9.0, data.Value.GetValue())

	data, err = s.localSut.GetDataForId(model.MeasurementIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForId(model.MeasurementIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetValuesForTypeCommodityScope() {
	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: &measurement,
		CommodityType:   &commodity,
		ScopeType:       &scope,
	}

	data, err := s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurement = model.MeasurementTypeTypeArea
	data, err = s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetDescriptionsForId() {
	measurementId := model.MeasurementIdType(0)
	data, err := s.localSut.GetDescriptionForId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionForId(measurementId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDescriptionForId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionForId(measurementId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetDescriptionsForScope() {
	filter := model.MeasurementDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	data, err := s.localSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetConstraints() {
	filter := model.MeasurementConstraintsDataType{}
	data, err := s.localSut.GetConstraintsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetConstraintsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	data, err = s.localSut.GetConstraintsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetConstraintsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *MeasurementSuite) Test_GetValues() {
	filter := model.MeasurementDescriptionDataType{}
	data, err := s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addConstraints()

	s.addDescription()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *MeasurementSuite) addDescription() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeMeasurementDescriptionListData, fData, nil, nil)
}

func (s *MeasurementSuite) addConstraints() {
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
	_ = s.localFeature.UpdateData(model.FunctionTypeMeasurementConstraintsListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeMeasurementConstraintsListData, fData, nil, nil)
}

func (s *MeasurementSuite) addData() {
	t := time.Now().UTC()
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
	_ = s.localFeature.UpdateData(model.FunctionTypeMeasurementListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeMeasurementListData, fData, nil, nil)
}
