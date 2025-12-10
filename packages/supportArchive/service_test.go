package supportArchive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
					ExcludedContents: &pbMaintenance.ExcludedContents{
						SystemState:   true,
						SensitiveData: true,
						Logs:          true,
						Events:        true,
						VolumeInfo:    true,
						SystemInfo:    true,
					},
					ContentTimeframe: &pbMaintenance.ContentTimeframe{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
					},
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
							Logs:          true,
							Events:        true,
							VolumeInfo:    true,
							SystemInfo:    true,
						},
						ContentTimeframe: v1.ContentTimeframe{
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
					ExcludedContents: &pbMaintenance.ExcludedContents{
						SystemState:   true,
						SensitiveData: true,
						Logs:          true,
						Events:        true,
						VolumeInfo:    true,
						SystemInfo:    true,
					},
					ContentTimeframe: &pbMaintenance.ContentTimeframe{
						EndDateTime:   &timestamppb.Timestamp{Seconds: int64(1600)},
						StartDateTime: &timestamppb.Timestamp{Seconds: int64(32000)},
					},
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
			ExcludedContents: &pbMaintenance.ExcludedContents{},
			ContentTimeframe: &pbMaintenance.ContentTimeframe{
				StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
			},
		}

		supportArchiveClientMock := newMockSupportArchiveClient(t)
		d := NewSupportArchiveService(supportArchiveClientMock, &http.Client{})

		beforeTime := metav1.Now()
		got, err := d.mapRequestSettingsToSupportArchive(request)

		assert.NoError(t, err)
		assert.True(t, got.Spec.ContentTimeframe.EndTime.Time.Equal(beforeTime.Time) ||
			got.Spec.ContentTimeframe.EndTime.Time.After(beforeTime.Time))
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
	testCtx := context.Background()

	tests := []struct {
		name                   string
		supportArchiveClientFn func(t *testing.T) supportArchiveClient
		req                    *pbMaintenance.CreateSupportArchiveRequest
		wantErrMessage         string
	}{
		{
			name: "should fail to map request settings",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				ExcludedContents: &pbMaintenance.ExcludedContents{},
				ContentTimeframe: &pbMaintenance.ContentTimeframe{
					EndDateTime:   &timestamppb.Timestamp{Seconds: int64(16000)},
					StartDateTime: &timestamppb.Timestamp{Seconds: int64(32000)},
				},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				return clientMock
			},
			wantErrMessage: "failed to map support archive settings: ",
		},
		{
			name: "should fail to create support archive CR",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				ExcludedContents: &pbMaintenance.ExcludedContents{},
				ContentTimeframe: &pbMaintenance.ContentTimeframe{
					EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
					StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
				},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					Return(nil, assert.AnError)
				return clientMock
			},
			wantErrMessage: "failed to create support archive: ",
		},
		{
			name: "should succeed",
			req: &pbMaintenance.CreateSupportArchiveRequest{
				ExcludedContents: &pbMaintenance.ExcludedContents{},
				ContentTimeframe: &pbMaintenance.ContentTimeframe{
					EndDateTime:   &timestamppb.Timestamp{Seconds: int64(32000)},
					StartDateTime: &timestamppb.Timestamp{Seconds: int64(16000)},
				},
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, mock.AnythingOfType("*v1.SupportArchive"), metav1.CreateOptions{}).
					RunAndReturn(func(ctx context.Context, archive *v1.SupportArchive, options metav1.CreateOptions) (*v1.SupportArchive, error) {
						assert.Equal(t, archive.Spec.ExcludedContents, v1.ExcludedContents{})
						assert.Equal(t, archive.Spec.ContentTimeframe.StartTime, metav1.NewTime(timestampStart.AsTime()))
						assert.Equal(t, archive.Spec.ContentTimeframe.EndTime, metav1.NewTime(timestampEnd.AsTime()))
						return nil, nil
					})

				return clientMock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			clientmock := tt.supportArchiveClientFn(t)
			service := &supportArchiveService{
				supportArchiveClient: clientmock,
			}

			// when
			resp, err := service.Create(testCtx, tt.req)

			//then
			if tt.wantErrMessage != "" {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.wantErrMessage))
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func Test_supportArchiveService_AllSupportArchives(t *testing.T) {
	testCtx := context.Background()
	t.Run("should return all support archives", func(t *testing.T) {
		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&v1.SupportArchiveList{
			Items: []v1.SupportArchive{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "archive-1",
						CreationTimestamp: metav1.Time{Time: time.Unix(16000, 0)},
					},
					Status: v1.SupportArchiveStatus{
						Errors:       []string{},
						DownloadPath: "",
						Conditions:   []metav1.Condition{},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "archive-2",
						CreationTimestamp: metav1.Time{Time: time.Unix(24000, 0)},
					},
					Status: v1.SupportArchiveStatus{
						Errors:       []string{},
						DownloadPath: "",
						Conditions:   []metav1.Condition{},
					},
				},
			},
		}, nil)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
		}

		archives, err := sut.AllSupportArchives(testCtx, nil)

		require.NoError(t, err)
		assert.Len(t, archives.SupportArchives, 2)

		assert.Equal(t, "archive-1", archives.SupportArchives[0].Name)
		assert.Equal(t, timestamppb.New(time.Unix(16000, 0)), archives.SupportArchives[0].CreatedDateTime)
		assert.Equal(t, pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_IN_PROGRESS, archives.SupportArchives[0].Status)

		assert.Equal(t, "archive-2", archives.SupportArchives[1].Name)
		assert.Equal(t, timestamppb.New(time.Unix(24000, 0)), archives.SupportArchives[1].CreatedDateTime)
		assert.Equal(t, pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_IN_PROGRESS, archives.SupportArchives[1].Status)
	})

	t.Run("should return error for error listing archives", func(t *testing.T) {
		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
		}

		_, err := sut.AllSupportArchives(testCtx, nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to list support archives:")
	})
}

