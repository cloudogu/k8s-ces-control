package backup

import (
	"context"
	"testing"

	componentV1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		policy, err := getRetentionPolicy(testCtx, mComponentClient)

		require.NoError(t, err)
		assert.Equal(t, "removeAllButKeepLatest", policy)
	})

	t.Run("should get empty string as retention policy when no valuesYamlOverwrite exists", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(&componentV1.Component{
			Spec: componentV1.ComponentSpec{},
		}, nil)

		policy, err := getRetentionPolicy(testCtx, mComponentClient)

		require.NoError(t, err)
		assert.Equal(t, "", policy)
	})

	t.Run("should get empty string as retention policy when no value is configured", func(t *testing.T) {
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

		policy, err := getRetentionPolicy(testCtx, mComponentClient)

		require.NoError(t, err)
		assert.Equal(t, "", policy)
	})

	t.Run("should fail get retention policy for error unmarshalling yaml", func(t *testing.T) {
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

		_, err := getRetentionPolicy(testCtx, mComponentClient)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get configured retention policy: failed to unmarshal backup-operator config from valuesYamlOverwrite:")
	})

	t.Run("should fail get retention policy for error getting component", func(t *testing.T) {
		mComponentClient := newMockComponentClient(t)
		mComponentClient.EXPECT().Get(testCtx, "k8s-backup-operator", metav1.GetOptions{}).Return(nil, assert.AnError)

		_, err := getRetentionPolicy(testCtx, mComponentClient)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get backup-operator component:")
	})
}
