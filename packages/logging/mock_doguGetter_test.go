// Code generated by mockery v2.42.1. DO NOT EDIT.

package logging

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
)

// mockDoguGetter is an autogenerated mock type for the doguGetter type
type mockDoguGetter struct {
	mock.Mock
}

type mockDoguGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguGetter) EXPECT() *mockDoguGetter_Expecter {
	return &mockDoguGetter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockDoguGetter) Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.Dogu, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v2.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) (*v2.Dogu, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) *v2.Dogu); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, v1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguGetter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDoguGetter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts v1.GetOptions
func (_e *mockDoguGetter_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockDoguGetter_Get_Call {
	return &mockDoguGetter_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockDoguGetter_Get_Call) Run(run func(ctx context.Context, name string, opts v1.GetOptions)) *mockDoguGetter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(v1.GetOptions))
	})
	return _c
}

func (_c *mockDoguGetter_Get_Call) Return(_a0 *v2.Dogu, _a1 error) *mockDoguGetter_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguGetter_Get_Call) RunAndReturn(run func(context.Context, string, v1.GetOptions) (*v2.Dogu, error)) *mockDoguGetter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguGetter creates a new instance of mockDoguGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguGetter {
	mock := &mockDoguGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
