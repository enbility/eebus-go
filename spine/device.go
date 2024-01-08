package spine

import "github.com/enbility/eebus-go/spine/model"

type DeviceImpl struct {
	address    *model.AddressDeviceType
	dType      *model.DeviceTypeType
	featureSet *model.NetworkManagementFeatureSetType
}

// Initialize a new device
// Both values are required for a local device but provided as empty strings for a remote device
// as the address is only provided via detailed discovery response
func NewDeviceImpl(address *model.AddressDeviceType, dType *model.DeviceTypeType, featureSet *model.NetworkManagementFeatureSetType) *DeviceImpl {
	deviceImpl := &DeviceImpl{}

	if dType != nil {
		deviceImpl.dType = dType
	}

	if address != nil {
		deviceImpl.address = address
	}

	if featureSet != nil {
		deviceImpl.featureSet = featureSet
	}

	return deviceImpl
}

func (r *DeviceImpl) Address() *model.AddressDeviceType {
	return r.address
}

func (r *DeviceImpl) DeviceType() *model.DeviceTypeType {
	return r.dType
}

func (r *DeviceImpl) FeatureSet() *model.NetworkManagementFeatureSetType {
	return r.featureSet
}

func (r *DeviceImpl) DestinationData() model.NodeManagementDestinationDataType {
	return model.NodeManagementDestinationDataType{
		DeviceDescription: &model.NetworkManagementDeviceDescriptionDataType{
			DeviceAddress: &model.DeviceAddressType{
				Device: r.Address(),
			},
			DeviceType:        r.DeviceType(),
			NetworkFeatureSet: r.FeatureSet(),
		},
	}
}
