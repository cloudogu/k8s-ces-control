package backup

import (
	"context"
	"fmt"

	"github.com/cloudogu/k8s-component-lib/api/v1"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// retentionPolicyName is an enum identifying a retention policy.
type retentionPolicyName string

const (
	keepAllPolicy                                           retentionPolicyName = "keepAll"
	removeAllButKeepLatestPolicy                            retentionPolicyName = "removeAllButKeepLatest"
	keepLastSevenDaysPolicy                                 retentionPolicyName = "keepLastSevenDays"
	keepLast7DaysOldestOf1Month1Quarter1HalfYear1YearPolicy retentionPolicyName = "keep7Days1Month1Quarter1Year"
)

const backupOperatorComponentName = "k8s-backup-operator"

type BackupOperatorConfig struct {
	Retention struct {
		Strategy string `yaml:"strategy"`
	} `yaml:"retention"`
}

func getRetentionPolicy(ctx context.Context, client componentClient) (string, error) {
	backupOpComponent, err := client.Get(ctx, backupOperatorComponentName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get backup-operator component: %w", err)
	}

	policy, err := getConfiguredRetentionPolicy(backupOpComponent)
	if err != nil {
		return "", fmt.Errorf("failed to get configured retention policy: %w", err)
	}

	return policy, nil
}

func getConfiguredRetentionPolicy(backupOpComponent *v1.Component) (string, error) {
	yamlString := backupOpComponent.Spec.ValuesYamlOverwrite

	if yamlString == "" {
		return "", nil
	}

	var cfg BackupOperatorConfig
	if err := yaml.Unmarshal([]byte(yamlString), &cfg); err != nil {
		return "", fmt.Errorf("failed to unmarshal backup-operator config from valuesYamlOverwrite: %w", err)
	}

	return cfg.Retention.Strategy, nil
}
