package spine

import (
	"fmt"

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

func (r *FeatureImpl) Operations() map[model.FunctionType]*Operations {
	return r.operations
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

func (r *FeatureImpl) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("Id: %d (%s)", *r.Address().Feature, *r.Description())
}

func featureAddressType(id uint, entityAddress *model.EntityAddressType) *model.FeatureAddressType {
	res := model.FeatureAddressType{
		Device:  entityAddress.Device,
		Entity:  entityAddress.Entity,
		Feature: util.Ptr(model.AddressFeatureType(id)),
	}

	return &res
}
