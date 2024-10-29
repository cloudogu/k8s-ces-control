// Code generated by mockery v2.30.1. DO NOT EDIT.

package main

import (
	appsv1 "k8s.io/api/apps/v1"
	apiautoscalingv1 "k8s.io/api/autoscaling/v1"

	autoscalingv1 "k8s.io/client-go/applyconfigurations/autoscaling/v1"

	context "context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mock "github.com/stretchr/testify/mock"

	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/client-go/applyconfigurations/apps/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// mockDeploymentInterface is an autogenerated mock type for the deploymentInterface type
type mockDeploymentInterface struct {
	mock.Mock
}

type mockDeploymentInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDeploymentInterface) EXPECT() *mockDeploymentInterface_Expecter {
	return &mockDeploymentInterface_Expecter{mock: &_m.Mock}
}

// Apply provides a mock function with given fields: ctx, deployment, opts
func (_m *mockDeploymentInterface) Apply(ctx context.Context, deployment *v1.DeploymentApplyConfiguration, opts metav1.ApplyOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, deployment, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, deployment, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, deployment, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, deployment, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Apply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Apply'
type mockDeploymentInterface_Apply_Call struct {
	*mock.Call
}

// Apply is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment *v1.DeploymentApplyConfiguration
//   - opts metav1.ApplyOptions
func (_e *mockDeploymentInterface_Expecter) Apply(ctx interface{}, deployment interface{}, opts interface{}) *mockDeploymentInterface_Apply_Call {
	return &mockDeploymentInterface_Apply_Call{Call: _e.mock.On("Apply", ctx, deployment, opts)}
}

func (_c *mockDeploymentInterface_Apply_Call) Run(run func(ctx context.Context, deployment *v1.DeploymentApplyConfiguration, opts metav1.ApplyOptions)) *mockDeploymentInterface_Apply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.DeploymentApplyConfiguration), args[2].(metav1.ApplyOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Apply_Call) Return(result *appsv1.Deployment, err error) *mockDeploymentInterface_Apply_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDeploymentInterface_Apply_Call) RunAndReturn(run func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_Apply_Call {
	_c.Call.Return(run)
	return _c
}

// ApplyScale provides a mock function with given fields: ctx, deploymentName, scale, opts
func (_m *mockDeploymentInterface) ApplyScale(ctx context.Context, deploymentName string, scale *autoscalingv1.ScaleApplyConfiguration, opts metav1.ApplyOptions) (*apiautoscalingv1.Scale, error) {
	ret := _m.Called(ctx, deploymentName, scale, opts)

	var r0 *apiautoscalingv1.Scale
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *autoscalingv1.ScaleApplyConfiguration, metav1.ApplyOptions) (*apiautoscalingv1.Scale, error)); ok {
		return rf(ctx, deploymentName, scale, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *autoscalingv1.ScaleApplyConfiguration, metav1.ApplyOptions) *apiautoscalingv1.Scale); ok {
		r0 = rf(ctx, deploymentName, scale, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apiautoscalingv1.Scale)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *autoscalingv1.ScaleApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, deploymentName, scale, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_ApplyScale_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyScale'
type mockDeploymentInterface_ApplyScale_Call struct {
	*mock.Call
}

// ApplyScale is a helper method to define mock.On call
//   - ctx context.Context
//   - deploymentName string
//   - scale *autoscalingv1.ScaleApplyConfiguration
//   - opts metav1.ApplyOptions
func (_e *mockDeploymentInterface_Expecter) ApplyScale(ctx interface{}, deploymentName interface{}, scale interface{}, opts interface{}) *mockDeploymentInterface_ApplyScale_Call {
	return &mockDeploymentInterface_ApplyScale_Call{Call: _e.mock.On("ApplyScale", ctx, deploymentName, scale, opts)}
}

func (_c *mockDeploymentInterface_ApplyScale_Call) Run(run func(ctx context.Context, deploymentName string, scale *autoscalingv1.ScaleApplyConfiguration, opts metav1.ApplyOptions)) *mockDeploymentInterface_ApplyScale_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*autoscalingv1.ScaleApplyConfiguration), args[3].(metav1.ApplyOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_ApplyScale_Call) Return(_a0 *apiautoscalingv1.Scale, _a1 error) *mockDeploymentInterface_ApplyScale_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_ApplyScale_Call) RunAndReturn(run func(context.Context, string, *autoscalingv1.ScaleApplyConfiguration, metav1.ApplyOptions) (*apiautoscalingv1.Scale, error)) *mockDeploymentInterface_ApplyScale_Call {
	_c.Call.Return(run)
	return _c
}

