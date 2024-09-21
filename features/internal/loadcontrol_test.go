package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestLoadControlSuite(t *testing.T) {
	suite.Run(t, new(LoadControlSuite))
}

type LoadControlSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.LoadControlCommon
}

func (s *LoadControlSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeLoadControl,
				functions: []model.FunctionType{
					model.FunctionTypeLoadControlLimitDescriptionListData,
					model.FunctionTypeLoadControlLimitConstraintsListData,
					model.FunctionTypeLoadControlLimitListData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalLoadControl(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteLoadControl(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *LoadControlSuite) Test_CheckEventPayloadDataForFilter() {
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     &limitType,
		LimitCategory: &category,
		ScopeType:     &scope,
	}
	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	exists = s.localSut.CheckEventPayloadDataForFilter(limitType, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(limitType, filter)
	assert.False(s.T(), exists)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				ScopeType: util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: util.Ptr(category),
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				LimitType:     util.Ptr(limitType),
				ScopeType:     util.Ptr(scope),
			},
		},
	}

	fErr := s.localFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.False(s.T(), exists)

	limitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.True(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.True(s.T(), exists)
}

func (s *LoadControlSuite) Test_GetLimitDescriptions() {
	filter := model.LoadControlLimitDescriptionDataType{}
	data, err := s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDescriptionsForId() {
	limitId := model.LoadControlLimitIdType(0)
	data, err := s.localSut.GetLimitDescriptionForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDescriptionForId(limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionForId(limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDescriptionsForFilter() {
	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	data, err := s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	filter.LimitType = util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit)
	data, err = s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDescriptionsForFilterMeasurementId() {
	measurementId := util.Ptr(model.MeasurementIdType(0))
	filter := model.LoadControlLimitDescriptionDataType{
		MeasurementId: measurementId,
	}

	data, err := s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	measurementId = util.Ptr(model.MeasurementIdType(10))
	filter = model.LoadControlLimitDescriptionDataType{
		MeasurementId: measurementId,
	}

	data, err = s.localSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitData() {
	filter := model.LoadControlLimitDescriptionDataType{}
	data, err := s.localSut.GetLimitDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetLimitDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *LoadControlSuite) Test_GetLimitDataForLimitId() {
	limitId := model.LoadControlLimitIdType(0)
	data, err := s.localSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetLimitDataForId(limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForId(limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	limitId = model.LoadControlLimitIdType(10)
	data, err = s.localSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetLimitDataForId(limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
}

// helper

func (s *LoadControlSuite) addDescription() {
	fData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				MeasurementId:  util.Ptr(model.MeasurementIdType(0)),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, fData, nil, nil)
}

func (s *LoadControlSuite) addData() {
	fData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(true),
				Value:             model.NewScaledNumberType(12),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, fData, nil, nil)
}
