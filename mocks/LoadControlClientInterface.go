// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	model "github.com/enbility/spine-go/model"
	mock "github.com/stretchr/testify/mock"
)

// LoadControlClientInterface is an autogenerated mock type for the LoadControlClientInterface type
type LoadControlClientInterface struct {
	mock.Mock
}

type LoadControlClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *LoadControlClientInterface) EXPECT() *LoadControlClientInterface_Expecter {
	return &LoadControlClientInterface_Expecter{mock: &_m.Mock}
}

// RequestLimitConstraints provides a mock function with given fields: selector, elements
func (_m *LoadControlClientInterface) RequestLimitConstraints(selector *model.LoadControlLimitConstraintsListDataSelectorsType, elements *model.LoadControlLimitConstraintsDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestLimitConstraints")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitConstraintsListDataSelectorsType, *model.LoadControlLimitConstraintsDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitConstraintsListDataSelectorsType, *model.LoadControlLimitConstraintsDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.LoadControlLimitConstraintsListDataSelectorsType, *model.LoadControlLimitConstraintsDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadControlClientInterface_RequestLimitConstraints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestLimitConstraints'
type LoadControlClientInterface_RequestLimitConstraints_Call struct {
	*mock.Call
}

// RequestLimitConstraints is a helper method to define mock.On call
//   - selector *model.LoadControlLimitConstraintsListDataSelectorsType
//   - elements *model.LoadControlLimitConstraintsDataElementsType
func (_e *LoadControlClientInterface_Expecter) RequestLimitConstraints(selector interface{}, elements interface{}) *LoadControlClientInterface_RequestLimitConstraints_Call {
	return &LoadControlClientInterface_RequestLimitConstraints_Call{Call: _e.mock.On("RequestLimitConstraints", selector, elements)}
}

func (_c *LoadControlClientInterface_RequestLimitConstraints_Call) Run(run func(selector *model.LoadControlLimitConstraintsListDataSelectorsType, elements *model.LoadControlLimitConstraintsDataElementsType)) *LoadControlClientInterface_RequestLimitConstraints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.LoadControlLimitConstraintsListDataSelectorsType), args[1].(*model.LoadControlLimitConstraintsDataElementsType))
	})
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitConstraints_Call) Return(_a0 *model.MsgCounterType, _a1 error) *LoadControlClientInterface_RequestLimitConstraints_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitConstraints_Call) RunAndReturn(run func(*model.LoadControlLimitConstraintsListDataSelectorsType, *model.LoadControlLimitConstraintsDataElementsType) (*model.MsgCounterType, error)) *LoadControlClientInterface_RequestLimitConstraints_Call {
	_c.Call.Return(run)
	return _c
}

// RequestLimitData provides a mock function with given fields: selector, elements
func (_m *LoadControlClientInterface) RequestLimitData(selector *model.LoadControlLimitListDataSelectorsType, elements *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestLimitData")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadControlClientInterface_RequestLimitData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestLimitData'
type LoadControlClientInterface_RequestLimitData_Call struct {
	*mock.Call
}

// RequestLimitData is a helper method to define mock.On call
//   - selector *model.LoadControlLimitListDataSelectorsType
//   - elements *model.LoadControlLimitDataElementsType
func (_e *LoadControlClientInterface_Expecter) RequestLimitData(selector interface{}, elements interface{}) *LoadControlClientInterface_RequestLimitData_Call {
	return &LoadControlClientInterface_RequestLimitData_Call{Call: _e.mock.On("RequestLimitData", selector, elements)}
}

func (_c *LoadControlClientInterface_RequestLimitData_Call) Run(run func(selector *model.LoadControlLimitListDataSelectorsType, elements *model.LoadControlLimitDataElementsType)) *LoadControlClientInterface_RequestLimitData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.LoadControlLimitListDataSelectorsType), args[1].(*model.LoadControlLimitDataElementsType))
	})
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitData_Call) Return(_a0 *model.MsgCounterType, _a1 error) *LoadControlClientInterface_RequestLimitData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitData_Call) RunAndReturn(run func(*model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error)) *LoadControlClientInterface_RequestLimitData_Call {
	_c.Call.Return(run)
	return _c
}

