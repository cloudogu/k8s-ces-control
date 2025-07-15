package supportArchive

import (
	"fmt"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	pbTypes "github.com/cloudogu/ces-control-api/generated/types"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
	"testing"
	"time"
)

func TestNewSupportArchiveService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//given
		supportArchiveClientMock := newMockSupportArchiveClient(t)

		//when
		service := NewSupportArchiveService(supportArchiveClientMock)

		require.NotNil(t, service)
	})
}

func Test_defaultSupportArchive_mapRequestSettingsToSupportArchive(t *testing.T) {
	tests := []struct {
		name    string
		reqFn   func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest
		wantFn  func(t *testing.T) *v1.SupportArchive
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successfully map settings",
			reqFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest {

				return &pbMaintenance.CreateSupportArchiveRequest{
					Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
						ExcludedContents: &pbMaintenance.ExcludedContents{
							SystemState:   true,
							SensitiveData: true,
							LogsAndEvents: true,
							VolumeInfo:    true,
							SystemInfo:    true,
						},
						LoggingConfig: &pbMaintenance.LoggingConfig{
							EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
							StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
						},
					}},
				}
			},
			wantFn: func(t *testing.T) *v1.SupportArchive {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				return &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{
							SystemState:   true,
							SensitiveData: true,
							LogsAndEvents: true,
							VolumeInfo:    true,
							SystemInfo:    true,
						},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "failed to map settings because end time is before start time",
			reqFn: func(t *testing.T) *pbMaintenance.CreateSupportArchiveRequest {

				return &pbMaintenance.CreateSupportArchiveRequest{
					Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
						ExcludedContents: &pbMaintenance.ExcludedContents{
							SystemState:   true,
							SensitiveData: true,
							LogsAndEvents: true,
							VolumeInfo:    true,
							SystemInfo:    true,
						},
						LoggingConfig: &pbMaintenance.LoggingConfig{
							EndDateTime:   &timestamppb.Timestamp{Seconds: int64(1600)},
							StartDateTime: &timestamppb.Timestamp{Seconds: int64(32000)},
						},
					}},
				}
			},
			wantFn: func(t *testing.T) *v1.SupportArchive {
				return &v1.SupportArchive{}
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "end time is before start time", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supportArchiveClientMock := newMockSupportArchiveClient(t)
			d := NewSupportArchiveService(supportArchiveClientMock)

			got, err := d.mapRequestSettingsToSupportArchive(tt.reqFn(t))

			if !tt.wantErr(t, err, fmt.Sprintf("mapRequestSettingsToSupportArchive(%v)", tt.reqFn(t))) {
				return
			}
			assert.Equal(t, tt.wantFn(t), got)
		})
	}
}

func Test_defaultSupportArchive_Create(t *testing.T) {
	tests := []struct {
		name                   string
		supportArchiveClientFn func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver)
		req                    *pbMaintenance.CreateSupportArchiveRequest
		wantErrMessage         string
	}{
		{
			name: "should fail to create support archive CR",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, assert.AnError)
				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				return clientMock, serviceMock
			},
			wantErrMessage: "failed to create support archive: ",
		},
		{
			name: "should fail to create watch interface",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, nil)

				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(nil, assert.AnError)
				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				return clientMock, serviceMock
			},
			wantErrMessage: "failed to create watch interface:",
		},
		{
			name: "should fail when watch stops without result",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Stop()
				}()

				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				return clientMock, serviceMock
			},
			wantErrMessage: "failed to create or watch support archive:",
		},
		{
			name: "should fail when watch returns nil",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, nil)
				}()

				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				return clientMock, serviceMock
			},
			wantErrMessage: "failed to create or watch support archive:",
		},
		{
			name: "should fail when download path cannot be sent",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				downloadPath := "testDownloadPath"
				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: downloadPath,
						},
					})
				}()

				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				resp := &pbTypes.ChunkedDataResponse{}
				resp.Data = []byte(downloadPath)
				serviceMock.EXPECT().Send(resp).Return(assert.AnError)

				return clientMock, serviceMock
			},
			wantErrMessage: "failed to send response:",
		},
		{
			name: "should succeed",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) (supportArchiveClient, supportArchiveCreateserver) {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
					Spec: v1.SupportArchiveSpec{
						ExcludedContents: v1.ExcludedContents{},
						LoggingConfig: v1.LoggingConfig{
							StartTime: metav1.NewTime(timestampStart.AsTime()),
							EndTime:   metav1.NewTime(timestampEnd.AsTime()),
						},
					},
				}, metav1.CreateOptions{}).Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				downloadPath := "testDownloadPath"
				go func() {
					time.Sleep(1 * time.Second)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: downloadPath,
						},
					})
				}()

				serviceMock := newMockSupportArchiveCreateserver(t)
				serviceMock.EXPECT().Context().Return(t.Context())
				resp := &pbTypes.ChunkedDataResponse{}
				resp.Data = []byte(downloadPath)
				serviceMock.EXPECT().Send(resp).Return(nil)

				return clientMock, serviceMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			clientmock, servicMock := tt.supportArchiveClientFn(t)
			service := NewSupportArchiveService(clientmock)

			// when
			err := service.Create(tt.req, servicMock)

			//then
			if tt.wantErrMessage != "" {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.wantErrMessage))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
