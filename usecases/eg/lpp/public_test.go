package lpp

import (
	"time"

	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EgLPPSuite) Test_LoadControlLimit() {
	data, err := s.sut.ProductionLimit(nil)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
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
				LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
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

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6000.0, data.Value)
	assert.Equal(s.T(), true, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)
}

func (s *EgLPPSuite) Test_WriteLoadControlLimit() {
	limit := ucapi.LoadLimit{
		Value:    6000,
		IsActive: true,
		Duration: 0,
	}
	_, err := s.sut.WriteProductionLimit(s.mockRemoteEntity, limit, nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit, nil)
	assert.NotNil(s.T(), err)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit, nil)
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

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit, nil)
	assert.Nil(s.T(), err)

	limit.Duration = time.Duration(time.Hour * 2)
	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit, func(result model.ResultDataType) {})
	assert.Nil(s.T(), err)
}

func (s *EgLPPSuite) Test_FailsafeProductionActivePowerLimit() {
	data, err := s.sut.FailsafeProductionActivePowerLimit(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
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

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
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

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 4000.0, data)
}

func (s *EgLPPSuite) Test_WriteFailsafeProductionActivePowerLimit() {
	_, err := s.sut.WriteFailsafeProductionActivePowerLimit(s.mockRemoteEntity, 6000)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
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

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)
}

func (s *EgLPPSuite) Test_FailsafeDurationMinimum() {
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

func (s *EgLPPSuite) Test_WriteFailsafeDurationMinimum() {
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

func (s *EgLPPSuite) Test_Heartbeat() {
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

func (s *EgLPPSuite) Test_PowerProductionNominalMax() {
	data, err := s.sut.ProductionNominalMax(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.ProductionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ProductionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData = &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       util.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualProductionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ProductionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 8000.0, data)
}
