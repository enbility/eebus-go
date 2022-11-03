package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
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

	supported := util.Implements[T, model.UpdaterFactory[T]]()
	if !supported {
		return NewErrorTypeFromString(fmt.Sprintf("partial updates are not supported for type '%s'", util.Type[T]().Name()))
	}

	if r.data == nil {
		r.data = new(T)
	}

	updater := any(r.data).(model.UpdaterFactory[T])
	updater.NewUpdater(newData, filterPartial, filterDelete).DoUpdate()
	return nil
}

func (r *FunctionDataImpl[T]) DataAny() any {
	return r.Data()
}

func (r *FunctionDataImpl[T]) UpdateDataAny(newData any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	err := r.UpdateData(newData.(*T), filterPartial, filterDelete)
	if err != nil {
		// TODO: log error
		log.Error(err.String())
	}
}
