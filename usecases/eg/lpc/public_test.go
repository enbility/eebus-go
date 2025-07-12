package lpc

import (
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EgLPCSuite) Test_ConsumptionLimit() {
	data, err := s.sut.ConsumptionLimit(nil)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(6000),
				TimePeriod: &model.TimePeriodType{
					EndTime: model.NewAbsoluteOrRelativeTimeType("PT2H"),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6000.0, data.Value)
	assert.Equal(s.T(), true, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)
}

func (s *EgLPCSuite) Test_WriteLoadControlLimit() {
	limit := ucapi.LoadLimit{
		Value:    6000,
		IsActive: true,
		Duration: 0,
	}
	_, err := s.sut.WriteConsumptionLimit(s.mockRemoteEntity, limit, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteConsumptionLimit(s.monitoredEntity, limit, nil)
	assert.NotNil(s.T(), err)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteConsumptionLimit(s.monitoredEntity, limit, nil)
	assert.NotNil(s.T(), err)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(6000),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteConsumptionLimit(s.monitoredEntity, limit, nil)
	assert.Nil(s.T(), err)

	limit.Duration = time.Duration(time.Hour * 2)
	_, err = s.sut.WriteConsumptionLimit(s.monitoredEntity, limit, func(result model.ResultDataType) {})
	assert.Nil(s.T(), err)
}

func (s *EgLPCSuite) Test_FailsafeConsumptionActivePowerLimit() {
	data, err := s.sut.FailsafeConsumptionActivePowerLimit(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(4000),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 4000.0, data)
}

func (s *EgLPCSuite) Test_WriteFailsafeConsumptionActivePowerLimit() {
	_, err := s.sut.WriteFailsafeConsumptionActivePowerLimit(s.mockRemoteEntity, 6000)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)
}

func (s *EgLPCSuite) Test_FailsafeDurationMinimum() {
	data, err := s.sut.FailsafeDurationMinimum(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 2),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Duration(time.Hour*2), data)
}

func (s *EgLPCSuite) Test_WriteFailsafeDurationMinimum() {
	_, err := s.sut.WriteFailsafeDurationMinimum(s.mockRemoteEntity, time.Duration(time.Hour*2))
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.NotNil(s.T(), err)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.Nil(s.T(), err)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*1))
	assert.NotNil(s.T(), err)
}

func (s *EgLPCSuite) Test_Heartbeat() {
	remoteDiagServer := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), remoteDiagServer)

	var err error
	heartbeatDiag, err := client.NewDeviceDiagnosis(s.sut.LocalEntity, s.monitoredEntity)
	assert.NotNil(s.T(), heartbeatDiag)
	assert.Nil(s.T(), err)

	// add heartbeat data to the remoteDiagServer
	timestamp := time.Now().Add(-time.Second * 121)
	data := &model.DeviceDiagnosisHeartbeatDataType{
		Timestamp:        model.NewAbsoluteOrRelativeTimeTypeFromTime(timestamp),
		HeartbeatCounter: util.Ptr(uint64(1)),
		HeartbeatTimeout: model.NewDurationType(time.Second * 120),
	}
	_, err1 := remoteDiagServer.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	assert.Nil(s.T(), err1)

	value := s.sut.IsHeartbeatWithinDuration(s.monitoredEntity)
	assert.False(s.T(), value)

	timestamp = time.Now()
	data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(timestamp)

	_, err1 = remoteDiagServer.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	assert.Nil(s.T(), err1)

	value = s.sut.IsHeartbeatWithinDuration(s.monitoredEntity)
	assert.True(s.T(), value)

	s.sut.StopHeartbeat()
	s.sut.StartHeartbeat()
}

func (s *EgLPCSuite) Test_PowerConsumptionNominalMax() {
	data, err := s.sut.ConsumptionNominalMax(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 8000.0, data)
}

// Integration tests for EG LPC validation with malformed data from CS

func (s *EgLPCSuite) Test_ConsumptionLimit_ValidationErrors() {
	// Setup LoadControl description first
	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test 1: Missing LimitId - should return ErrDataInvalid when validator runs
	// Note: Data must include Value to pass the `value.Value == nil` check in public.go:59
	invalidLimitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				// LimitId missing - this should trigger validation error
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(6000), // Include value to pass initial check
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, invalidLimitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.ConsumptionLimit(s.monitoredEntity)
	// This will actually return ErrDataNotAvailable because GetLimitDataForId won't find matching ID
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	// Test 2: Valid ID but missing Value - should return ErrDataNotAvailable from public.go:59
	invalidLimitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				// Value missing - will cause ErrDataNotAvailable from public.go:59
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, invalidLimitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data.Value)

	// Test 3: Negative Value - per spec, EG should accept values from CS without range validation
	invalidLimitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(-5000), // Negative value
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, invalidLimitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept negative value per spec
	assert.Equal(s.T(), -5000.0, data.Value)

	// Test 4: Excessive Value - per spec, EG should accept values from CS without range validation
	invalidLimitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(2000000), // Excessive value (2MW > 1MW limit)
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, invalidLimitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept excessive value per spec
	assert.Equal(s.T(), 2000000.0, data.Value)

	// Test 5: Valid data after validation errors - should work correctly
	validLimitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: util.Ptr(true),
				IsLimitActive:     util.Ptr(false),
				Value:             model.NewScaledNumberType(6000),
				TimePeriod: &model.TimePeriodType{
					EndTime: model.NewAbsoluteOrRelativeTimeType("PT2H"),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitListData, validLimitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6000.0, data.Value)
	assert.Equal(s.T(), true, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)
}

