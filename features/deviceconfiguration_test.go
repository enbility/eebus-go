package features_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
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

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	deviceConfiguration *features.DeviceConfiguration
}

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

	var err error
	s.deviceConfiguration, err = features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, s.localEntity, s.remoteEntity)
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

func (s *DeviceConfigurationSuite) Test_GetDescriptionForKeyId() {
	keyId := model.DeviceConfigurationKeyIdType(0)
	desc, err := s.deviceConfiguration.GetDescriptionForKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescription()

	desc, err = s.deviceConfiguration.GetDescriptionForKeyId(keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
}

func (s *DeviceConfigurationSuite) Test_GetDescriptionForKeyName() {
	desc, err := s.deviceConfiguration.GetDescriptionForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescription()

	desc, err = s.deviceConfiguration.GetDescriptionForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
}

func (s *DeviceConfigurationSuite) Test_GetValueForKey() {
	key := model.DeviceConfigurationKeyNameTypeCommunicationsStandard
	valueType := model.DeviceConfigurationKeyValueTypeTypeString

	value, err := s.deviceConfiguration.GetKeyValueForKeyName(key, valueType)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addDescription()

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(key, valueType)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addData()

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(key, valueType)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported, model.DeviceConfigurationKeyValueTypeTypeBoolean)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeAzimuth, model.DeviceConfigurationKeyValueTypeTypeDate)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeBatteryType, model.DeviceConfigurationKeyValueTypeTypeDateTime)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeTimeToAcDischargePowerMax, model.DeviceConfigurationKeyValueTypeTypeDuration)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeIncentivesWaitIncentiveWriteable, model.DeviceConfigurationKeyValueTypeTypeTime)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeIncentivesWaitIncentiveWriteable, model.DeviceConfigurationKeyValueTypeType("invalid"))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	value, err = s.deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameType("invalid"), model.DeviceConfigurationKeyValueTypeType("invalid"))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
}

func (s *DeviceConfigurationSuite) Test_GetValues() {
	data, err := s.deviceConfiguration.GetKeyValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.deviceConfiguration.GetKeyValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.deviceConfiguration.GetKeyValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *DeviceConfigurationSuite) addDescription() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeString),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeBoolean),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
				Unit:      util.Ptr(model.UnitOfMeasurementTypepct),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(3)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeAzimuth),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDate),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(4)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeBatteryType),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDateTime),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(5)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeTimeToAcDischargePowerMax),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(6)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeIncentivesWaitIncentiveWriteable),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeTime),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, fData, nil, nil)
}

func (s *DeviceConfigurationSuite) addData() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringType("test")),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(50),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(3)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Date: model.NewDateType("01.01.2023"),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(4)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					DateTime: model.NewDateTimeTypeFromTime(time.Now()),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(5)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Second * 4),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(6)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Time: model.NewTimeType("13:05"),
				},
			},
		},
	}
	rF.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, fData, nil, nil)
}
