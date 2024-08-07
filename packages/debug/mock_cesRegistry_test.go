// Code generated by mockery v2.20.0. DO NOT EDIT.

package debug

import (
	registry "github.com/cloudogu/cesapp-lib/registry"
	mock "github.com/stretchr/testify/mock"
)

// mockCesRegistry is an autogenerated mock type for the cesRegistry type
type mockCesRegistry struct {
	mock.Mock
}

type mockCesRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCesRegistry) EXPECT() *mockCesRegistry_Expecter {
	return &mockCesRegistry_Expecter{mock: &_m.Mock}
}

// BlueprintRegistry provides a mock function with given fields:
func (_m *mockCesRegistry) BlueprintRegistry() registry.ConfigurationContext {
	ret := _m.Called()

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.ConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockCesRegistry_BlueprintRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlueprintRegistry'
type mockCesRegistry_BlueprintRegistry_Call struct {
	*mock.Call
}

// BlueprintRegistry is a helper method to define mock.On call
func (_e *mockCesRegistry_Expecter) BlueprintRegistry() *mockCesRegistry_BlueprintRegistry_Call {
	return &mockCesRegistry_BlueprintRegistry_Call{Call: _e.mock.On("BlueprintRegistry")}
}

func (_c *mockCesRegistry_BlueprintRegistry_Call) Run(run func()) *mockCesRegistry_BlueprintRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCesRegistry_BlueprintRegistry_Call) Return(_a0 registry.ConfigurationContext) *mockCesRegistry_BlueprintRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_BlueprintRegistry_Call) RunAndReturn(run func() registry.ConfigurationContext) *mockCesRegistry_BlueprintRegistry_Call {
	_c.Call.Return(run)
	return _c
}

// DoguConfig provides a mock function with given fields: dogu
func (_m *mockCesRegistry) DoguConfig(dogu string) registry.ConfigurationContext {
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

// mockCesRegistry_DoguConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoguConfig'
type mockCesRegistry_DoguConfig_Call struct {
	*mock.Call
}

// DoguConfig is a helper method to define mock.On call
//   - dogu string
func (_e *mockCesRegistry_Expecter) DoguConfig(dogu interface{}) *mockCesRegistry_DoguConfig_Call {
	return &mockCesRegistry_DoguConfig_Call{Call: _e.mock.On("DoguConfig", dogu)}
}

func (_c *mockCesRegistry_DoguConfig_Call) Run(run func(dogu string)) *mockCesRegistry_DoguConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCesRegistry_DoguConfig_Call) Return(_a0 registry.ConfigurationContext) *mockCesRegistry_DoguConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_DoguConfig_Call) RunAndReturn(run func(string) registry.ConfigurationContext) *mockCesRegistry_DoguConfig_Call {
	_c.Call.Return(run)
	return _c
}

// DoguRegistry provides a mock function with given fields:
func (_m *mockCesRegistry) DoguRegistry() registry.DoguRegistry {
	ret := _m.Called()

	var r0 registry.DoguRegistry
	if rf, ok := ret.Get(0).(func() registry.DoguRegistry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.DoguRegistry)
		}
	}

	return r0
}

// mockCesRegistry_DoguRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoguRegistry'
type mockCesRegistry_DoguRegistry_Call struct {
	*mock.Call
}

// DoguRegistry is a helper method to define mock.On call
func (_e *mockCesRegistry_Expecter) DoguRegistry() *mockCesRegistry_DoguRegistry_Call {
	return &mockCesRegistry_DoguRegistry_Call{Call: _e.mock.On("DoguRegistry")}
}

func (_c *mockCesRegistry_DoguRegistry_Call) Run(run func()) *mockCesRegistry_DoguRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCesRegistry_DoguRegistry_Call) Return(_a0 registry.DoguRegistry) *mockCesRegistry_DoguRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_DoguRegistry_Call) RunAndReturn(run func() registry.DoguRegistry) *mockCesRegistry_DoguRegistry_Call {
	_c.Call.Return(run)
	return _c
}

// GetNode provides a mock function with given fields:
func (_m *mockCesRegistry) GetNode() (registry.Node, error) {
	ret := _m.Called()

	var r0 registry.Node
	var r1 error
	if rf, ok := ret.Get(0).(func() (registry.Node, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() registry.Node); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(registry.Node)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCesRegistry_GetNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNode'
type mockCesRegistry_GetNode_Call struct {
	*mock.Call
}

// GetNode is a helper method to define mock.On call
func (_e *mockCesRegistry_Expecter) GetNode() *mockCesRegistry_GetNode_Call {
	return &mockCesRegistry_GetNode_Call{Call: _e.mock.On("GetNode")}
}

func (_c *mockCesRegistry_GetNode_Call) Run(run func()) *mockCesRegistry_GetNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCesRegistry_GetNode_Call) Return(_a0 registry.Node, _a1 error) *mockCesRegistry_GetNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCesRegistry_GetNode_Call) RunAndReturn(run func() (registry.Node, error)) *mockCesRegistry_GetNode_Call {
	_c.Call.Return(run)
	return _c
}

