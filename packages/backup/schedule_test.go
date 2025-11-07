package backup

import (
	"context"
	"testing"

	backupV1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_getBackupSchedule(t *testing.T) {
	testCtx := context.Background()
	schedule := &backupV1.BackupSchedule{Spec: backupV1.BackupScheduleSpec{Schedule: "1 2 * 4 *"}}

	t.Run("should get schedule", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(schedule, nil)

		schedule, err := getBackupSchedule(testCtx, mBackupScheduleClient)

		require.NoError(t, err)
		assert.Equal(t, "1 2 * 4 *", schedule)
	})

	t.Run("should get empty schedule if non exists", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, "not found"))

		schedule, err := getBackupSchedule(testCtx, mBackupScheduleClient)

		require.NoError(t, err)
		assert.Equal(t, "", schedule)
	})

	t.Run("should fail to get empty schedule with error", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, assert.AnError)

		_, err := getBackupSchedule(testCtx, mBackupScheduleClient)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get backup schedule:")
	})
}

func Test_setBackupSchedule(t *testing.T) {
	testCtx := context.Background()
	schedule := &backupV1.BackupSchedule{
		ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
		Spec:       backupV1.BackupScheduleSpec{Schedule: "1 2 * 4 *"},
	}

	t.Run("should create new schedule", func(t *testing.T) {
		expectedSchedule := &backupV1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
			Spec:       backupV1.BackupScheduleSpec{Schedule: "* 2 3 * *"},
		}

		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, "not found"))
		mBackupScheduleClient.EXPECT().Create(testCtx, expectedSchedule, metav1.CreateOptions{}).Return(expectedSchedule, nil)

		err := setBackupSchedule(testCtx, mBackupScheduleClient, "* 2 3 * *")

		require.NoError(t, err)
	})

	t.Run("should fail to create new schedule", func(t *testing.T) {
		expectedSchedule := &backupV1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
			Spec:       backupV1.BackupScheduleSpec{Schedule: "* 2 3 * *"},
		}

		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, "not found"))
		mBackupScheduleClient.EXPECT().Create(testCtx, expectedSchedule, metav1.CreateOptions{}).Return(schedule, assert.AnError)

		err := setBackupSchedule(testCtx, mBackupScheduleClient, "* 2 3 * *")

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create backup schedule:")
	})

	t.Run("should update schedule", func(t *testing.T) {
		expectedSchedule := &backupV1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
			Spec:       backupV1.BackupScheduleSpec{Schedule: "* 2 3 * *"},
		}

		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(schedule, nil)
		mBackupScheduleClient.EXPECT().Update(testCtx, expectedSchedule, metav1.UpdateOptions{}).Return(schedule, nil)

		err := setBackupSchedule(testCtx, mBackupScheduleClient, "* 2 3 * *")

		require.NoError(t, err)
	})

	t.Run("should fail to update schedule", func(t *testing.T) {
		expectedSchedule := &backupV1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{Name: "ces-schedule"},
			Spec:       backupV1.BackupScheduleSpec{Schedule: "* 2 3 * *"},
		}

		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(schedule, nil)
		mBackupScheduleClient.EXPECT().Update(testCtx, expectedSchedule, metav1.UpdateOptions{}).Return(nil, assert.AnError)

		err := setBackupSchedule(testCtx, mBackupScheduleClient, "* 2 3 * *")

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update backup schedule:")
	})

	t.Run("should fail to set schedule for error while getting schedule", func(t *testing.T) {
		mBackupScheduleClient := newMockBackupScheduleClient(t)
		mBackupScheduleClient.EXPECT().Get(testCtx, "ces-schedule", metav1.GetOptions{}).Return(nil, assert.AnError)

		err := setBackupSchedule(testCtx, mBackupScheduleClient, "* 2 3 * *")

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get existing backup schedule:")
	})
}
