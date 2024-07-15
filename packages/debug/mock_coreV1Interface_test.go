// Code generated by mockery v2.20.0. DO NOT EDIT.

package debug

import (
	mock "github.com/stretchr/testify/mock"
	rest "k8s.io/client-go/rest"

	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// mockCoreV1Interface is an autogenerated mock type for the coreV1Interface type
type mockCoreV1Interface struct {
	mock.Mock
}

type mockCoreV1Interface_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCoreV1Interface) EXPECT() *mockCoreV1Interface_Expecter {
	return &mockCoreV1Interface_Expecter{mock: &_m.Mock}
}

// ComponentStatuses provides a mock function with given fields:
func (_m *mockCoreV1Interface) ComponentStatuses() v1.ComponentStatusInterface {
	ret := _m.Called()

	var r0 v1.ComponentStatusInterface
	if rf, ok := ret.Get(0).(func() v1.ComponentStatusInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ComponentStatusInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_ComponentStatuses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ComponentStatuses'
type mockCoreV1Interface_ComponentStatuses_Call struct {
	*mock.Call
}

// ComponentStatuses is a helper method to define mock.On call
func (_e *mockCoreV1Interface_Expecter) ComponentStatuses() *mockCoreV1Interface_ComponentStatuses_Call {
	return &mockCoreV1Interface_ComponentStatuses_Call{Call: _e.mock.On("ComponentStatuses")}
}

func (_c *mockCoreV1Interface_ComponentStatuses_Call) Run(run func()) *mockCoreV1Interface_ComponentStatuses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCoreV1Interface_ComponentStatuses_Call) Return(_a0 v1.ComponentStatusInterface) *mockCoreV1Interface_ComponentStatuses_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_ComponentStatuses_Call) RunAndReturn(run func() v1.ComponentStatusInterface) *mockCoreV1Interface_ComponentStatuses_Call {
	_c.Call.Return(run)
	return _c
}

// ConfigMaps provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) ConfigMaps(namespace string) v1.ConfigMapInterface {
	ret := _m.Called(namespace)

	var r0 v1.ConfigMapInterface
	if rf, ok := ret.Get(0).(func(string) v1.ConfigMapInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ConfigMapInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_ConfigMaps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigMaps'
type mockCoreV1Interface_ConfigMaps_Call struct {
	*mock.Call
}

// ConfigMaps is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) ConfigMaps(namespace interface{}) *mockCoreV1Interface_ConfigMaps_Call {
	return &mockCoreV1Interface_ConfigMaps_Call{Call: _e.mock.On("ConfigMaps", namespace)}
}

func (_c *mockCoreV1Interface_ConfigMaps_Call) Run(run func(namespace string)) *mockCoreV1Interface_ConfigMaps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_ConfigMaps_Call) Return(_a0 v1.ConfigMapInterface) *mockCoreV1Interface_ConfigMaps_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_ConfigMaps_Call) RunAndReturn(run func(string) v1.ConfigMapInterface) *mockCoreV1Interface_ConfigMaps_Call {
	_c.Call.Return(run)
	return _c
}

// Endpoints provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) Endpoints(namespace string) v1.EndpointsInterface {
	ret := _m.Called(namespace)

	var r0 v1.EndpointsInterface
	if rf, ok := ret.Get(0).(func(string) v1.EndpointsInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.EndpointsInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Endpoints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Endpoints'
type mockCoreV1Interface_Endpoints_Call struct {
	*mock.Call
}

// Endpoints is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) Endpoints(namespace interface{}) *mockCoreV1Interface_Endpoints_Call {
	return &mockCoreV1Interface_Endpoints_Call{Call: _e.mock.On("Endpoints", namespace)}
}

func (_c *mockCoreV1Interface_Endpoints_Call) Run(run func(namespace string)) *mockCoreV1Interface_Endpoints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_Endpoints_Call) Return(_a0 v1.EndpointsInterface) *mockCoreV1Interface_Endpoints_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Endpoints_Call) RunAndReturn(run func(string) v1.EndpointsInterface) *mockCoreV1Interface_Endpoints_Call {
	_c.Call.Return(run)
	return _c
}

