// Code generated by mockery v2.20.0. DO NOT EDIT.

package logging

import (
	registry "github.com/cloudogu/cesapp-lib/registry"
	mock "github.com/stretchr/testify/mock"
)

// mockConfigProvider is an autogenerated mock type for the configProvider type
type mockConfigProvider struct {
	mock.Mock
}

type mockConfigProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *mockConfigProvider) EXPECT() *mockConfigProvider_Expecter {
	return &mockConfigProvider_Expecter{mock: &_m.Mock}
}

// DoguConfig provides a mock function with given fields: dogu
func (_m *mockConfigProvider) DoguConfig(dogu string) registry.ConfigurationContext {
	ret := _m.Called(dogu)

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func(string) registry.ConfigurationContext); ok {
		r0 = rf(dogu)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockConfigProvider_DoguConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoguConfig'
type mockConfigProvider_DoguConfig_Call struct {
	*mock.Call
}

// DoguConfig is a helper method to define mock.On call
//   - dogu string
func (_e *mockConfigProvider_Expecter) DoguConfig(dogu interface{}) *mockConfigProvider_DoguConfig_Call {
	return &mockConfigProvider_DoguConfig_Call{Call: _e.mock.On("DoguConfig", dogu)}
}

func (_c *mockConfigProvider_DoguConfig_Call) Run(run func(dogu string)) *mockConfigProvider_DoguConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockConfigProvider_DoguConfig_Call) Return(_a0 registry.ConfigurationContext) *mockConfigProvider_DoguConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockConfigProvider_DoguConfig_Call) RunAndReturn(run func(string) registry.ConfigurationContext) *mockConfigProvider_DoguConfig_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockConfigProvider interface {
	mock.TestingT
	Cleanup(func())
}

// newMockConfigProvider creates a new instance of mockConfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockConfigProvider(t mockConstructorTestingTnewMockConfigProvider) *mockConfigProvider {
	mock := &mockConfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
