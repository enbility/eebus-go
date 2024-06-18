// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	model "github.com/enbility/spine-go/model"
	mock "github.com/stretchr/testify/mock"
)

// IncentiveTableCommonInterface is an autogenerated mock type for the IncentiveTableCommonInterface type
type IncentiveTableCommonInterface struct {
	mock.Mock
}

type IncentiveTableCommonInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *IncentiveTableCommonInterface) EXPECT() *IncentiveTableCommonInterface_Expecter {
	return &IncentiveTableCommonInterface_Expecter{mock: &_m.Mock}
}

// GetConstraints provides a mock function with given fields:
func (_m *IncentiveTableCommonInterface) GetConstraints() ([]model.IncentiveTableConstraintsType, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConstraints")
	}

	var r0 []model.IncentiveTableConstraintsType
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.IncentiveTableConstraintsType, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.IncentiveTableConstraintsType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncentiveTableConstraintsType)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncentiveTableCommonInterface_GetConstraints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConstraints'
type IncentiveTableCommonInterface_GetConstraints_Call struct {
	*mock.Call
}

// GetConstraints is a helper method to define mock.On call
func (_e *IncentiveTableCommonInterface_Expecter) GetConstraints() *IncentiveTableCommonInterface_GetConstraints_Call {
	return &IncentiveTableCommonInterface_GetConstraints_Call{Call: _e.mock.On("GetConstraints")}
}

func (_c *IncentiveTableCommonInterface_GetConstraints_Call) Run(run func()) *IncentiveTableCommonInterface_GetConstraints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IncentiveTableCommonInterface_GetConstraints_Call) Return(_a0 []model.IncentiveTableConstraintsType, _a1 error) *IncentiveTableCommonInterface_GetConstraints_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IncentiveTableCommonInterface_GetConstraints_Call) RunAndReturn(run func() ([]model.IncentiveTableConstraintsType, error)) *IncentiveTableCommonInterface_GetConstraints_Call {
	_c.Call.Return(run)
	return _c
}

// GetData provides a mock function with given fields:
func (_m *IncentiveTableCommonInterface) GetData() ([]model.IncentiveTableType, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetData")
	}

	var r0 []model.IncentiveTableType
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.IncentiveTableType, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.IncentiveTableType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncentiveTableType)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncentiveTableCommonInterface_GetData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetData'
type IncentiveTableCommonInterface_GetData_Call struct {
	*mock.Call
}

// GetData is a helper method to define mock.On call
func (_e *IncentiveTableCommonInterface_Expecter) GetData() *IncentiveTableCommonInterface_GetData_Call {
	return &IncentiveTableCommonInterface_GetData_Call{Call: _e.mock.On("GetData")}
}

func (_c *IncentiveTableCommonInterface_GetData_Call) Run(run func()) *IncentiveTableCommonInterface_GetData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IncentiveTableCommonInterface_GetData_Call) Return(_a0 []model.IncentiveTableType, _a1 error) *IncentiveTableCommonInterface_GetData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IncentiveTableCommonInterface_GetData_Call) RunAndReturn(run func() ([]model.IncentiveTableType, error)) *IncentiveTableCommonInterface_GetData_Call {
	_c.Call.Return(run)
	return _c
}

// GetDescriptionsForFilter provides a mock function with given fields: filter
func (_m *IncentiveTableCommonInterface) GetDescriptionsForFilter(filter model.TariffDescriptionDataType) ([]model.IncentiveTableDescriptionType, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for GetDescriptionsForFilter")
	}

	var r0 []model.IncentiveTableDescriptionType
	var r1 error
	if rf, ok := ret.Get(0).(func(model.TariffDescriptionDataType) ([]model.IncentiveTableDescriptionType, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(model.TariffDescriptionDataType) []model.IncentiveTableDescriptionType); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncentiveTableDescriptionType)
		}
	}

	if rf, ok := ret.Get(1).(func(model.TariffDescriptionDataType) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncentiveTableCommonInterface_GetDescriptionsForFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDescriptionsForFilter'
type IncentiveTableCommonInterface_GetDescriptionsForFilter_Call struct {
	*mock.Call
}

// GetDescriptionsForFilter is a helper method to define mock.On call
//   - filter model.TariffDescriptionDataType
func (_e *IncentiveTableCommonInterface_Expecter) GetDescriptionsForFilter(filter interface{}) *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call {
	return &IncentiveTableCommonInterface_GetDescriptionsForFilter_Call{Call: _e.mock.On("GetDescriptionsForFilter", filter)}
}

func (_c *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call) Run(run func(filter model.TariffDescriptionDataType)) *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.TariffDescriptionDataType))
	})
	return _c
}

func (_c *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call) Return(_a0 []model.IncentiveTableDescriptionType, _a1 error) *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call) RunAndReturn(run func(model.TariffDescriptionDataType) ([]model.IncentiveTableDescriptionType, error)) *IncentiveTableCommonInterface_GetDescriptionsForFilter_Call {
	_c.Call.Return(run)
	return _c
}

// NewIncentiveTableCommonInterface creates a new instance of IncentiveTableCommonInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIncentiveTableCommonInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *IncentiveTableCommonInterface {
	mock := &IncentiveTableCommonInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
