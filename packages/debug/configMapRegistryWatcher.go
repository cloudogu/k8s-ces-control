package debug

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-ces-control/generated/debug"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"time"
)

var tickerInterval = time.Second * 30

type defaultConfigMapRegistryWatcher struct {
	configMapInterface         configMapInterface
	configMapDebugModeRegistry debugModeRegistry
	debugModeService           debugModeServer
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

	for range ticker.C {
		err := w.checkDisableRegistry(ctx)
		if err != nil {
			logrus.Error(fmt.Errorf("watch debug mode registry: %w", err))
		}
	}
}

func (w *defaultConfigMapRegistryWatcher) checkDisableRegistry(ctx context.Context) error {
	registryConfigMap, err := w.configMapInterface.Get(ctx, registryName, v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to get debug mode registry %s: %w", registryName, err)
	}

	if registryConfigMap.Data == nil {
		return nil
	}

	timestamp, ok := registryConfigMap.Data[keyDisableAtTimestamp]
	if !ok {
		return nil
	}

	disableAtTimestamp, err := time.Parse(timestampFormat, timestamp)
	if err != nil {
		return fmt.Errorf("failed parse disableAtTimestamp %s: %w", timestamp, err)
	}

	logrus.Info("disable debug mode registry")
	if time.Now().After(disableAtTimestamp) {
		_, err = w.debugModeService.Disable(ctx, &debug.ToggleDebugModeRequest{})
		if err != nil {
			return fmt.Errorf("failed to disable debug mode: %w", err)
		}
	}

	return nil
}
