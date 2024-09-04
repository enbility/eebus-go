package client

import (
	"testing"

	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDeviceConfigurationSuite(t *testing.T) {
	suite.Run(t, new(DeviceConfigurationSuite))
}

type DeviceConfigurationSuite struct {
	suite.Suite

	localEntity        spineapi.EntityLocalInterface
	localEntityPartial spineapi.EntityLocalInterface

	remoteEntity        spineapi.EntityRemoteInterface
	remoteEntityPartial spineapi.EntityRemoteInterface

	mockRemoteEntity *mocks.EntityRemoteInterface

	deviceConfiguration        *DeviceConfiguration
	deviceConfigurationPartial *DeviceConfiguration
}

const remoteSki string = "testremoteski"

func (s *DeviceConfigurationSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceConfiguration,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
					model.FunctionTypeDeviceConfigurationKeyValueListData,
				},
			},
		},
	)

	s.localEntityPartial, s.remoteEntityPartial = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceConfiguration,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
					model.FunctionTypeDeviceConfigurationKeyValueListData,
				},
				partial: true,
			},
		},
	)

	mockRemoteDevice := mocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = mocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := mocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()

	var err error
	s.deviceConfiguration, err = NewDeviceConfiguration(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.deviceConfiguration)

	s.deviceConfiguration, err = NewDeviceConfiguration(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceConfiguration)

	s.deviceConfigurationPartial, err = NewDeviceConfiguration(s.localEntityPartial, s.remoteEntityPartial)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceConfiguration)
}

func (s *DeviceConfigurationSuite) Test_RequestKeyValueDescriptions() {
	counter, err := s.deviceConfiguration.RequestKeyValueDescriptions(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.deviceConfiguration.RequestKeyValueDescriptions(
		&model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType{},
		&model.DeviceConfigurationKeyValueDescriptionDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *DeviceConfigurationSuite) Test_RequestKeyValueList() {
	counter, err := s.deviceConfiguration.RequestKeyValues(nil, nil)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)

	counter, err = s.deviceConfiguration.RequestKeyValues(
		&model.DeviceConfigurationKeyValueListDataSelectorsType{},
		&model.DeviceConfigurationKeyValueDataElementsType{},
	)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *DeviceConfigurationSuite) Test_WriteValues() {
	counter, err := s.deviceConfiguration.WriteKeyValues(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.DeviceConfigurationKeyValueDataType{}
	counter, err = s.deviceConfiguration.WriteKeyValues(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	rF := s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	data1 := rF.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueListData).(*model.DeviceConfigurationKeyValueListDataType)
	assert.Nil(s.T(), data1)

	defaultData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId:             util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				IsValueChangeable: util.Ptr(true),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(16),
				},
			},
			{
				KeyId:             util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				IsValueChangeable: util.Ptr(true),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(32),
				},
			},
		},
	}
	_, err1 := rF.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, defaultData, nil, nil)
	assert.Nil(s.T(), err1)
	data1 = rF.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueListData).(*model.DeviceConfigurationKeyValueListDataType)
	assert.NotNil(s.T(), data1)
	assert.Equal(s.T(), 2, len(data1.DeviceConfigurationKeyValueData))

	data = []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(10),
			},
		},
	}
	counter, err = s.deviceConfiguration.WriteKeyValues(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

// test with partial support
func (s *DeviceConfigurationSuite) Test_WriteValues_Partial() {
	counter, err := s.deviceConfigurationPartial.WriteKeyValues(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.DeviceConfigurationKeyValueDataType{}
	counter, err = s.deviceConfigurationPartial.WriteKeyValues(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	rF := s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	data1 := rF.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueListData).(*model.DeviceConfigurationKeyValueListDataType)
	assert.Nil(s.T(), data1)

	defaultData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId:             util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				IsValueChangeable: util.Ptr(true),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(16),
				},
			},
			{
				KeyId:             util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				IsValueChangeable: util.Ptr(true),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(32),
				},
			},
		},
	}
	_, err1 := rF.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, defaultData, nil, nil)
	assert.Nil(s.T(), err1)
	data1 = rF.DataCopy(model.FunctionTypeDeviceConfigurationKeyValueListData).(*model.DeviceConfigurationKeyValueListDataType)
	assert.NotNil(s.T(), data1)
	assert.Equal(s.T(), 2, len(data1.DeviceConfigurationKeyValueData))

	data = []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(10),
			},
		},
	}
	counter, err = s.deviceConfigurationPartial.WriteKeyValues(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