// Events provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) Events(namespace string) v1.EventInterface {
	ret := _m.Called(namespace)

	var r0 v1.EventInterface
	if rf, ok := ret.Get(0).(func(string) v1.EventInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.EventInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Events_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Events'
type mockCoreV1Interface_Events_Call struct {
	*mock.Call
}

// Events is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) Events(namespace interface{}) *mockCoreV1Interface_Events_Call {
	return &mockCoreV1Interface_Events_Call{Call: _e.mock.On("Events", namespace)}
}

func (_c *mockCoreV1Interface_Events_Call) Run(run func(namespace string)) *mockCoreV1Interface_Events_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_Events_Call) Return(_a0 v1.EventInterface) *mockCoreV1Interface_Events_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Events_Call) RunAndReturn(run func(string) v1.EventInterface) *mockCoreV1Interface_Events_Call {
	_c.Call.Return(run)
	return _c
}

// LimitRanges provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) LimitRanges(namespace string) v1.LimitRangeInterface {
	ret := _m.Called(namespace)

	var r0 v1.LimitRangeInterface
	if rf, ok := ret.Get(0).(func(string) v1.LimitRangeInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.LimitRangeInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_LimitRanges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LimitRanges'
type mockCoreV1Interface_LimitRanges_Call struct {
	*mock.Call
}

// LimitRanges is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) LimitRanges(namespace interface{}) *mockCoreV1Interface_LimitRanges_Call {
	return &mockCoreV1Interface_LimitRanges_Call{Call: _e.mock.On("LimitRanges", namespace)}
}

func (_c *mockCoreV1Interface_LimitRanges_Call) Run(run func(namespace string)) *mockCoreV1Interface_LimitRanges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_LimitRanges_Call) Return(_a0 v1.LimitRangeInterface) *mockCoreV1Interface_LimitRanges_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_LimitRanges_Call) RunAndReturn(run func(string) v1.LimitRangeInterface) *mockCoreV1Interface_LimitRanges_Call {
	_c.Call.Return(run)
	return _c
}

// Namespaces provides a mock function with given fields:
func (_m *mockCoreV1Interface) Namespaces() v1.NamespaceInterface {
	ret := _m.Called()

	var r0 v1.NamespaceInterface
	if rf, ok := ret.Get(0).(func() v1.NamespaceInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.NamespaceInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Namespaces_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Namespaces'
type mockCoreV1Interface_Namespaces_Call struct {
	*mock.Call
}

// Namespaces is a helper method to define mock.On call
func (_e *mockCoreV1Interface_Expecter) Namespaces() *mockCoreV1Interface_Namespaces_Call {
	return &mockCoreV1Interface_Namespaces_Call{Call: _e.mock.On("Namespaces")}
}

func (_c *mockCoreV1Interface_Namespaces_Call) Run(run func()) *mockCoreV1Interface_Namespaces_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCoreV1Interface_Namespaces_Call) Return(_a0 v1.NamespaceInterface) *mockCoreV1Interface_Namespaces_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Namespaces_Call) RunAndReturn(run func() v1.NamespaceInterface) *mockCoreV1Interface_Namespaces_Call {
	_c.Call.Return(run)
	return _c
}

// Nodes provides a mock function with given fields:
func (_m *mockCoreV1Interface) Nodes() v1.NodeInterface {
	ret := _m.Called()

	var r0 v1.NodeInterface
	if rf, ok := ret.Get(0).(func() v1.NodeInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.NodeInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Nodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Nodes'
type mockCoreV1Interface_Nodes_Call struct {
	*mock.Call
}

// Nodes is a helper method to define mock.On call
func (_e *mockCoreV1Interface_Expecter) Nodes() *mockCoreV1Interface_Nodes_Call {
	return &mockCoreV1Interface_Nodes_Call{Call: _e.mock.On("Nodes")}
}

func (_c *mockCoreV1Interface_Nodes_Call) Run(run func()) *mockCoreV1Interface_Nodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCoreV1Interface_Nodes_Call) Return(_a0 v1.NodeInterface) *mockCoreV1Interface_Nodes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Nodes_Call) RunAndReturn(run func() v1.NodeInterface) *mockCoreV1Interface_Nodes_Call {
	_c.Call.Return(run)
	return _c
}

