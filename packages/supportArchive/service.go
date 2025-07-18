package supportArchive

import (
	"fmt"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	v1 "github.com/cloudogu/k8s-support-archive-lib/api/v1"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
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

func (d *supportArchiveService) Create(req *pbMaintenance.CreateSupportArchiveRequest, server pbMaintenance.SupportArchive_CreateServer) error {
	supportArchive, err := d.mapRequestSettingsToSupportArchive(req)
	if err != nil {
		return fmt.Errorf("failed to map support archive settings: %q", err)
	}

	downloadPath, err := d.createAndWatchSupportArchive(supportArchive, server)
	if err != nil {
		return fmt.Errorf("failed to create or watch support archive: %q", err)
	}

	// Download the ZIP file from the download path
	zipContent, err := d.downloadFile(downloadPath)
	if err != nil {
		log.Log.Info("Error in downloadFile", "err", err)
		return fmt.Errorf("failed to download ZIP file: %q", err)
	}

	log.Log.Info("Write to Stream", "zipContent", zipContent)
	err = d.writeToStream(zipContent, server)
	if err != nil {
		return fmt.Errorf("failed to send response: %q", err)
	}

	return nil
}

func (d *supportArchiveService) mapRequestSettingsToSupportArchive(req *pbMaintenance.CreateSupportArchiveRequest) (*v1.SupportArchive, error) {
	log.Log.Info("mapRequestSettingsToSupportArchive", "request", req)
	//TODO: potential nil pointer!
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

	timestamp := metav1.Now().Format("20060102150405")
	return &v1.SupportArchive{
		ObjectMeta: metav1.ObjectMeta{
			Name: "support-archive-" + timestamp + "z",
		},
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
	ctx := server.Context()
	_, err := d.supportArchiveClient.Create(ctx, supportArchive, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create support archive: %q", err)
	}

	watchInterface, err := d.supportArchiveClient.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create watch interface: %q", err)
	}

	for {
		select {
		case event, channelOk := <-watchInterface.ResultChan():
			if !channelOk {
				return "", fmt.Errorf("watch channel closed unexpectedly")
			}

			if event.Type == watch.Added || event.Type == watch.Modified {
				eventSupportArchive, typeOk := event.Object.(*v1.SupportArchive)
				if !typeOk {
					return "", fmt.Errorf("unexpected type")
				}
				// ignore changes to other supportArchives
				if eventSupportArchive.Name != supportArchive.Name {
					continue
				}
				// TODO watch Conditions instead of Phase
				if eventSupportArchive.Status.Phase == v1.StatusPhaseCreated {
					return eventSupportArchive.Status.DownloadPath, nil
				}
			}
		case <-server.Context().Done():
			return "", fmt.Errorf("timed out waiting for support archive to be created")
		}
	}
}

// downloadFile downloads a file from the given URL and returns its content as a byte slice
func (d *supportArchiveService) downloadFile(url string) ([]byte, error) {
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
	log.Log.Info("Request send", "Response", resp)

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing response body: %w", closeErr)
		}
	}()

	log.Log.Info("Check Status code", "status", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	log.Log.Info("Read Body", "body", resp.Body)
	return io.ReadAll(resp.Body)
}
