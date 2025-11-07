package backup

import (
	"context"
	"fmt"

	v1 "github.com/cloudogu/k8s-backup-lib/api/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const backupScheduleName = "ces-schedule"

func getBackupSchedule(ctx context.Context, client backupScheduleClient) (string, error) {
	schedule, err := client.Get(ctx, backupScheduleName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get backup schedule: %w", err)
	}

	return schedule.Spec.Schedule, nil
}

func setBackupSchedule(ctx context.Context, client backupScheduleClient, schedule string) error {
	backupSchedule, err := client.Get(ctx, backupScheduleName, metav1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return createSchedule(ctx, client, schedule)
		} else {
			return fmt.Errorf("failed to get existing backup schedule: %w", err)
		}
	}

	return updateSchedule(ctx, client, backupSchedule, schedule)
}

func createSchedule(ctx context.Context, client backupScheduleClient, schedule string) error {
	backupSchedule := &v1.BackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name: backupScheduleName,
		},
		Spec: v1.BackupScheduleSpec{
			Schedule: schedule,
		},
	}

	_, err := client.Create(ctx, backupSchedule, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create backup schedule: %w", err)
	}

	return nil
}

func updateSchedule(ctx context.Context, client backupScheduleClient, backupSchedule *v1.BackupSchedule, schedule string) error {
	backupSchedule.Spec.Schedule = schedule

	_, err := client.Update(ctx, backupSchedule, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update backup schedule: %w", err)
	}

	return nil
}