// PersistentVolumeClaims provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) PersistentVolumeClaims(namespace string) v1.PersistentVolumeClaimInterface {
	ret := _m.Called(namespace)

	var r0 v1.PersistentVolumeClaimInterface
	if rf, ok := ret.Get(0).(func(string) v1.PersistentVolumeClaimInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.PersistentVolumeClaimInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_PersistentVolumeClaims_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PersistentVolumeClaims'
type mockCoreV1Interface_PersistentVolumeClaims_Call struct {
	*mock.Call
}

// PersistentVolumeClaims is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) PersistentVolumeClaims(namespace interface{}) *mockCoreV1Interface_PersistentVolumeClaims_Call {
	return &mockCoreV1Interface_PersistentVolumeClaims_Call{Call: _e.mock.On("PersistentVolumeClaims", namespace)}
}

func (_c *mockCoreV1Interface_PersistentVolumeClaims_Call) Run(run func(namespace string)) *mockCoreV1Interface_PersistentVolumeClaims_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_PersistentVolumeClaims_Call) Return(_a0 v1.PersistentVolumeClaimInterface) *mockCoreV1Interface_PersistentVolumeClaims_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_PersistentVolumeClaims_Call) RunAndReturn(run func(string) v1.PersistentVolumeClaimInterface) *mockCoreV1Interface_PersistentVolumeClaims_Call {
	_c.Call.Return(run)
	return _c
}

// PersistentVolumes provides a mock function with given fields:
func (_m *mockCoreV1Interface) PersistentVolumes() v1.PersistentVolumeInterface {
	ret := _m.Called()

	var r0 v1.PersistentVolumeInterface
	if rf, ok := ret.Get(0).(func() v1.PersistentVolumeInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.PersistentVolumeInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_PersistentVolumes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PersistentVolumes'
type mockCoreV1Interface_PersistentVolumes_Call struct {
	*mock.Call
}

// PersistentVolumes is a helper method to define mock.On call
func (_e *mockCoreV1Interface_Expecter) PersistentVolumes() *mockCoreV1Interface_PersistentVolumes_Call {
	return &mockCoreV1Interface_PersistentVolumes_Call{Call: _e.mock.On("PersistentVolumes")}
}

func (_c *mockCoreV1Interface_PersistentVolumes_Call) Run(run func()) *mockCoreV1Interface_PersistentVolumes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCoreV1Interface_PersistentVolumes_Call) Return(_a0 v1.PersistentVolumeInterface) *mockCoreV1Interface_PersistentVolumes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_PersistentVolumes_Call) RunAndReturn(run func() v1.PersistentVolumeInterface) *mockCoreV1Interface_PersistentVolumes_Call {
	_c.Call.Return(run)
	return _c
}

// PodTemplates provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) PodTemplates(namespace string) v1.PodTemplateInterface {
	ret := _m.Called(namespace)

	var r0 v1.PodTemplateInterface
	if rf, ok := ret.Get(0).(func(string) v1.PodTemplateInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.PodTemplateInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_PodTemplates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PodTemplates'
type mockCoreV1Interface_PodTemplates_Call struct {
	*mock.Call
}

// PodTemplates is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) PodTemplates(namespace interface{}) *mockCoreV1Interface_PodTemplates_Call {
	return &mockCoreV1Interface_PodTemplates_Call{Call: _e.mock.On("PodTemplates", namespace)}
}

func (_c *mockCoreV1Interface_PodTemplates_Call) Run(run func(namespace string)) *mockCoreV1Interface_PodTemplates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_PodTemplates_Call) Return(_a0 v1.PodTemplateInterface) *mockCoreV1Interface_PodTemplates_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_PodTemplates_Call) RunAndReturn(run func(string) v1.PodTemplateInterface) *mockCoreV1Interface_PodTemplates_Call {
	_c.Call.Return(run)
	return _c
}

