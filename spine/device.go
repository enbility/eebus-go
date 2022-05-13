package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type DeviceImpl struct {
	address        model.AddressDeviceType
	dType          model.DeviceTypeType
	useCaseManager *UseCaseManager
}

// Initialize a new device
// Both values are required for a local device but provided as empty strings for a remote device
// as the address is only provided via detailed discovery response
func NewDeviceImpl(address model.AddressDeviceType, dType model.DeviceTypeType) *DeviceImpl {
	deviceImpl := &DeviceImpl{
		useCaseManager: NewUseCaseManager(),
	}

	if dType != "" {
		deviceImpl.dType = dType
	}

	if address != "" {
		deviceImpl.address = address
	}

	return deviceImpl
}

func (r *DeviceImpl) Address() *model.AddressDeviceType {
	return &r.address
}

func (r *DeviceImpl) UseCaseManager() *UseCaseManager {
	return r.useCaseManager
}

func (r *DeviceImpl) DeviceType() model.DeviceTypeType {
	return r.dType
}