func Test_supportArchiveService_DeleteSupportArchive(t *testing.T) {
	testCtx := context.Background()
	t.Run("should delete archive", func(t *testing.T) {
		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Delete(testCtx, "archive-1", metav1.DeleteOptions{}).Return(nil)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
		}

		resp, err := sut.DeleteSupportArchive(testCtx, &pbMaintenance.DeleteSupportArchiveRequest{Name: "archive-1"})

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("should fail delete archive for error deleting", func(t *testing.T) {
		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Delete(testCtx, "archive-1", metav1.DeleteOptions{}).Return(assert.AnError)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
		}

		_, err := sut.DeleteSupportArchive(testCtx, &pbMaintenance.DeleteSupportArchiveRequest{Name: "archive-1"})

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to delete support archive:")
	})
}

func Test_supportArchiveService_DownloadSupportArchive(t *testing.T) {
	testCtx := context.Background()
	t.Run("should download archive", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, request.Method, "GET")
			assert.Equal(t, request.URL.Path, "/download/archive-1")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("myTest.zip")),
			}, nil
		})

		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Get(testCtx, "archive-1", metav1.GetOptions{}).Return(&v1.SupportArchive{
			Status: v1.SupportArchiveStatus{
				DownloadPath: "http://localhost:8080/download/archive-1",
			},
		}, nil)

		mServer := newMockSupportArchiveDownloadServer(t)
		mServer.EXPECT().Send(mock.Anything).Return(nil)
		mServer.EXPECT().Context().Return(testCtx)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
			httpClient:           mHttpClient,
			writeToStream:        nil,
		}

		err := sut.DownloadSupportArchive(&pbMaintenance.DownloadSupportArchiveRequest{Name: "archive-1"}, mServer)

		require.NoError(t, err)
	})

	t.Run("should fail to download archive for error streaming", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, request.Method, "GET")
			assert.Equal(t, request.URL.Path, "/download/archive-1")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("myTest.zip")),
			}, nil
		})

		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Get(testCtx, "archive-1", metav1.GetOptions{}).Return(&v1.SupportArchive{
			Status: v1.SupportArchiveStatus{
				DownloadPath: "http://localhost:8080/download/archive-1",
			},
		}, nil)

		mServer := newMockSupportArchiveDownloadServer(t)
		mServer.EXPECT().Send(mock.Anything).Return(assert.AnError)
		mServer.EXPECT().Context().Return(testCtx)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
			httpClient:           mHttpClient,
			writeToStream:        nil,
		}

		err := sut.DownloadSupportArchive(&pbMaintenance.DownloadSupportArchiveRequest{Name: "archive-1"}, mServer)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed stream support-archive file:")
	})

	t.Run("should fail to download archive for error getting file", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, request.Method, "GET")
			assert.Equal(t, request.URL.Path, "/download/archive-1")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("myTest.zip")),
			}, assert.AnError
		})

		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Get(testCtx, "archive-1", metav1.GetOptions{}).Return(&v1.SupportArchive{
			Status: v1.SupportArchiveStatus{
				DownloadPath: "http://localhost:8080/download/archive-1",
			},
		}, nil)

		mServer := newMockSupportArchiveDownloadServer(t)
		mServer.EXPECT().Context().Return(testCtx)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
			httpClient:           mHttpClient,
			writeToStream:        nil,
		}

		err := sut.DownloadSupportArchive(&pbMaintenance.DownloadSupportArchiveRequest{Name: "archive-1"}, mServer)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to download ZIP file:")
	})

	t.Run("should fail to download archive when archive has now downloadLink", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)

		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Get(testCtx, "archive-1", metav1.GetOptions{}).Return(&v1.SupportArchive{
			Status: v1.SupportArchiveStatus{},
		}, nil)

		mServer := newMockSupportArchiveDownloadServer(t)
		mServer.EXPECT().Context().Return(testCtx)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
			httpClient:           mHttpClient,
			writeToStream:        nil,
		}

		err := sut.DownloadSupportArchive(&pbMaintenance.DownloadSupportArchiveRequest{Name: "archive-1"}, mServer)

		require.Error(t, err)
		assert.ErrorContains(t, err, "support archive is not ready yet")
	})

	t.Run("should fail to download archive for error getting archive", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mSupportArchiveClient := newMockSupportArchiveClient(t)
		mSupportArchiveClient.EXPECT().Get(testCtx, "archive-1", metav1.GetOptions{}).Return(nil, assert.AnError)

		mServer := newMockSupportArchiveDownloadServer(t)
		mServer.EXPECT().Context().Return(testCtx)

		sut := &supportArchiveService{
			supportArchiveClient: mSupportArchiveClient,
			httpClient:           mHttpClient,
			writeToStream:        nil,
		}

		err := sut.DownloadSupportArchive(&pbMaintenance.DownloadSupportArchiveRequest{Name: "archive-1"}, mServer)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get support archive:")
	})
}

