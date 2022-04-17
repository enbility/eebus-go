package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type FunctionDataImpl[T any] struct {
	functionType model.FunctionEnumType
	data         *T
}

func NewFunctionData[T any](function model.FunctionEnumType) *FunctionDataImpl[T] {
	return &FunctionDataImpl[T]{
		functionType: function,
	}
}

func (r *FunctionDataImpl[T]) Function() model.FunctionEnumType {
	return r.functionType
}

func (r *FunctionDataImpl[T]) SetData(newData *T) {
	r.data = newData
}

func (r *FunctionDataImpl[T]) Data() *T {
	return r.data
}