// ApplyStatus provides a mock function with given fields: ctx, deployment, opts
func (_m *mockDeploymentInterface) ApplyStatus(ctx context.Context, deployment *v1.DeploymentApplyConfiguration, opts metav1.ApplyOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, deployment, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, deployment, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, deployment, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, deployment, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_ApplyStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyStatus'
type mockDeploymentInterface_ApplyStatus_Call struct {
	*mock.Call
}

// ApplyStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment *v1.DeploymentApplyConfiguration
//   - opts metav1.ApplyOptions
func (_e *mockDeploymentInterface_Expecter) ApplyStatus(ctx interface{}, deployment interface{}, opts interface{}) *mockDeploymentInterface_ApplyStatus_Call {
	return &mockDeploymentInterface_ApplyStatus_Call{Call: _e.mock.On("ApplyStatus", ctx, deployment, opts)}
}

func (_c *mockDeploymentInterface_ApplyStatus_Call) Run(run func(ctx context.Context, deployment *v1.DeploymentApplyConfiguration, opts metav1.ApplyOptions)) *mockDeploymentInterface_ApplyStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.DeploymentApplyConfiguration), args[2].(metav1.ApplyOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_ApplyStatus_Call) Return(result *appsv1.Deployment, err error) *mockDeploymentInterface_ApplyStatus_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDeploymentInterface_ApplyStatus_Call) RunAndReturn(run func(context.Context, *v1.DeploymentApplyConfiguration, metav1.ApplyOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_ApplyStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, deployment, opts
func (_m *mockDeploymentInterface) Create(ctx context.Context, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, deployment, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.CreateOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, deployment, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.CreateOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, deployment, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.Deployment, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, deployment, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDeploymentInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment *appsv1.Deployment
//   - opts metav1.CreateOptions
func (_e *mockDeploymentInterface_Expecter) Create(ctx interface{}, deployment interface{}, opts interface{}) *mockDeploymentInterface_Create_Call {
	return &mockDeploymentInterface_Create_Call{Call: _e.mock.On("Create", ctx, deployment, opts)}
}

func (_c *mockDeploymentInterface_Create_Call) Run(run func(ctx context.Context, deployment *appsv1.Deployment, opts metav1.CreateOptions)) *mockDeploymentInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.Deployment), args[2].(metav1.CreateOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Create_Call) Return(_a0 *appsv1.Deployment, _a1 error) *mockDeploymentInterface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_Create_Call) RunAndReturn(run func(context.Context, *appsv1.Deployment, metav1.CreateOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *mockDeploymentInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDeploymentInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDeploymentInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.DeleteOptions
func (_e *mockDeploymentInterface_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *mockDeploymentInterface_Delete_Call {
	return &mockDeploymentInterface_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *mockDeploymentInterface_Delete_Call) Run(run func(ctx context.Context, name string, opts metav1.DeleteOptions)) *mockDeploymentInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.DeleteOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Delete_Call) Return(_a0 error) *mockDeploymentInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDeploymentInterface_Delete_Call) RunAndReturn(run func(context.Context, string, metav1.DeleteOptions) error) *mockDeploymentInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *mockDeploymentInterface) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDeploymentInterface_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type mockDeploymentInterface_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.DeleteOptions
//   - listOpts metav1.ListOptions
func (_e *mockDeploymentInterface_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *mockDeploymentInterface_DeleteCollection_Call {
	return &mockDeploymentInterface_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *mockDeploymentInterface_DeleteCollection_Call) Run(run func(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions)) *mockDeploymentInterface_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.DeleteOptions), args[2].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_DeleteCollection_Call) Return(_a0 error) *mockDeploymentInterface_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDeploymentInterface_DeleteCollection_Call) RunAndReturn(run func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error) *mockDeploymentInterface_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockDeploymentInterface) Get(ctx context.Context, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, name, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDeploymentInterface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.GetOptions
func (_e *mockDeploymentInterface_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockDeploymentInterface_Get_Call {
	return &mockDeploymentInterface_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockDeploymentInterface_Get_Call) Run(run func(ctx context.Context, name string, opts metav1.GetOptions)) *mockDeploymentInterface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Get_Call) Return(_a0 *appsv1.Deployment, _a1 error) *mockDeploymentInterface_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_Get_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetScale provides a mock function with given fields: ctx, deploymentName, options
func (_m *mockDeploymentInterface) GetScale(ctx context.Context, deploymentName string, options metav1.GetOptions) (*apiautoscalingv1.Scale, error) {
	ret := _m.Called(ctx, deploymentName, options)

	var r0 *apiautoscalingv1.Scale
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*apiautoscalingv1.Scale, error)); ok {
		return rf(ctx, deploymentName, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *apiautoscalingv1.Scale); ok {
		r0 = rf(ctx, deploymentName, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apiautoscalingv1.Scale)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, deploymentName, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_GetScale_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetScale'
type mockDeploymentInterface_GetScale_Call struct {
	*mock.Call
}

// GetScale is a helper method to define mock.On call
//   - ctx context.Context
//   - deploymentName string
//   - options metav1.GetOptions
func (_e *mockDeploymentInterface_Expecter) GetScale(ctx interface{}, deploymentName interface{}, options interface{}) *mockDeploymentInterface_GetScale_Call {
	return &mockDeploymentInterface_GetScale_Call{Call: _e.mock.On("GetScale", ctx, deploymentName, options)}
}

func (_c *mockDeploymentInterface_GetScale_Call) Run(run func(ctx context.Context, deploymentName string, options metav1.GetOptions)) *mockDeploymentInterface_GetScale_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_GetScale_Call) Return(_a0 *apiautoscalingv1.Scale, _a1 error) *mockDeploymentInterface_GetScale_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_GetScale_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*apiautoscalingv1.Scale, error)) *mockDeploymentInterface_GetScale_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *mockDeploymentInterface) List(ctx context.Context, opts metav1.ListOptions) (*appsv1.DeploymentList, error) {
	ret := _m.Called(ctx, opts)

	var r0 *appsv1.DeploymentList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*appsv1.DeploymentList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *appsv1.DeploymentList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DeploymentList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockDeploymentInterface_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDeploymentInterface_Expecter) List(ctx interface{}, opts interface{}) *mockDeploymentInterface_List_Call {
	return &mockDeploymentInterface_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *mockDeploymentInterface_List_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDeploymentInterface_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_List_Call) Return(_a0 *appsv1.DeploymentList, _a1 error) *mockDeploymentInterface_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_List_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (*appsv1.DeploymentList, error)) *mockDeploymentInterface_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *mockDeploymentInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*appsv1.Deployment, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*appsv1.Deployment, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *appsv1.Deployment); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockDeploymentInterface_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts metav1.PatchOptions