func Test_supportArchiveService_downloadFile(t *testing.T) {
	t.Run("should successfully download file", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).RunAndReturn(func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "http://example.com/file.zip", request.URL.String())
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("file content")),
			}, nil
		})

		sut := &supportArchiveService{
			httpClient: mHttpClient,
		}

		reader, err := sut.getDownloadFile("http://example.com/file.zip")

		require.NoError(t, err)
		require.NotNil(t, reader)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "file content", string(content))
	})

	t.Run("should fail when request creation fails", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)

		sut := &supportArchiveService{
			httpClient: mHttpClient,
		}

		reader, err := sut.getDownloadFile(":")

		require.Error(t, err)
		assert.Nil(t, reader)
		assert.ErrorContains(t, err, "failed to create request:")
	})

	t.Run("should fail when HTTP request fails", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)

		sut := &supportArchiveService{
			httpClient: mHttpClient,
		}

		reader, err := sut.getDownloadFile("http://example.com/file.zip")

		require.Error(t, err)
		assert.Nil(t, reader)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to send request:")
	})

	t.Run("should fail when status code is not OK", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		sut := &supportArchiveService{
			httpClient: mHttpClient,
		}

		reader, err := sut.getDownloadFile("http://example.com/file.zip")

		require.Error(t, err)
		assert.Nil(t, reader)
		assert.ErrorContains(t, err, "bad status: 404 Not Found")
	})

	t.Run("should fail when status code is 500", func(t *testing.T) {
		mHttpClient := newMockHttpClient(t)
		mHttpClient.EXPECT().Do(mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Status:     "500 Internal Server Error",
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil)

		sut := &supportArchiveService{
			httpClient: mHttpClient,
		}

		reader, err := sut.getDownloadFile("http://example.com/file.zip")

		require.Error(t, err)
		assert.Nil(t, reader)
		assert.ErrorContains(t, err, "bad status: 500 Internal Server Error")
	})
}

func Test_getStatus(t *testing.T) {
	tests := []struct {
		name    string
		archive v1.SupportArchive
		want    pbMaintenance.SupportArchiveStatus
	}{
		{
			name: "should return FAILED when archive has errors",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors: []string{"error 1", "error 2"},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_FAILED,
		},
		{
			name: "should return COMPLETED when archive is created with download path",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors:       []string{},
					DownloadPath: "http://localhost:8080/download/archive-1",
					Conditions: []metav1.Condition{
						{
							Type:   v1.ConditionSupportArchiveCreated,
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_COMPLETED,
		},
		{
			name: "should return CREATED when archive is created without download path",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors:       []string{},
					DownloadPath: "",
					Conditions: []metav1.Condition{
						{
							Type:   v1.ConditionSupportArchiveCreated,
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_CREATED,
		},
		{
			name: "should return IN_PROGRESS when archive has no conditions",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors:     []string{},
					Conditions: []metav1.Condition{},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_IN_PROGRESS,
		},
		{
			name: "should return IN_PROGRESS when archive condition is not true",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors: []string{},
					Conditions: []metav1.Condition{
						{
							Type:   v1.ConditionSupportArchiveCreated,
							Status: metav1.ConditionFalse,
						},
					},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_IN_PROGRESS,
		},
		{
			name: "should return FAILED when archive has errors even with completed condition",
			archive: v1.SupportArchive{
				Status: v1.SupportArchiveStatus{
					Errors:       []string{"error"},
					DownloadPath: "http://localhost:8080/download/archive-1",
					Conditions: []metav1.Condition{
						{
							Type:   v1.ConditionSupportArchiveCreated,
							Status: metav1.ConditionTrue,
						},
					},
				},
			},
			want: pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_FAILED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatus(tt.archive)
			assert.Equal(t, tt.want, got)
		})
	}
}
