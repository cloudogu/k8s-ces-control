package support

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/cesapp-lib/archive"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
)

type defaultSupportArchiveService struct {
	pbMaintenance.UnimplementedSupportArchiveServer
	createArchive func(ctx context.Context, environment *pbMaintenance.CreateSupportArchiveRequest) ([]byte, error)
}

type archiveManager interface {
	GetContent() []byte
	AddContentAsFile(content string, fileNameInArchive string) error
	AddContentAsFileWithModifiedDate(content string, fileNameInArchive string, modified time.Time) error
	AddFileToArchive(file archive.File) error
	AddFilesToArchive(files []archive.File, closeAfterFinish bool) error
	SaveArchiveAsFile(archivePath string) error
	Close() error
}

func (d *defaultSupportArchiveService) Create(request *pbMaintenance.CreateSupportArchiveRequest, server pbMaintenance.SupportArchive_CreateServer) error {

	dataToStream, err := d.createArchive(server.Context(), request)
	if err != nil {
		return fmt.Errorf("create support archive failed: %w", err)
	}

	return stream.WriteToStream(dataToStream, server)
}

type SupportArchiveCreator struct {
	environment       *pbMaintenance.CreateSupportArchiveRequest_Common
	ArchiveManager    archiveManager
	resourceCollector resourceCollector
}

func (sac *SupportArchiveCreator) CreateSupportArchive(ctx context.Context) ([]byte, error) {
	manager, err := sac.collectArchiveData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect support archive data: %w", err)
	}
	return manager.GetContent(), nil
}

func (sac *SupportArchiveCreator) collectArchiveData(ctx context.Context) (archive.Manager, error) {
	manager := sac.ArchiveManager
	defer func() {
		_ = manager.Close()
	}()

	if !sac.environment.Common.ExcludedContents.SystemState {
		labelSelector := &metav1.LabelSelector{MatchLabels: map[string]string{"app": "ces"}}
		excludedGVKs := []gvkMatcher{{
			Version: "v1",
			Kind:    "Secret",
		}}
		resources, err := sac.resourceCollector.Collect(ctx, labelSelector, excludedGVKs)
		if err != nil {
			return archive.Manager{}, fmt.Errorf("collect resources failed: %w", err)
		}
		for _, item := range resources {
			gvk := item.GroupVersionKind()
			resource, err := yaml.Marshal(item)
			if err != nil {
				return archive.Manager{}, fmt.Errorf("failed to marshal resource to yaml: %w", err)
			}
			// File name should be like this k8s/k8s.cloudogu.com/backups/MyBackup.yaml
			err = sac.ArchiveManager.AddContentAsFile(string(resource), fmt.Sprintf("k8s/%s/%s/%s.yaml", gvk.Group, gvk.Kind, item.GetName()))
			if err != nil {
				return archive.Manager{}, fmt.Errorf("failed to add resource to archive: %w", err)
			}
		}

	}

	return *archive.NewManager(), nil
}

func NewSupportArchiveService(client k8sClient, discoveryClient discoveryInterface) *defaultSupportArchiveService {
	createArchiveFunc := func(ctx context.Context, request *pbMaintenance.CreateSupportArchiveRequest) ([]byte, error) {
		environment, ok := request.GetEnvironment().(*pbMaintenance.CreateSupportArchiveRequest_Common)
		if !ok {
			return nil, fmt.Errorf("unsupported environment; the sent environment is probably legacy but must be common")
		}
		archiveCreator := &SupportArchiveCreator{
			environment:    environment,
			ArchiveManager: *archive.NewManager(),
			resourceCollector: &defaultResourceCollector{
				client:          client,
				discoveryClient: discoveryClient,
			},
		}
		return archiveCreator.CreateSupportArchive(ctx)
	}
	return &defaultSupportArchiveService{
		createArchive: createArchiveFunc,
	}
}