func (s *EgLPCSuite) Test_FailsafeConsumptionActivePowerLimit_ValidationErrors() {
	// Setup DeviceConfiguration description first
	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test 1: Missing KeyId - should return ErrDataNotAvailable (can't match filter)
	invalidKeyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				// KeyId missing - GetKeyValueDataForFilter won't find matching data
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(4000),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data)

	// Test 2: Missing Value - should return ErrDataNotAvailable (public.go:129 check)
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				// Value missing - public.go:129 will return ErrDataNotAvailable
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data)

	// Test 3: Empty Value object - should return ErrDataNotAvailable (public.go:129 check)
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					// No ScaledNumber - public.go:129 will return ErrDataNotAvailable
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), 0.0, data)

	// Test 4: Negative ScaledNumber value - per spec, EG should accept values from CS without range validation
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(-1000), // Negative value
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept negative value per spec
	assert.Equal(s.T(), -1000.0, data)

	// Test 5: Excessive ScaledNumber value - per spec, EG should accept values from CS without range validation
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(150000), // Excessive value (>100kW)
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept excessive value per spec
	assert.Equal(s.T(), 150000.0, data)

	// Test 6: Valid data after validation errors - should work correctly
	validKeyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(4000),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, validKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 4000.0, data)
}

func (s *EgLPCSuite) Test_FailsafeDurationMinimum_ValidationErrors() {
	// Setup DeviceConfiguration description first
	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	// Test 1: Missing KeyId - should return ErrDataNotAvailable (can't match filter)
	invalidKeyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				// KeyId missing - GetKeyValueDataForFilter won't find matching data
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 2),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), time.Duration(0), data)

	// Test 2: Missing Value - should return ErrDataNotAvailable (public.go:200 check)
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				// Value missing - public.go:200 will return ErrDataNotAvailable
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), time.Duration(0), data)

	// Test 3: Empty Value object - should return ErrDataNotAvailable (public.go:200 check)
	invalidKeyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					// No Duration - public.go:200 will return ErrDataNotAvailable
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, invalidKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err)
	assert.Equal(s.T(), time.Duration(0), data)

	// Test 4: Duration below 2h - should be accepted when reading (EG trusts CS data)
	shortDurationData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 1), // 1 hour
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, shortDurationData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Hour*1, data)

	// Test 5: Duration above 24h - should be accepted when reading (EG trusts CS data)
	longDurationData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 25), // 25 hours
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, longDurationData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Hour*25, data)

	// Test 6: Valid data after validation errors - should work correctly
	validKeyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 2),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, validKeyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Duration(time.Hour*2), data)
}

func (s *EgLPCSuite) Test_ConsumptionNominalMax_ValidationErrors() {
	// Test 1: Missing CharacteristicId - should return ErrDataInvalid
	invalidCharData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				// CharacteristicId missing
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax), // Use correct type for EMS
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, invalidCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataInvalid, err)
	assert.Equal(s.T(), 0.0, data)

	// Test 2: Wrong CharacteristicType - should return ErrDataNotAvailable (filtered out)
	invalidCharData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax), // Wrong type - won't match filter
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, invalidCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err) // Filtered out, so no data found
	assert.Equal(s.T(), 0.0, data)

	// Test 3: Wrong CharacteristicContext - should return ErrDataNotAvailable (filtered out)
	invalidCharData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeDevice), // Wrong context - won't match filter
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, invalidCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Equal(s.T(), api.ErrDataNotAvailable, err) // Filtered out, so no data found
	assert.Equal(s.T(), 0.0, data)

	// Test 4: Negative Value - per spec, EG should accept values from CS without range validation
	invalidCharData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax), // Use correct type for EMS
				Value:                  model.NewScaledNumberType(-5000), // Negative value
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, invalidCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept negative value per spec
	assert.Equal(s.T(), -5000.0, data)

	// Test 5: Excessive Value - per spec, EG should accept values from CS without range validation
	invalidCharData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax), // Use correct type for EMS
				Value:                  model.NewScaledNumberType(2000000), // Excessive value (2MW > 1MW limit)
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, invalidCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err) // Should accept excessive value per spec
	assert.Equal(s.T(), 2000000.0, data)

	// Test 6: Valid ContractualConsumptionNominalMax data - should work correctly
	// Since the test device has no device type (nil), it uses ContractualConsumptionNominalMax
	validCharData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, validCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 8000.0, data)

	// Test 7: Valid ContractualConsumptionNominalMax data - should work correctly
	// Since the test device has no device type (nil), it uses ContractualConsumptionNominalMax
	validCharData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(12000),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, validCharData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 12000.0, data)
}
