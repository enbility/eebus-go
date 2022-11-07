package features

import (
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type ManufacturerType struct {
	BrandName                      string
	VendorName                     string
	VendorCode                     string
	DeviceName                     string
	DeviceCode                     string
	SerialNumber                   string
	SoftwareRevision               string
	HardwareRevision               string
	PowerSource                    string
	ManufacturerNodeIdentification string
	ManufacturerLabel              string
	ManufacturerDescription        string
}

type DeviceClassification struct {
	*FeatureImpl
}

func NewDeviceClassification(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*DeviceClassification, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceClassification, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceClassification{
		FeatureImpl: feature,
	}

	return dc, nil
}

// request DeviceClassificationManufacturerData from a remote device entity
func (d *DeviceClassification) RequestManufacturerDetailsForEntity() (*model.MsgCounterType, error) {
	// request DeviceClassificationManufacturer from a remote entity
	msgCounter, err := d.requestData(model.FunctionTypeDeviceClassificationManufacturerData)
	if err != nil {
		logging.Log.Error(err)
		return nil, err
	}

	return msgCounter, nil
}

// get the current manufacturer details for a remote device entity
func (d *DeviceClassification) GetManufacturerDetails() (*ManufacturerType, error) {
	rData := d.featureRemote.Data(model.FunctionTypeDeviceClassificationManufacturerData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceClassificationManufacturerDataType)
	details := &ManufacturerType{}

	if data.BrandName != nil {
		details.BrandName = string(*data.BrandName)
	}
	if data.VendorName != nil {
		details.VendorName = string(*data.VendorName)
	}
	if data.VendorCode != nil {
		details.VendorCode = string(*data.VendorCode)
	}
	if data.DeviceName != nil {
		details.DeviceName = string(*data.DeviceName)
	}
	if data.DeviceCode != nil {
		details.DeviceCode = string(*data.DeviceCode)
	}
	if data.SerialNumber != nil {
		details.SerialNumber = string(*data.SerialNumber)
	}
	if data.SoftwareRevision != nil {
		details.SoftwareRevision = string(*data.SoftwareRevision)
	}
	if data.HardwareRevision != nil {
		details.HardwareRevision = string(*data.HardwareRevision)
	}
	if data.PowerSource != nil {
		details.PowerSource = string(*data.PowerSource)
	}
	if data.ManufacturerNodeIdentification != nil {
		details.ManufacturerNodeIdentification = string(*data.ManufacturerNodeIdentification)
	}
	if data.ManufacturerLabel != nil {
		details.ManufacturerLabel = string(*data.ManufacturerLabel)
	}
	if data.ManufacturerDescription != nil {
		details.ManufacturerDescription = string(*data.ManufacturerDescription)
	}

	return details, nil
}
