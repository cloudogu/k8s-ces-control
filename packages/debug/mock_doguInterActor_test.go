// Code generated by mockery v2.20.0. DO NOT EDIT.

package debug

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguInterActor is an autogenerated mock type for the doguInterActor type
type mockDoguInterActor struct {
	mock.Mock
}

type mockDoguInterActor_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguInterActor) EXPECT() *mockDoguInterActor_Expecter {
	return &mockDoguInterActor_Expecter{mock: &_m.Mock}
}

// RestartDoguWithWait provides a mock function with given fields: ctx, doguName, waitForRollout
func (_m *mockDoguInterActor) RestartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	ret := _m.Called(ctx, doguName, waitForRollout)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, doguName, waitForRollout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInterActor_RestartDoguWithWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RestartDoguWithWait'
type mockDoguInterActor_RestartDoguWithWait_Call struct {
	*mock.Call
}

// RestartDoguWithWait is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName string
//   - waitForRollout bool
func (_e *mockDoguInterActor_Expecter) RestartDoguWithWait(ctx interface{}, doguName interface{}, waitForRollout interface{}) *mockDoguInterActor_RestartDoguWithWait_Call {
	return &mockDoguInterActor_RestartDoguWithWait_Call{Call: _e.mock.On("RestartDoguWithWait", ctx, doguName, waitForRollout)}
}

func (_c *mockDoguInterActor_RestartDoguWithWait_Call) Run(run func(ctx context.Context, doguName string, waitForRollout bool)) *mockDoguInterActor_RestartDoguWithWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *mockDoguInterActor_RestartDoguWithWait_Call) Return(_a0 error) *mockDoguInterActor_RestartDoguWithWait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInterActor_RestartDoguWithWait_Call) RunAndReturn(run func(context.Context, string, bool) error) *mockDoguInterActor_RestartDoguWithWait_Call {
	_c.Call.Return(run)
	return _c
}

// StartDoguWithWait provides a mock function with given fields: ctx, doguName, waitForRollout
func (_m *mockDoguInterActor) StartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	ret := _m.Called(ctx, doguName, waitForRollout)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, doguName, waitForRollout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInterActor_StartDoguWithWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartDoguWithWait'
type mockDoguInterActor_StartDoguWithWait_Call struct {
	*mock.Call
}

// StartDoguWithWait is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName string
//   - waitForRollout bool
func (_e *mockDoguInterActor_Expecter) StartDoguWithWait(ctx interface{}, doguName interface{}, waitForRollout interface{}) *mockDoguInterActor_StartDoguWithWait_Call {
	return &mockDoguInterActor_StartDoguWithWait_Call{Call: _e.mock.On("StartDoguWithWait", ctx, doguName, waitForRollout)}
}

func (_c *mockDoguInterActor_StartDoguWithWait_Call) Run(run func(ctx context.Context, doguName string, waitForRollout bool)) *mockDoguInterActor_StartDoguWithWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *mockDoguInterActor_StartDoguWithWait_Call) Return(_a0 error) *mockDoguInterActor_StartDoguWithWait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInterActor_StartDoguWithWait_Call) RunAndReturn(run func(context.Context, string, bool) error) *mockDoguInterActor_StartDoguWithWait_Call {
	_c.Call.Return(run)
	return _c
}

// StopDoguWithWait provides a mock function with given fields: ctx, doguName, waitForRollout
func (_m *mockDoguInterActor) StopDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	ret := _m.Called(ctx, doguName, waitForRollout)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, doguName, waitForRollout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInterActor_StopDoguWithWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopDoguWithWait'
type mockDoguInterActor_StopDoguWithWait_Call struct {
	*mock.Call
}

// StopDoguWithWait is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName string
//   - waitForRollout bool
func (_e *mockDoguInterActor_Expecter) StopDoguWithWait(ctx interface{}, doguName interface{}, waitForRollout interface{}) *mockDoguInterActor_StopDoguWithWait_Call {
	return &mockDoguInterActor_StopDoguWithWait_Call{Call: _e.mock.On("StopDoguWithWait", ctx, doguName, waitForRollout)}
}

func (_c *mockDoguInterActor_StopDoguWithWait_Call) Run(run func(ctx context.Context, doguName string, waitForRollout bool)) *mockDoguInterActor_StopDoguWithWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *mockDoguInterActor_StopDoguWithWait_Call) Return(_a0 error) *mockDoguInterActor_StopDoguWithWait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInterActor_StopDoguWithWait_Call) RunAndReturn(run func(context.Context, string, bool) error) *mockDoguInterActor_StopDoguWithWait_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguInterActor interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguInterActor creates a new instance of mockDoguInterActor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguInterActor(t mockConstructorTestingTnewMockDoguInterActor) *mockDoguInterActor {
	mock := &mockDoguInterActor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
