// Code generated by mockery v2.20.0. DO NOT EDIT.

package doguAdministration

import (
	mock "github.com/stretchr/testify/mock"
	rest "k8s.io/client-go/rest"

	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

// mockAppsV1Client is an autogenerated mock type for the appsV1Client type
type mockAppsV1Client struct {
	mock.Mock
}

type mockAppsV1Client_Expecter struct {
	mock *mock.Mock
}

func (_m *mockAppsV1Client) EXPECT() *mockAppsV1Client_Expecter {
	return &mockAppsV1Client_Expecter{mock: &_m.Mock}
}

// ControllerRevisions provides a mock function with given fields: namespace
func (_m *mockAppsV1Client) ControllerRevisions(namespace string) v1.ControllerRevisionInterface {
	ret := _m.Called(namespace)

	var r0 v1.ControllerRevisionInterface
	if rf, ok := ret.Get(0).(func(string) v1.ControllerRevisionInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ControllerRevisionInterface)
		}
	}

	return r0
}

// mockAppsV1Client_ControllerRevisions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ControllerRevisions'
type mockAppsV1Client_ControllerRevisions_Call struct {
	*mock.Call
}

// ControllerRevisions is a helper method to define mock.On call
//   - namespace string
func (_e *mockAppsV1Client_Expecter) ControllerRevisions(namespace interface{}) *mockAppsV1Client_ControllerRevisions_Call {
	return &mockAppsV1Client_ControllerRevisions_Call{Call: _e.mock.On("ControllerRevisions", namespace)}
}

func (_c *mockAppsV1Client_ControllerRevisions_Call) Run(run func(namespace string)) *mockAppsV1Client_ControllerRevisions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockAppsV1Client_ControllerRevisions_Call) Return(_a0 v1.ControllerRevisionInterface) *mockAppsV1Client_ControllerRevisions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_ControllerRevisions_Call) RunAndReturn(run func(string) v1.ControllerRevisionInterface) *mockAppsV1Client_ControllerRevisions_Call {
	_c.Call.Return(run)
	return _c
}

// DaemonSets provides a mock function with given fields: namespace
func (_m *mockAppsV1Client) DaemonSets(namespace string) v1.DaemonSetInterface {
	ret := _m.Called(namespace)

	var r0 v1.DaemonSetInterface
	if rf, ok := ret.Get(0).(func(string) v1.DaemonSetInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.DaemonSetInterface)
		}
	}

	return r0
}

// mockAppsV1Client_DaemonSets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DaemonSets'
type mockAppsV1Client_DaemonSets_Call struct {
	*mock.Call
}

// DaemonSets is a helper method to define mock.On call
//   - namespace string
func (_e *mockAppsV1Client_Expecter) DaemonSets(namespace interface{}) *mockAppsV1Client_DaemonSets_Call {
	return &mockAppsV1Client_DaemonSets_Call{Call: _e.mock.On("DaemonSets", namespace)}
}

func (_c *mockAppsV1Client_DaemonSets_Call) Run(run func(namespace string)) *mockAppsV1Client_DaemonSets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockAppsV1Client_DaemonSets_Call) Return(_a0 v1.DaemonSetInterface) *mockAppsV1Client_DaemonSets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_DaemonSets_Call) RunAndReturn(run func(string) v1.DaemonSetInterface) *mockAppsV1Client_DaemonSets_Call {
	_c.Call.Return(run)
	return _c
}

// Deployments provides a mock function with given fields: namespace
func (_m *mockAppsV1Client) Deployments(namespace string) v1.DeploymentInterface {
	ret := _m.Called(namespace)

	var r0 v1.DeploymentInterface
	if rf, ok := ret.Get(0).(func(string) v1.DeploymentInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.DeploymentInterface)
		}
	}

	return r0
}

// mockAppsV1Client_Deployments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deployments'
type mockAppsV1Client_Deployments_Call struct {
	*mock.Call
}

// Deployments is a helper method to define mock.On call
//   - namespace string
func (_e *mockAppsV1Client_Expecter) Deployments(namespace interface{}) *mockAppsV1Client_Deployments_Call {
	return &mockAppsV1Client_Deployments_Call{Call: _e.mock.On("Deployments", namespace)}
}

