// Code generated by mockery v2.42.0. DO NOT EDIT.

package servicemocks

import (
	mock "github.com/stretchr/testify/mock"
	health "github.com/willie68/go-micro/internal/services/health"

	time "time"
)

// SHealth is an autogenerated mock type for the SHealth type
type SHealth struct {
	mock.Mock
}

type SHealth_Expecter struct {
	mock *mock.Mock
}

func (_m *SHealth) EXPECT() *SHealth_Expecter {
	return &SHealth_Expecter{mock: &_m.Mock}
}

// CheckHealthCheckTimer provides a mock function with given fields:
func (_m *SHealth) CheckHealthCheckTimer() {
	_m.Called()
}

// SHealth_CheckHealthCheckTimer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckHealthCheckTimer'
type SHealth_CheckHealthCheckTimer_Call struct {
	*mock.Call
}

// CheckHealthCheckTimer is a helper method to define mock.On call
func (_e *SHealth_Expecter) CheckHealthCheckTimer() *SHealth_CheckHealthCheckTimer_Call {
	return &SHealth_CheckHealthCheckTimer_Call{Call: _e.mock.On("CheckHealthCheckTimer")}
}

func (_c *SHealth_CheckHealthCheckTimer_Call) Run(run func()) *SHealth_CheckHealthCheckTimer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SHealth_CheckHealthCheckTimer_Call) Return() *SHealth_CheckHealthCheckTimer_Call {
	_c.Call.Return()
	return _c
}

func (_c *SHealth_CheckHealthCheckTimer_Call) RunAndReturn(run func()) *SHealth_CheckHealthCheckTimer_Call {
	_c.Call.Return(run)
	return _c
}

// Healthy provides a mock function with given fields:
func (_m *SHealth) Healthy() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Healthy")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SHealth_Healthy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Healthy'
type SHealth_Healthy_Call struct {
	*mock.Call
}

// Healthy is a helper method to define mock.On call
func (_e *SHealth_Expecter) Healthy() *SHealth_Healthy_Call {
	return &SHealth_Healthy_Call{Call: _e.mock.On("Healthy")}
}

func (_c *SHealth_Healthy_Call) Run(run func()) *SHealth_Healthy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SHealth_Healthy_Call) Return(_a0 bool) *SHealth_Healthy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SHealth_Healthy_Call) RunAndReturn(run func() bool) *SHealth_Healthy_Call {
	_c.Call.Return(run)
	return _c
}

// LastChecked provides a mock function with given fields:
func (_m *SHealth) LastChecked() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LastChecked")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// SHealth_LastChecked_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LastChecked'
type SHealth_LastChecked_Call struct {
	*mock.Call
}

// LastChecked is a helper method to define mock.On call
func (_e *SHealth_Expecter) LastChecked() *SHealth_LastChecked_Call {
	return &SHealth_LastChecked_Call{Call: _e.mock.On("LastChecked")}
}

func (_c *SHealth_LastChecked_Call) Run(run func()) *SHealth_LastChecked_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SHealth_LastChecked_Call) Return(_a0 time.Time) *SHealth_LastChecked_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SHealth_LastChecked_Call) RunAndReturn(run func() time.Time) *SHealth_LastChecked_Call {
	_c.Call.Return(run)
	return _c
}

// Message provides a mock function with given fields:
func (_m *SHealth) Message() health.Message {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Message")
	}

	var r0 health.Message
	if rf, ok := ret.Get(0).(func() health.Message); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(health.Message)
	}

	return r0
}

// SHealth_Message_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Message'
type SHealth_Message_Call struct {
	*mock.Call
}

// Message is a helper method to define mock.On call
func (_e *SHealth_Expecter) Message() *SHealth_Message_Call {
	return &SHealth_Message_Call{Call: _e.mock.On("Message")}
}

func (_c *SHealth_Message_Call) Run(run func()) *SHealth_Message_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SHealth_Message_Call) Return(_a0 health.Message) *SHealth_Message_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SHealth_Message_Call) RunAndReturn(run func() health.Message) *SHealth_Message_Call {
	_c.Call.Return(run)
	return _c
}

// Readyz provides a mock function with given fields:
func (_m *SHealth) Readyz() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Readyz")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SHealth_Readyz_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Readyz'
type SHealth_Readyz_Call struct {
	*mock.Call
}

// Readyz is a helper method to define mock.On call
func (_e *SHealth_Expecter) Readyz() *SHealth_Readyz_Call {
	return &SHealth_Readyz_Call{Call: _e.mock.On("Readyz")}
}

func (_c *SHealth_Readyz_Call) Run(run func()) *SHealth_Readyz_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SHealth_Readyz_Call) Return(_a0 bool) *SHealth_Readyz_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SHealth_Readyz_Call) RunAndReturn(run func() bool) *SHealth_Readyz_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: check
func (_m *SHealth) Register(check health.Check) {
	_m.Called(check)
}

// SHealth_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type SHealth_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - check health.Check
func (_e *SHealth_Expecter) Register(check interface{}) *SHealth_Register_Call {
	return &SHealth_Register_Call{Call: _e.mock.On("Register", check)}
}

func (_c *SHealth_Register_Call) Run(run func(check health.Check)) *SHealth_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(health.Check))
	})
	return _c
}

func (_c *SHealth_Register_Call) Return() *SHealth_Register_Call {
	_c.Call.Return()
	return _c
}

func (_c *SHealth_Register_Call) RunAndReturn(run func(health.Check)) *SHealth_Register_Call {
	_c.Call.Return(run)
	return _c
}

// Unregister provides a mock function with given fields: checkname
func (_m *SHealth) Unregister(checkname string) bool {
	ret := _m.Called(checkname)

	if len(ret) == 0 {
		panic("no return value specified for Unregister")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(checkname)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SHealth_Unregister_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unregister'
type SHealth_Unregister_Call struct {
	*mock.Call
}

// Unregister is a helper method to define mock.On call
//   - checkname string
func (_e *SHealth_Expecter) Unregister(checkname interface{}) *SHealth_Unregister_Call {
	return &SHealth_Unregister_Call{Call: _e.mock.On("Unregister", checkname)}
}

func (_c *SHealth_Unregister_Call) Run(run func(checkname string)) *SHealth_Unregister_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *SHealth_Unregister_Call) Return(_a0 bool) *SHealth_Unregister_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SHealth_Unregister_Call) RunAndReturn(run func(string) bool) *SHealth_Unregister_Call {
	_c.Call.Return(run)
	return _c
}

// NewSHealth creates a new instance of SHealth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSHealth(t interface {
	mock.TestingT
	Cleanup(func())
}) *SHealth {
	mock := &SHealth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
