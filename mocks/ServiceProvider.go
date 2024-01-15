// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	api "github.com/enbility/eebus-go/api"
	mock "github.com/stretchr/testify/mock"
)

// ServiceProvider is an autogenerated mock type for the ServiceProvider type
type ServiceProvider struct {
	mock.Mock
}

type ServiceProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *ServiceProvider) EXPECT() *ServiceProvider_Expecter {
	return &ServiceProvider_Expecter{mock: &_m.Mock}
}

// AllowWaitingForTrust provides a mock function with given fields: ski
func (_m *ServiceProvider) AllowWaitingForTrust(ski string) bool {
	ret := _m.Called(ski)

	if len(ret) == 0 {
		panic("no return value specified for AllowWaitingForTrust")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(ski)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ServiceProvider_AllowWaitingForTrust_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AllowWaitingForTrust'
type ServiceProvider_AllowWaitingForTrust_Call struct {
	*mock.Call
}

// AllowWaitingForTrust is a helper method to define mock.On call
//   - ski string
func (_e *ServiceProvider_Expecter) AllowWaitingForTrust(ski interface{}) *ServiceProvider_AllowWaitingForTrust_Call {
	return &ServiceProvider_AllowWaitingForTrust_Call{Call: _e.mock.On("AllowWaitingForTrust", ski)}
}

func (_c *ServiceProvider_AllowWaitingForTrust_Call) Run(run func(ski string)) *ServiceProvider_AllowWaitingForTrust_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ServiceProvider_AllowWaitingForTrust_Call) Return(_a0 bool) *ServiceProvider_AllowWaitingForTrust_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ServiceProvider_AllowWaitingForTrust_Call) RunAndReturn(run func(string) bool) *ServiceProvider_AllowWaitingForTrust_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteSKIConnected provides a mock function with given fields: ski
func (_m *ServiceProvider) RemoteSKIConnected(ski string) {
	_m.Called(ski)
}

// ServiceProvider_RemoteSKIConnected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteSKIConnected'
type ServiceProvider_RemoteSKIConnected_Call struct {
	*mock.Call
}

// RemoteSKIConnected is a helper method to define mock.On call
//   - ski string
func (_e *ServiceProvider_Expecter) RemoteSKIConnected(ski interface{}) *ServiceProvider_RemoteSKIConnected_Call {
	return &ServiceProvider_RemoteSKIConnected_Call{Call: _e.mock.On("RemoteSKIConnected", ski)}
}

func (_c *ServiceProvider_RemoteSKIConnected_Call) Run(run func(ski string)) *ServiceProvider_RemoteSKIConnected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ServiceProvider_RemoteSKIConnected_Call) Return() *ServiceProvider_RemoteSKIConnected_Call {
	_c.Call.Return()
	return _c
}

func (_c *ServiceProvider_RemoteSKIConnected_Call) RunAndReturn(run func(string)) *ServiceProvider_RemoteSKIConnected_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteSKIDisconnected provides a mock function with given fields: ski
func (_m *ServiceProvider) RemoteSKIDisconnected(ski string) {
	_m.Called(ski)
}

// ServiceProvider_RemoteSKIDisconnected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteSKIDisconnected'
type ServiceProvider_RemoteSKIDisconnected_Call struct {
	*mock.Call
}

// RemoteSKIDisconnected is a helper method to define mock.On call
//   - ski string
func (_e *ServiceProvider_Expecter) RemoteSKIDisconnected(ski interface{}) *ServiceProvider_RemoteSKIDisconnected_Call {
	return &ServiceProvider_RemoteSKIDisconnected_Call{Call: _e.mock.On("RemoteSKIDisconnected", ski)}
}

func (_c *ServiceProvider_RemoteSKIDisconnected_Call) Run(run func(ski string)) *ServiceProvider_RemoteSKIDisconnected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ServiceProvider_RemoteSKIDisconnected_Call) Return() *ServiceProvider_RemoteSKIDisconnected_Call {
	_c.Call.Return()
	return _c
}

func (_c *ServiceProvider_RemoteSKIDisconnected_Call) RunAndReturn(run func(string)) *ServiceProvider_RemoteSKIDisconnected_Call {
	_c.Call.Return(run)
	return _c
}

// ServicePairingDetailUpdate provides a mock function with given fields: ski, detail
func (_m *ServiceProvider) ServicePairingDetailUpdate(ski string, detail *api.ConnectionStateDetail) {
	_m.Called(ski, detail)
}

