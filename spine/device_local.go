package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type DeviceLocalImpl struct {
	*DeviceImpl
	entities []*EntityLocalImpl
}

func NewDeviceLocalImpl(address model.AddressDeviceType) *DeviceLocalImpl {
	return &DeviceLocalImpl{
		DeviceImpl: NewDeviceImpl(address),
	}
}

func (d *DeviceLocalImpl) AddEntity(entity *EntityLocalImpl) {
	d.entities = append(d.entities, entity)
	// if nmf := feature.GetNodeManagementLocal(d); nmf != nil {
	// 	nmf.NotifySubscribersOfEntity(entity, model.NetworkManagementStateChangeTypeAdded)
	// }
}
