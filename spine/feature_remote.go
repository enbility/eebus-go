package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type FeatureRemoteImpl struct {
	functionDataMap map[model.FunctionType]FunctionData
}

func NewFeatureRemoteImpl(ftype model.FeatureTypeType) *FeatureRemoteImpl {
	result := &FeatureRemoteImpl{
		functionDataMap: make(map[model.FunctionType]FunctionData),
	}
	for _, fd := range CreateFunctionData[FunctionData](ftype) {
		result.functionDataMap[fd.Function()] = fd
	}

	return result
}

func (r *FeatureRemoteImpl) Data(function model.FunctionType) any {
	return r.functionData(function).DataAny()
}

func (r *FeatureRemoteImpl) SetData(function model.FunctionType, data any) {
	r.functionData(function).SetDataAny(data)

	// TODO: fire event
}

func (r *FeatureRemoteImpl) functionData(function model.FunctionType) FunctionData {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
