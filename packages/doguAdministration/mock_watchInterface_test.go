// Code generated by mockery v2.20.0. DO NOT EDIT.

package doguAdministration

import (
	mock "github.com/stretchr/testify/mock"
	watch "k8s.io/apimachinery/pkg/watch"
)

// mockWatchInterface is an autogenerated mock type for the watchInterface type
type mockWatchInterface struct {
	mock.Mock
}

type mockWatchInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *mockWatchInterface) EXPECT() *mockWatchInterface_Expecter {
	return &mockWatchInterface_Expecter{mock: &_m.Mock}
}

// ResultChan provides a mock function with given fields:
func (_m *mockWatchInterface) ResultChan() <-chan watch.Event {
	ret := _m.Called()

	var r0 <-chan watch.Event
	if rf, ok := ret.Get(0).(func() <-chan watch.Event); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan watch.Event)
		}
	}

	return r0
}

// mockWatchInterface_ResultChan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResultChan'
type mockWatchInterface_ResultChan_Call struct {
	*mock.Call
}

// ResultChan is a helper method to define mock.On call
func (_e *mockWatchInterface_Expecter) ResultChan() *mockWatchInterface_ResultChan_Call {
	return &mockWatchInterface_ResultChan_Call{Call: _e.mock.On("ResultChan")}
}

func (_c *mockWatchInterface_ResultChan_Call) Run(run func()) *mockWatchInterface_ResultChan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockWatchInterface_ResultChan_Call) Return(_a0 <-chan watch.Event) *mockWatchInterface_ResultChan_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockWatchInterface_ResultChan_Call) RunAndReturn(run func() <-chan watch.Event) *mockWatchInterface_ResultChan_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *mockWatchInterface) Stop() {
	_m.Called()
}

// mockWatchInterface_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type mockWatchInterface_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *mockWatchInterface_Expecter) Stop() *mockWatchInterface_Stop_Call {
	return &mockWatchInterface_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *mockWatchInterface_Stop_Call) Run(run func()) *mockWatchInterface_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockWatchInterface_Stop_Call) Return() *mockWatchInterface_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockWatchInterface_Stop_Call) RunAndReturn(run func()) *mockWatchInterface_Stop_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockWatchInterface interface {
	mock.TestingT
	Cleanup(func())
}

// newMockWatchInterface creates a new instance of mockWatchInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockWatchInterface(t mockConstructorTestingTnewMockWatchInterface) *mockWatchInterface {
	mock := &mockWatchInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}