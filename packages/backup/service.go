package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	pbBackup "github.com/cloudogu/ces-control-api/generated/backup"
	v1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	v3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	backupStatusInProgress = "inProgress"
	backupStatusCompleted  = "completed"
	backupStatusFailed     = "failed"
	backupStatusDeleting   = "deleting"
)

const (
	blueprintIdAnnotation = "backup.cloudogu.com/blueprintId"
	dogusAnnotation       = "backup.cloudogu.com/dogus"
)

// annotationDogus are found in the annotation "backup.cloudogu.com/dogus" of a backup
type annotationDogus struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DefaultBackupService struct {
	pbBackup.UnimplementedBackupManagementServer
	backupClient         backupInterface
	restoreClient        restoreInterface
	backupScheduleClient backupScheduleClient
	componentClient      componentClient
	blueprintLister      blueprintLister
}

// NewBackupService returns an instance of defaultBackupService.
func NewBackupService(backupClient backupInterface, restoreClient restoreInterface, backupScheduleClient backupScheduleClient, componentClient componentClient, blueprintLister blueprintLister) *DefaultBackupService {
	return &DefaultBackupService{
		backupClient:         backupClient,
		restoreClient:        restoreClient,
		backupScheduleClient: backupScheduleClient,
		componentClient:      componentClient,
		blueprintLister:      blueprintLister,
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

	blueprintList, err := s.blueprintLister.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint: %w", err)
	}
	if len(blueprintList.Items) == 0 {
		return nil, fmt.Errorf("failed to get blueprint: no blueprints available")
	}

	backups := s.mapBackups(list, &blueprintList.Items[0])
	return &pbBackup.GetAllBackupsResponse{Backups: backups}, nil
}

// CreateRestore creates a restore for the given backup.
// The restore is only created if the backup is restorable.
func (s *DefaultBackupService) CreateRestore(ctx context.Context, request *pbBackup.CreateRestoreRequest) (*pbBackup.CreateRestoreResponse, error) {
	backup, err := s.backupClient.Get(ctx, request.BackupId, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get backup: %w", err)
	}

	list, err := s.blueprintLister.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint: %w", err)
	}
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("failed to get blueprint: blueprint not found")
	}

	// only create restore if the backup is restorable
	restorable, err := s.isBackupRestorable(backup, &list.Items[0])
	if err != nil {
		return nil, fmt.Errorf("failed to check if backup is restorable: %w", err)
	}
	if restorable {
		timestamp := time.Now().Format("20060102-1504")
		restoreName := fmt.Sprintf("restore-%s", timestamp)
		restore := &v1.Restore{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: restoreName,
			},
			Spec: v1.RestoreSpec{
				BackupName: request.BackupId,
			},
			Status: v1.RestoreStatus{},
		}
		_, err := s.restoreClient.Create(ctx, restore, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create restore: %w", err)
		}
	} else {
		return nil, fmt.Errorf("backup is not restorable")
	}

	return &pbBackup.CreateRestoreResponse{}, nil
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

func (s *DefaultBackupService) mapBackups(backupList *v1.BackupList, blueprint *v3.Blueprint) []*pbBackup.BackupResponse {
	backupResponseList := make([]*pbBackup.BackupResponse, 0, 5)
	for _, backup := range backupList.Items {
		// skip backups in deleting state
		if backup.Status.Status == backupStatusDeleting {
			continue
		}
		restorable, err := s.isBackupRestorable(&backup, blueprint)
		if err != nil {
			// There might be backups that do not have the annotations. In this case we just log the error and continue.
			slog.Error(fmt.Sprintf("failed to check if backup is restorable: %v", err))
			restorable = false
		}
		backupResponse := pbBackup.BackupResponse{
			Id:             backup.Name,
			StartTime:      backup.Status.StartTimestamp.String(),
			EndTime:        backup.Status.CompletionTimestamp.String(),
			Status:         backupStatus(&backup),
			CurrentVersion: true,
			Restorable:     restorable && backup.Status.Status == backupStatusCompleted,
			BlueprintId:    backup.GetAnnotations()[blueprintIdAnnotation],
		}
		backupResponseList = append(backupResponseList, &backupResponse)
	}

	return backupResponseList
}

func (s *DefaultBackupService) mapRestores(ctx context.Context, restoreList *v1.RestoreList) ([]*pbBackup.RestoreResponse, error) {
	restoreResponseList := make([]*pbBackup.RestoreResponse, 0, 5)
	for _, restore := range restoreList.Items {
		blueprintId := ""
		backup, err := s.backupClient.Get(ctx, restore.Spec.BackupName, metav1.GetOptions{})
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				slog.Warn(fmt.Sprintf("could not find backup for restore: %v", err))
			} else {
				return nil, fmt.Errorf("failed to get backup for restore %s: %w", restore.Name, err)
			}
		} else {
			blueprintId = backup.GetAnnotations()[blueprintIdAnnotation]
		}

		restoreResponse := pbBackup.RestoreResponse{
			BackupId:    restore.Spec.BackupName,
			StartTime:   restore.CreationTimestamp.String(),
			Success:     restore.Status.Status == "completed",
			BlueprintId: blueprintId,
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

// a backup is restorable if it is from the same blueprint and the dogus are matching
func (s *DefaultBackupService) isBackupRestorable(backup *v1.Backup, blueprint *v3.Blueprint) (bool, error) {
	ans := backup.GetAnnotations()

	// get all dogus from backup annotations
	backupDogus := make([]annotationDogus, 0, 5)
	err := json.Unmarshal([]byte(ans[dogusAnnotation]), &backupDogus)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal dogus from backup: %w", err)
	}

	blueprintNameIsMatching := ans != nil && ans[blueprintIdAnnotation] == blueprint.Spec.DisplayName

	if blueprintNameIsMatching && s.isDoguListMatching(backupDogus, blueprint.Spec.Blueprint.Dogus) {
		return true, nil
	}
	return false, nil
}

// isDoguListMatching checks if the given list of backup dogus is matching the given list of dogus from the blueprint.
func (s *DefaultBackupService) isDoguListMatching(backupDogus []annotationDogus, blueprintDogus []v3.Dogu) bool {
	if len(backupDogus) != len(blueprintDogus) {
		return false
	}

	backupDoguMap := make(map[string]annotationDogus)
	for _, v := range backupDogus {
		backupDoguMap[v.Name] = v
	}
	for _, blueprintDogu := range blueprintDogus {
		backupDogu, ok := backupDoguMap[blueprintDogu.Name]
		if !ok {
			return false
		}
		if backupDogu.Version != *blueprintDogu.Version {
			return false
		}
	}
	return true
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
