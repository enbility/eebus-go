package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type DeviceImpl struct {
	address model.AddressDeviceType
	dType   model.DeviceTypeType
}

func NewDeviceImpl(address model.AddressDeviceType) *DeviceImpl {
	return &DeviceImpl{
		address: address,
	}
}

func (r *DeviceImpl) Address() *model.AddressDeviceType {
	return &r.address
}

func (r *DeviceImpl) SetDeviceType(dType model.DeviceTypeType) {
	r.dType = dType
}
