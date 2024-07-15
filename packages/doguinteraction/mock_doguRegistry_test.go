// Code generated by mockery v2.20.0. DO NOT EDIT.

package doguinteraction

import (
	context "context"

	core "github.com/cloudogu/cesapp-lib/core"
	mock "github.com/stretchr/testify/mock"
)

// mockDoguRegistry is an autogenerated mock type for the doguRegistry type
type mockDoguRegistry struct {
	mock.Mock
}

type mockDoguRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguRegistry) EXPECT() *mockDoguRegistry_Expecter {
	return &mockDoguRegistry_Expecter{mock: &_m.Mock}
}

// GetCurrentOfAll provides a mock function with given fields: ctx
func (_m *mockDoguRegistry) GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error) {
	ret := _m.Called(ctx)

	var r0 []*core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*core.Dogu, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*core.Dogu); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguRegistry_GetCurrentOfAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentOfAll'
type mockDoguRegistry_GetCurrentOfAll_Call struct {
	*mock.Call
}

// GetCurrentOfAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockDoguRegistry_Expecter) GetCurrentOfAll(ctx interface{}) *mockDoguRegistry_GetCurrentOfAll_Call {
	return &mockDoguRegistry_GetCurrentOfAll_Call{Call: _e.mock.On("GetCurrentOfAll", ctx)}
}

func (_c *mockDoguRegistry_GetCurrentOfAll_Call) Run(run func(ctx context.Context)) *mockDoguRegistry_GetCurrentOfAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockDoguRegistry_GetCurrentOfAll_Call) Return(_a0 []*core.Dogu, _a1 error) *mockDoguRegistry_GetCurrentOfAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguRegistry_GetCurrentOfAll_Call) RunAndReturn(run func(context.Context) ([]*core.Dogu, error)) *mockDoguRegistry_GetCurrentOfAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguRegistry creates a new instance of mockDoguRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguRegistry(t mockConstructorTestingTnewMockDoguRegistry) *mockDoguRegistry {
	mock := &mockDoguRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
