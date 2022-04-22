package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

const DeviceInformationEntityId uint = 0

type EntityImpl struct {
	eType        model.EntityTypeType
	address      []model.AddressEntityType
	fIdGenerator func() uint
}

func NewEntity(eType model.EntityTypeType, address []model.AddressEntityType) *EntityImpl {
	return &EntityImpl{
		eType:        eType,
		address:      address,
		fIdGenerator: newFeatureIdGenerator(0),
	}
}

func (r *EntityImpl) Address() []model.AddressEntityType {
	return r.address
}

func (r *EntityImpl) NextFeatureId() uint {
	return r.fIdGenerator()
}

func newFeatureIdGenerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
