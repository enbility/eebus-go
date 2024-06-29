package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceClassificationCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalDeviceClassification(featureLocal spineapi.FeatureLocalInterface) *DeviceClassificationCommon {
	return &DeviceClassificationCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteDeviceClassification(featureRemote spineapi.FeatureRemoteInterface) *DeviceClassificationCommon {
	return &DeviceClassificationCommon{
		featureRemote: featureRemote,
	}
}

var _ api.DeviceClassificationCommonInterface = (*DeviceClassificationCommon)(nil)

// get the current manufacturer details for a remote device entity
func (d *DeviceClassificationCommon) GetManufacturerDetails() (*model.DeviceClassificationManufacturerDataType, error) {
	function := model.FunctionTypeDeviceClassificationManufacturerData
	data, err := featureDataCopyOfType[model.DeviceClassificationManufacturerDataType](d.featureLocal, d.featureRemote, function)
	if err != nil || data == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data, nil
}
