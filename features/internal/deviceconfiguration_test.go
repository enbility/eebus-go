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

func TestDeviceConfigurationSuite(t *testing.T) {
	suite.Run(t, new(DeviceConfigurationSuite))
}

type DeviceConfigurationSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.DeviceConfigurationCommon
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

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalDeviceConfiguration(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteDeviceConfiguration(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *DeviceConfigurationSuite) Test_CheckEventPayloadDataForFilter() {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit
	keyName2 := model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyName,
	}
	exists := s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	exists = s.localSut.CheckEventPayloadDataForFilter(keyName, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(keyName, filter)
	assert.False(s.T(), exists)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyName: util.Ptr(keyName),
			},
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(keyName2),
			},
		},
	}

	_, fErr := s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	fErr = s.localFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(keyData, filter)
	assert.False(s.T(), exists)

	descData = &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(keyName),
			},
		},
	}

	_, fErr = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	fErr = s.localFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = s.localSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringTypeIEC61851),
				},
			},
		},
	}

	exists = s.localSut.CheckEventPayloadDataForFilter(keyData, filter)
	assert.True(s.T(), exists)
	exists = s.remoteSut.CheckEventPayloadDataForFilter(keyData, filter)
	assert.True(s.T(), exists)
}

func (s *DeviceConfigurationSuite) Test_DescriptionForKeyId() {
	keyId := model.DeviceConfigurationKeyIdType(0)
	desc, err := s.localSut.GetKeyValueDescriptionFoKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionFoKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescription()

	desc, err = s.localSut.GetKeyValueDescriptionFoKeyId(keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionFoKeyId(keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
}

func (s *DeviceConfigurationSuite) Test_Description() {
	keyId := model.DeviceConfigurationKeyIdType(0)
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyId: util.Ptr(keyId),
	}
	desc, err := s.localSut.GetKeyValueDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescription()

	desc, err = s.localSut.GetKeyValueDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
}

func (s *DeviceConfigurationSuite) Test_GetDescriptionForKeyName() {
	keyName := model.DeviceConfigurationKeyNameTypeCommunicationsStandard
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyName,
	}
	desc, err := s.localSut.GetKeyValueDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), desc)

	s.addDescription()

	desc, err = s.localSut.GetKeyValueDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
	desc, err = s.remoteSut.GetKeyValueDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), desc)
}

func (s *DeviceConfigurationSuite) Test_GetValueForKey() {
	keyName := model.DeviceConfigurationKeyNameTypeCommunicationsStandard
	valueType := model.DeviceConfigurationKeyValueTypeTypeString
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   &keyName,
		ValueType: &valueType,
	}

	value, err := s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addDescription()

	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addData()

	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeBoolean)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeAzimuth)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDate)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeBatteryType)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDateTime)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeTimeToAcDischargePowerMax)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeIncentivesWaitIncentiveWriteable)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeTypeTime)
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeIncentivesWaitIncentiveWriteable)
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeType("invalid"))
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameType("invalid"))
	filter.ValueType = util.Ptr(model.DeviceConfigurationKeyValueTypeType("invalid"))
	value, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
	value, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
}

func (s *DeviceConfigurationSuite) Test_GetData() {
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{}

	data, err := s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *DeviceConfigurationSuite) Test_GetDataForKeyId() {
	keyId := model.DeviceConfigurationKeyIdType(0)

	data, err := s.localSut.GetKeyValueDataForKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addDescription()

	data, err = s.localSut.GetKeyValueDataForKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForKeyId(keyId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.localSut.GetKeyValueDataForKeyId(keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	data, err = s.remoteSut.GetKeyValueDataForKeyId(keyId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

// helper

func (s *DeviceConfigurationSuite) addDescription() {
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
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, fData, nil, nil)
	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, fData, nil, nil)
}

func (s *DeviceConfigurationSuite) addData() {
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
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, fData, nil, nil)
	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, fData, nil, nil)
}
