package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type DeviceRemoteImpl struct {
	*DeviceImpl
	//	entities []*EntityRemoteImpl
}

func NewDeviceRemoteImpl(address model.AddressDeviceType) *DeviceRemoteImpl {
	return &DeviceRemoteImpl{
		DeviceImpl: NewDeviceImpl(address),
	}
}
