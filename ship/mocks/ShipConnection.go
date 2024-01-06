// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	ship "github.com/enbility/eebus-go/ship"
	mock "github.com/stretchr/testify/mock"
)

// ShipConnection is an autogenerated mock type for the ShipConnection type
type ShipConnection struct {
	mock.Mock
}

// AbortPendingHandshake provides a mock function with given fields:
func (_m *ShipConnection) AbortPendingHandshake() {
	_m.Called()
}

// ApprovePendingHandshake provides a mock function with given fields:
func (_m *ShipConnection) ApprovePendingHandshake() {
	_m.Called()
}

// CloseConnection provides a mock function with given fields: safe, code, reason
func (_m *ShipConnection) CloseConnection(safe bool, code int, reason string) {
	_m.Called(safe, code, reason)
}

// DataHandler provides a mock function with given fields:
func (_m *ShipConnection) DataHandler() ship.ShipDataConnection {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DataHandler")
	}

	var r0 ship.ShipDataConnection
	if rf, ok := ret.Get(0).(func() ship.ShipDataConnection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ship.ShipDataConnection)
		}
	}

	return r0
}

// RemoteSKI provides a mock function with given fields:
func (_m *ShipConnection) RemoteSKI() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RemoteSKI")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ShipHandshakeState provides a mock function with given fields:
func (_m *ShipConnection) ShipHandshakeState() (ship.ShipMessageExchangeState, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ShipHandshakeState")
	}

	var r0 ship.ShipMessageExchangeState
	var r1 error
	if rf, ok := ret.Get(0).(func() (ship.ShipMessageExchangeState, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() ship.ShipMessageExchangeState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ship.ShipMessageExchangeState)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShipConnection creates a new instance of ShipConnection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShipConnection(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShipConnection {
	mock := &ShipConnection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}