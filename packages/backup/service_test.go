package backup

import (
	"context"
	"testing"

	"github.com/cloudogu/ces-control-api/generated/backup"
	backupV1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_getAllBackups(t *testing.T) {
	t.Run("should return all backups", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
		}

		backupTwo := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_two",
			},
		}

		backupThree := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_three",
			},
		}

		backups := make([]backupV1.Backup, 0)
		backups = append(backups, backupOne, backupTwo, backupThree)

		backupClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.BackupList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    backups,
		}, nil)

		sut := DefaultBackupService{
			backupClient:  backupClientMock,
			restoreClient: nil,
		}

		// when
		allBackups, err := sut.AllBackups(testCtx, &backup.GetAllBackupsRequest{})
		// then
		require.NoError(t, err)
		assert.Equal(t, 3, len(allBackups.Backups))
	})

	t.Run("should return no backups when there are none", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)

		backupClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.BackupList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    []backupV1.Backup{},
		}, nil)

		sut := DefaultBackupService{
			backupClient:  backupClientMock,
			restoreClient: nil,
		}

		// when
		allBackups, err := sut.AllBackups(testCtx, &backup.GetAllBackupsRequest{})
		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(allBackups.Backups))
	})
}

func Test_getAllRestores(t *testing.T) {
	t.Run("should return all restores", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		restoreClientMock := newMockRestoreInterface(t)
		backupClientMock := newMockBackupInterface(t)

		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
		}

		backupTwo := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_two",
			},
		}

		backupThree := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_three",
			},
		}

		restoreOne := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_one",
			},
			Spec: backupV1.RestoreSpec{
				BackupName: "backup_one",
			},
		}

		restoreTwo := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_two",
			},
			Spec: backupV1.RestoreSpec{
				BackupName: "backup_two",
			},
		}

		restoreThree := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_three",
			},
			Spec: backupV1.RestoreSpec{
				BackupName: "backup_three",
			},
		}

		restores := make([]backupV1.Restore, 0)
		restores = append(restores, restoreOne, restoreTwo, restoreThree)

		restoreClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.RestoreList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    restores,
		}, nil)

		backupClientMock.EXPECT().Get(testCtx, "backup_one", metav1.GetOptions{}).Return(&backupOne, nil)
		backupClientMock.EXPECT().Get(testCtx, "backup_two", metav1.GetOptions{}).Return(&backupTwo, nil)
		backupClientMock.EXPECT().Get(testCtx, "backup_three", metav1.GetOptions{}).Return(&backupThree, nil)

		sut := DefaultBackupService{
			backupClient:  backupClientMock,
			restoreClient: restoreClientMock,
		}

		// when
		allRestores, err := sut.AllRestores(testCtx, &backup.GetAllRestoresRequest{})
		// then
		require.NoError(t, err)
		assert.Equal(t, 3, len(allRestores.Restores))
	})

	t.Run("should return no restores when there are none", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		restoreClientMock := newMockRestoreInterface(t)
		restoreOne := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_one",
			},
		}

		restoreTwo := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_two",
			},
		}

		restoreThree := backupV1.Restore{
			ObjectMeta: metav1.ObjectMeta{
				Name: "restore_three",
			},
		}

		restores := make([]backupV1.Restore, 3)
		restores = append(restores, restoreOne, restoreTwo, restoreThree)

		restoreClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.RestoreList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    []backupV1.Restore{},
		}, nil)

		sut := DefaultBackupService{
			backupClient:  nil,
			restoreClient: restoreClientMock,
		}

		// when
		allRestores, err := sut.AllRestores(testCtx, &backup.GetAllRestoresRequest{})
		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(allRestores.Restores))
	})
}

func TestDefaultBackupService_GetSchedule(t *testing.T) {
	testCtx := context.Background()
	schedule := &backupV1.BackupSchedule{Spec: backupV1.BackupScheduleSpec{Schedule: "1 2 * 4 *"}}

	t.Run("should get schedule", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(schedule, nil)

		svc := &DefaultBackupService{
			backupScheduleClient: mBackupScheduleClient,
		}

		response, err := svc.GetSchedule(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, "1 2 * 4 *", response.Schedule)
	})

	t.Run("should fail to get empty schedule with error", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, assert.AnError)

		svc := &DefaultBackupService{
			backupScheduleClient: mBackupScheduleClient,
		}

		_, err := svc.GetSchedule(testCtx, nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get backup schedule:")
	})
}

func TestDefaultBackupService_SetSchedule(t *testing.T) {
	testCtx := context.Background()

	t.Run("should set schedule", func(t *testing.T) {
		expectedSchedule := &backupV1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
			Spec:       backupV1.BackupScheduleSpec{Schedule: "* 2 3 * *"},
		}

		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, "not found"))
		mBackupScheduleClient.EXPECT().Create(testCtx, expectedSchedule, metav1.CreateOptions{}).Return(expectedSchedule, nil)

		svc := &DefaultBackupService{
			backupScheduleClient: mBackupScheduleClient,
		}

		response, err := svc.SetSchedule(testCtx, &backup.SetBackupScheduleRequest{Schedule: "* 2 3 * *"})

		require.NoError(t, err)
		assert.Equal(t, &backup.SetBackupScheduleResponse{}, response)
	})

	t.Run("should fail to set schedule", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, assert.AnError)

		svc := &DefaultBackupService{
			backupScheduleClient: mBackupScheduleClient,
		}

		_, err := svc.SetSchedule(testCtx, &backup.SetBackupScheduleRequest{Schedule: "* 2 3 * *"})

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get existing backup schedule:")
	})
}
