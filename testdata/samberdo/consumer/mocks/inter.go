// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Inter is an autogenerated mock type for the Inter type
type Inter struct {
	mock.Mock
}

type Inter_Expecter struct {
	mock *mock.Mock
}

func (_m *Inter) EXPECT() *Inter_Expecter {
	return &Inter_Expecter{mock: &_m.Mock}
}

// Doit provides a mock function with given fields: in
func (_m *Inter) Doit(in string) string {
	ret := _m.Called(in)

	if len(ret) == 0 {
		panic("no return value specified for Doit")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Inter_Doit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Doit'
type Inter_Doit_Call struct {
	*mock.Call
}

// Doit is a helper method to define mock.On call
//   - in string
func (_e *Inter_Expecter) Doit(in interface{}) *Inter_Doit_Call {
	return &Inter_Doit_Call{Call: _e.mock.On("Doit", in)}
}

func (_c *Inter_Doit_Call) Run(run func(in string)) *Inter_Doit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Inter_Doit_Call) Return(_a0 string) *Inter_Doit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Inter_Doit_Call) RunAndReturn(run func(string) string) *Inter_Doit_Call {
	_c.Call.Return(run)
	return _c
}

// NewInter creates a new instance of Inter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Inter {
	mock := &Inter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
