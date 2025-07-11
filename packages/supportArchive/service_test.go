package supportArchive

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	pbTypes "github.com/cloudogu/ces-control-api/generated/types"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net"
	"strings"
	"testing"
)

var (
	listener *bufconn.Listener
)

const bufferSize1MB = 1024 * 1024

func setupSupportArchiveTestServer(
	writer stream.Writer,
	client supportArchiveClient) {
	s := grpc.NewServer()
	listener = bufconn.Listen(bufferSize1MB)
	srv := &supportArchiveService{supportArchiveClient: client, writeToStream: writer}
	pbMaintenance.RegisterSupportArchiveServer(s, srv)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}

type clientStreamReader interface {
	Recv() (*pbTypes.ChunkedDataResponse, error)
}

func getDataFromStream(t *testing.T, client clientStreamReader) ([]byte, error) {
	t.Helper()
	data, err := readStreamData(t, client)
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)
	file := reader.File[0]
	f, err := file.Open()
	require.NoError(t, err)
	defer func() { _ = f.Close() }()
	unzippedData, err := ioutil.ReadAll(f)
	require.NoError(t, err)
	return unzippedData, nil
}

func readStreamData(t *testing.T, client clientStreamReader) ([]byte, error) {
	t.Helper()
	var binaryArchiveData []byte
	for {
		response, err := client.Recv()
		if err != nil {
			if err == io.EOF {
				t.Logf("Transfer of %d bytes successful", len(binaryArchiveData))
				break
			}

			return nil, err
		}

		binaryArchiveData = append(binaryArchiveData, response.Data...)
	}
	return binaryArchiveData, nil
}

func getSupportArchiveGrpcClient(t *testing.T) (pbMaintenance.SupportArchiveClient, *context.Context, *grpc.ClientConn) {
	ctx := context.Background()
	conn, clientErr := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if clientErr != nil {
		t.Fatalf("Failed to dial bufnet: %v", clientErr)
	}
	return pbMaintenance.NewSupportArchiveClient(conn), &ctx, conn
}

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
		supportArchiveClientFn func(t *testing.T) supportArchiveClient
		req                    *pbMaintenance.CreateSupportArchiveRequest
		wantErrMessage         string
	}{
		{
			name: "should fail to create support archive CR",
			req: &pbMaintenance.CreateSupportArchiveRequest{
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
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
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
				}, metav1.CreateOptions{}).Return(nil, assert.AnError)
				return clientMock
			},
			wantErrMessage: "failed to create support archive: ",
		},
		{
			name: "should fail to create watch interface",
			req: &pbMaintenance.CreateSupportArchiveRequest{
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
			},
			supportArchiveClientFn: func(t *testing.T) supportArchiveClient {
				timestampStart := &timestamppb.Timestamp{Seconds: int64(16000)}
				timestampEnd := &timestamppb.Timestamp{Seconds: int64(32000)}
				clientMock := newMockSupportArchiveClient(t)
				clientMock.EXPECT().Create(mock.Anything, &v1.SupportArchive{
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
				}, metav1.CreateOptions{}).Return(nil, nil)

				clientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{}).Return(nil, assert.AnError)
				return clientMock
			},
			wantErrMessage: "failed to create watch interface:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := func(data []byte, streamer stream.GRPCStreamServer) error {
				return streamer.Send(&pbTypes.ChunkedDataResponse{Data: data})
			}
			setupSupportArchiveTestServer(writer, tt.supportArchiveClientFn(t))

			client, ctx, conn := getSupportArchiveGrpcClient(t)
			defer func() { _ = conn.Close() }()

			archiveDataClient, err := client.Create(*ctx, tt.req)

			data, err := getDataFromStream(t, archiveDataClient)
			if tt.wantErrMessage != "" {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.wantErrMessage))
			} else {
				require.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}
