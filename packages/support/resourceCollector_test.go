package support

import (
	"context"
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
		name    string
		fields  fields
		args    args
		want    []*unstructured.Unstructured
		wantErr bool
	}{
		{
			name: "should fail to get resource kind lists from server",
			fields: fields{
				clientFn:          nil,
				discoveryClientFn: nil,
			},
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &resourceCollector{
				client:          tt.fields.clientFn(t),
				discoveryClient: tt.fields.discoveryClientFn(t),
			}
			got, err := rc.Collect(tt.args.ctx, tt.args.labelSelector, tt.args.excludedGVKs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collect() got = %v, want %v", got, tt.want)
			}
		})
	}
}