// RequestLimitDescriptions provides a mock function with given fields: selector, elements
func (_m *LoadControlClientInterface) RequestLimitDescriptions(selector *model.LoadControlLimitDescriptionListDataSelectorsType, elements *model.LoadControlLimitDescriptionDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestLimitDescriptions")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitDescriptionListDataSelectorsType, *model.LoadControlLimitDescriptionDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.LoadControlLimitDescriptionListDataSelectorsType, *model.LoadControlLimitDescriptionDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.LoadControlLimitDescriptionListDataSelectorsType, *model.LoadControlLimitDescriptionDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadControlClientInterface_RequestLimitDescriptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestLimitDescriptions'
type LoadControlClientInterface_RequestLimitDescriptions_Call struct {
	*mock.Call
}

// RequestLimitDescriptions is a helper method to define mock.On call
//   - selector *model.LoadControlLimitDescriptionListDataSelectorsType
//   - elements *model.LoadControlLimitDescriptionDataElementsType
func (_e *LoadControlClientInterface_Expecter) RequestLimitDescriptions(selector interface{}, elements interface{}) *LoadControlClientInterface_RequestLimitDescriptions_Call {
	return &LoadControlClientInterface_RequestLimitDescriptions_Call{Call: _e.mock.On("RequestLimitDescriptions", selector, elements)}
}

func (_c *LoadControlClientInterface_RequestLimitDescriptions_Call) Run(run func(selector *model.LoadControlLimitDescriptionListDataSelectorsType, elements *model.LoadControlLimitDescriptionDataElementsType)) *LoadControlClientInterface_RequestLimitDescriptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.LoadControlLimitDescriptionListDataSelectorsType), args[1].(*model.LoadControlLimitDescriptionDataElementsType))
	})
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitDescriptions_Call) Return(_a0 *model.MsgCounterType, _a1 error) *LoadControlClientInterface_RequestLimitDescriptions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LoadControlClientInterface_RequestLimitDescriptions_Call) RunAndReturn(run func(*model.LoadControlLimitDescriptionListDataSelectorsType, *model.LoadControlLimitDescriptionDataElementsType) (*model.MsgCounterType, error)) *LoadControlClientInterface_RequestLimitDescriptions_Call {
	_c.Call.Return(run)
	return _c
}

// WriteLimitData provides a mock function with given fields: data, deleteSelectors, deleteElements
func (_m *LoadControlClientInterface) WriteLimitData(data []model.LoadControlLimitDataType, deleteSelectors *model.LoadControlLimitListDataSelectorsType, deleteElements *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(data, deleteSelectors, deleteElements)

	if len(ret) == 0 {
		panic("no return value specified for WriteLimitData")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func([]model.LoadControlLimitDataType, *model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(data, deleteSelectors, deleteElements)
	}
	if rf, ok := ret.Get(0).(func([]model.LoadControlLimitDataType, *model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(data, deleteSelectors, deleteElements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func([]model.LoadControlLimitDataType, *model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) error); ok {
		r1 = rf(data, deleteSelectors, deleteElements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadControlClientInterface_WriteLimitData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteLimitData'
type LoadControlClientInterface_WriteLimitData_Call struct {
	*mock.Call
}

// WriteLimitData is a helper method to define mock.On call
//   - data []model.LoadControlLimitDataType
//   - deleteSelectors *model.LoadControlLimitListDataSelectorsType
//   - deleteElements *model.LoadControlLimitDataElementsType
func (_e *LoadControlClientInterface_Expecter) WriteLimitData(data interface{}, deleteSelectors interface{}, deleteElements interface{}) *LoadControlClientInterface_WriteLimitData_Call {
	return &LoadControlClientInterface_WriteLimitData_Call{Call: _e.mock.On("WriteLimitData", data, deleteSelectors, deleteElements)}
}

func (_c *LoadControlClientInterface_WriteLimitData_Call) Run(run func(data []model.LoadControlLimitDataType, deleteSelectors *model.LoadControlLimitListDataSelectorsType, deleteElements *model.LoadControlLimitDataElementsType)) *LoadControlClientInterface_WriteLimitData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]model.LoadControlLimitDataType), args[1].(*model.LoadControlLimitListDataSelectorsType), args[2].(*model.LoadControlLimitDataElementsType))
	})
	return _c
}

func (_c *LoadControlClientInterface_WriteLimitData_Call) Return(_a0 *model.MsgCounterType, _a1 error) *LoadControlClientInterface_WriteLimitData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LoadControlClientInterface_WriteLimitData_Call) RunAndReturn(run func([]model.LoadControlLimitDataType, *model.LoadControlLimitListDataSelectorsType, *model.LoadControlLimitDataElementsType) (*model.MsgCounterType, error)) *LoadControlClientInterface_WriteLimitData_Call {
	_c.Call.Return(run)
	return _c
}

// NewLoadControlClientInterface creates a new instance of LoadControlClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoadControlClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoadControlClientInterface {
	mock := &LoadControlClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}