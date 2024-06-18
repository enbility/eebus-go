// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	model "github.com/enbility/spine-go/model"
	mock "github.com/stretchr/testify/mock"
)

// MeasurementClientInterface is an autogenerated mock type for the MeasurementClientInterface type
type MeasurementClientInterface struct {
	mock.Mock
}

type MeasurementClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MeasurementClientInterface) EXPECT() *MeasurementClientInterface_Expecter {
	return &MeasurementClientInterface_Expecter{mock: &_m.Mock}
}

// RequestConstraints provides a mock function with given fields: selector, elements
func (_m *MeasurementClientInterface) RequestConstraints(selector *model.MeasurementConstraintsListDataSelectorsType, elements *model.MeasurementConstraintsDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestConstraints")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.MeasurementConstraintsListDataSelectorsType, *model.MeasurementConstraintsDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.MeasurementConstraintsListDataSelectorsType, *model.MeasurementConstraintsDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.MeasurementConstraintsListDataSelectorsType, *model.MeasurementConstraintsDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MeasurementClientInterface_RequestConstraints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestConstraints'
type MeasurementClientInterface_RequestConstraints_Call struct {
	*mock.Call
}

// RequestConstraints is a helper method to define mock.On call
//   - selector *model.MeasurementConstraintsListDataSelectorsType
//   - elements *model.MeasurementConstraintsDataElementsType
func (_e *MeasurementClientInterface_Expecter) RequestConstraints(selector interface{}, elements interface{}) *MeasurementClientInterface_RequestConstraints_Call {
	return &MeasurementClientInterface_RequestConstraints_Call{Call: _e.mock.On("RequestConstraints", selector, elements)}
}

func (_c *MeasurementClientInterface_RequestConstraints_Call) Run(run func(selector *model.MeasurementConstraintsListDataSelectorsType, elements *model.MeasurementConstraintsDataElementsType)) *MeasurementClientInterface_RequestConstraints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.MeasurementConstraintsListDataSelectorsType), args[1].(*model.MeasurementConstraintsDataElementsType))
	})
	return _c
}

func (_c *MeasurementClientInterface_RequestConstraints_Call) Return(_a0 *model.MsgCounterType, _a1 error) *MeasurementClientInterface_RequestConstraints_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MeasurementClientInterface_RequestConstraints_Call) RunAndReturn(run func(*model.MeasurementConstraintsListDataSelectorsType, *model.MeasurementConstraintsDataElementsType) (*model.MsgCounterType, error)) *MeasurementClientInterface_RequestConstraints_Call {
	_c.Call.Return(run)
	return _c
}

// RequestData provides a mock function with given fields: selector, elements
func (_m *MeasurementClientInterface) RequestData(selector *model.MeasurementListDataSelectorsType, elements *model.MeasurementDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestData")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MeasurementClientInterface_RequestData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestData'
type MeasurementClientInterface_RequestData_Call struct {
	*mock.Call
}

// RequestData is a helper method to define mock.On call
//   - selector *model.MeasurementListDataSelectorsType
//   - elements *model.MeasurementDataElementsType
func (_e *MeasurementClientInterface_Expecter) RequestData(selector interface{}, elements interface{}) *MeasurementClientInterface_RequestData_Call {
	return &MeasurementClientInterface_RequestData_Call{Call: _e.mock.On("RequestData", selector, elements)}
}

func (_c *MeasurementClientInterface_RequestData_Call) Run(run func(selector *model.MeasurementListDataSelectorsType, elements *model.MeasurementDataElementsType)) *MeasurementClientInterface_RequestData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.MeasurementListDataSelectorsType), args[1].(*model.MeasurementDataElementsType))
	})
	return _c
}

func (_c *MeasurementClientInterface_RequestData_Call) Return(_a0 *model.MsgCounterType, _a1 error) *MeasurementClientInterface_RequestData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MeasurementClientInterface_RequestData_Call) RunAndReturn(run func(*model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) (*model.MsgCounterType, error)) *MeasurementClientInterface_RequestData_Call {
	_c.Call.Return(run)
	return _c
}

// RequestDescriptions provides a mock function with given fields: selector, elements
func (_m *MeasurementClientInterface) RequestDescriptions(selector *model.MeasurementDescriptionListDataSelectorsType, elements *model.MeasurementDescriptionDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestDescriptions")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.MeasurementDescriptionListDataSelectorsType, *model.MeasurementDescriptionDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.MeasurementDescriptionListDataSelectorsType, *model.MeasurementDescriptionDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.MeasurementDescriptionListDataSelectorsType, *model.MeasurementDescriptionDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MeasurementClientInterface_RequestDescriptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestDescriptions'
type MeasurementClientInterface_RequestDescriptions_Call struct {
	*mock.Call
}

// RequestDescriptions is a helper method to define mock.On call
//   - selector *model.MeasurementDescriptionListDataSelectorsType
//   - elements *model.MeasurementDescriptionDataElementsType
func (_e *MeasurementClientInterface_Expecter) RequestDescriptions(selector interface{}, elements interface{}) *MeasurementClientInterface_RequestDescriptions_Call {
	return &MeasurementClientInterface_RequestDescriptions_Call{Call: _e.mock.On("RequestDescriptions", selector, elements)}
}

func (_c *MeasurementClientInterface_RequestDescriptions_Call) Run(run func(selector *model.MeasurementDescriptionListDataSelectorsType, elements *model.MeasurementDescriptionDataElementsType)) *MeasurementClientInterface_RequestDescriptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.MeasurementDescriptionListDataSelectorsType), args[1].(*model.MeasurementDescriptionDataElementsType))
	})
	return _c
}

func (_c *MeasurementClientInterface_RequestDescriptions_Call) Return(_a0 *model.MsgCounterType, _a1 error) *MeasurementClientInterface_RequestDescriptions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MeasurementClientInterface_RequestDescriptions_Call) RunAndReturn(run func(*model.MeasurementDescriptionListDataSelectorsType, *model.MeasurementDescriptionDataElementsType) (*model.MsgCounterType, error)) *MeasurementClientInterface_RequestDescriptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewMeasurementClientInterface creates a new instance of MeasurementClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMeasurementClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MeasurementClientInterface {
	mock := &MeasurementClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
