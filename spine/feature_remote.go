package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type FeatureRemoteImpl struct {
	functionDataMap map[model.FunctionEnumType]FunctionData
}

func NewFeatureRemoteImpl(ftype model.FeatureTypeEnumType) *FeatureRemoteImpl {
	result := &FeatureRemoteImpl{
		functionDataMap: make(map[model.FunctionEnumType]FunctionData),
	}
	for _, fd := range CreateFunctionData[FunctionData](ftype) {
		result.functionDataMap[fd.Function()] = fd
	}

	return result
}

func (r *FeatureRemoteImpl) Data(function model.FunctionEnumType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureRemoteImpl) SetData(function model.FunctionEnumType, data any) {
	r.functionData(function).SetDataAny(data)

	// TODO: fire event
}

func (r *FeatureRemoteImpl) functionData(function model.FunctionEnumType) FunctionData {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
