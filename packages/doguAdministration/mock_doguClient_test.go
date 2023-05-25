// Code generated by mockery v2.20.0. DO NOT EDIT.

package doguAdministration

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "k8s.io/apimachinery/pkg/types"

	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// mockDoguClient is an autogenerated mock type for the doguClient type
type mockDoguClient struct {
	mock.Mock
}

type mockDoguClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguClient) EXPECT() *mockDoguClient_Expecter {
	return &mockDoguClient_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, dogu, opts
func (_m *mockDoguClient) Create(ctx context.Context, dogu *v1.Dogu, opts metav1.CreateOptions) (*v1.Dogu, error) {
	ret := _m.Called(ctx, dogu, opts)

	var r0 *v1.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.CreateOptions) (*v1.Dogu, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.CreateOptions) *v1.Dogu); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Dogu, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDoguClient_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v1.Dogu
//   - opts metav1.CreateOptions
func (_e *mockDoguClient_Expecter) Create(ctx interface{}, dogu interface{}, opts interface{}) *mockDoguClient_Create_Call {
	return &mockDoguClient_Create_Call{Call: _e.mock.On("Create", ctx, dogu, opts)}
}

func (_c *mockDoguClient_Create_Call) Run(run func(ctx context.Context, dogu *v1.Dogu, opts metav1.CreateOptions)) *mockDoguClient_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Dogu), args[2].(metav1.CreateOptions))
	})
	return _c
}

func (_c *mockDoguClient_Create_Call) Return(_a0 *v1.Dogu, _a1 error) *mockDoguClient_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_Create_Call) RunAndReturn(run func(context.Context, *v1.Dogu, metav1.CreateOptions) (*v1.Dogu, error)) *mockDoguClient_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *mockDoguClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguClient_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDoguClient_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.DeleteOptions
func (_e *mockDoguClient_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *mockDoguClient_Delete_Call {
	return &mockDoguClient_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *mockDoguClient_Delete_Call) Run(run func(ctx context.Context, name string, opts metav1.DeleteOptions)) *mockDoguClient_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.DeleteOptions))
	})
	return _c
}

func (_c *mockDoguClient_Delete_Call) Return(_a0 error) *mockDoguClient_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguClient_Delete_Call) RunAndReturn(run func(context.Context, string, metav1.DeleteOptions) error) *mockDoguClient_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *mockDoguClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguClient_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type mockDoguClient_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.DeleteOptions
//   - listOpts metav1.ListOptions
func (_e *mockDoguClient_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *mockDoguClient_DeleteCollection_Call {
	return &mockDoguClient_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *mockDoguClient_DeleteCollection_Call) Run(run func(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions)) *mockDoguClient_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.DeleteOptions), args[2].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDoguClient_DeleteCollection_Call) Return(_a0 error) *mockDoguClient_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguClient_DeleteCollection_Call) RunAndReturn(run func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error) *mockDoguClient_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockDoguClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Dogu, error) {
	ret := _m.Called(ctx, name, opts)

	var r0 *v1.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*v1.Dogu, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *v1.Dogu); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDoguClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.GetOptions
func (_e *mockDoguClient_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockDoguClient_Get_Call {
	return &mockDoguClient_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockDoguClient_Get_Call) Run(run func(ctx context.Context, name string, opts metav1.GetOptions)) *mockDoguClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockDoguClient_Get_Call) Return(_a0 *v1.Dogu, _a1 error) *mockDoguClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_Get_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*v1.Dogu, error)) *mockDoguClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *mockDoguClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.DoguList, error) {
	ret := _m.Called(ctx, opts)

	var r0 *v1.DoguList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*v1.DoguList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *v1.DoguList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.DoguList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockDoguClient_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDoguClient_Expecter) List(ctx interface{}, opts interface{}) *mockDoguClient_List_Call {
	return &mockDoguClient_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *mockDoguClient_List_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDoguClient_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDoguClient_List_Call) Return(_a0 *v1.DoguList, _a1 error) *mockDoguClient_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_List_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (*v1.DoguList, error)) *mockDoguClient_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *mockDoguClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1.Dogu, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *v1.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*v1.Dogu, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *v1.Dogu); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockDoguClient_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts metav1.PatchOptions
//   - subresources ...string
func (_e *mockDoguClient_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *mockDoguClient_Patch_Call {
	return &mockDoguClient_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *mockDoguClient_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string)) *mockDoguClient_Patch_Call {
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

func (_c *mockDoguClient_Patch_Call) Return(result *v1.Dogu, err error) *mockDoguClient_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDoguClient_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*v1.Dogu, error)) *mockDoguClient_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, dogu, opts
func (_m *mockDoguClient) Update(ctx context.Context, dogu *v1.Dogu, opts metav1.UpdateOptions) (*v1.Dogu, error) {
	ret := _m.Called(ctx, dogu, opts)

	var r0 *v1.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) (*v1.Dogu, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) *v1.Dogu); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockDoguClient_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v1.Dogu
//   - opts metav1.UpdateOptions
func (_e *mockDoguClient_Expecter) Update(ctx interface{}, dogu interface{}, opts interface{}) *mockDoguClient_Update_Call {
	return &mockDoguClient_Update_Call{Call: _e.mock.On("Update", ctx, dogu, opts)}
}

func (_c *mockDoguClient_Update_Call) Run(run func(ctx context.Context, dogu *v1.Dogu, opts metav1.UpdateOptions)) *mockDoguClient_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Dogu), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDoguClient_Update_Call) Return(_a0 *v1.Dogu, _a1 error) *mockDoguClient_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_Update_Call) RunAndReturn(run func(context.Context, *v1.Dogu, metav1.UpdateOptions) (*v1.Dogu, error)) *mockDoguClient_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, dogu, opts
func (_m *mockDoguClient) UpdateStatus(ctx context.Context, dogu *v1.Dogu, opts metav1.UpdateOptions) (*v1.Dogu, error) {
	ret := _m.Called(ctx, dogu, opts)

	var r0 *v1.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) (*v1.Dogu, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) *v1.Dogu); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Dogu, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguClient_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type mockDoguClient_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v1.Dogu
//   - opts metav1.UpdateOptions
func (_e *mockDoguClient_Expecter) UpdateStatus(ctx interface{}, dogu interface{}, opts interface{}) *mockDoguClient_UpdateStatus_Call {
	return &mockDoguClient_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, dogu, opts)}
}

func (_c *mockDoguClient_UpdateStatus_Call) Run(run func(ctx context.Context, dogu *v1.Dogu, opts metav1.UpdateOptions)) *mockDoguClient_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Dogu), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDoguClient_UpdateStatus_Call) Return(_a0 *v1.Dogu, _a1 error) *mockDoguClient_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_UpdateStatus_Call) RunAndReturn(run func(context.Context, *v1.Dogu, metav1.UpdateOptions) (*v1.Dogu, error)) *mockDoguClient_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *mockDoguClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
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

// mockDoguClient_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type mockDoguClient_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDoguClient_Expecter) Watch(ctx interface{}, opts interface{}) *mockDoguClient_Watch_Call {
	return &mockDoguClient_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *mockDoguClient_Watch_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDoguClient_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDoguClient_Watch_Call) Return(_a0 watch.Interface, _a1 error) *mockDoguClient_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguClient_Watch_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (watch.Interface, error)) *mockDoguClient_Watch_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguClient creates a new instance of mockDoguClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguClient(t mockConstructorTestingTnewMockDoguClient) *mockDoguClient {
	mock := &mockDoguClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}