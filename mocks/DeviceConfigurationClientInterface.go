// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	model "github.com/enbility/spine-go/model"
	mock "github.com/stretchr/testify/mock"
)

// DeviceConfigurationClientInterface is an autogenerated mock type for the DeviceConfigurationClientInterface type
type DeviceConfigurationClientInterface struct {
	mock.Mock
}

type DeviceConfigurationClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *DeviceConfigurationClientInterface) EXPECT() *DeviceConfigurationClientInterface_Expecter {
	return &DeviceConfigurationClientInterface_Expecter{mock: &_m.Mock}
}

// CheckEventPayloadDataForFilter provides a mock function with given fields: payloadData, filter
func (_m *DeviceConfigurationClientInterface) CheckEventPayloadDataForFilter(payloadData interface{}, filter interface{}) bool {
	ret := _m.Called(payloadData, filter)

	if len(ret) == 0 {
		panic("no return value specified for CheckEventPayloadDataForFilter")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) bool); ok {
		r0 = rf(payloadData, filter)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckEventPayloadDataForFilter'
type DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call struct {
	*mock.Call
}

// CheckEventPayloadDataForFilter is a helper method to define mock.On call
//   - payloadData interface{}
//   - filter interface{}
func (_e *DeviceConfigurationClientInterface_Expecter) CheckEventPayloadDataForFilter(payloadData interface{}, filter interface{}) *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call {
	return &DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call{Call: _e.mock.On("CheckEventPayloadDataForFilter", payloadData, filter)}
}

func (_c *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call) Run(run func(payloadData interface{}, filter interface{})) *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(interface{}))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call) Return(_a0 bool) *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call) RunAndReturn(run func(interface{}, interface{}) bool) *DeviceConfigurationClientInterface_CheckEventPayloadDataForFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueDataForFilter provides a mock function with given fields: filter
func (_m *DeviceConfigurationClientInterface) GetKeyValueDataForFilter(filter model.DeviceConfigurationKeyValueDescriptionDataType) (*model.DeviceConfigurationKeyValueDataType, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for GetKeyValueDataForFilter")
	}

	var r0 *model.DeviceConfigurationKeyValueDataType
	var r1 error
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyValueDescriptionDataType) (*model.DeviceConfigurationKeyValueDataType, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyValueDescriptionDataType) *model.DeviceConfigurationKeyValueDataType); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DeviceConfigurationKeyValueDataType)
		}
	}

	if rf, ok := ret.Get(1).(func(model.DeviceConfigurationKeyValueDescriptionDataType) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueDataForFilter'
type DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call struct {
	*mock.Call
}

// GetKeyValueDataForFilter is a helper method to define mock.On call
//   - filter model.DeviceConfigurationKeyValueDescriptionDataType
func (_e *DeviceConfigurationClientInterface_Expecter) GetKeyValueDataForFilter(filter interface{}) *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call {
	return &DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call{Call: _e.mock.On("GetKeyValueDataForFilter", filter)}
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call) Run(run func(filter model.DeviceConfigurationKeyValueDescriptionDataType)) *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.DeviceConfigurationKeyValueDescriptionDataType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call) Return(_a0 *model.DeviceConfigurationKeyValueDataType, _a1 error) *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call) RunAndReturn(run func(model.DeviceConfigurationKeyValueDescriptionDataType) (*model.DeviceConfigurationKeyValueDataType, error)) *DeviceConfigurationClientInterface_GetKeyValueDataForFilter_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueDataForKeyId provides a mock function with given fields: keyId
func (_m *DeviceConfigurationClientInterface) GetKeyValueDataForKeyId(keyId model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDataType, error) {
	ret := _m.Called(keyId)

	if len(ret) == 0 {
		panic("no return value specified for GetKeyValueDataForKeyId")
	}

	var r0 *model.DeviceConfigurationKeyValueDataType
	var r1 error
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDataType, error)); ok {
		return rf(keyId)
	}
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyIdType) *model.DeviceConfigurationKeyValueDataType); ok {
		r0 = rf(keyId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DeviceConfigurationKeyValueDataType)
		}
	}

	if rf, ok := ret.Get(1).(func(model.DeviceConfigurationKeyIdType) error); ok {
		r1 = rf(keyId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueDataForKeyId'
type DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call struct {
	*mock.Call
}

// GetKeyValueDataForKeyId is a helper method to define mock.On call
//   - keyId model.DeviceConfigurationKeyIdType
func (_e *DeviceConfigurationClientInterface_Expecter) GetKeyValueDataForKeyId(keyId interface{}) *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call {
	return &DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call{Call: _e.mock.On("GetKeyValueDataForKeyId", keyId)}
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call) Run(run func(keyId model.DeviceConfigurationKeyIdType)) *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.DeviceConfigurationKeyIdType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call) Return(_a0 *model.DeviceConfigurationKeyValueDataType, _a1 error) *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call) RunAndReturn(run func(model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDataType, error)) *DeviceConfigurationClientInterface_GetKeyValueDataForKeyId_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueDescriptionFoKeyId provides a mock function with given fields: keyId
func (_m *DeviceConfigurationClientInterface) GetKeyValueDescriptionFoKeyId(keyId model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	ret := _m.Called(keyId)

	if len(ret) == 0 {
		panic("no return value specified for GetKeyValueDescriptionFoKeyId")
	}

	var r0 *model.DeviceConfigurationKeyValueDescriptionDataType
	var r1 error
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error)); ok {
		return rf(keyId)
	}
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyIdType) *model.DeviceConfigurationKeyValueDescriptionDataType); ok {
		r0 = rf(keyId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.DeviceConfigurationKeyValueDescriptionDataType)
		}
	}

	if rf, ok := ret.Get(1).(func(model.DeviceConfigurationKeyIdType) error); ok {
		r1 = rf(keyId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueDescriptionFoKeyId'
type DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call struct {
	*mock.Call
}

// GetKeyValueDescriptionFoKeyId is a helper method to define mock.On call
//   - keyId model.DeviceConfigurationKeyIdType
func (_e *DeviceConfigurationClientInterface_Expecter) GetKeyValueDescriptionFoKeyId(keyId interface{}) *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call {
	return &DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call{Call: _e.mock.On("GetKeyValueDescriptionFoKeyId", keyId)}
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call) Run(run func(keyId model.DeviceConfigurationKeyIdType)) *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.DeviceConfigurationKeyIdType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call) Return(_a0 *model.DeviceConfigurationKeyValueDescriptionDataType, _a1 error) *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call) RunAndReturn(run func(model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error)) *DeviceConfigurationClientInterface_GetKeyValueDescriptionFoKeyId_Call {
	_c.Call.Return(run)
	return _c
}

