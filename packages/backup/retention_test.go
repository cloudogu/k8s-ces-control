package backup

import (
	"context"
	"testing"

	componentV1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getRetentionPolicy(t *testing.T) {
	testCtx := context.Background()

	t.Run("should get retention policy", func(t *testing.T) {
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

		policy, err := getRetentionPolicy(testCtx, mComponentClient, nil)

		require.NoError(t, err)
		assert.Equal(t, "removeAllButKeepLatest", policy)
	})

	t.Run("should get default string as retention policy when no valuesYamlOverwrite exists", func(t *testing.T) {
		mCronJobClient := newMockCronJobClient(t)
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{},
		}, nil)
		cronJob := &v1.CronJob{
			Spec: v1.CronJobSpec{
				JobTemplate: v1.JobTemplateSpec{
					Spec: v1.JobSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Args: []string{"--strategy=keepAll"},
									},
								},
							},
						},
					},
				},
			},
		}
		mCronJobClient.EXPECT().Get(testCtx, backupGarbageCollectorCronJobName, metav1.GetOptions{}).Return(cronJob, nil)

		policy, err := getRetentionPolicy(testCtx, mComponentClient, mCronJobClient)

		require.NoError(t, err)
		assert.Equal(t, "keepAll", policy)
	})

	t.Run("should get default string as retention policy when no value is configured", func(t *testing.T) {
		mCronJobClient := newMockCronJobClient(t)
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
cleanup:
  exclude: foo
retention:
  garbageCollectionCron: "0 * * * *"
`,
			},
		}, nil)

		cronJob := &v1.CronJob{
			Spec: v1.CronJobSpec{
				JobTemplate: v1.JobTemplateSpec{
					Spec: v1.JobSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Args: []string{"--strategy=keepAll"},
									},
								},
							},
						},
					},
				},
			},
		}
		mCronJobClient.EXPECT().Get(testCtx, backupGarbageCollectorCronJobName, metav1.GetOptions{}).Return(cronJob, nil)

		policy, err := getRetentionPolicy(testCtx, mComponentClient, mCronJobClient)

		require.NoError(t, err)
		assert.Equal(t, "keepAll", policy)
	})

	t.Run("should fail get retention policy for error unmarshalling yaml", func(t *testing.T) {
		mCronJobClient := newMockCronJobClient(t)
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{
				ValuesYamlOverwrite: `
key: value
  invalid: indentation
    another: level
`,
			},
		}, nil)

		_, err := getRetentionPolicy(testCtx, mComponentClient, mCronJobClient)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get configured retention policy: failed to unmarshal backup-operator config from valuesYamlOverwrite:")
	})

	t.Run("should fail get retention policy for error getting component", func(t *testing.T) {
		mCronJobClient := newMockCronJobClient(t)
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(nil, assert.AnError)

		_, err := getRetentionPolicy(testCtx, mComponentClient, mCronJobClient)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get backup-operator component:")
	})
	t.Run("should fail get retention policy for error getting default from cronjob", func(t *testing.T) {
		mCronJobClient := newMockCronJobClient(t)
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{},
		}, nil)
		mCronJobClient.EXPECT().Get(testCtx, backupGarbageCollectorCronJobName, metav1.GetOptions{}).Return(nil, assert.AnError)

		policy, err := getRetentionPolicy(testCtx, mComponentClient, mCronJobClient)

		require.NoError(t, err)
		assert.Equal(t, "", policy)
	})
}
