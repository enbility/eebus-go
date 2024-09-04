package evcc

import (
	"testing"

	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemEVCCSuite) Test_ChargeState() {
	data, err := s.sut.ChargeState(s.mockRemoteEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeUnplugged, data)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeUnknown, data)

	stateData := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeActive, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypePaused, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFailure),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeError, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFinished),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeFinished, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeInAlarm),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ucapi.EVChargeStateTypeUnknown, data)
}

func (s *CemEVCCSuite) Test_EVConnected() {
	data := s.sut.EVConnected(nil)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.mockRemoteEntity)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.evEntity)
	assert.Equal(s.T(), false, data)

	stateData := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.EVConnected(s.evEntity)
	assert.Equal(s.T(), true, data)
}

func (s *CemEVCCSuite) Test_EVCommunicationStandard() {
	data, err := s.sut.CommunicationStandard(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), EVCCCommunicationStandardUnknown, data)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), EVCCCommunicationStandardUnknown, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), EVCCCommunicationStandardUnknown, data)

	descData = &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeString),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), EVCCCommunicationStandardUnknown, data)

	devData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringTypeISO151182ED2),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, devData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.DeviceConfigurationKeyValueStringTypeISO151182ED2, data)
}

func (s *CemEVCCSuite) Test_EVAsymmetricChargingSupport() {
	data, err := s.sut.AsymmetricChargingSupport(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.AsymmetricChargingSupport(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupport(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData = &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeBoolean),
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupport(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	devData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceConfigurationKeyValueListData, devData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupport(s.evEntity)
	assert.Nil(s.T(), err)
	assert.True(s.T(), data)
}

func (s *CemEVCCSuite) Test_EVIdentification() {
	data, err := s.sut.Identifications(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []ucapi.IdentificationItem(nil), data)

	data, err = s.sut.Identifications(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []ucapi.IdentificationItem(nil), data)

	data, err = s.sut.Identifications(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []ucapi.IdentificationItem(nil), data)

	idData := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(0)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIdentification, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeIdentificationListData, idData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Identifications(s.evEntity)
	assert.Nil(s.T(), err)
	resultData := []ucapi.IdentificationItem{{Value: "test", ValueType: model.IdentificationTypeTypeEui64}}
	assert.Equal(s.T(), resultData, data)
}

func (s *CemEVCCSuite) Test_EVManufacturerData() {
	_, err := s.sut.ManufacturerData(nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ManufacturerData(s.evEntity)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.ManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), "", data.DeviceName)
	assert.Equal(s.T(), "", data.SerialNumber)

	descData = &model.DeviceClassificationManufacturerDataType{
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("test")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("12345")),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), "test", data.DeviceName)
	assert.Equal(s.T(), "12345", data.SerialNumber)
	assert.Equal(s.T(), "", data.BrandName)
}

func (s *CemEVCCSuite) Test_EVChargingPowerLimits() {
	minData, maxData, standByData, err := s.sut.ChargingPowerLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, minData)
	assert.Equal(s.T(), 0.0, maxData)
	assert.Equal(s.T(), 0.0, standByData)

	minData, maxData, standByData, err = s.sut.ChargingPowerLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, minData)
	assert.Equal(s.T(), 0.0, maxData)
	assert.Equal(s.T(), 0.0, standByData)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	minData, maxData, standByData, err = s.sut.ChargingPowerLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, minData)
	assert.Equal(s.T(), 0.0, maxData)
	assert.Equal(s.T(), 0.0, standByData)

	type permittedStruct struct {
		standByValue, expectedStandByValue float64
		minValue, expectedMinValue         float64
		maxValue, expectedMaxValue         float64
	}

	tests := []struct {
		name      string
		permitted permittedStruct
	}{
		{
			"IEC 3 Phase",
			permittedStruct{0.1, 0.1, 4287600, 4287600, 11433600, 11433600},
		},
		{
			"ISO15118 VW",
			permittedStruct{0.1, 0.1, 800, 800, 11433600, 11433600},
		},
		{
			"ISO15118 Taycan",
			permittedStruct{0.1, 0.1, 400, 400, 11433600, 11433600},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			item := model.ScaledNumberSetType{
				Range: []model.ScaledNumberRangeType{
					{
						Min: model.NewScaledNumberType(tc.permitted.minValue),
						Max: model.NewScaledNumberType(tc.permitted.maxValue),
					},
				},
				Value: []model.ScaledNumberType{*model.NewScaledNumberType(tc.permitted.standByValue)},
			}
			permittedData = append(permittedData, item)

			permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet:      permittedData,
			}
			dataSet = append(dataSet, permittedItem)

			permData := &model.ElectricalConnectionPermittedValueSetListDataType{
				ElectricalConnectionPermittedValueSetData: dataSet,
			}

			_, fErr := rFeature.UpdateData(true, model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
			assert.Nil(s.T(), fErr)

			minData, maxData, standByData, err = s.sut.ChargingPowerLimits(s.evEntity)
			assert.Nil(s.T(), err)

			assert.Nil(s.T(), err)
			assert.Equal(s.T(), tc.permitted.expectedMinValue, minData)
			assert.Equal(s.T(), tc.permitted.expectedMaxValue, maxData)
			assert.Equal(s.T(), tc.permitted.expectedStandByValue, standByData)
		})
	}
}

func (s *CemEVCCSuite) Test_EVInSleepMode() {
	data, err := s.sut.IsInSleepMode(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsInSleepMode(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.DeviceDiagnosisStateDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsInSleepMode(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsInSleepMode(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
