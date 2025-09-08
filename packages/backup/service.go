package backup

import (
	"context"

	pbBackup "github.com/cloudogu/ces-control-api/generated/backup"
	v1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type defaultBackupService struct {
	pbBackup.BackupManagementServer
	backupClient  backupInterface
	restoreClient restoreInterface
}

// NewBackupService returns an instance of defaultBackupService.
func NewBackupService(backupClient backupInterface, restoreClient restoreInterface) *defaultBackupService {
	return &defaultBackupService{
		backupClient:  backupClient,
		restoreClient: restoreClient,
	}
}

func (s *defaultBackupService) AllBackups(ctx context.Context, _ *pbBackup.GetAllBackupsRequest) (*pbBackup.GetAllBackupsResponse, error) {
	list, err := s.backupClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &pbBackup.GetAllBackupsResponse{Backups: s.mapBackups(list)}, nil
}

func (s *defaultBackupService) AllRestores(ctx context.Context, _ *pbBackup.GetAllRestoresRequest) (*pbBackup.GetAllRestoresResponse, error) {
	list, err := s.restoreClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	restores, err := s.mapRestores(ctx, list)
	if err != nil {
		return nil, err
	}

	return &pbBackup.GetAllRestoresResponse{Restores: restores}, nil
}

func (s *defaultBackupService) mapBackups(backupList *v1.BackupList) []*pbBackup.BackupResponse {
	backupResponseList := make([]*pbBackup.BackupResponse, 0, 5)
	for _, backup := range backupList.Items {
		backupResponse := pbBackup.BackupResponse{
			Id:             string(backup.UID),
			StartTime:      backup.Status.StartTimestamp.String(),
			EndTime:        backup.Status.CompletionTimestamp.String(),
			Success:        backup.Status.Status == "completed",
			CurrentVersion: true,
		}
		backupResponseList = append(backupResponseList, &backupResponse)
	}

	return backupResponseList
}

func (s *defaultBackupService) mapRestores(ctx context.Context, restoreList *v1.RestoreList) ([]*pbBackup.RestoreResponse, error) {
	restoreResponseList := make([]*pbBackup.RestoreResponse, 0, 5)
	for _, restore := range restoreList.Items {
		backup, err := s.backupClient.Get(ctx, restore.Spec.BackupName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		restoreResponse := pbBackup.RestoreResponse{
			BackupId:    string(backup.UID),
			StartTime:   backup.Status.StartTimestamp.String(),
			EndTime:     backup.Status.CompletionTimestamp.String(),
			Success:     restore.Status.Status == "completed",
			BlueprintId: "Unknown", // FIXME: fix once backup POC is completed and BlueprintId has been added in some shape or form.
		}
		restoreResponseList = append(restoreResponseList, &restoreResponse)
	}

	return restoreResponseList, nil
}
