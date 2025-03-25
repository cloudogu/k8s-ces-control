package support

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestGVKMatcher_Matches(t *testing.T) {
	tests := []struct {
		name string
		m    gvkMatcher
		gvk  schema.GroupVersionKind
		want bool
	}{
		{
			name: "match group",
			m:    gvkMatcher{Group: "", Version: "*", Kind: "*"},
			gvk:  schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			want: true,
		},
		{
			name: "match group and version",
			m:    gvkMatcher{Group: "", Version: "v1", Kind: "*"},
			gvk:  schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			want: true,
		},
		{
			name: "match group and version and kind",
			m:    gvkMatcher{Group: "", Version: "v1", Kind: "Pod"},
			gvk:  schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
			want: true,
		},
		{
			name: "mismatch group",
			m:    gvkMatcher{Group: "", Version: "*", Kind: "*"},
			gvk:  schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"},
			want: false,
		},
		{
			name: "mismatch version",
			m:    gvkMatcher{Group: "", Version: "v1", Kind: "*"},
			gvk:  schema.GroupVersionKind{Group: "", Version: "v2", Kind: "Pod"},
			want: false,
		},
		{
			name: "mismatch kind",
			m:    gvkMatcher{Group: "", Version: "v1", Kind: "Pod"},
			gvk:  schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Matches(tt.gvk); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceCollector_Collect(t *testing.T) {
	type fields struct {
		clientFn          func(t *testing.T) k8sClient
		discoveryClientFn func(t *testing.T) discoveryInterface
	}
	type args struct {
		ctx           context.Context
		labelSelector *metav1.LabelSelector
		excludedGVKs  []gvkMatcher
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*unstructured.Unstructured
		wantErrFn assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get resource kind lists from server",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(nil, assert.AnError)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want: nil,
			wantErrFn: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get resource kind lists from server", i)
			},
		},
		{
			name: "should fail to create label selector",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return([]*metav1.APIResourceList{}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"_invalid": "_invalid"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want: nil,
			wantErrFn: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to create selector from given label selector", i)
			},
		},
		{
			name: "should skip empty api resources list",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return([]*metav1.APIResourceList{{}}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want:      nil,
			wantErrFn: assert.NoError,
		},
		{
			name: "should fail to parse group version",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(
						[]*metav1.APIResourceList{
							{},
							{APIResources: []metav1.APIResource{{}}, GroupVersion: "invalid/invalid/invalid"},
						}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want: nil,
			wantErrFn: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to list api resources with group version \"invalid/invalid/invalid\"", i)
			},
		},
		{
			name: "should skip resources with no verbs or without list",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(
						[]*metav1.APIResourceList{
							{},
							{APIResources: []metav1.APIResource{{Verbs: make(metav1.Verbs, 0)}}},
							{APIResources: []metav1.APIResource{{Verbs: []string{"get"}}}},
						}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want:      nil,
			wantErrFn: assert.NoError,
		},
		{
			name: "should skip resource if gvk matcher matches",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(
						[]*metav1.APIResourceList{
							{},
							{APIResources: []metav1.APIResource{{Verbs: make(metav1.Verbs, 0)}}},
							{APIResources: []metav1.APIResource{{Verbs: []string{"get"}}}},
							{GroupVersion: "v1", APIResources: []metav1.APIResource{{Verbs: []string{"get", "list"}, Version: "v1", Kind: "Secret"}}},
						}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want:      nil,
			wantErrFn: assert.NoError,
		},
		{
			name: "should fail to list resource",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					m := newMockK8sClient(t)
					m.EXPECT().List(
						context.Background(),
						&unstructured.UnstructuredList{Object: map[string]interface{}{"apiVersion": "v1", "kind": "Pod"}},
						mock.Anything,
					).Return(assert.AnError)
					return m
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(
						[]*metav1.APIResourceList{
							{},
							{APIResources: []metav1.APIResource{{Verbs: make(metav1.Verbs, 0)}}},
							{APIResources: []metav1.APIResource{{Verbs: []string{"get"}}}},
							{GroupVersion: "v1", APIResources: []metav1.APIResource{{Verbs: []string{"get", "list"}, Version: "v1", Kind: "Secret"}}},
							{GroupVersion: "v1", APIResources: []metav1.APIResource{{Verbs: []string{"get", "list"}, Version: "v1", Kind: "Pod"}}},
						}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want: nil,
			wantErrFn: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to list objects in /v1, Kind=Pod", i) &&
					assert.ErrorContains(t, err, "failed to list api resources with label selector \"app=ces\"", i)
			},
		},
		{
			name: "should succeed to list resource",
			fields: fields{
				clientFn: func(t *testing.T) k8sClient {
					m := newMockK8sClient(t)
					m.EXPECT().List(
						context.Background(),
						&unstructured.UnstructuredList{Object: map[string]interface{}{"apiVersion": "v1", "kind": "Pod"}},
						mock.Anything,
					).RunAndReturn(func(ctx context.Context, list client.ObjectList, option ...client.ListOption) error {
						list.(*unstructured.UnstructuredList).Items = []unstructured.Unstructured{
							{Object: map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{"name": "test"}}},
						}
						return nil
					})
					return m
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					m := newMockDiscoveryInterface(t)
					m.EXPECT().ServerPreferredResources().Return(
						[]*metav1.APIResourceList{
							{},
							{APIResources: []metav1.APIResource{{Verbs: make(metav1.Verbs, 0)}}},
							{APIResources: []metav1.APIResource{{Verbs: []string{"get"}}}},
							{GroupVersion: "v1", APIResources: []metav1.APIResource{{Verbs: []string{"get", "list"}, Version: "v1", Kind: "Secret"}}},
							{GroupVersion: "v1", APIResources: []metav1.APIResource{{Verbs: []string{"get", "list"}, Version: "v1", Kind: "Pod"}}},
						}, nil)
					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				labelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}},
				excludedGVKs: []gvkMatcher{{
					Version: "v1",
					Kind:    "Secret",
				}},
			},
			want:      []*unstructured.Unstructured{{Object: map[string]interface{}{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]interface{}{"name": "test"}}}},
			wantErrFn: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &resourceCollector{
				client:          tt.fields.clientFn(t),
				discoveryClient: tt.fields.discoveryClientFn(t),
			}
			got, err := rc.Collect(tt.args.ctx, tt.args.labelSelector, tt.args.excludedGVKs)
			_ = tt.wantErrFn(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collect() got = %v, want %v", got, tt.want)
			}
		})
	}
}
