// Code generated by mockery v2.42.1. DO NOT EDIT.

package logging

import mock "github.com/stretchr/testify/mock"

// mockLogProvider is an autogenerated mock type for the logProvider type
type mockLogProvider struct {
	mock.Mock
}

type mockLogProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *mockLogProvider) EXPECT() *mockLogProvider_Expecter {
	return &mockLogProvider_Expecter{mock: &_m.Mock}
}

// getLogs provides a mock function with given fields: doguName, linesCount
func (_m *mockLogProvider) getLogs(doguName string, linesCount int) ([]logLine, error) {
	ret := _m.Called(doguName, linesCount)

	if len(ret) == 0 {
		panic("no return value specified for getLogs")
	}

	var r0 []logLine
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int) ([]logLine, error)); ok {
		return rf(doguName, linesCount)
	}
	if rf, ok := ret.Get(0).(func(string, int) []logLine); ok {
		r0 = rf(doguName, linesCount)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logLine)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(doguName, linesCount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockLogProvider_getLogs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'getLogs'
type mockLogProvider_getLogs_Call struct {
	*mock.Call
}

// getLogs is a helper method to define mock.On call
//   - doguName string
//   - linesCount int
func (_e *mockLogProvider_Expecter) getLogs(doguName interface{}, linesCount interface{}) *mockLogProvider_getLogs_Call {
	return &mockLogProvider_getLogs_Call{Call: _e.mock.On("getLogs", doguName, linesCount)}
}

func (_c *mockLogProvider_getLogs_Call) Run(run func(doguName string, linesCount int)) *mockLogProvider_getLogs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int))
	})
	return _c
}

func (_c *mockLogProvider_getLogs_Call) Return(_a0 []logLine, _a1 error) *mockLogProvider_getLogs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockLogProvider_getLogs_Call) RunAndReturn(run func(string, int) ([]logLine, error)) *mockLogProvider_getLogs_Call {
	_c.Call.Return(run)
	return _c
}

// newMockLogProvider creates a new instance of mockLogProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockLogProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockLogProvider {
	mock := &mockLogProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
