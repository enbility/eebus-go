package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type FeatureImpl struct {
	address     *model.FeatureAddressType
	ftype       model.FeatureTypeType
	description *model.DescriptionType
	role        model.RoleType
	//functions        map[model.FunctionEnumType]spine.RW
	//maxResponseDelay *model.MaxResponseDelayType
}

func NewFeatureImpl(address *model.FeatureAddressType, ftype model.FeatureTypeType, role model.RoleType) *FeatureImpl {
	res := &FeatureImpl{
		address: address,
		ftype:   ftype,
		role:    role,
	}

	return res
}

func (r *FeatureImpl) Address() *model.FeatureAddressType {
	return r.address
}

func (r *FeatureImpl) Description() *model.DescriptionType {
	return r.description
}

func (r *FeatureImpl) SetDescription(d *model.DescriptionType) {
	r.description = d
}

func (r *FeatureImpl) SetDescriptionString(s string) {
	r.description = util.Ptr(model.DescriptionType(s))
}

func featureAddressType(id uint, deviceAddress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *model.FeatureAddressType {
	res := model.FeatureAddressType{
		Device:  deviceAddress,
		Entity:  entityAddress,
		Feature: util.Ptr(model.AddressFeatureType(id)),
	}

	return &res
}