func (_c *mockAppsV1Client_Deployments_Call) Run(run func(namespace string)) *mockAppsV1Client_Deployments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockAppsV1Client_Deployments_Call) Return(_a0 v1.DeploymentInterface) *mockAppsV1Client_Deployments_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_Deployments_Call) RunAndReturn(run func(string) v1.DeploymentInterface) *mockAppsV1Client_Deployments_Call {
	_c.Call.Return(run)
	return _c
}

// RESTClient provides a mock function with given fields:
func (_m *mockAppsV1Client) RESTClient() rest.Interface {
	ret := _m.Called()

	var r0 rest.Interface
	if rf, ok := ret.Get(0).(func() rest.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rest.Interface)
		}
	}

	return r0
}

// mockAppsV1Client_RESTClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RESTClient'
type mockAppsV1Client_RESTClient_Call struct {
	*mock.Call
}

// RESTClient is a helper method to define mock.On call
func (_e *mockAppsV1Client_Expecter) RESTClient() *mockAppsV1Client_RESTClient_Call {
	return &mockAppsV1Client_RESTClient_Call{Call: _e.mock.On("RESTClient")}
}

func (_c *mockAppsV1Client_RESTClient_Call) Run(run func()) *mockAppsV1Client_RESTClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockAppsV1Client_RESTClient_Call) Return(_a0 rest.Interface) *mockAppsV1Client_RESTClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_RESTClient_Call) RunAndReturn(run func() rest.Interface) *mockAppsV1Client_RESTClient_Call {
	_c.Call.Return(run)
	return _c
}

// ReplicaSets provides a mock function with given fields: namespace
func (_m *mockAppsV1Client) ReplicaSets(namespace string) v1.ReplicaSetInterface {
	ret := _m.Called(namespace)

	var r0 v1.ReplicaSetInterface
	if rf, ok := ret.Get(0).(func(string) v1.ReplicaSetInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ReplicaSetInterface)
		}
	}

	return r0
}

// mockAppsV1Client_ReplicaSets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReplicaSets'
type mockAppsV1Client_ReplicaSets_Call struct {
	*mock.Call
}

// ReplicaSets is a helper method to define mock.On call
//   - namespace string
func (_e *mockAppsV1Client_Expecter) ReplicaSets(namespace interface{}) *mockAppsV1Client_ReplicaSets_Call {
	return &mockAppsV1Client_ReplicaSets_Call{Call: _e.mock.On("ReplicaSets", namespace)}
}

func (_c *mockAppsV1Client_ReplicaSets_Call) Run(run func(namespace string)) *mockAppsV1Client_ReplicaSets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockAppsV1Client_ReplicaSets_Call) Return(_a0 v1.ReplicaSetInterface) *mockAppsV1Client_ReplicaSets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_ReplicaSets_Call) RunAndReturn(run func(string) v1.ReplicaSetInterface) *mockAppsV1Client_ReplicaSets_Call {
	_c.Call.Return(run)
	return _c
}

// StatefulSets provides a mock function with given fields: namespace
func (_m *mockAppsV1Client) StatefulSets(namespace string) v1.StatefulSetInterface {
	ret := _m.Called(namespace)

	var r0 v1.StatefulSetInterface
	if rf, ok := ret.Get(0).(func(string) v1.StatefulSetInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.StatefulSetInterface)
		}
	}

	return r0
}

// mockAppsV1Client_StatefulSets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StatefulSets'
type mockAppsV1Client_StatefulSets_Call struct {
	*mock.Call
}

// StatefulSets is a helper method to define mock.On call
//   - namespace string
func (_e *mockAppsV1Client_Expecter) StatefulSets(namespace interface{}) *mockAppsV1Client_StatefulSets_Call {
	return &mockAppsV1Client_StatefulSets_Call{Call: _e.mock.On("StatefulSets", namespace)}
}

func (_c *mockAppsV1Client_StatefulSets_Call) Run(run func(namespace string)) *mockAppsV1Client_StatefulSets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockAppsV1Client_StatefulSets_Call) Return(_a0 v1.StatefulSetInterface) *mockAppsV1Client_StatefulSets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockAppsV1Client_StatefulSets_Call) RunAndReturn(run func(string) v1.StatefulSetInterface) *mockAppsV1Client_StatefulSets_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockAppsV1Client interface {
	mock.TestingT
	Cleanup(func())
}

// newMockAppsV1Client creates a new instance of mockAppsV1Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockAppsV1Client(t mockConstructorTestingTnewMockAppsV1Client) *mockAppsV1Client {
	mock := &mockAppsV1Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
