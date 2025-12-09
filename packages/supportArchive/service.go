package supportArchive

import (
	"context"
	"fmt"
	"io"
	"net/http"

	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type supportArchiveService struct {
	pbMaintenance.UnimplementedSupportArchiveServer
	supportArchiveClient supportArchiveClient
	httpClient           httpClient
	writeToStream        stream.Writer
}

func NewSupportArchiveService(client supportArchiveClient, http httpClient) *supportArchiveService {
	return &supportArchiveService{
		supportArchiveClient: client,
		httpClient:           http,
		writeToStream:        stream.WriteToStream,
	}
}

func (d *supportArchiveService) Create(ctx context.Context, req *pbMaintenance.CreateSupportArchiveRequest) (*pbMaintenance.CreateSupportArchiveResponse, error) {
	supportArchive, err := d.mapRequestSettingsToSupportArchive(req)
	if err != nil {
		return nil, fmt.Errorf("failed to map support archive settings: %q", err)
	}

	_, err = d.supportArchiveClient.Create(ctx, supportArchive, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create support archive: %q", err)
	}

	return &pbMaintenance.CreateSupportArchiveResponse{}, nil
}

func (d *supportArchiveService) AllSupportArchives(ctx context.Context, _ *pbMaintenance.GetAllSupportArchivesRequest) (*pbMaintenance.GetAllSupportArchivesResponse, error) {
	list, err := d.supportArchiveClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list support archives: %w", err)
	}

	archives := make([]*pbMaintenance.Archive, len(list.Items))
	for i, item := range list.Items {
		archives[i] = &pbMaintenance.Archive{
			Name:            item.Name,
			CreatedDateTime: timestamppb.New(item.CreationTimestamp.Time),
			Status:          getStatus(item),
		}
	}

	return &pbMaintenance.GetAllSupportArchivesResponse{SupportArchives: archives}, nil
}

func (d *supportArchiveService) DeleteSupportArchive(ctx context.Context, req *pbMaintenance.DeleteSupportArchiveRequest) (*pbMaintenance.DeleteSupportArchiveResponse, error) {
	err := d.supportArchiveClient.Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete support archive: %w", err)
	}

	return &pbMaintenance.DeleteSupportArchiveResponse{}, nil
}

func (d *supportArchiveService) DownloadSupportArchive(req *pbMaintenance.DownloadSupportArchiveRequest, server pbMaintenance.SupportArchive_DownloadSupportArchiveServer) error {
	ctx := server.Context()
	archive, err := d.supportArchiveClient.Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get support archive: %w", err)
	}

	if archive.Status.DownloadPath == "" {
		return fmt.Errorf("support archive is not ready yet")
	}

	// Download the ZIP file from the download path and stream it
	reader, err := d.getDownloadFile(archive.Status.DownloadPath)
	if err != nil {
		log.Log.Info("Error in getDownloadFile", "err", err)
		return fmt.Errorf("failed to download ZIP file: %w", err)
	}
	defer reader.Close()

	err = stream.WriteReaderToStream(reader, server)
	if err != nil {
		return fmt.Errorf("failed stream support-archive file: %w", err)
	}

	return nil
}

func getStatus(archive v1.SupportArchive) pbMaintenance.SupportArchiveStatus {
	archiveStatus := archive.Status

	if len(archiveStatus.Errors) > 0 {
		return pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_FAILED
	}

	for _, condition := range archiveStatus.Conditions {
		if condition.Type == v1.ConditionSupportArchiveCreated && condition.Status == metav1.ConditionTrue {
			if archiveStatus.DownloadPath == "" {
				return pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_CREATED
			}

			return pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_COMPLETED
		}
	}

	return pbMaintenance.SupportArchiveStatus_SUPPORT_ARCHIVE_STATUS_IN_PROGRESS
}

func (d *supportArchiveService) mapRequestSettingsToSupportArchive(req *pbMaintenance.CreateSupportArchiveRequest) (*v1.SupportArchive, error) {
	exclude := req.ExcludedContents
	logConf := req.ContentTimeframe

	startTime := setDefaultTimeIfEmpty(logConf.StartDateTime, metav1.NewTime(metav1.Now().AddDate(0, 0, -4)))
	endTime := setDefaultTimeIfEmpty(logConf.EndDateTime, metav1.Now())

	if endTime.Before(&startTime) {
		return &v1.SupportArchive{}, fmt.Errorf("end time is before start time")
	}

	timestamp := metav1.Now().Format("20060102150405")
	return &v1.SupportArchive{
		ObjectMeta: metav1.ObjectMeta{
			Name: "support-archive-" + timestamp + "z",
		},
		Spec: v1.SupportArchiveSpec{
			ExcludedContents: v1.ExcludedContents{
				SystemState:   exclude.SystemState,
				SensitiveData: exclude.SensitiveData,
				Logs:          exclude.Events,
				Events:        exclude.Logs,
				VolumeInfo:    exclude.VolumeInfo,
				SystemInfo:    exclude.SystemInfo,
			},
			ContentTimeframe: v1.ContentTimeframe{
				StartTime: startTime,
				EndTime:   endTime,
			},
		},
	}, nil
}

func setDefaultTimeIfEmpty(timestamp *timestamppb.Timestamp, defaultTime metav1.Time) metav1.Time {
	var timeObj metav1.Time

	if timestamp.AsTime().IsZero() || (timestamp.GetSeconds() == 0 && timestamp.GetNanos() == 0) {
		timeObj = defaultTime
	} else {
		timeObj = metav1.NewTime(timestamp.AsTime())
	}
	return timeObj
}

// getDownloadFile downloads a file from the given URL and returns its content as an io.ReadCloser for streaming
func (d *supportArchiveService) getDownloadFile(url string) (io.ReadCloser, error) {
	client := d.httpClient
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, nil
}
