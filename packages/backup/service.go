package backup

import (
	"context"
	"fmt"
	"time"

	pbBackup "github.com/cloudogu/ces-control-api/generated/backup"
	v1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	backupStatusInProgress = "inProgress"
	backupStatusCompleted  = "completed"
	backupStatusFailed     = "failed"
	backupStatusDeleting   = "deleting"
)

type DefaultBackupService struct {
	pbBackup.UnimplementedBackupManagementServer
	backupClient         backupInterface
	restoreClient        restoreInterface
	backupScheduleClient backupScheduleClient
	componentClient      componentClient
}

// NewBackupService returns an instance of defaultBackupService.
func NewBackupService(backupClient backupInterface, restoreClient restoreInterface, backupScheduleClient backupScheduleClient, componentClient componentClient) *DefaultBackupService {
	return &DefaultBackupService{
		backupClient:         backupClient,
		restoreClient:        restoreClient,
		backupScheduleClient: backupScheduleClient,
		componentClient:      componentClient,
	}
}

func (s *DefaultBackupService) DeleteBackup(ctx context.Context, req *pbBackup.DeleteBackupRequest) (*pbBackup.DeleteBackupResponse, error) {
	err := s.backupClient.Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete backup: %w", err)
	}
	return &pbBackup.DeleteBackupResponse{}, nil
}

func (s *DefaultBackupService) CreateBackup(ctx context.Context, _ *pbBackup.CreateBackupRequest) (*pbBackup.CreateBackupResponse, error) {
	timestamp := time.Now().Format("20060102-1504")
	backupName := fmt.Sprintf("backup-%s", timestamp)
	backup := &v1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name: backupName,
		},
		Spec: v1.BackupSpec{
			SyncedFromProvider: false,
		},
	}
	_, err := s.backupClient.Create(ctx, backup, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &pbBackup.CreateBackupResponse{}, nil
}

func (s *DefaultBackupService) AllBackups(ctx context.Context, _ *pbBackup.GetAllBackupsRequest) (*pbBackup.GetAllBackupsResponse, error) {
	list, err := s.backupClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	return &pbBackup.GetAllBackupsResponse{Backups: s.mapBackups(list)}, nil
}

func (s *DefaultBackupService) AllRestores(ctx context.Context, _ *pbBackup.GetAllRestoresRequest) (*pbBackup.GetAllRestoresResponse, error) {
	list, err := s.restoreClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list restores: %w", err)
	}

	restores, err := s.mapRestores(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to map restores to dto: %w", err)
	}

	return &pbBackup.GetAllRestoresResponse{Restores: restores}, nil
}

func (s *DefaultBackupService) mapBackups(backupList *v1.BackupList) []*pbBackup.BackupResponse {
	backupResponseList := make([]*pbBackup.BackupResponse, 0, 5)
	for _, backup := range backupList.Items {
		// skip backups in deleting state
		if backup.Status.Status == backupStatusDeleting {
			continue
		}
		backupResponse := pbBackup.BackupResponse{
			Id:             backup.Name,
			StartTime:      backup.Status.StartTimestamp.String(),
			EndTime:        backup.Status.CompletionTimestamp.String(),
			Status:         backupStatus(&backup),
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
			return nil, fmt.Errorf("failed to get backup for restore %s: %w", restore.Name, err)
		}

		restoreResponse := pbBackup.RestoreResponse{
			BackupId:    backup.Name,
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
		return nil, fmt.Errorf("failed to get backup schedule: %w", err)
	}

	return &pbBackup.GetBackupScheduleResponse{Schedule: schedule}, nil
}

func (s *DefaultBackupService) SetSchedule(ctx context.Context, req *pbBackup.SetBackupScheduleRequest) (*pbBackup.SetBackupScheduleResponse, error) {
	err := setBackupSchedule(ctx, s.backupScheduleClient, req.Schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to set backup schedule: %w", err)
	}

	return &pbBackup.SetBackupScheduleResponse{}, nil
}

func (s *DefaultBackupService) GetRetentionPolicy(ctx context.Context, _ *pbBackup.GetRetentionPolicyRequest) (*pbBackup.GetRetentionPolicyResponse, error) {
	policy, err := getRetentionPolicy(ctx, s.componentClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get retention policy: %w", err)
	}

	// map policy to protobuf enum
	var retentionPolicy pbBackup.RetentionPolicy
	switch policy {
	case string(keepAllPolicy):
		retentionPolicy = pbBackup.RetentionPolicy_RETENTION_POLICY_KEEP_ALL
	case string(removeAllButKeepLatestPolicy):
		retentionPolicy = pbBackup.RetentionPolicy_RETENTION_POLICY_REMOVE_ALL_BUT_KEEP_LATEST
	case string(keepLastSevenDaysPolicy):
		retentionPolicy = pbBackup.RetentionPolicy_RETENTION_POLICY_KEEP_LAST_SEVEN_DAYS
	case string(keepLast7DaysOldestOf1Month1Quarter1HalfYear1YearPolicy):
		retentionPolicy = pbBackup.RetentionPolicy_RETENTION_POLICY_KEEP_LAST_7_DAYS_OLDEST_OF_1_MONTH_1_QUARTER_1_HALF_YEAR_1_YEAR
	default:
		retentionPolicy = pbBackup.RetentionPolicy_RETENTION_POLICY_UNSPECIFIED
	}

	return &pbBackup.GetRetentionPolicyResponse{Policy: retentionPolicy}, nil
}

func backupStatus(backup *v1.Backup) string {
	switch backup.Status.Status {
	case backupStatusCompleted:
		return backupStatusCompleted
	case backupStatusFailed:
		return backupStatusFailed
	default:
		return backupStatusInProgress
	}
}