// GlobalConfig provides a mock function with given fields:
func (_m *mockCesRegistry) GlobalConfig() registry.ConfigurationContext {
	ret := _m.Called()

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.ConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockCesRegistry_GlobalConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GlobalConfig'
type mockCesRegistry_GlobalConfig_Call struct {
	*mock.Call
}

// GlobalConfig is a helper method to define mock.On call
func (_e *mockCesRegistry_Expecter) GlobalConfig() *mockCesRegistry_GlobalConfig_Call {
	return &mockCesRegistry_GlobalConfig_Call{Call: _e.mock.On("GlobalConfig")}
}

func (_c *mockCesRegistry_GlobalConfig_Call) Run(run func()) *mockCesRegistry_GlobalConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCesRegistry_GlobalConfig_Call) Return(_a0 registry.ConfigurationContext) *mockCesRegistry_GlobalConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_GlobalConfig_Call) RunAndReturn(run func() registry.ConfigurationContext) *mockCesRegistry_GlobalConfig_Call {
	_c.Call.Return(run)
	return _c
}

// HostConfig provides a mock function with given fields: hostService
func (_m *mockCesRegistry) HostConfig(hostService string) registry.ConfigurationContext {
	ret := _m.Called(hostService)

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func(string) registry.ConfigurationContext); ok {
		r0 = rf(hostService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockCesRegistry_HostConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HostConfig'
type mockCesRegistry_HostConfig_Call struct {
	*mock.Call
}

// HostConfig is a helper method to define mock.On call
//   - hostService string
func (_e *mockCesRegistry_Expecter) HostConfig(hostService interface{}) *mockCesRegistry_HostConfig_Call {
	return &mockCesRegistry_HostConfig_Call{Call: _e.mock.On("HostConfig", hostService)}
}

func (_c *mockCesRegistry_HostConfig_Call) Run(run func(hostService string)) *mockCesRegistry_HostConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCesRegistry_HostConfig_Call) Return(_a0 registry.ConfigurationContext) *mockCesRegistry_HostConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_HostConfig_Call) RunAndReturn(run func(string) registry.ConfigurationContext) *mockCesRegistry_HostConfig_Call {
	_c.Call.Return(run)
	return _c
}

// RootConfig provides a mock function with given fields:
func (_m *mockCesRegistry) RootConfig() registry.WatchConfigurationContext {
	ret := _m.Called()

	var r0 registry.WatchConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.WatchConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.WatchConfigurationContext)
		}
	}

	return r0
}

// mockCesRegistry_RootConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RootConfig'
type mockCesRegistry_RootConfig_Call struct {
	*mock.Call
}

// RootConfig is a helper method to define mock.On call
func (_e *mockCesRegistry_Expecter) RootConfig() *mockCesRegistry_RootConfig_Call {
	return &mockCesRegistry_RootConfig_Call{Call: _e.mock.On("RootConfig")}
}

func (_c *mockCesRegistry_RootConfig_Call) Run(run func()) *mockCesRegistry_RootConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCesRegistry_RootConfig_Call) Return(_a0 registry.WatchConfigurationContext) *mockCesRegistry_RootConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_RootConfig_Call) RunAndReturn(run func() registry.WatchConfigurationContext) *mockCesRegistry_RootConfig_Call {
	_c.Call.Return(run)
	return _c
}

// State provides a mock function with given fields: dogu
func (_m *mockCesRegistry) State(dogu string) registry.State {
	ret := _m.Called(dogu)

	var r0 registry.State
	if rf, ok := ret.Get(0).(func(string) registry.State); ok {
		r0 = rf(dogu)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.State)
		}
	}

	return r0
}

// mockCesRegistry_State_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'State'
type mockCesRegistry_State_Call struct {
	*mock.Call
}

// State is a helper method to define mock.On call
//   - dogu string
func (_e *mockCesRegistry_Expecter) State(dogu interface{}) *mockCesRegistry_State_Call {
	return &mockCesRegistry_State_Call{Call: _e.mock.On("State", dogu)}
}

func (_c *mockCesRegistry_State_Call) Run(run func(dogu string)) *mockCesRegistry_State_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCesRegistry_State_Call) Return(_a0 registry.State) *mockCesRegistry_State_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCesRegistry_State_Call) RunAndReturn(run func(string) registry.State) *mockCesRegistry_State_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockCesRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// newMockCesRegistry creates a new instance of mockCesRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockCesRegistry(t mockConstructorTestingTnewMockCesRegistry) *mockCesRegistry {
	mock := &mockCesRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
