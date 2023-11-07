// Code generated by mockery v2.20.0. DO NOT EDIT.

package stream

import (
	types "github.com/cloudogu/k8s-ces-control/generated/types"
	mock "github.com/stretchr/testify/mock"
)

// MockGRPCStreamServer is an autogenerated mock type for the GRPCStreamServer type
type MockGRPCStreamServer struct {
	mock.Mock
}

type MockGRPCStreamServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGRPCStreamServer) EXPECT() *MockGRPCStreamServer_Expecter {
	return &MockGRPCStreamServer_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: response
func (_m *MockGRPCStreamServer) Send(response *types.ChunkedDataResponse) error {
	ret := _m.Called(response)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.ChunkedDataResponse) error); ok {
		r0 = rf(response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockGRPCStreamServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockGRPCStreamServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - response *types.ChunkedDataResponse
func (_e *MockGRPCStreamServer_Expecter) Send(response interface{}) *MockGRPCStreamServer_Send_Call {
	return &MockGRPCStreamServer_Send_Call{Call: _e.mock.On("Send", response)}
}

func (_c *MockGRPCStreamServer_Send_Call) Run(run func(response *types.ChunkedDataResponse)) *MockGRPCStreamServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*types.ChunkedDataResponse))
	})
	return _c
}

func (_c *MockGRPCStreamServer_Send_Call) Return(_a0 error) *MockGRPCStreamServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGRPCStreamServer_Send_Call) RunAndReturn(run func(*types.ChunkedDataResponse) error) *MockGRPCStreamServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockGRPCStreamServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGRPCStreamServer creates a new instance of MockGRPCStreamServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGRPCStreamServer(t mockConstructorTestingTNewMockGRPCStreamServer) *MockGRPCStreamServer {
	mock := &MockGRPCStreamServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
