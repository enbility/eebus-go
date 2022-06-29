package spine

import (
	"fmt"
	"reflect"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type FunctionData interface {
	Function() model.FunctionType
	DataAny() any
	UpdateDataAny(data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
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

func (r *FunctionDataImpl[T]) UpdateData(newData *T, filterPartial *model.FilterType, filterDelete *model.FilterType) *ErrorType {
	if filterPartial == nil && filterDelete == nil {
		// just set the data
		r.data = newData
		return nil
	}

	newT := new(T)
	_, supported := any(newT).(model.Updater[T])
	if !supported {
		return NewErrorTypeFromString(fmt.Sprintf("partial updates are not supported for type '%s'", reflect.TypeOf(*newT).Name()))
	}

	if r.data == nil {
		r.data = newT
	}

	updater, _ := any(r.data).(model.Updater[T])
	updater.Update(newData, filterPartial, filterDelete)
	return nil
}

func (r *FunctionDataImpl[T]) DataAny() any {
	return r.Data()
}

func (r *FunctionDataImpl[T]) UpdateDataAny(newData any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	err := r.UpdateData(newData.(*T), filterPartial, filterDelete)
	if err != nil {
		// TODO: log error
		fmt.Print(err.String())
	}
}
