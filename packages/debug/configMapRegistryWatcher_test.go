package debug

import (
	"github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
	"time"
)

func Test_defaultConfigMapRegistryWatcher_StartWatch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		disableAtTimestamp := time.Now().Add(time.Second * -2)
		disableAtTimestampStrFormat := disableAtTimestamp.Format(time.RFC822)
		registryCm := &v1.ConfigMap{Data: map[string]string{"disable-at-timestamp": disableAtTimestampStrFormat}}

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(registryCm, nil)
		debugModeService := newMockDebugModeServer(t)
		debugModeService.EXPECT().Disable(testCtx, &maintenance.ToggleDebugModeRequest{}).Return(nil, nil)

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
			debugModeService:   debugModeService,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil if registry is not found", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on error getting debug config map", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(nil, assert.AnError)

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get debug mode registry debug-mode-registry")
	})

	t.Run("should return nil on empty registry", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		registryCm := &v1.ConfigMap{}

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(registryCm, nil)

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil on missing key", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		registryCm := &v1.ConfigMap{Data: map[string]string{}}

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(registryCm, nil)

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return error on invalid timestamp format", func(t *testing.T) {
		// given
		oldTickerInterval := tickerInterval
		tickerInterval = time.Millisecond
		defer func() { tickerInterval = oldTickerInterval }()

		registryCm := &v1.ConfigMap{Data: map[string]string{"disable-at-timestamp": "invalid timestamp"}}

		configMapMock := newMockConfigMapInterface(t)
		configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(registryCm, nil)

		sut := defaultConfigMapRegistryWatcher{
			configMapInterface: configMapMock,
		}

		// when
		err := sut.checkDisableRegistry(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed parse disableAtTimestamp")
	})

	t.Run("should return error on disable error", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// given
			oldTickerInterval := tickerInterval
			tickerInterval = time.Millisecond
			defer func() { tickerInterval = oldTickerInterval }()

			disableAtTimestamp := time.Now().Add(time.Second * -2)
			disableAtTimestampStrFormat := disableAtTimestamp.Format(time.RFC822)
			registryCm := &v1.ConfigMap{Data: map[string]string{"disable-at-timestamp": disableAtTimestampStrFormat}}

			configMapMock := newMockConfigMapInterface(t)
			configMapMock.EXPECT().Get(testCtx, "debug-mode-registry", v12.GetOptions{}).Return(registryCm, nil)
			debugModeServiceMock := newMockDebugModeServer(t)
			debugModeServiceMock.EXPECT().Disable(testCtx, &maintenance.ToggleDebugModeRequest{}).Return(nil, assert.AnError)

			sut := defaultConfigMapRegistryWatcher{
				configMapInterface: configMapMock,
				debugModeService:   debugModeServiceMock,
			}

			// when
			err := sut.checkDisableRegistry(testCtx)

			// then
			require.Error(t, err)
			assert.ErrorIs(t, err, assert.AnError)
			assert.ErrorContains(t, err, "failed to disable debug mode")
		})
	})
}