// Pods provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) Pods(namespace string) v1.PodInterface {
	ret := _m.Called(namespace)

	var r0 v1.PodInterface
	if rf, ok := ret.Get(0).(func(string) v1.PodInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.PodInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Pods_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Pods'
type mockCoreV1Interface_Pods_Call struct {
	*mock.Call
}

// Pods is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) Pods(namespace interface{}) *mockCoreV1Interface_Pods_Call {
	return &mockCoreV1Interface_Pods_Call{Call: _e.mock.On("Pods", namespace)}
}

func (_c *mockCoreV1Interface_Pods_Call) Run(run func(namespace string)) *mockCoreV1Interface_Pods_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_Pods_Call) Return(_a0 v1.PodInterface) *mockCoreV1Interface_Pods_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Pods_Call) RunAndReturn(run func(string) v1.PodInterface) *mockCoreV1Interface_Pods_Call {
	_c.Call.Return(run)
	return _c
}

// RESTClient provides a mock function with given fields:
func (_m *mockCoreV1Interface) RESTClient() rest.Interface {
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

// mockCoreV1Interface_RESTClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RESTClient'
type mockCoreV1Interface_RESTClient_Call struct {
	*mock.Call
}

// RESTClient is a helper method to define mock.On call
func (_e *mockCoreV1Interface_Expecter) RESTClient() *mockCoreV1Interface_RESTClient_Call {
	return &mockCoreV1Interface_RESTClient_Call{Call: _e.mock.On("RESTClient")}
}

func (_c *mockCoreV1Interface_RESTClient_Call) Run(run func()) *mockCoreV1Interface_RESTClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockCoreV1Interface_RESTClient_Call) Return(_a0 rest.Interface) *mockCoreV1Interface_RESTClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_RESTClient_Call) RunAndReturn(run func() rest.Interface) *mockCoreV1Interface_RESTClient_Call {
	_c.Call.Return(run)
	return _c
}

// ReplicationControllers provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) ReplicationControllers(namespace string) v1.ReplicationControllerInterface {
	ret := _m.Called(namespace)

	var r0 v1.ReplicationControllerInterface
	if rf, ok := ret.Get(0).(func(string) v1.ReplicationControllerInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ReplicationControllerInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_ReplicationControllers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReplicationControllers'
type mockCoreV1Interface_ReplicationControllers_Call struct {
	*mock.Call
}

// ReplicationControllers is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) ReplicationControllers(namespace interface{}) *mockCoreV1Interface_ReplicationControllers_Call {
	return &mockCoreV1Interface_ReplicationControllers_Call{Call: _e.mock.On("ReplicationControllers", namespace)}
}

func (_c *mockCoreV1Interface_ReplicationControllers_Call) Run(run func(namespace string)) *mockCoreV1Interface_ReplicationControllers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_ReplicationControllers_Call) Return(_a0 v1.ReplicationControllerInterface) *mockCoreV1Interface_ReplicationControllers_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_ReplicationControllers_Call) RunAndReturn(run func(string) v1.ReplicationControllerInterface) *mockCoreV1Interface_ReplicationControllers_Call {
	_c.Call.Return(run)
	return _c
}

// ResourceQuotas provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) ResourceQuotas(namespace string) v1.ResourceQuotaInterface {
	ret := _m.Called(namespace)

	var r0 v1.ResourceQuotaInterface
	if rf, ok := ret.Get(0).(func(string) v1.ResourceQuotaInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ResourceQuotaInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_ResourceQuotas_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResourceQuotas'
type mockCoreV1Interface_ResourceQuotas_Call struct {
	*mock.Call
}

// ResourceQuotas is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) ResourceQuotas(namespace interface{}) *mockCoreV1Interface_ResourceQuotas_Call {
	return &mockCoreV1Interface_ResourceQuotas_Call{Call: _e.mock.On("ResourceQuotas", namespace)}
}

