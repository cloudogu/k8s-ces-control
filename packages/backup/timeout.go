package backup

import (
	"context"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configMapName     = "k8s-backup-operator-backup-config"
	retryTimeLimitKey = "retryTimeLimit"
)

func getBackupTimeout(ctx context.Context, client configmapClient) (int, error) {
	cm, err := client.Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get configmap [%s]: %w", configMapName, err)
	}
	backupRetryTimeLimitStr := cm.Data[retryTimeLimitKey]

	retryLimit, err := strconv.Atoi(backupRetryTimeLimitStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert [%s]: %w", backupRetryTimeLimitStr, err)
	}
	return retryLimit, nil
}

func setBackupTimeout(ctx context.Context, client configmapClient, retryLimit int) error {
	cm, err := client.Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get configmap [%s]: %w", configMapName, err)
	}
	cm.Data[retryTimeLimitKey] = fmt.Sprintf("%d", retryLimit)

	_, err = client.Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update backup timeout: %w", err)
	}
	return nil
}
