package backup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getBackupTimeout(t *testing.T) {
	testCtx := context.Background()

	t.Run("should get timeout", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		cm := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "10"}}
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(cm, nil)

		timeout, err := getBackupTimeout(testCtx, mConfigmapClient)

		require.NoError(t, err)
		assert.Equal(t, 10, timeout)
	})

	t.Run("should fail without config map", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(nil, assert.AnError)

		timeout, err := getBackupTimeout(testCtx, mConfigmapClient)

		require.Error(t, err)
		assert.Equal(t, 0, timeout)
	})

	t.Run("should fail to parse value", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		cm := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "invalid"}}
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(cm, nil)

		timeout, err := getBackupTimeout(testCtx, mConfigmapClient)

		require.Error(t, err)
		assert.Equal(t, 0, timeout)
	})
}

func Test_setBackupTimeout(t *testing.T) {
	testCtx := context.Background()

	t.Run("should set timeout", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		cm1 := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "10"}}
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(cm1, nil)

		cm2 := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "20"}}
		mConfigmapClient.EXPECT().Update(testCtx, cm2, metav1.UpdateOptions{}).Return(cm2, nil)

		err := setBackupTimeout(testCtx, mConfigmapClient, 20)

		require.NoError(t, err)
	})
	t.Run("should fail without config map", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(nil, assert.AnError)

		err := setBackupTimeout(testCtx, mConfigmapClient, 10)

		require.Error(t, err)
	})
	t.Run("should fail on update", func(t *testing.T) {
		mConfigmapClient := newMockConfigmapClient(t)
		cm1 := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "10"}}
		mConfigmapClient.EXPECT().Get(testCtx, configMapName, metav1.GetOptions{}).Return(cm1, nil)

		cm2 := &corev1.ConfigMap{Data: map[string]string{retryTimeLimitKey: "20"}}
		mConfigmapClient.EXPECT().Update(testCtx, cm2, metav1.UpdateOptions{}).Return(cm2, assert.AnError)

		err := setBackupTimeout(testCtx, mConfigmapClient, 20)

		require.Error(t, err)
	})
}
