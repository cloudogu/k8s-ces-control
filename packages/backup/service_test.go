package backup

import (
	"context"
	"testing"

	"github.com/cloudogu/ces-control-api/generated/backup"
	backupV1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	v2 "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	componentV1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_getAllBackups(t *testing.T) {
	backupAnnotations := make(map[string]string)
	backupAnnotations["backup.cloudogu.com/dogus"] = "[{\"name\": \"hallowelt/bluespice\", \"version\": \"4.5.5-3\"},{\"name\": \"official\", \"version\": \"7.2.6-3\"}]"
	backupAnnotations["backup.cloudogu.com/blueprintId"] = "all-dogus-sample"

	t.Run("should return all backups", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_one",
				Annotations: backupAnnotations,
			},
		}

		backupTwo := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_two",
				Annotations: backupAnnotations,
			},
		}

		backupThree := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_three",
				Annotations: backupAnnotations,
			},
		}

		backups := make([]backupV1.Backup, 0)
		backups = append(backups, backupOne, backupTwo, backupThree)

		backupClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.BackupList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    backups,
		}, nil)

		bps := v2.BlueprintList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items: []v2.Blueprint{{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
			}},
		}

		lister := newMockBlueprintLister(t)
		lister.EXPECT().List(testCtx, mock.Anything).Return(&bps, nil)

		sut := DefaultBackupService{
			backupClient:    backupClientMock,
			blueprintLister: lister,
			restoreClient:   nil,
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

		bps := v2.BlueprintList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items: []v2.Blueprint{{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
			}},
		}

		lister := newMockBlueprintLister(t)
		lister.EXPECT().List(testCtx, mock.Anything).Return(&bps, nil)

		sut := DefaultBackupService{
			backupClient:    backupClientMock,
			restoreClient:   nil,
			blueprintLister: lister,
		}

		// when
		allBackups, err := sut.AllBackups(testCtx, &backup.GetAllBackupsRequest{})
		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(allBackups.Backups))
	})

	t.Run("should error when blueprintclient is not available", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_one",
				Annotations: backupAnnotations,
			},
		}

		backups := make([]backupV1.Backup, 0)
		backups = append(backups, backupOne)

		backupClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&backupV1.BackupList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    backups,
		}, nil)

		lister := newMockBlueprintLister(t)
		lister.EXPECT().List(testCtx, mock.Anything).Return(nil, assert.AnError)

		sut := DefaultBackupService{
			backupClient:    backupClientMock,
			blueprintLister: lister,
			restoreClient:   nil,
		}

		// when
		_, err := sut.AllBackups(testCtx, &backup.GetAllBackupsRequest{})
		// then
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_getAllRestores(t *testing.T) {
	backupAnnotations := make(map[string]string)
	backupAnnotations["backup.cloudogu.com/dogus"] = "[{\"name\": \"hallowelt/bluespice\", \"version\": \"4.5.5-3\"},{\"name\": \"official\", \"version\": \"7.2.6-3\"}]"
	backupAnnotations["backup.cloudogu.com/blueprintId"] = "all-dogus-sample"

	t.Run("should return all restores", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		restoreClientMock := newMockRestoreInterface(t)
		backupClientMock := newMockBackupInterface(t)

		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_one",
				Annotations: backupAnnotations,
			},
		}

		backupTwo := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_two",
				Annotations: backupAnnotations,
			},
		}

		backupThree := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "backup_three",
				Annotations: backupAnnotations,
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

func TestDefaultBackupService_GetRetentionPolicy(t *testing.T) {
	testCtx := context.Background()

	t.Run("should get default policy", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{},
		}, nil)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		response, err := svc.GetRetentionPolicy(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, backup.RetentionPolicy_RETENTION_POLICY_UNSPECIFIED, response.Policy)
	})

	t.Run("should get policy", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
cleanup:
  exclude: foo
retention:
  strategy: "removeAllButKeepLatest"
  garbageCollectionCron: "0 * * * *"
`,
			},
		}, nil)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		response, err := svc.GetRetentionPolicy(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, backup.RetentionPolicy_RETENTION_POLICY_REMOVE_ALL_BUT_KEEP_LATEST, response.Policy)
	})

	t.Run("should get policy keepAll", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
cleanup:
  exclude: foo
retention:
  strategy: "keepAll"
  garbageCollectionCron: "0 * * * *"
`,
			},
		}, nil)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		response, err := svc.GetRetentionPolicy(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, backup.RetentionPolicy_RETENTION_POLICY_KEEP_ALL, response.Policy)
	})

	t.Run("should get policy keepLastSevenDays", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
cleanup:
  exclude: foo
retention:
  strategy: "keepLastSevenDays"
  garbageCollectionCron: "0 * * * *"
`,
			},
		}, nil)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		response, err := svc.GetRetentionPolicy(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, backup.RetentionPolicy_RETENTION_POLICY_KEEP_LAST_SEVEN_DAYS, response.Policy)
	})

	t.Run("should get policy keep7Days1Month1Quarter1Year", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
cleanup:
  exclude: foo
retention:
  strategy: "keep7Days1Month1Quarter1Year"
  garbageCollectionCron: "0 * * * *"
`,
			},
		}, nil)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		response, err := svc.GetRetentionPolicy(testCtx, nil)

		require.NoError(t, err)
		assert.Equal(t, backup.RetentionPolicy_RETENTION_POLICY_KEEP_LAST_7_DAYS_OLDEST_OF_1_MONTH_1_QUARTER_1_HALF_YEAR_1_YEAR, response.Policy)
	})

	t.Run("should fail to get empty policy with error", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(nil, assert.AnError)

		svc := &DefaultBackupService{
			componentClient: mComponentClient,
		}

		_, err := svc.GetRetentionPolicy(testCtx, nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get backup-operator component:")
	})
}

