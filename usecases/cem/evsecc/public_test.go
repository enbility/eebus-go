package evsecc

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CemEVSECCSuite) Test_EVSEManufacturerData() {
	_, err := s.sut.ManufacturerData(nil)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ManufacturerData(s.evseEntity)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.ManufacturerData(s.evseEntity)
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

	data, err = s.sut.ManufacturerData(s.evseEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), "test", data.DeviceName)
	assert.Equal(s.T(), "12345", data.SerialNumber)
	assert.Equal(s.T(), "", data.BrandName)
}

func (s *CemEVSECCSuite) Test_EVSEOperatingState() {
	data, errCode, err := s.sut.OperatingState(nil)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), nil, err)

	data, errCode, err = s.sut.OperatingState(s.mockRemoteEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceDiagnosisStateDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), err)

	descData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
		LastErrorCode:  util.Ptr(model.LastErrorCodeType("error")),
	}

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeStandby, data)
	assert.Equal(s.T(), "error", errCode)
	assert.Nil(s.T(), err)
}