// GetKeyValueDescriptionsForFilter provides a mock function with given fields: filter
func (_m *DeviceConfigurationClientInterface) GetKeyValueDescriptionsForFilter(filter model.DeviceConfigurationKeyValueDescriptionDataType) ([]model.DeviceConfigurationKeyValueDescriptionDataType, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for GetKeyValueDescriptionsForFilter")
	}

	var r0 []model.DeviceConfigurationKeyValueDescriptionDataType
	var r1 error
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyValueDescriptionDataType) ([]model.DeviceConfigurationKeyValueDescriptionDataType, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(model.DeviceConfigurationKeyValueDescriptionDataType) []model.DeviceConfigurationKeyValueDescriptionDataType); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.DeviceConfigurationKeyValueDescriptionDataType)
		}
	}

	if rf, ok := ret.Get(1).(func(model.DeviceConfigurationKeyValueDescriptionDataType) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKeyValueDescriptionsForFilter'
type DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call struct {
	*mock.Call
}

// GetKeyValueDescriptionsForFilter is a helper method to define mock.On call
//   - filter model.DeviceConfigurationKeyValueDescriptionDataType
func (_e *DeviceConfigurationClientInterface_Expecter) GetKeyValueDescriptionsForFilter(filter interface{}) *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call {
	return &DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call{Call: _e.mock.On("GetKeyValueDescriptionsForFilter", filter)}
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call) Run(run func(filter model.DeviceConfigurationKeyValueDescriptionDataType)) *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.DeviceConfigurationKeyValueDescriptionDataType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call) Return(_a0 []model.DeviceConfigurationKeyValueDescriptionDataType, _a1 error) *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call) RunAndReturn(run func(model.DeviceConfigurationKeyValueDescriptionDataType) ([]model.DeviceConfigurationKeyValueDescriptionDataType, error)) *DeviceConfigurationClientInterface_GetKeyValueDescriptionsForFilter_Call {
	_c.Call.Return(run)
	return _c
}

// RequestKeyValueDescriptions provides a mock function with given fields: selector, elements
func (_m *DeviceConfigurationClientInterface) RequestKeyValueDescriptions(selector *model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, elements *model.DeviceConfigurationKeyValueDescriptionDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestKeyValueDescriptions")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, *model.DeviceConfigurationKeyValueDescriptionDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, *model.DeviceConfigurationKeyValueDescriptionDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, *model.DeviceConfigurationKeyValueDescriptionDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestKeyValueDescriptions'
type DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call struct {
	*mock.Call
}

// RequestKeyValueDescriptions is a helper method to define mock.On call
//   - selector *model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType
//   - elements *model.DeviceConfigurationKeyValueDescriptionDataElementsType
func (_e *DeviceConfigurationClientInterface_Expecter) RequestKeyValueDescriptions(selector interface{}, elements interface{}) *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call {
	return &DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call{Call: _e.mock.On("RequestKeyValueDescriptions", selector, elements)}
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call) Run(run func(selector *model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, elements *model.DeviceConfigurationKeyValueDescriptionDataElementsType)) *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType), args[1].(*model.DeviceConfigurationKeyValueDescriptionDataElementsType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call) Return(_a0 *model.MsgCounterType, _a1 error) *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call) RunAndReturn(run func(*model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType, *model.DeviceConfigurationKeyValueDescriptionDataElementsType) (*model.MsgCounterType, error)) *DeviceConfigurationClientInterface_RequestKeyValueDescriptions_Call {
	_c.Call.Return(run)
	return _c
}