func (_c *mockCoreV1Interface_ResourceQuotas_Call) Run(run func(namespace string)) *mockCoreV1Interface_ResourceQuotas_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_ResourceQuotas_Call) Return(_a0 v1.ResourceQuotaInterface) *mockCoreV1Interface_ResourceQuotas_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_ResourceQuotas_Call) RunAndReturn(run func(string) v1.ResourceQuotaInterface) *mockCoreV1Interface_ResourceQuotas_Call {
	_c.Call.Return(run)
	return _c
}

// Secrets provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) Secrets(namespace string) v1.SecretInterface {
	ret := _m.Called(namespace)

	var r0 v1.SecretInterface
	if rf, ok := ret.Get(0).(func(string) v1.SecretInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.SecretInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Secrets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Secrets'
type mockCoreV1Interface_Secrets_Call struct {
	*mock.Call
}

// Secrets is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) Secrets(namespace interface{}) *mockCoreV1Interface_Secrets_Call {
	return &mockCoreV1Interface_Secrets_Call{Call: _e.mock.On("Secrets", namespace)}
}

func (_c *mockCoreV1Interface_Secrets_Call) Run(run func(namespace string)) *mockCoreV1Interface_Secrets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_Secrets_Call) Return(_a0 v1.SecretInterface) *mockCoreV1Interface_Secrets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Secrets_Call) RunAndReturn(run func(string) v1.SecretInterface) *mockCoreV1Interface_Secrets_Call {
	_c.Call.Return(run)
	return _c
}

// ServiceAccounts provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) ServiceAccounts(namespace string) v1.ServiceAccountInterface {
	ret := _m.Called(namespace)

	var r0 v1.ServiceAccountInterface
	if rf, ok := ret.Get(0).(func(string) v1.ServiceAccountInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ServiceAccountInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_ServiceAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServiceAccounts'
type mockCoreV1Interface_ServiceAccounts_Call struct {
	*mock.Call
}

// ServiceAccounts is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) ServiceAccounts(namespace interface{}) *mockCoreV1Interface_ServiceAccounts_Call {
	return &mockCoreV1Interface_ServiceAccounts_Call{Call: _e.mock.On("ServiceAccounts", namespace)}
}

func (_c *mockCoreV1Interface_ServiceAccounts_Call) Run(run func(namespace string)) *mockCoreV1Interface_ServiceAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_ServiceAccounts_Call) Return(_a0 v1.ServiceAccountInterface) *mockCoreV1Interface_ServiceAccounts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_ServiceAccounts_Call) RunAndReturn(run func(string) v1.ServiceAccountInterface) *mockCoreV1Interface_ServiceAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// Services provides a mock function with given fields: namespace
func (_m *mockCoreV1Interface) Services(namespace string) v1.ServiceInterface {
	ret := _m.Called(namespace)

	var r0 v1.ServiceInterface
	if rf, ok := ret.Get(0).(func(string) v1.ServiceInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.ServiceInterface)
		}
	}

	return r0
}

// mockCoreV1Interface_Services_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Services'
type mockCoreV1Interface_Services_Call struct {
	*mock.Call
}

// Services is a helper method to define mock.On call
//   - namespace string
func (_e *mockCoreV1Interface_Expecter) Services(namespace interface{}) *mockCoreV1Interface_Services_Call {
	return &mockCoreV1Interface_Services_Call{Call: _e.mock.On("Services", namespace)}
}

func (_c *mockCoreV1Interface_Services_Call) Run(run func(namespace string)) *mockCoreV1Interface_Services_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockCoreV1Interface_Services_Call) Return(_a0 v1.ServiceInterface) *mockCoreV1Interface_Services_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockCoreV1Interface_Services_Call) RunAndReturn(run func(string) v1.ServiceInterface) *mockCoreV1Interface_Services_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockCoreV1Interface interface {
	mock.TestingT
	Cleanup(func())
}

// newMockCoreV1Interface creates a new instance of mockCoreV1Interface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockCoreV1Interface(t mockConstructorTestingTnewMockCoreV1Interface) *mockCoreV1Interface {
	mock := &mockCoreV1Interface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