func Test_createBackups(t *testing.T) {
	t.Run("should create backup", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
		}

		backupClientMock.EXPECT().Create(testCtx, mock.Anything, metav1.CreateOptions{}).Return(&backupOne, nil)

		sut := DefaultBackupService{
			backupClient:  backupClientMock,
			restoreClient: nil,
		}

		// when
		_, err := sut.CreateBackup(testCtx, &backup.CreateBackupRequest{})
		// then
		require.NoError(t, err)
	})
	t.Run("should error on creating backup", func(t *testing.T) {
		// given
		testCtx := context.TODO()
		backupClientMock := newMockBackupInterface(t)

		backupClientMock.EXPECT().Create(testCtx, mock.Anything, metav1.CreateOptions{}).Return(nil, assert.AnError)

		sut := DefaultBackupService{
			backupClient:  backupClientMock,
			restoreClient: nil,
		}

		// when
		_, err := sut.CreateBackup(testCtx, &backup.CreateBackupRequest{})
		// then
		require.ErrorIs(t, err, assert.AnError)
	})
}

func TestIsDoguListMatching(t *testing.T) {
	t.Run("should return true when both lists contain the same dogus", func(t *testing.T) {
		version1 := "1.2.3"
		version2 := "4.5.6"
		dogu1 := v2.Dogu{Name: "test/1", Version: &version1}
		dogu2 := v2.Dogu{Name: "test/2", Version: &version2}

		annotationDogus := []annotationDogus{
			{Name: "test/1", Version: "1.2.3"},
			{Name: "test/2", Version: "4.5.6"},
		}

		bpDogus := []v2.Dogu{
			dogu1, dogu2,
		}

		svc := &DefaultBackupService{}

		ret := svc.isDoguListMatching(annotationDogus, bpDogus)
		assert.True(t, ret)
	})

	t.Run("should return false when both lists are not matching", func(t *testing.T) {
		version1 := "9.9.9"
		version2 := "6.6.6"
		dogu1 := v2.Dogu{Name: "test/1", Version: &version1}
		dogu2 := v2.Dogu{Name: "test/2", Version: &version2}

		annotationDogus := []annotationDogus{
			{Name: "test/1", Version: "1.2.3"},
			{Name: "test/2", Version: "4.5.6"},
		}

		bpDogus := []v2.Dogu{
			dogu1, dogu2,
		}

		svc := &DefaultBackupService{}

		ret := svc.isDoguListMatching(annotationDogus, bpDogus)
		assert.False(t, ret)
	})
}

func TestBackupStatus(t *testing.T) {
	t.Run("should return completed backup status", func(t *testing.T) {
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
			Status: backupV1.BackupStatus{
				Status: "completed",
			},
		}

		status := backupStatus(&backupOne)
		assert.Equal(t, backupOne.Status.Status, status)
	})
	t.Run("should return failed backup status", func(t *testing.T) {
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
			Status: backupV1.BackupStatus{
				Status: "failed",
			},
		}

		status := backupStatus(&backupOne)
		assert.Equal(t, backupOne.Status.Status, status)
	})
	t.Run("should return inProgress backup status", func(t *testing.T) {
		backupOne := backupV1.Backup{
			ObjectMeta: metav1.ObjectMeta{
				Name: "backup_one",
			},
			Status: backupV1.BackupStatus{
				Status: "inProgress",
			},
		}

		status := backupStatus(&backupOne)
		assert.Equal(t, backupOne.Status.Status, status)
	})
}
