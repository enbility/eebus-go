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

func TestSetpointSuite(t *testing.T) {
	suite.Run(t, new(SetpointSuite))
}

type SetpointSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.SetpointCommon
}

func (s *SetpointSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeSetpoint,
				functions: []model.FunctionType{
					model.FunctionTypeSetpointDescriptionListData,
					model.FunctionTypeSetpointConstraintsListData,
					model.FunctionTypeSetpointListData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeSetpoint, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalSetpoint(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeSetpoint, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteSetpoint(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *SetpointSuite) Test_GetSetpointDescriptions() {
	filter := model.SetpointDescriptionDataType{}
	data, err := s.localSut.GetSetpointDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetSetpointDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc, err := s.localSut.GetSetpointDescriptionForId(model.SetpointIdType(0))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescriptions()

	data, err = s.localSut.GetSetpointDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetSetpointDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	filter.ScopeType = util.Ptr(model.ScopeTypeTypeDhwTemperature)
	data, err = s.localSut.GetSetpointDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))

	desc, err = s.localSut.GetSetpointDescriptionForId(model.SetpointIdType(0))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetSetpointDescriptionForId(model.SetpointIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
}

func (s *SetpointSuite) Test_GetSetpointData() {
	filter := model.SetpointDataType{}
	data, err := s.localSut.GetSetpointDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetSetpointDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	sp, err := s.localSut.GetSetpointForId(model.SetpointIdType(0))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), sp)

	s.addData()

	data, err = s.localSut.GetSetpointDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))
	data, err = s.remoteSut.GetSetpointDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2, len(data))

	sp, err = s.localSut.GetSetpointForId(model.SetpointIdType(0))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), sp)
	assert.Equal(s.T(), 21.0, sp.Value.GetValue())

	sp, err = s.remoteSut.GetSetpointForId(model.SetpointIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), sp)
}

func (s *SetpointSuite) Test_GetSetpointConstraints() {
	filter := model.SetpointConstraintsDataType{}
	data, err := s.localSut.GetSetpointConstraintsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetSetpointConstraintsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	constraints, err := s.localSut.GetSetpointConstraintsForId(model.SetpointIdType(0))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), constraints)

	s.addConstraints()

	data, err = s.localSut.GetSetpointConstraintsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	data, err = s.remoteSut.GetSetpointConstraintsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))

	constraints, err = s.localSut.GetSetpointConstraintsForId(model.SetpointIdType(0))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), constraints)
	assert.Equal(s.T(), 16.0, constraints.SetpointRangeMin.GetValue())
	assert.Equal(s.T(), 25.0, constraints.SetpointRangeMax.GetValue())

	constraints, err = s.remoteSut.GetSetpointConstraintsForId(model.SetpointIdType(10))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), constraints)
}

func (s *SetpointSuite) Test_CheckEventPayloadDataForFilter() {
	filter := model.SetpointDataType{
		SetpointId: util.Ptr(model.SetpointIdType(0)),
	}

	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	temp := true
	exists = s.localSut.CheckEventPayloadDataForFilter(temp, filter)
	assert.False(s.T(), exists)

	data := &model.SetpointListDataType{}
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)

	data = &model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId: util.Ptr(model.SetpointIdType(0)),
				Value:      model.NewScaledNumberType(21),
			},
		},
	}
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.True(s.T(), exists)

	filter.SetpointId = util.Ptr(model.SetpointIdType(10))
	exists = s.localSut.CheckEventPayloadDataForFilter(data, filter)
	assert.False(s.T(), exists)
}

// helpers

func (s *SetpointSuite) addDescriptions() {
	fData := &model.SetpointDescriptionListDataType{
		SetpointDescriptionData: []model.SetpointDescriptionDataType{
			{
				SetpointId:   util.Ptr(model.SetpointIdType(0)),
				SetpointType: util.Ptr(model.SetpointTypeTypeValueAbsolute),
				ScopeType:    util.Ptr(model.ScopeTypeTypeDhwTemperature),
				Unit:         util.Ptr(model.UnitOfMeasurementTypedegC),
			},
			{
				SetpointId:   util.Ptr(model.SetpointIdType(1)),
				SetpointType: util.Ptr(model.SetpointTypeTypeValueAbsolute),
				ScopeType:    util.Ptr(model.ScopeTypeTypeRoomAirTemperature),
				Unit:         util.Ptr(model.UnitOfMeasurementTypedegC),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeSetpointDescriptionListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeSetpointDescriptionListData, fData, nil, nil)
}

func (s *SetpointSuite) addData() {
	fData := &model.SetpointListDataType{
		SetpointData: []model.SetpointDataType{
			{
				SetpointId:           util.Ptr(model.SetpointIdType(0)),
				Value:                model.NewScaledNumberType(21),
				IsSetpointChangeable: util.Ptr(true),
			},
			{
				SetpointId:       util.Ptr(model.SetpointIdType(1)),
				Value:            model.NewScaledNumberType(45),
				IsSetpointActive: util.Ptr(true),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeSetpointListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeSetpointListData, fData, nil, nil)
}

func (s *SetpointSuite) addConstraints() {
	fData := &model.SetpointConstraintsListDataType{
		SetpointConstraintsData: []model.SetpointConstraintsDataType{
			{
				SetpointId:       util.Ptr(model.SetpointIdType(0)),
				SetpointRangeMin: model.NewScaledNumberType(16),
				SetpointRangeMax: model.NewScaledNumberType(25),
				SetpointStepSize: model.NewScaledNumberType(0.5),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeSetpointConstraintsListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeSetpointConstraintsListData, fData, nil, nil)
}
