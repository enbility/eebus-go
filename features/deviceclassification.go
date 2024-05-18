package features

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type DeviceClassification struct {
	*Feature
}

// Get a new DeviceClassification features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewDeviceClassification(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*DeviceClassification, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceClassification, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceClassification{
		Feature: feature,
	}

	return dc, nil
}

// request DeviceClassificationManufacturerData from a remote device entity
func (d *DeviceClassification) RequestManufacturerDetails() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceClassificationManufacturerData, nil, nil)
}

// get the current manufacturer details for a remote device entity
func (d *DeviceClassification) GetManufacturerDetails() (*model.DeviceClassificationManufacturerDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.DeviceClassificationManufacturerDataType](d.featureRemote, model.FunctionTypeDeviceClassificationManufacturerData)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data, nil
}
