package supportArchive

import (
	"fmt"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type supportArchiveService struct {
	pbMaintenance.UnimplementedSupportArchiveServer
	supportArchiveClient supportArchiveClient
	writeToStream        stream.Writer
}

func NewSupportArchiveService(client supportArchiveClient) *supportArchiveService {
	return &supportArchiveService{
		supportArchiveClient: client,
		writeToStream:        stream.WriteToStream,
	}
}

func (d *supportArchiveService) Create(req *pbMaintenance.CreateSupportArchiveRequest, server pbMaintenance.SupportArchive_CreateServer) error {
	supportArchive, err := d.mapRequestSettingsToSupportArchive(req)
	if err != nil {
		return fmt.Errorf("failed to map support archive settings: %q", err)
	}

	downloadPath, err := d.createAndWatchSupportArchive(supportArchive, server)
	if err != nil {
		return fmt.Errorf("failed to create or watch support archive: %q", err)
	}

	err = d.writeToStream([]byte(downloadPath), server)
	if err != nil {
		return fmt.Errorf("failed to send response: %q", err)
	}

	return nil
}

func (d *supportArchiveService) mapRequestSettingsToSupportArchive(req *pbMaintenance.CreateSupportArchiveRequest) (*v1.SupportArchive, error) {
	exclude := req.GetCommon().ExcludedContents
	logConf := req.GetCommon().LoggingConfig

	startTime := metav1.NewTime(logConf.StartDateTime.AsTime())
	var endTime metav1.Time
	// set endTime to now if it was not set
	if logConf.EndDateTime.AsTime().IsZero() || (logConf.EndDateTime.GetSeconds() == 0 && logConf.EndDateTime.GetNanos() == 0) {
		endTime = metav1.Now()
	} else {
		endTime = metav1.NewTime(logConf.EndDateTime.AsTime())
	}

	if endTime.Before(&startTime) {
		return &v1.SupportArchive{}, fmt.Errorf("end time is before start time")
	}

	return &v1.SupportArchive{
		Spec: v1.SupportArchiveSpec{
			ExcludedContents: v1.ExcludedContents{
				SystemState:   exclude.SystemState,
				SensitiveData: exclude.SensitiveData,
				LogsAndEvents: exclude.LogsAndEvents,
				VolumeInfo:    exclude.VolumeInfo,
				SystemInfo:    exclude.SystemInfo,
			},
			LoggingConfig: v1.LoggingConfig{
				StartTime: startTime,
				EndTime:   endTime,
			},
		},
	}, nil
}

func (d *supportArchiveService) createAndWatchSupportArchive(supportArchive *v1.SupportArchive, server pbMaintenance.SupportArchive_CreateServer) (string, error) {
	_, err := d.supportArchiveClient.Create(server.Context(), supportArchive, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create support archive: %q", err)
	}

	watchInterface, err := d.supportArchiveClient.Watch(server.Context(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create watch interface: %q", err)
	}

	for event := range watchInterface.ResultChan() {
		if event.Type == watch.Added || event.Type == watch.Modified {
			supportArchive, ok := event.Object.(*v1.SupportArchive)
			if !ok {
				return "", fmt.Errorf("unexpected type")
			}
			// TODO watch Conditions instead of Phase
			if supportArchive.Status.Phase == v1.StatusPhaseCreated {
				return supportArchive.Status.DownloadPath, nil
			}
		}
	}
	return "", fmt.Errorf("failed to watch support Archive")
}
