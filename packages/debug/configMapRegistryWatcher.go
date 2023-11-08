package debug

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

const tickerInterval = time.Second * 30

type defaultConfigMapRegistryWatcher struct {
	configMapInterface         configMapInterface
	configMapDebugModeRegistry debugModeRegistry
}

// NewDefaultConfigMapRegistryWatcher creates an instance of defaultConfigMapRegistryWatcher.
func NewDefaultConfigMapRegistryWatcher(configMapInterface configMapInterface, registry debugModeRegistry) *defaultConfigMapRegistryWatcher {
	return &defaultConfigMapRegistryWatcher{
		configMapInterface:         configMapInterface,
		configMapDebugModeRegistry: registry,
	}
}

// StartWatch checks if the disableAtTimestamp in the registry is after now and if yes it disables the debug mode.
func (w *defaultConfigMapRegistryWatcher) StartWatch(ctx context.Context) {
	go func() {
		w.doWatch(ctx)
	}()
}

func (w *defaultConfigMapRegistryWatcher) doWatch(ctx context.Context) {
	ticker := time.NewTicker(tickerInterval)
	logger := log.FromContext(ctx)

	for {
		select {
		case <-ticker.C:
			registryConfigMap, err := w.configMapInterface.Get(ctx, registryName, v1.GetOptions{})
			errorMsg := "watch debug mode registry error"
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				}
				logger.Error(fmt.Errorf("failed to get debug mode registry %s: %w", registryName, err), errorMsg)
				continue
			}

			if registryConfigMap.Data == nil {
				continue
			}

			timestamp, ok := registryConfigMap.Data[keyDisableAtTimestamp]
			if !ok {
				logger.Error(fmt.Errorf("failed to get debug mode disableAtTimestamp"), errorMsg)
				continue
			}

			disableAtTimestamp, err := time.Parse(timestampFormat, timestamp)
			if err != nil {
				logger.Error(fmt.Errorf("failed parse disableAtTimestamp %s: %w", timestamp, err), errorMsg)
				continue
			}

			logger.Info("disable debug mode registry")
			if time.Now().After(disableAtTimestamp) {
				err = w.configMapDebugModeRegistry.Disable(ctx)
				if err != nil {
					logger.Error(fmt.Errorf("failed to disable debug mode: %w", err), errorMsg)
				}
			}
		}
	}
}
