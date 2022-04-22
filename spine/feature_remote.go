package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type FeatureRemoteImpl struct {
	*FeatureImpl
	entity          *EntityRemoteImpl
	sender          Sender
	functionDataMap map[model.FunctionType]FunctionData
}

func NewFeatureRemoteImpl(id uint, entity *EntityRemoteImpl, ftype model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	res := &FeatureRemoteImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Device().Address(), entity.Address()),
			ftype,
			role),
		entity:          entity,
		sender:          sender,
		functionDataMap: make(map[model.FunctionType]FunctionData),
	}
	for _, fd := range CreateFunctionData[FunctionData](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}

	return res
}

func (r *FeatureRemoteImpl) Data(function model.FunctionType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureRemoteImpl) SetData(function model.FunctionType, data any) {
	r.functionData(function).SetDataAny(data)

	// TODO: fire event
}

func (r *FeatureRemoteImpl) Sender() Sender {
	return r.sender
}

func (r *FeatureRemoteImpl) functionData(function model.FunctionType) FunctionData {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
