package supportArchive

import (
	"context"
	"fmt"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	pbTypes "github.com/cloudogu/ces-control-api/generated/types"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewSupportArchiveService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//given
		supportArchiveClientMock := newMockSupportArchiveClient(t)

		//when
		service := NewSupportArchiveService(supportArchiveClientMock, &http.Client{})

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
			wantErr: nil,
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
			d := NewSupportArchiveService(supportArchiveClientMock, &http.Client{})

			beforeTime := metav1.Now()
			got, err := d.mapRequestSettingsToSupportArchive(tt.reqFn(t))
			afterTime := metav1.Now()

			if tt.wantErr != nil {
				tt.wantErr(t, err, fmt.Sprintf("mapRequestSettingsToSupportArchive(%v)", tt.reqFn(t)))
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantFn(t).Spec, got.Spec)
			assertArchiveName(t, got, beforeTime, afterTime)
		})
	}

	t.Run("successfully take endtime of now, if not set", func(t *testing.T) {
		request := &pbMaintenance.CreateSupportArchiveRequest{
			Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
				ExcludedContents: &pbMaintenance.ExcludedContents{},
				LoggingConfig: &pbMaintenance.LoggingConfig{
					StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
				},
			}},
		}

		supportArchiveClientMock := newMockSupportArchiveClient(t)
		d := NewSupportArchiveService(supportArchiveClientMock, &http.Client{})

		beforeTime := metav1.Now()
		got, err := d.mapRequestSettingsToSupportArchive(request)

		assert.NoError(t, err)
		assert.True(t, got.Spec.LoggingConfig.EndTime.Time.Equal(beforeTime.Time) ||
			got.Spec.LoggingConfig.EndTime.Time.After(beforeTime.Time))
	})
}

func assertArchiveName(t *testing.T, got *v1.SupportArchive, beforeTime metav1.Time, afterTime metav1.Time) {
	nameParts := strings.Split(got.Name, "-")
	assert.Equal(t, len(nameParts), 3)
	timestampStr, _ := strings.CutSuffix(nameParts[2], "z")
	// Parse the timestamp string
	archiveTime, err := time.Parse("20060102150405", timestampStr)
	assert.NoError(t, err)

	// Assert that the archive name timestamp is between before and after times
	assert.True(t, archiveTime.Before(beforeTime.Time) || archiveTime.After(afterTime.Time))
}

