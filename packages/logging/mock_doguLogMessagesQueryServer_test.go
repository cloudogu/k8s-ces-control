// Code generated by mockery v2.42.1. DO NOT EDIT.

package logging

import (
	context "context"

	generatedlogging "github.com/cloudogu/ces-control-api/generated/logging"
	metadata "google.golang.org/grpc/metadata"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguLogMessagesQueryServer is an autogenerated mock type for the doguLogMessagesQueryServer type
type mockDoguLogMessagesQueryServer struct {
	mock.Mock
}

type mockDoguLogMessagesQueryServer_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguLogMessagesQueryServer) EXPECT() *mockDoguLogMessagesQueryServer_Expecter {
	return &mockDoguLogMessagesQueryServer_Expecter{mock: &_m.Mock}
}

// Context provides a mock function with given fields:
func (_m *mockDoguLogMessagesQueryServer) Context() context.Context {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Context")
	}

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// mockDoguLogMessagesQueryServer_Context_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Context'
type mockDoguLogMessagesQueryServer_Context_Call struct {
	*mock.Call
}

// Context is a helper method to define mock.On call
func (_e *mockDoguLogMessagesQueryServer_Expecter) Context() *mockDoguLogMessagesQueryServer_Context_Call {
	return &mockDoguLogMessagesQueryServer_Context_Call{Call: _e.mock.On("Context")}
}

func (_c *mockDoguLogMessagesQueryServer_Context_Call) Run(run func()) *mockDoguLogMessagesQueryServer_Context_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_Context_Call) Return(_a0 context.Context) *mockDoguLogMessagesQueryServer_Context_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_Context_Call) RunAndReturn(run func() context.Context) *mockDoguLogMessagesQueryServer_Context_Call {
	_c.Call.Return(run)
	return _c
}

// RecvMsg provides a mock function with given fields: m
func (_m *mockDoguLogMessagesQueryServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for RecvMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogMessagesQueryServer_RecvMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecvMsg'
type mockDoguLogMessagesQueryServer_RecvMsg_Call struct {
	*mock.Call
}

// RecvMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *mockDoguLogMessagesQueryServer_Expecter) RecvMsg(m interface{}) *mockDoguLogMessagesQueryServer_RecvMsg_Call {
	return &mockDoguLogMessagesQueryServer_RecvMsg_Call{Call: _e.mock.On("RecvMsg", m)}
}

func (_c *mockDoguLogMessagesQueryServer_RecvMsg_Call) Run(run func(m interface{})) *mockDoguLogMessagesQueryServer_RecvMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_RecvMsg_Call) Return(_a0 error) *mockDoguLogMessagesQueryServer_RecvMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_RecvMsg_Call) RunAndReturn(run func(interface{}) error) *mockDoguLogMessagesQueryServer_RecvMsg_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0
func (_m *mockDoguLogMessagesQueryServer) Send(_a0 *generatedlogging.DoguLogMessage) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*generatedlogging.DoguLogMessage) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogMessagesQueryServer_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type mockDoguLogMessagesQueryServer_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 *generatedlogging.DoguLogMessage
func (_e *mockDoguLogMessagesQueryServer_Expecter) Send(_a0 interface{}) *mockDoguLogMessagesQueryServer_Send_Call {
	return &mockDoguLogMessagesQueryServer_Send_Call{Call: _e.mock.On("Send", _a0)}
}

func (_c *mockDoguLogMessagesQueryServer_Send_Call) Run(run func(_a0 *generatedlogging.DoguLogMessage)) *mockDoguLogMessagesQueryServer_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*generatedlogging.DoguLogMessage))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_Send_Call) Return(_a0 error) *mockDoguLogMessagesQueryServer_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_Send_Call) RunAndReturn(run func(*generatedlogging.DoguLogMessage) error) *mockDoguLogMessagesQueryServer_Send_Call {
	_c.Call.Return(run)
	return _c
}

