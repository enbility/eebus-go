package internal

import (
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *InternalSuite) Test_ManufacturerData() {
	_, err := ManufacturerData(nil, nil)
	assert.NotNil(s.T(), err)

	_, err = ManufacturerData(s.localEntity, s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = ManufacturerData(s.localEntity, s.monitoredEntity)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceClassificationManufacturerDataType{

		DeviceName:   util.Ptr(model.DeviceClassificationStringType("deviceName")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("deviceCode")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serialNumber")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	assert.NotNil(s.T(), rFeature)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	data, err := ManufacturerData(s.localEntity, s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), "deviceName", data.DeviceName)
	assert.Equal(s.T(), "deviceCode", data.DeviceCode)
	assert.Equal(s.T(), "serialNumber", data.SerialNumber)
	assert.Equal(s.T(), "", data.SoftwareRevision)
}
