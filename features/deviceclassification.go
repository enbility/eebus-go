package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type DeviceClassification struct {
	*FeatureImpl
}

func NewDeviceClassification(
	localRole, remoteRole model.RoleType,
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*DeviceClassification, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceClassification, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceClassification{
		FeatureImpl: feature,
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
		return nil, ErrDataNotAvailable
	}

	return data, nil
}
