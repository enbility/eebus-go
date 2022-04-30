package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type DeviceImpl struct {
	address        model.AddressDeviceType
	dType          model.DeviceTypeType
	useCaseManager *UseCaseManager
}

func NewDeviceImpl(address model.AddressDeviceType, dType model.DeviceTypeType) *DeviceImpl {
	return &DeviceImpl{
		address:        address,
		dType:          dType,
		useCaseManager: NewUseCaseManager(),
	}
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