// ServiceProvider_ServicePairingDetailUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServicePairingDetailUpdate'
type ServiceProvider_ServicePairingDetailUpdate_Call struct {
	*mock.Call
}

// ServicePairingDetailUpdate is a helper method to define mock.On call
//   - ski string
//   - detail *api.ConnectionStateDetail
func (_e *ServiceProvider_Expecter) ServicePairingDetailUpdate(ski interface{}, detail interface{}) *ServiceProvider_ServicePairingDetailUpdate_Call {
	return &ServiceProvider_ServicePairingDetailUpdate_Call{Call: _e.mock.On("ServicePairingDetailUpdate", ski, detail)}
}

func (_c *ServiceProvider_ServicePairingDetailUpdate_Call) Run(run func(ski string, detail *api.ConnectionStateDetail)) *ServiceProvider_ServicePairingDetailUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*api.ConnectionStateDetail))
	})
	return _c
}

func (_c *ServiceProvider_ServicePairingDetailUpdate_Call) Return() *ServiceProvider_ServicePairingDetailUpdate_Call {
	_c.Call.Return()
	return _c
}

func (_c *ServiceProvider_ServicePairingDetailUpdate_Call) RunAndReturn(run func(string, *api.ConnectionStateDetail)) *ServiceProvider_ServicePairingDetailUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// ServiceShipIDUpdate provides a mock function with given fields: ski, shipID
func (_m *ServiceProvider) ServiceShipIDUpdate(ski string, shipID string) {
	_m.Called(ski, shipID)
}

// ServiceProvider_ServiceShipIDUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServiceShipIDUpdate'
type ServiceProvider_ServiceShipIDUpdate_Call struct {
	*mock.Call
}

// ServiceShipIDUpdate is a helper method to define mock.On call
//   - ski string
//   - shipID string
func (_e *ServiceProvider_Expecter) ServiceShipIDUpdate(ski interface{}, shipID interface{}) *ServiceProvider_ServiceShipIDUpdate_Call {
	return &ServiceProvider_ServiceShipIDUpdate_Call{Call: _e.mock.On("ServiceShipIDUpdate", ski, shipID)}
}

func (_c *ServiceProvider_ServiceShipIDUpdate_Call) Run(run func(ski string, shipID string)) *ServiceProvider_ServiceShipIDUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *ServiceProvider_ServiceShipIDUpdate_Call) Return() *ServiceProvider_ServiceShipIDUpdate_Call {
	_c.Call.Return()
	return _c
}

func (_c *ServiceProvider_ServiceShipIDUpdate_Call) RunAndReturn(run func(string, string)) *ServiceProvider_ServiceShipIDUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// VisibleMDNSRecordsUpdated provides a mock function with given fields: entries
func (_m *ServiceProvider) VisibleMDNSRecordsUpdated(entries []*api.MdnsEntry) {
	_m.Called(entries)
}

// ServiceProvider_VisibleMDNSRecordsUpdated_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VisibleMDNSRecordsUpdated'
type ServiceProvider_VisibleMDNSRecordsUpdated_Call struct {
	*mock.Call
}

// VisibleMDNSRecordsUpdated is a helper method to define mock.On call
//   - entries []*api.MdnsEntry
func (_e *ServiceProvider_Expecter) VisibleMDNSRecordsUpdated(entries interface{}) *ServiceProvider_VisibleMDNSRecordsUpdated_Call {
	return &ServiceProvider_VisibleMDNSRecordsUpdated_Call{Call: _e.mock.On("VisibleMDNSRecordsUpdated", entries)}
}

func (_c *ServiceProvider_VisibleMDNSRecordsUpdated_Call) Run(run func(entries []*api.MdnsEntry)) *ServiceProvider_VisibleMDNSRecordsUpdated_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*api.MdnsEntry))
	})
	return _c
}

func (_c *ServiceProvider_VisibleMDNSRecordsUpdated_Call) Return() *ServiceProvider_VisibleMDNSRecordsUpdated_Call {
	_c.Call.Return()
	return _c
}

func (_c *ServiceProvider_VisibleMDNSRecordsUpdated_Call) RunAndReturn(run func([]*api.MdnsEntry)) *ServiceProvider_VisibleMDNSRecordsUpdated_Call {
	_c.Call.Return(run)
	return _c
}

// NewServiceProvider creates a new instance of ServiceProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceProvider {
	mock := &ServiceProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}