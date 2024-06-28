// Code generated by mockery v2.42.1. DO NOT EDIT.

package debug

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

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

// MarshalFromCesRegistryToString provides a mock function with given fields: ctx
func (_m *mockDoguLogLevelRegistry) MarshalFromCesRegistryToString(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for MarshalFromCesRegistryToString")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
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
//   - ctx context.Context
func (_e *mockDoguLogLevelRegistry_Expecter) MarshalFromCesRegistryToString(ctx interface{}) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	return &mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call{Call: _e.mock.On("MarshalFromCesRegistryToString", ctx)}
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) Run(run func(ctx context.Context)) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) Return(_a0 string, _a1 error) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call) RunAndReturn(run func(context.Context) (string, error)) *mockDoguLogLevelRegistry_MarshalFromCesRegistryToString_Call {
	_c.Call.Return(run)
	return _c
}

// UnMarshalFromStringToCesRegistry provides a mock function with given fields: unmarshal
func (_m *mockDoguLogLevelRegistry) UnMarshalFromStringToCesRegistry(unmarshal string) error {
	ret := _m.Called(unmarshal)

	if len(ret) == 0 {
		panic("no return value specified for UnMarshalFromStringToCesRegistry")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(unmarshal)
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
//   - unmarshal string
func (_e *mockDoguLogLevelRegistry_Expecter) UnMarshalFromStringToCesRegistry(unmarshal interface{}) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	return &mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call{Call: _e.mock.On("UnMarshalFromStringToCesRegistry", unmarshal)}
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) Run(run func(unmarshal string)) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) Return(_a0 error) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call) RunAndReturn(run func(string) error) *mockDoguLogLevelRegistry_UnMarshalFromStringToCesRegistry_Call {
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
