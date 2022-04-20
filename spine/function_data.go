package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type FunctionData interface {
	Function() model.FunctionType
	DataAny() any
	SetDataAny(data any)
}

var _ FunctionData = (*FunctionDataImpl[int])(nil)

type FunctionDataImpl[T any] struct {
	functionType model.FunctionType
	data         *T
}

func NewFunctionData[T any](function model.FunctionType) *FunctionDataImpl[T] {
	return &FunctionDataImpl[T]{
		functionType: function,
	}
}

func (r *FunctionDataImpl[T]) Function() model.FunctionType {
	return r.functionType
}

func (r *FunctionDataImpl[T]) Data() *T {
	return r.data
}

func (r *FunctionDataImpl[T]) SetData(newData *T) {
	r.data = newData
}

func (r *FunctionDataImpl[T]) DataAny() any {
	return r.Data()
}

func (r *FunctionDataImpl[T]) SetDataAny(newData any) {
	r.SetData(newData.(*T))
}
