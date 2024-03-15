// Code generated by mockery v2.42.1. DO NOT EDIT.

package debug

import mock "github.com/stretchr/testify/mock"

// mockMaintenanceModeSwitch is an autogenerated mock type for the maintenanceModeSwitch type
type mockMaintenanceModeSwitch struct {
	mock.Mock
}

type mockMaintenanceModeSwitch_Expecter struct {
	mock *mock.Mock
}

func (_m *mockMaintenanceModeSwitch) EXPECT() *mockMaintenanceModeSwitch_Expecter {
	return &mockMaintenanceModeSwitch_Expecter{mock: &_m.Mock}
}

// ActivateMaintenanceMode provides a mock function with given fields: title, text
func (_m *mockMaintenanceModeSwitch) ActivateMaintenanceMode(title string, text string) error {
	ret := _m.Called(title, text)

	if len(ret) == 0 {
		panic("no return value specified for ActivateMaintenanceMode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(title, text)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ActivateMaintenanceMode'
type mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call struct {
	*mock.Call
}

// ActivateMaintenanceMode is a helper method to define mock.On call
//   - title string
//   - text string
func (_e *mockMaintenanceModeSwitch_Expecter) ActivateMaintenanceMode(title interface{}, text interface{}) *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call {
	return &mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call{Call: _e.mock.On("ActivateMaintenanceMode", title, text)}
}

func (_c *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call) Run(run func(title string, text string)) *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call) Return(_a0 error) *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call) RunAndReturn(run func(string, string) error) *mockMaintenanceModeSwitch_ActivateMaintenanceMode_Call {
	_c.Call.Return(run)
	return _c
}

// DeactivateMaintenanceMode provides a mock function with given fields:
func (_m *mockMaintenanceModeSwitch) DeactivateMaintenanceMode() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DeactivateMaintenanceMode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeactivateMaintenanceMode'
type mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call struct {
	*mock.Call
}

// DeactivateMaintenanceMode is a helper method to define mock.On call
func (_e *mockMaintenanceModeSwitch_Expecter) DeactivateMaintenanceMode() *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call {
	return &mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call{Call: _e.mock.On("DeactivateMaintenanceMode")}
}

func (_c *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call) Run(run func()) *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call) Return(_a0 error) *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call) RunAndReturn(run func() error) *mockMaintenanceModeSwitch_DeactivateMaintenanceMode_Call {
	_c.Call.Return(run)
	return _c
}

// newMockMaintenanceModeSwitch creates a new instance of mockMaintenanceModeSwitch. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockMaintenanceModeSwitch(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockMaintenanceModeSwitch {
	mock := &mockMaintenanceModeSwitch{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
