package features_test

import (
	"testing"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceConfigurationSuite(t *testing.T) {
	suite.Run(t, new(DeviceConfigurationSuite))
}

type DeviceConfigurationSuite struct {
	suite.Suite

	localEntity  *spine.EntityLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	deviceConfiguration *features.DeviceConfiguration
	sentMessage         []byte
}

var _ spine.SpineDataConnection = (*DeviceConfigurationSuite)(nil)

func (s *DeviceConfigurationSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *DeviceConfigurationSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
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
	s.deviceConfiguration, err = features.NewDeviceConfiguration(model.RoleTypeServer, model.RoleTypeClient, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceConfiguration)
}

func (s *DeviceConfigurationSuite) Test_RequestDescriptions() {
	err := s.deviceConfiguration.RequestDescriptions()
	assert.Nil(s.T(), err)
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
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
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
		},
	}
	rF.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, fData, nil, nil)
}

func (s *DeviceConfigurationSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
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
		},
	}
	rF.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, fData, nil, nil)
}
