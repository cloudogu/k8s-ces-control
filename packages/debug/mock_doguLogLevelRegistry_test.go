// Code generated by mockery v2.42.1. DO NOT EDIT.

package debug

import mock "github.com/stretchr/testify/mock"

// mockDoguLogLevelRegistry is an autogenerated mock type for the doguLogLevelRegistry type
type mockDoguLogLevelRegistry struct {
	mock.Mock
}

type mockDoguLogLevelRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguLogLevelRegistry) EXPECT() *mockDoguLogLevelRegistry_Expecter {
	return &mockDoguLogLevelRegistry_Expecter{mock: &_m.Mock}
}

// MarshalFromCesRegistryToString provides a mock function with given fields: registry
func (_m *mockDoguLogLevelRegistry) MarshalFromCesRegistryToString(registry cesRegistry) (string, error) {
	ret := _m.Called(registry)

	if len(ret) == 0 {
		panic("no return value specified for MarshalFromCesRegistryToString")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(cesRegistry) (string, error)); ok {
		return rf(registry)
	}
	if rf, ok := ret.Get(0).(func(cesRegistry) string); ok {
		r0 = rf(registry)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(cesRegistry) error); ok {
		r1 = rf(registry)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalFromCesRegistryToString'
type mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call struct {
	*mock.Call
}

// MarshalFromCesRegistryToString is a helper method to define mock.On call
//   - registry cesRegistry
func (_e *mockDoguLogLevelRegistry_Expecter) MarshalFromCesRegistryToString(registry interface{}) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	return &mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call{Call: _e.mock.On("MarshalFromCesRegistryToString", registry)}
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) Run(run func(registry cesRegistry)) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(cesRegistry))
	})
	return _c
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) Return(_a0 string, _a1 error) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) RunAndReturn(run func(cesRegistry) (string, error)) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Return(run)
	return _c
}

// UnMarshalFromStringToCesRegistry provides a mock function with given fields: registry, unmarshal
func (_m *mockDoguLogLevelRegistry) UnMarshalFromStringToCesRegistry(registry cesRegistry, unmarshal string) error {
	ret := _m.Called(registry, unmarshal)

	if len(ret) == 0 {
		panic("no return value specified for UnMarshalFromStringToCesRegistry")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cesRegistry, string) error); ok {
		r0 = rf(registry, unmarshal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnMarshalFromStringToCesRegistry'
type mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call struct {
	*mock.Call
}

// UnMarshalFromStringToCesRegistry is a helper method to define mock.On call
//   - registry cesRegistry
//   - unmarshal string
func (_e *mockDoguLogLevelRegistry_Expecter) UnMarshalFromStringToCesRegistry(registry interface{}, unmarshal interface{}) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	return &mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call{Call: _e.mock.On("UnMarshalFromStringToCesRegistry", registry, unmarshal)}
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) Run(run func(registry cesRegistry, unmarshal string)) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(cesRegistry), args[1].(string))
	})
	return _c
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) Return(_a0 error) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) RunAndReturn(run func(cesRegistry, string) error) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguLogLevelRegistry creates a new instance of mockDoguLogLevelRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguLogLevelRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguLogLevelRegistry {
	mock := &mockDoguLogLevelRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