func Test_defaultSupportArchive_Create(t *testing.T) {
	tests := []struct {
		name                   string
		supportArchiveClientFn func(t *testing.T) supportArchiveClient
		supportArchiveServerFn func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc)
		supportArchiveHttpFn   func(t *testing.T) httpClient
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
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, assert.AnError)
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(t.Context())
				return serverMock, nil
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "failed to create support archive: ",
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
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, nil)
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(nil, assert.AnError)
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(t.Context())
				return serverMock, nil
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "failed to create watch interface:",
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
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Stop()
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(t.Context())
				return serverMock, nil
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "failed to create or watch support archive:",
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
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, nil)
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(t.Context())
				return serverMock, nil
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "failed to create or watch support archive:",
		},
		{
			name: "should fail when download path is not a url",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						return nil, nil
					})
				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "§$(§/$$=§%)(",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "failed to create request: ",
		},
		{
			name: "should fail when http client throws error",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						return nil, nil
					})
				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "testDownloadPath",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient {
				httpMock := newMockHttpClient(t)
				httpMock.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)
				return httpMock
			},
			wantErrMessage: "failed to send request: ",
		},
		{
			name: "should fail when http status code is not ok",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						return nil, nil
					})
				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "testDownloadPath",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient {
				httpMock := newMockHttpClient(t)
				httpMock.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(
					&http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader("myTest.zip")),
					}, nil)
				return httpMock
			},
			wantErrMessage: "bad status: ",
		},
		{
			name: "should fail when download archive cannot be sent",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						return nil, nil
					})
				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "testDownloadPath",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				serverMock := newMockSupportArchiveCreateserver(t)
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				resp := &pbTypes.ChunkedDataResponse{}
				resp.Data = []byte("myTest.zip")
				serverMock.EXPECT().Send(resp).Return(assert.AnError)

				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient {
				httpMock := newMockHttpClient(t)
				httpMock.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("myTest.zip")),
					}, nil)
				return httpMock
			},
			wantErrMessage: "failed to send response:",
		},
		{
			name: "should fail on context timeout",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig:    &pbMaintenance.LoggingConfig{},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, nil)

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(timeoutCtx)

				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient { return nil },
			wantErrMessage:       "timed out waiting for support archive to be created",
		},
		{
			name: "should ignore archives with different name",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				Environment: &pbMaintenance.CreateSupportArchiveRequest_Common{Common: &pbMaintenance.CommonSupportArchiveRequest{
					ExcludedContents: &pbMaintenance.ExcludedContents{},
					LoggingConfig: &pbMaintenance.LoggingConfig{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
				}},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						assert.Equal(t, archive.Spec.ExcludedContents, v1.ExcludedContents{})
						assert.Equal(t, archive.Spec.LoggingConfig.StartTime, metav1.NewTime(timestampStart.AsTime()))
						assert.Equal(t, archive.Spec.LoggingConfig.EndTime, metav1.NewTime(timestampEnd.AsTime()))
						return nil, nil
					})

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(100 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: "wrongName",
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "differentPath", // this would let the Send-Mock fail
						},
					})
					time.Sleep(500 * time.Millisecond)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "testDownloadPath",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				resp := &pbTypes.ChunkedDataResponse{}
				resp.Data = []byte("myTest.zip")
				serverMock.EXPECT().Send(resp).Return(nil)

				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient {
				httpMock := newMockHttpClient(t)
				httpMock.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
					assert.Equal(t, request.Method, "GET")
					assert.Equal(t, request.URL.Path, "testDownloadPath")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("myTest.zip")),
					}, nil
				})
				return httpMock
			},
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
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				var archiveName string
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						archiveName = archive.Name
						assert.Equal(t, archive.Spec.ExcludedContents, v1.ExcludedContents{})
						assert.Equal(t, archive.Spec.LoggingConfig.StartTime, metav1.NewTime(timestampStart.AsTime()))
						assert.Equal(t, archive.Spec.LoggingConfig.EndTime, metav1.NewTime(timestampEnd.AsTime()))
						return nil, nil
					})

				watcher := watch.NewFake()
				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(watcher, nil)

				go func() {
					time.Sleep(1 * time.Second)
					watcher.Action(watch.Modified, &v1.SupportArchive{
						ObjectMeta: metav1.ObjectMeta{
							Name: archiveName,
						},
						Status: v1.SupportArchiveStatus{
							Phase:        v1.StatusPhaseCreated,
							DownloadPath: "testDownloadPath",
						},
					})
				}()
				return clientMock
			},
			supportArchiveServerFn: func(t *testing.T) (supportArchiveCreateserver, context.CancelFunc) {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				serverMock := newMockSupportArchiveCreateserver(t)
				serverMock.EXPECT().Context().Return(timeoutCtx)
				resp := &pbTypes.ChunkedDataResponse{}
				resp.Data = []byte(("myTest.zip"))
				serverMock.EXPECT().Send(resp).Return(nil)

				return serverMock, cancel
			},
			supportArchiveHttpFn: func(t *testing.T) httpClient {
				httpMock := newMockHttpClient(t)
				httpMock.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
					assert.Equal(t, request.Method, "GET")
					assert.Equal(t, request.URL.Path, "testDownloadPath")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("myTest.zip")),
					}, nil
				})
				return httpMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			clientmock := tt.supportArchiveClientFn(t)
			serverMock, cancelWatch := tt.supportArchiveServerFn(t)
			if cancelWatch != nil {
				defer cancelWatch()
			}
			httpmock := tt.supportArchiveHttpFn(t)
			service := NewSupportArchiveService(clientmock, httpmock)

			// when
			err := service.Create(tt.req, serverMock)

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
