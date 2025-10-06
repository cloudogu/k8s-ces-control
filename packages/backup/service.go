package backup

import (
	"context"
	"fmt"

	pbBackup "github.com/cloudogu/ces-control-api/generated/backup"
	v1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DefaultBackupService struct {
	pbBackup.UnimplementedBackupManagementServer
	backupClient         backupInterface
	restoreClient        restoreInterface
	backupScheduleClient backupScheduleClient
}

// NewBackupService returns an instance of defaultBackupService.
func NewBackupService(backupClient backupInterface, restoreClient restoreInterface, backupScheduleClient backupScheduleClient) *DefaultBackupService {
	return &DefaultBackupService{
		backupClient:         backupClient,
		restoreClient:        restoreClient,
		backupScheduleClient: backupScheduleClient,
	}
}

func (s *DefaultBackupService) AllBackups(ctx context.Context, _ *pbBackup.GetAllBackupsRequest) (*pbBackup.GetAllBackupsResponse, error) {
	list, err := s.backupClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &pbBackup.GetAllBackupsResponse{Backups: s.mapBackups(list)}, nil
}

func (s *DefaultBackupService) AllRestores(ctx context.Context, _ *pbBackup.GetAllRestoresRequest) (*pbBackup.GetAllRestoresResponse, error) {
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

func (s *DefaultBackupService) mapBackups(backupList *v1.BackupList) []*pbBackup.BackupResponse {
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

func (s *DefaultBackupService) mapRestores(ctx context.Context, restoreList *v1.RestoreList) ([]*pbBackup.RestoreResponse, error) {
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

func (s *DefaultBackupService) GetSchedule(ctx context.Context, _ *pbBackup.GetBackupScheduleRequest) (*pbBackup.GetBackupScheduleResponse, error) {
	schedule, err := getBackupSchedule(ctx, s.backupScheduleClient)
	if err != nil {
		return nil, err
	}

	return &pbBackup.GetBackupScheduleResponse{Schedule: schedule}, nil
}

func (s *DefaultBackupService) SetSchedule(ctx context.Context, req *pbBackup.SetBackupScheduleRequest) (*pbBackup.SetBackupScheduleResponse, error) {
	err := setBackupSchedule(ctx, s.backupScheduleClient, req.Schedule)
	if err != nil {
		return nil, err
	}

	return &pbBackup.SetBackupScheduleResponse{}, nil
}

func (s *DefaultBackupService) GetRetentionPolicy(context.Context, *pbBackup.RetentionPolicyRequest) (*pbBackup.RetentionPolicyResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
