// Code generated by mockery v2.20.0. DO NOT EDIT.

package debug

import (
	context "context"

	generateddebug "github.com/cloudogu/k8s-ces-control/generated/debug"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cloudogu/k8s-ces-control/generated/types"
)

// mockDebugModeServer is an autogenerated mock type for the debugModeServer type
type mockDebugModeServer struct {
	mock.Mock
}

type mockDebugModeServer_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDebugModeServer) EXPECT() *mockDebugModeServer_Expecter {
	return &mockDebugModeServer_Expecter{mock: &_m.Mock}
}

// Disable provides a mock function with given fields: _a0, _a1
func (_m *mockDebugModeServer) Disable(_a0 context.Context, _a1 *generateddebug.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *types.BasicResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *generateddebug.ToggleDebugModeRequest) (*types.BasicResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *generateddebug.ToggleDebugModeRequest) *types.BasicResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.BasicResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *generateddebug.ToggleDebugModeRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDebugModeServer_Disable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disable'
type mockDebugModeServer_Disable_Call struct {
	*mock.Call
}

// Disable is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *generateddebug.ToggleDebugModeRequest
func (_e *mockDebugModeServer_Expecter) Disable(_a0 interface{}, _a1 interface{}) *mockDebugModeServer_Disable_Call {
	return &mockDebugModeServer_Disable_Call{Call: _e.mock.On("Disable", _a0, _a1)}
}

func (_c *mockDebugModeServer_Disable_Call) Run(run func(_a0 context.Context, _a1 *generateddebug.ToggleDebugModeRequest)) *mockDebugModeServer_Disable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*generateddebug.ToggleDebugModeRequest))
	})
	return _c
}

func (_c *mockDebugModeServer_Disable_Call) Return(_a0 *types.BasicResponse, _a1 error) *mockDebugModeServer_Disable_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDebugModeServer_Disable_Call) RunAndReturn(run func(context.Context, *generateddebug.ToggleDebugModeRequest) (*types.BasicResponse, error)) *mockDebugModeServer_Disable_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDebugModeServer interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDebugModeServer creates a new instance of mockDebugModeServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDebugModeServer(t mockConstructorTestingTnewMockDebugModeServer) *mockDebugModeServer {
	mock := &mockDebugModeServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}