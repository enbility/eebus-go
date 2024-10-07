package server_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
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

	sut *server.DeviceConfiguration

	service api.ServiceInterface

	localEntity spineapi.EntityLocalInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	remoteEntity     spineapi.EntityRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
}

func (s *DeviceConfigurationSuite) BeforeTest(suiteName, testName string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()
	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	mockRemoteDevice := spinemocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = spinemocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := spinemocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()

	var entities []spineapi.EntityRemoteInterface

	s.remoteDevice, entities = setupFeatures(s.service, s.T())
	s.remoteEntity = entities[1]

	var err error
	s.sut, err = server.NewDeviceConfiguration(nil)
	assert.NotNil(s.T(), err)

	s.sut, err = server.NewDeviceConfiguration(s.localEntity)
	assert.Nil(s.T(), err)
}

func (s *DeviceConfigurationSuite) Test_CheckEventPayloadDataForFilter() {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit

	exists := s.sut.CheckEventPayloadDataForFilter(nil, keyName)
	assert.False(s.T(), exists)

	exists = s.sut.CheckEventPayloadDataForFilter(keyName, keyName)
	assert.False(s.T(), exists)

	descData := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(keyName),
	}

	keyId := s.sut.AddKeyValueDescription(descData)
	assert.NotNil(s.T(), keyId)

	exists = s.sut.CheckEventPayloadDataForFilter(nil, descData)
	assert.False(s.T(), exists)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
	}

	exists = s.sut.CheckEventPayloadDataForFilter(keyData, descData)
	assert.False(s.T(), exists)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: keyId,
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringTypeIEC61851),
				},
			},
		},
	}

	exists = s.sut.CheckEventPayloadDataForFilter(keyData, descData)
	assert.True(s.T(), exists)
}

func (s *DeviceConfigurationSuite) Test_Description() {
	filter1 := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
	}

	filter2 := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
	}

	data, err := s.sut.GetKeyValueDescriptionsForFilter(filter1)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.GetKeyValueDescriptionsForFilter(filter2)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   filter1.KeyName,
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
	}
	keyId := s.sut.AddKeyValueDescription(desc)
	assert.NotNil(s.T(), keyId)

	data, err = s.sut.GetKeyValueDescriptionsForFilter(filter1)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), *keyId, *data[0].KeyId)

	result, err := s.sut.GetKeyValueDescriptionFoKeyId(*keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	data, err = s.sut.GetKeyValueDescriptionsForFilter(filter2)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), data)

	desc.KeyName = filter2.KeyName
	keyId = s.sut.AddKeyValueDescription(desc)
	assert.NotNil(s.T(), keyId)
}

func (s *DeviceConfigurationSuite) Test_GetKeyValue() {
	filter1 := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
	}
	filter2 := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
	}

	data := model.DeviceConfigurationKeyValueDataType{
		IsValueChangeable: util.Ptr(false),
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(10),
		},
	}

	result, err := s.sut.GetKeyValueDataForFilter(filter1)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	desc := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   filter1.KeyName,
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
	}
	keyId := s.sut.AddKeyValueDescription(desc)
	assert.NotNil(s.T(), keyId)

	err = s.sut.UpdateKeyValueDataForFilter(data, nil, filter2)
	assert.NotNil(s.T(), err)

	err = s.sut.UpdateKeyValueDataForKeyId(data, nil, *keyId)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetKeyValueDataForKeyId(*keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), *keyId, *result.KeyId)
}

func (s *DeviceConfigurationSuite) Test_UpdateKeyValueDataForFilter() {
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
	}

	data := model.DeviceConfigurationKeyValueDataType{
		IsValueChangeable: util.Ptr(false),
		Value: &model.DeviceConfigurationKeyValueValueType{
			ScaledNumber: model.NewScaledNumberType(10),
		},
	}

	err := s.sut.UpdateKeyValueDataForFilter(data, nil, filter)
	assert.NotNil(s.T(), err)

	desc := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: filter.KeyName,
	}

	keyId := s.sut.AddKeyValueDescription(desc)
	assert.NotNil(s.T(), keyId)

	err = s.sut.UpdateKeyValueDataForFilter(data, nil, filter)
	assert.Nil(s.T(), err)

	result, err := s.sut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result.KeyId)
	assert.Equal(s.T(), *keyId, *result.KeyId)
	assert.NotNil(s.T(), result.Value)
	assert.NotNil(s.T(), result.Value.ScaledNumber)
	assert.Equal(s.T(), 10.0, result.Value.ScaledNumber.GetValue())

	err = s.sut.UpdateKeyValueDataForFilter(data, nil, filter)
	assert.Nil(s.T(), err)

	deleteElements := &model.DeviceConfigurationKeyValueDataElementsType{
		Value: &model.DeviceConfigurationKeyValueValueElementsType{},
	}
	data = model.DeviceConfigurationKeyValueDataType{}
	err = s.sut.UpdateKeyValueDataForFilter(data, deleteElements, filter)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result.KeyId)
	assert.Equal(s.T(), *keyId, *result.KeyId)
	assert.Nil(s.T(), result.Value)
}