// SendHeader provides a mock function with given fields: _a0
func (_m *mockDoguLogMessagesQueryServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogMessagesQueryServer_SendHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHeader'
type mockDoguLogMessagesQueryServer_SendHeader_Call struct {
	*mock.Call
}

// SendHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *mockDoguLogMessagesQueryServer_Expecter) SendHeader(_a0 interface{}) *mockDoguLogMessagesQueryServer_SendHeader_Call {
	return &mockDoguLogMessagesQueryServer_SendHeader_Call{Call: _e.mock.On("SendHeader", _a0)}
}

func (_c *mockDoguLogMessagesQueryServer_SendHeader_Call) Run(run func(_a0 metadata.MD)) *mockDoguLogMessagesQueryServer_SendHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SendHeader_Call) Return(_a0 error) *mockDoguLogMessagesQueryServer_SendHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SendHeader_Call) RunAndReturn(run func(metadata.MD) error) *mockDoguLogMessagesQueryServer_SendHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SendMsg provides a mock function with given fields: m
func (_m *mockDoguLogMessagesQueryServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for SendMsg")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogMessagesQueryServer_SendMsg_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMsg'
type mockDoguLogMessagesQueryServer_SendMsg_Call struct {
	*mock.Call
}

// SendMsg is a helper method to define mock.On call
//   - m interface{}
func (_e *mockDoguLogMessagesQueryServer_Expecter) SendMsg(m interface{}) *mockDoguLogMessagesQueryServer_SendMsg_Call {
	return &mockDoguLogMessagesQueryServer_SendMsg_Call{Call: _e.mock.On("SendMsg", m)}
}

func (_c *mockDoguLogMessagesQueryServer_SendMsg_Call) Run(run func(m interface{})) *mockDoguLogMessagesQueryServer_SendMsg_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SendMsg_Call) Return(_a0 error) *mockDoguLogMessagesQueryServer_SendMsg_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SendMsg_Call) RunAndReturn(run func(interface{}) error) *mockDoguLogMessagesQueryServer_SendMsg_Call {
	_c.Call.Return(run)
	return _c
}

// SetHeader provides a mock function with given fields: _a0
func (_m *mockDoguLogMessagesQueryServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetHeader")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguLogMessagesQueryServer_SetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeader'
type mockDoguLogMessagesQueryServer_SetHeader_Call struct {
	*mock.Call
}

// SetHeader is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *mockDoguLogMessagesQueryServer_Expecter) SetHeader(_a0 interface{}) *mockDoguLogMessagesQueryServer_SetHeader_Call {
	return &mockDoguLogMessagesQueryServer_SetHeader_Call{Call: _e.mock.On("SetHeader", _a0)}
}

func (_c *mockDoguLogMessagesQueryServer_SetHeader_Call) Run(run func(_a0 metadata.MD)) *mockDoguLogMessagesQueryServer_SetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SetHeader_Call) Return(_a0 error) *mockDoguLogMessagesQueryServer_SetHeader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SetHeader_Call) RunAndReturn(run func(metadata.MD) error) *mockDoguLogMessagesQueryServer_SetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *mockDoguLogMessagesQueryServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

// mockDoguLogMessagesQueryServer_SetTrailer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrailer'
type mockDoguLogMessagesQueryServer_SetTrailer_Call struct {
	*mock.Call
}

// SetTrailer is a helper method to define mock.On call
//   - _a0 metadata.MD
func (_e *mockDoguLogMessagesQueryServer_Expecter) SetTrailer(_a0 interface{}) *mockDoguLogMessagesQueryServer_SetTrailer_Call {
	return &mockDoguLogMessagesQueryServer_SetTrailer_Call{Call: _e.mock.On("SetTrailer", _a0)}
}

func (_c *mockDoguLogMessagesQueryServer_SetTrailer_Call) Run(run func(_a0 metadata.MD)) *mockDoguLogMessagesQueryServer_SetTrailer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.MD))
	})
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SetTrailer_Call) Return() *mockDoguLogMessagesQueryServer_SetTrailer_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockDoguLogMessagesQueryServer_SetTrailer_Call) RunAndReturn(run func(metadata.MD)) *mockDoguLogMessagesQueryServer_SetTrailer_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguLogMessagesQueryServer creates a new instance of mockDoguLogMessagesQueryServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguLogMessagesQueryServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguLogMessagesQueryServer {
	mock := &mockDoguLogMessagesQueryServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
