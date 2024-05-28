package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/util"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDeviceConfigurationSuite(t *testing.T) {
	suite.Run(t, new(DeviceConfigurationSuite))
}

type DeviceConfigurationSuite struct {
	suite.Suite

	localEntity      spineapi.EntityLocalInterface
	remoteEntity     spineapi.EntityRemoteInterface
	mockRemoteEntity *mocks.EntityRemoteInterface

	deviceConfiguration *features.DeviceConfiguration
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
	s.deviceConfiguration, err = features.NewDeviceConfiguration(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceConfiguration)
}

func (s *DeviceConfigurationSuite) Test_RequestDescriptions() {
	counter, err := s.deviceConfiguration.RequestDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *DeviceConfigurationSuite) Test_RequestKeyValueList() {
	counter, err := s.deviceConfiguration.RequestKeyValues()
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
