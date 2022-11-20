package spine

import (
	"fmt"
	"sync"

	"github.com/DerAndereAndi/eebus-go/logging"
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

	mux sync.Mutex
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
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.data
}

func (r *FunctionDataImpl[T]) UpdateData(newData *T, filterPartial *model.FilterType, filterDelete *model.FilterType) *ErrorType {
	r.mux.Lock()
	defer r.mux.Unlock()

	if filterPartial == nil && filterDelete == nil {
		// just set the data
		r.data = newData
		return nil
	}

	supported := util.Implements[T, model.Updater]()
	if !supported {
		return NewErrorTypeFromString(fmt.Sprintf("partial updates are not supported for type '%s'", util.Type[T]().Name()))
	}

	if r.data == nil {
		r.data = new(T)
	}

	updater := any(r.data).(model.Updater)
	updater.UpdateList(newData, filterPartial, filterDelete)

	return nil
}

func (r *FunctionDataImpl[T]) DataAny() any {
	return r.Data()
}

func (r *FunctionDataImpl[T]) UpdateDataAny(newData any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	err := r.UpdateData(newData.(*T), filterPartial, filterDelete)
	if err != nil {
		logging.Log.Error(err.String())
	}
}