// RequestKeyValues provides a mock function with given fields: selector, elements
func (_m *DeviceConfigurationClientInterface) RequestKeyValues(selector *model.DeviceConfigurationKeyValueListDataSelectorsType, elements *model.DeviceConfigurationKeyValueDataElementsType) (*model.MsgCounterType, error) {
	ret := _m.Called(selector, elements)

	if len(ret) == 0 {
		panic("no return value specified for RequestKeyValues")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.DeviceConfigurationKeyValueListDataSelectorsType, *model.DeviceConfigurationKeyValueDataElementsType) (*model.MsgCounterType, error)); ok {
		return rf(selector, elements)
	}
	if rf, ok := ret.Get(0).(func(*model.DeviceConfigurationKeyValueListDataSelectorsType, *model.DeviceConfigurationKeyValueDataElementsType) *model.MsgCounterType); ok {
		r0 = rf(selector, elements)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.DeviceConfigurationKeyValueListDataSelectorsType, *model.DeviceConfigurationKeyValueDataElementsType) error); ok {
		r1 = rf(selector, elements)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_RequestKeyValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestKeyValues'
type DeviceConfigurationClientInterface_RequestKeyValues_Call struct {
	*mock.Call
}

// RequestKeyValues is a helper method to define mock.On call
//   - selector *model.DeviceConfigurationKeyValueListDataSelectorsType
//   - elements *model.DeviceConfigurationKeyValueDataElementsType
func (_e *DeviceConfigurationClientInterface_Expecter) RequestKeyValues(selector interface{}, elements interface{}) *DeviceConfigurationClientInterface_RequestKeyValues_Call {
	return &DeviceConfigurationClientInterface_RequestKeyValues_Call{Call: _e.mock.On("RequestKeyValues", selector, elements)}
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValues_Call) Run(run func(selector *model.DeviceConfigurationKeyValueListDataSelectorsType, elements *model.DeviceConfigurationKeyValueDataElementsType)) *DeviceConfigurationClientInterface_RequestKeyValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.DeviceConfigurationKeyValueListDataSelectorsType), args[1].(*model.DeviceConfigurationKeyValueDataElementsType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValues_Call) Return(_a0 *model.MsgCounterType, _a1 error) *DeviceConfigurationClientInterface_RequestKeyValues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_RequestKeyValues_Call) RunAndReturn(run func(*model.DeviceConfigurationKeyValueListDataSelectorsType, *model.DeviceConfigurationKeyValueDataElementsType) (*model.MsgCounterType, error)) *DeviceConfigurationClientInterface_RequestKeyValues_Call {
	_c.Call.Return(run)
	return _c
}

// WriteKeyValues provides a mock function with given fields: data
func (_m *DeviceConfigurationClientInterface) WriteKeyValues(data []model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for WriteKeyValues")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func([]model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func([]model.DeviceConfigurationKeyValueDataType) *model.MsgCounterType); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func([]model.DeviceConfigurationKeyValueDataType) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeviceConfigurationClientInterface_WriteKeyValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteKeyValues'
type DeviceConfigurationClientInterface_WriteKeyValues_Call struct {
	*mock.Call
}

// WriteKeyValues is a helper method to define mock.On call
//   - data []model.DeviceConfigurationKeyValueDataType
func (_e *DeviceConfigurationClientInterface_Expecter) WriteKeyValues(data interface{}) *DeviceConfigurationClientInterface_WriteKeyValues_Call {
	return &DeviceConfigurationClientInterface_WriteKeyValues_Call{Call: _e.mock.On("WriteKeyValues", data)}
}

func (_c *DeviceConfigurationClientInterface_WriteKeyValues_Call) Run(run func(data []model.DeviceConfigurationKeyValueDataType)) *DeviceConfigurationClientInterface_WriteKeyValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]model.DeviceConfigurationKeyValueDataType))
	})
	return _c
}

func (_c *DeviceConfigurationClientInterface_WriteKeyValues_Call) Return(_a0 *model.MsgCounterType, _a1 error) *DeviceConfigurationClientInterface_WriteKeyValues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DeviceConfigurationClientInterface_WriteKeyValues_Call) RunAndReturn(run func([]model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error)) *DeviceConfigurationClientInterface_WriteKeyValues_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeviceConfigurationClientInterface creates a new instance of DeviceConfigurationClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeviceConfigurationClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeviceConfigurationClientInterface {
	mock := &DeviceConfigurationClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
