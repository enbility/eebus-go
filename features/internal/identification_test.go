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

func TestIdentificationSuite(t *testing.T) {
	suite.Run(t, new(IdentificationSuite))
}

type IdentificationSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.IdentificationCommon
}

func (s *IdentificationSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeIdentification,
				functions: []model.FunctionType{
					model.FunctionTypeIdentificationListData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeIdentification, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalIdentification(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeIdentification, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteIdentification(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *IdentificationSuite) Test_CheckEventPayloadDataForFilter() {
	exists := s.localSut.CheckEventPayloadDataForFilter(nil)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil)
	assert.False(s.T(), exists)

	temp := true
	exists = s.localSut.CheckEventPayloadDataForFilter(temp)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(temp)
	assert.False(s.T(), exists)

	data := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(data)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(data)
	assert.False(s.T(), exists)

	data = &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId: util.Ptr(model.IdentificationIdType(0)),
			},
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(1)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(data)
	assert.True(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(data)
	assert.True(s.T(), exists)
}

func (s *IdentificationSuite) Test_GetValues() {
	filter := model.IdentificationDataType{}
	data, err := s.localSut.GetDataForFilter(filter)
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

func (s *IdentificationSuite) addData() {
	fData := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(0)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeIdentificationListData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeIdentificationListData, fData, nil, nil)
}
