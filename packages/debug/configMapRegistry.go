package debug

import (
	"context"
	"fmt"
	cesregistry "github.com/cloudogu/cesapp-lib/registry"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"strconv"
	"time"
)

const (
	registryName          = "debug-mode-registry"
	keyDebugModeEnabled   = "enabled"
	keyDisableAtTimestamp = "disable-at-timestamp"
	keyDoguLogLevel       = "dogus"
)

var maxThirtySecondsBackoff = wait.Backoff{
	Duration: 1500 * time.Millisecond,
	Factor:   1.5,
	Jitter:   0,
	Steps:    9999,
	Cap:      30 * time.Second,
}

type configMapDebugModeRegistry struct {
	cesRegistry          cesregistry.Registry
	doguLogLevelRegistry doguLogLevelRegistry
	configMapInterface   configMapInterface
	namespace            string
}

// NewConfigMapDebugModeRegistry creates a new instance of configMapDebugModeRegistry.
func NewConfigMapDebugModeRegistry(cesRegistry cesregistry.Registry, clusterClientSet clusterClientSet, namespace string) *configMapDebugModeRegistry {
	return &configMapDebugModeRegistry{
		cesRegistry:          cesRegistry,
		configMapInterface:   clusterClientSet.CoreV1().ConfigMaps(namespace),
		namespace:            namespace,
		doguLogLevelRegistry: NewDoguLogLevelRegistryMap(cesRegistry),
	}
}

// Enable writes `enabled: true` in the registry.
func (c *configMapDebugModeRegistry) Enable(ctx context.Context, endTimestamp int64) error {
	cm, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	cm.Data[keyDebugModeEnabled] = "true"
	cm.Data[keyDisableAtTimestamp] = strconv.FormatInt(endTimestamp, 10)

	return c.updateConfigMap(ctx, cm)
}

func (c *configMapDebugModeRegistry) getRegistry(ctx context.Context) (*corev1.ConfigMap, error) {
	cm, err := c.configMapInterface.Get(ctx, registryName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return c.createRegistry(ctx)
		} else {
			return nil, wrapRegistryError(c.namespace, registryName, "get", err)
		}
	}

	return cm, nil
}

func (c *configMapDebugModeRegistry) createRegistry(ctx context.Context) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: registryName, Namespace: c.namespace},
	}

	create, err := c.configMapInterface.Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		return nil, wrapRegistryError(cm.Namespace, cm.Name, "create", err)
	}

	return create, nil
}

func (c *configMapDebugModeRegistry) updateConfigMap(ctx context.Context, cm *corev1.ConfigMap) error {
	err := retryOnConflict(func() error {
		registry, getErr := c.getRegistry(ctx)
		if getErr != nil {
			return getErr
		}

		registry.Data = cm.Data

		_, err := c.configMapInterface.Update(ctx, registry, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return wrapRegistryError(cm.Namespace, cm.Name, "update", err)
	}

	return nil
}

func retryOnConflict(fn func() error) error {
	return retry.RetryOnConflict(maxThirtySecondsBackoff, fn)
}

func wrapRegistryError(namespace, name, verb string, err error) error {
	return fmt.Errorf("failed to %s config map %s/%s: %w", verb, namespace, name, err)
}

var deleteRetriable = func(err error) bool {
	return err != nil
}

// Disable writes `enabled: false` in the registry.
func (c *configMapDebugModeRegistry) Disable(ctx context.Context) error {
	err := retry.OnError(maxThirtySecondsBackoff, deleteRetriable, func() error {
		err := c.configMapInterface.Delete(ctx, registryName, metav1.DeleteOptions{})
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	})

	if err != nil {
		return fmt.Errorf("failed to delete debug mode registry configmap %s: %w", registryName, err)
	}

	return nil
}

// Status parses the fields `enabled` and `delete_at_timestamp` and returns them.
func (c *configMapDebugModeRegistry) Status(ctx context.Context) (isEnabled bool, disableAtTimestamp int64, err error) {
	registry, err := c.getRegistry(ctx)
	if err != nil {
		return false, 0, err
	}

	isEnabledStr := registry.Data[keyDebugModeEnabled]
	isEnabled, err = strconv.ParseBool(isEnabledStr)
	if err != nil {
		return false, 0, fmt.Errorf("failed to parse bool %s: %w", isEnabledStr, err)
	}

	disableAtTimestampStr := registry.Data[keyDisableAtTimestamp]
	disableAtTimestamp, err = strconv.ParseInt(disableAtTimestampStr, 10, 32)
	if err != nil {
		return false, 0, fmt.Errorf("failed to parse %s to uint: %w", disableAtTimestampStr, err)
	}

	return
}

func (c *configMapDebugModeRegistry) BackupDoguLogLevels(ctx context.Context) error {
	registry, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	newRegistry, err := c.doguLogLevelRegistry.MarshalToString()
	if err != nil {
		// TODO Log errors here?
		return fmt.Errorf("failed to renew dogu log level registry: %w", err)
	}

	registry.Data[keyDoguLogLevel] = newRegistry

	return c.updateConfigMap(ctx, registry)
}

func (c *configMapDebugModeRegistry) RestoreDoguLogLevels(ctx context.Context) error {
	registry, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	doguLogLevelData := registry.Data[keyDoguLogLevel]
	doguLogLevelReg, err := c.doguLogLevelRegistry.UnMarshalFromString(doguLogLevelData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dogu log level registry from [%s]: %w", doguLogLevelData, err)
	}

	err = doguLogLevelReg.RestoreToCesRegistry()
	if err != nil {
		// TODO log errors here?
		return fmt.Errorf("failed to restore dogu log level to ces registry: %w", err)
	}

	return nil
}