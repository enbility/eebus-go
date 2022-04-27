package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

const DeviceInformationEntityId uint = 0

var DeviceInformationAddressEntity = []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)}

type EntityImpl struct {
	eType        model.EntityTypeType
	address      *model.EntityAddressType
	description  *model.DescriptionType
	fIdGenerator func() uint
}

func NewEntity(eType model.EntityTypeType, deviceAdress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *EntityImpl {
	return &EntityImpl{
		eType: eType,
		address: &model.EntityAddressType{
			Device: deviceAdress,
			Entity: entityAddress,
		},
		fIdGenerator: newFeatureIdGenerator(0),
	}
}

func (r *EntityImpl) Address() *model.EntityAddressType {
	return r.address
}

func (r *EntityImpl) SetDescription(d *model.DescriptionType) {
	r.description = d
}

func (r *EntityImpl) NextFeatureId() uint {
	return r.fIdGenerator()
}

func EntityAddressType(deviceAdress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *model.EntityAddressType {
	return &model.EntityAddressType{
		Device: deviceAdress,
		Entity: entityAddress,
	}
}

func newFeatureIdGenerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