//   - subresources ...string
func (_e *mockDeploymentInterface_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *mockDeploymentInterface_Patch_Call {
	return &mockDeploymentInterface_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *mockDeploymentInterface_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string)) *mockDeploymentInterface_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-5)
		for i, a := range args[5:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(types.PatchType), args[3].([]byte), args[4].(metav1.PatchOptions), variadicArgs...)
	})
	return _c
}

func (_c *mockDeploymentInterface_Patch_Call) Return(result *appsv1.Deployment, err error) *mockDeploymentInterface_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDeploymentInterface_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*appsv1.Deployment, error)) *mockDeploymentInterface_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, deployment, opts
func (_m *mockDeploymentInterface) Update(ctx context.Context, deployment *appsv1.Deployment, opts metav1.UpdateOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, deployment, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, deployment, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, deployment, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, deployment, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockDeploymentInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment *appsv1.Deployment
//   - opts metav1.UpdateOptions
func (_e *mockDeploymentInterface_Expecter) Update(ctx interface{}, deployment interface{}, opts interface{}) *mockDeploymentInterface_Update_Call {
	return &mockDeploymentInterface_Update_Call{Call: _e.mock.On("Update", ctx, deployment, opts)}
}

func (_c *mockDeploymentInterface_Update_Call) Run(run func(ctx context.Context, deployment *appsv1.Deployment, opts metav1.UpdateOptions)) *mockDeploymentInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.Deployment), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Update_Call) Return(_a0 *appsv1.Deployment, _a1 error) *mockDeploymentInterface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_Update_Call) RunAndReturn(run func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateScale provides a mock function with given fields: ctx, deploymentName, scale, opts
func (_m *mockDeploymentInterface) UpdateScale(ctx context.Context, deploymentName string, scale *apiautoscalingv1.Scale, opts metav1.UpdateOptions) (*apiautoscalingv1.Scale, error) {
	ret := _m.Called(ctx, deploymentName, scale, opts)

	var r0 *apiautoscalingv1.Scale
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *apiautoscalingv1.Scale, metav1.UpdateOptions) (*apiautoscalingv1.Scale, error)); ok {
		return rf(ctx, deploymentName, scale, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *apiautoscalingv1.Scale, metav1.UpdateOptions) *apiautoscalingv1.Scale); ok {
		r0 = rf(ctx, deploymentName, scale, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apiautoscalingv1.Scale)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *apiautoscalingv1.Scale, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, deploymentName, scale, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_UpdateScale_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateScale'
type mockDeploymentInterface_UpdateScale_Call struct {
	*mock.Call
}

// UpdateScale is a helper method to define mock.On call
//   - ctx context.Context
//   - deploymentName string
//   - scale *apiautoscalingv1.Scale
//   - opts metav1.UpdateOptions
func (_e *mockDeploymentInterface_Expecter) UpdateScale(ctx interface{}, deploymentName interface{}, scale interface{}, opts interface{}) *mockDeploymentInterface_UpdateScale_Call {
	return &mockDeploymentInterface_UpdateScale_Call{Call: _e.mock.On("UpdateScale", ctx, deploymentName, scale, opts)}
}

func (_c *mockDeploymentInterface_UpdateScale_Call) Run(run func(ctx context.Context, deploymentName string, scale *apiautoscalingv1.Scale, opts metav1.UpdateOptions)) *mockDeploymentInterface_UpdateScale_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*apiautoscalingv1.Scale), args[3].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_UpdateScale_Call) Return(_a0 *apiautoscalingv1.Scale, _a1 error) *mockDeploymentInterface_UpdateScale_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_UpdateScale_Call) RunAndReturn(run func(context.Context, string, *apiautoscalingv1.Scale, metav1.UpdateOptions) (*apiautoscalingv1.Scale, error)) *mockDeploymentInterface_UpdateScale_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, deployment, opts
func (_m *mockDeploymentInterface) UpdateStatus(ctx context.Context, deployment *appsv1.Deployment, opts metav1.UpdateOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, deployment, opts)

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, deployment, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, deployment, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, deployment, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type mockDeploymentInterface_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment *appsv1.Deployment
//   - opts metav1.UpdateOptions
func (_e *mockDeploymentInterface_Expecter) UpdateStatus(ctx interface{}, deployment interface{}, opts interface{}) *mockDeploymentInterface_UpdateStatus_Call {
	return &mockDeploymentInterface_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, deployment, opts)}
}

func (_c *mockDeploymentInterface_UpdateStatus_Call) Run(run func(ctx context.Context, deployment *appsv1.Deployment, opts metav1.UpdateOptions)) *mockDeploymentInterface_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.Deployment), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_UpdateStatus_Call) Return(_a0 *appsv1.Deployment, _a1 error) *mockDeploymentInterface_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_UpdateStatus_Call) RunAndReturn(run func(context.Context, *appsv1.Deployment, metav1.UpdateOptions) (*appsv1.Deployment, error)) *mockDeploymentInterface_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *mockDeploymentInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (watch.Interface, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentInterface_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type mockDeploymentInterface_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDeploymentInterface_Expecter) Watch(ctx interface{}, opts interface{}) *mockDeploymentInterface_Watch_Call {
	return &mockDeploymentInterface_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *mockDeploymentInterface_Watch_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDeploymentInterface_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDeploymentInterface_Watch_Call) Return(_a0 watch.Interface, _a1 error) *mockDeploymentInterface_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentInterface_Watch_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (watch.Interface, error)) *mockDeploymentInterface_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDeploymentInterface creates a new instance of mockDeploymentInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDeploymentInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDeploymentInterface {
	mock := &mockDeploymentInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
