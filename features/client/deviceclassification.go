package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceClassification struct {
	*Feature

	*internal.DeviceClassificationCommon
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
		Feature:                    feature,
		DeviceClassificationCommon: internal.NewRemoteDeviceClassification(feature.featureRemote),
	}

	return dc, nil
}

var _ api.DeviceClassificationClientInterface = (*DeviceClassification)(nil)

// request DeviceClassificationManufacturerData from a remote device entity
func (d *DeviceClassification) RequestManufacturerDetails() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceClassificationManufacturerData, nil, nil)
}
