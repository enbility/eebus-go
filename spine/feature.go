package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type Feature interface {
	Address() *model.FeatureAddressType
	Type() model.FeatureTypeType
	Role() model.RoleType
}

type FeatureImpl struct {
	address     *model.FeatureAddressType
	ftype       model.FeatureTypeType
	description *model.DescriptionType
	role        model.RoleType
	operations  map[model.FunctionType]*Operations
}

var _ Feature = (*FeatureImpl)(nil)

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

func (r *FeatureImpl) Type() model.FeatureTypeType {
	return r.ftype
}

func (r *FeatureImpl) Role() model.RoleType {
	return r.role
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

func featureAddressType(id uint, entityAddress *model.EntityAddressType) *model.FeatureAddressType {
	res := model.FeatureAddressType{
		Device:  entityAddress.Device,
		Entity:  entityAddress.Entity,
		Feature: util.Ptr(model.AddressFeatureType(id)),
	}

	return &res
}
