package support

import (
	"context"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/cesapp-lib/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"testing"
)

func TestNewSupportArchiveService(t *testing.T) {
	t.Run("create new supportArchiveService", func(t *testing.T) {
		// given
		k8sClientMock := newMockK8sClient(t)
		discoveryClientMock := newMockDiscoveryInterface(t)

		// when
		sut := NewSupportArchiveService(k8sClientMock, discoveryClientMock)

		// then
		require.NotNil(t, sut)
	})
}

func TestSupportArchiveCreator_CreateSupportArchive(t *testing.T) {
	type fields struct {
		environmentFn       func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest_Common
		ArchiveManagerFn    func(t *testing.T) archiveManager
		resourceCollectorFn func(t *testing.T) resourceCollector
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no k8s resources found",
			fields: fields{
				environmentFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest_Common {
					return &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
						ExcludedContents: &pbMaintenance.ExcludedContents{
							SystemState: false,
						},
					}}
				},
				ArchiveManagerFn: func(t *testing.T) archiveManager {
					return *archive.NewManager()
				},
				resourceCollectorFn: func(t *testing.T) resourceCollector {
					rc := newMockResourceCollector(t)
					rc.EXPECT().Collect(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
					return rc
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "collect resources failed:", i)
			},
		},
		{
			name: "failed to add k8s resource to archive",
			fields: fields{
				environmentFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest_Common {
					return &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
						ExcludedContents: &pbMaintenance.ExcludedContents{
							SystemState: false,
						},
					}}
				},
				ArchiveManagerFn: func(t *testing.T) archiveManager {
					manager := newMockArchiveManager(t)
					manager.EXPECT().AddContentAsFile(mock.Anything, mock.Anything).Return(assert.AnError)
					manager.EXPECT().Close().Return(nil)
					return manager
				},
				resourceCollectorFn: func(t *testing.T) resourceCollector {
					rc := newMockResourceCollector(t)
					rc.EXPECT().Collect(mock.Anything, mock.Anything, mock.Anything).Return([]*unstructured.Unstructured{
						{
							make(map[string]interface{}),
						},
					}, nil)
					return rc
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to add resource to archive:", i)
			},
		},
		{
			name: "create support archive successfully",
			fields: fields{
				environmentFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest_Common {
					return &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
						ExcludedContents: &pbMaintenance.ExcludedContents{
							SystemState: true,
						},
					}}
				},
				ArchiveManagerFn: func(t *testing.T) archiveManager {
					return *archive.NewManager()
				},
				resourceCollectorFn: func(t *testing.T) resourceCollector {
					return newMockResourceCollector(t)
				},
			},
			wantErr: assert.NoError,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sac := &SupportArchiveCreator{
				environment:       tt.fields.environmentFn(t),
				ArchiveManager:    tt.fields.ArchiveManagerFn(t),
				resourceCollector: tt.fields.resourceCollectorFn(t),
			}
			got, err := sac.CreateSupportArchive(context.Background())
			if !tt.wantErr(t, err) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSupportArchive() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultSupportArchiveService_Create(t *testing.T) {
	type args struct {
		k8sClientFn       func(t *testing.T) k8sClient
		discoveryClientFn func(t *testing.T) discoveryInterface
		requestFn         func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest
		serverFn          func(t *testing.T) supportArchive_CreateServer
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "used unsupported environment",
			args: args{
				k8sClientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					return newMockDiscoveryInterface(t)
				},
				requestFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest {
					r := &pbMaintenance.CreateSupportArchiveRequest{
						Environment: &pbMaintenance.CreateSupportArchiveRequest_Legacy{},
					}
					return r
				},
				serverFn: func(t *testing.T) supportArchive_CreateServer {
					server := newMockSupportArchive_CreateServer(t)
					server.EXPECT().Context().Return(context.Background())
					return server
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "unsupported environment; the sent environment is probably legacy but must be common", i)
			},
		},
		{
			name: "used supported environment",
			args: args{
				k8sClientFn: func(t *testing.T) k8sClient {
					return newMockK8sClient(t)
				},
				discoveryClientFn: func(t *testing.T) discoveryInterface {
					discoveryClient := newMockDiscoveryInterface(t)
					discoveryClient.EXPECT().ServerPreferredResources().Return([]*metav1.APIResourceList{}, nil)
					return discoveryClient
				},
				requestFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest {
					r := &pbMaintenance.CreateSupportArchiveRequest{
						Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{
							Common: &pbMaintenance.CommonSupportArchiveRequest{
								ExcludedContents: &pbMaintenance.ExcludedContents{
									SystemState: false,
								},
							},
						},
					}
					return r
				},
				serverFn: func(t *testing.T) supportArchive_CreateServer {
					server := newMockSupportArchive_CreateServer(t)
					server.EXPECT().Context().Return(context.Background())
					return server
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSupportArchiveService(tt.args.k8sClientFn(t), tt.args.discoveryClientFn(t))
			err := d.Create(tt.args.requestFn(t), tt.args.serverFn(t))
			if !tt.wantErr(t, err) {
				return
			}
		})
	}
}
