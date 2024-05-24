// Code generated by mockery v2.42.1. DO NOT EDIT.

package logging

import (
	context "context"

	appsv1 "k8s.io/api/apps/v1"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockDeploymentGetter is an autogenerated mock type for the deploymentGetter type
type mockDeploymentGetter struct {
	mock.Mock
}

type mockDeploymentGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDeploymentGetter) EXPECT() *mockDeploymentGetter_Expecter {
	return &mockDeploymentGetter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockDeploymentGetter) Get(ctx context.Context, name string, opts v1.GetOptions) (*appsv1.Deployment, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *appsv1.Deployment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) (*appsv1.Deployment, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) *appsv1.Deployment); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.Deployment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, v1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDeploymentGetter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDeploymentGetter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts v1.GetOptions
func (_e *mockDeploymentGetter_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockDeploymentGetter_Get_Call {
	return &mockDeploymentGetter_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockDeploymentGetter_Get_Call) Run(run func(ctx context.Context, name string, opts v1.GetOptions)) *mockDeploymentGetter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(v1.GetOptions))
	})
	return _c
}

func (_c *mockDeploymentGetter_Get_Call) Return(_a0 *appsv1.Deployment, _a1 error) *mockDeploymentGetter_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDeploymentGetter_Get_Call) RunAndReturn(run func(context.Context, string, v1.GetOptions) (*appsv1.Deployment, error)) *mockDeploymentGetter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDeploymentGetter creates a new instance of mockDeploymentGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDeploymentGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDeploymentGetter {
	mock := &mockDeploymentGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}