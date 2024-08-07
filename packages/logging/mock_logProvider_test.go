// Code generated by mockery v2.20.0. DO NOT EDIT.

package logging

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

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

// queryLogs provides a mock function with given fields: doguName, startDate, endDate, filter
func (_m *mockLogProvider) queryLogs(doguName string, startDate time.Time, endDate time.Time, filter string) ([]logLine, error) {
	ret := _m.Called(doguName, startDate, endDate, filter)

	var r0 []logLine
	var r1 error
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time, string) ([]logLine, error)); ok {
		return rf(doguName, startDate, endDate, filter)
	}
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time, string) []logLine); ok {
		r0 = rf(doguName, startDate, endDate, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logLine)
		}
	}

	if rf, ok := ret.Get(1).(func(string, time.Time, time.Time, string) error); ok {
		r1 = rf(doguName, startDate, endDate, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockLogProvider_queryLogs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'queryLogs'
type mockLogProvider_queryLogs_Call struct {
	*mock.Call
}

// queryLogs is a helper method to define mock.On call
//   - doguName string
//   - startDate time.Time
//   - endDate time.Time
//   - filter string
func (_e *mockLogProvider_Expecter) queryLogs(doguName interface{}, startDate interface{}, endDate interface{}, filter interface{}) *mockLogProvider_queryLogs_Call {
	return &mockLogProvider_queryLogs_Call{Call: _e.mock.On("queryLogs", doguName, startDate, endDate, filter)}
}

func (_c *mockLogProvider_queryLogs_Call) Run(run func(doguName string, startDate time.Time, endDate time.Time, filter string)) *mockLogProvider_queryLogs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Time), args[2].(time.Time), args[3].(string))
	})
	return _c
}

func (_c *mockLogProvider_queryLogs_Call) Return(_a0 []logLine, _a1 error) *mockLogProvider_queryLogs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockLogProvider_queryLogs_Call) RunAndReturn(run func(string, time.Time, time.Time, string) ([]logLine, error)) *mockLogProvider_queryLogs_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockLogProvider interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLogProvider creates a new instance of mockLogProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLogProvider(t mockConstructorTestingTnewMockLogProvider) *mockLogProvider {
	mock := &mockLogProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
