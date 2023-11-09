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
	timestampFormat       = time.RFC822
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
func (c *configMapDebugModeRegistry) Enable(ctx context.Context, timerInMinutes int32) error {
	cm, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	timerDuration := time.Duration(timerInMinutes)
	disableAtTimestamp := time.Now().Add(time.Minute * timerDuration)
	cm.Data[keyDisableAtTimestamp] = disableAtTimestamp.Format(timestampFormat)
	cm.Data[keyDebugModeEnabled] = "true"

	return c.updateConfigMap(ctx, cm)
}

func (c *configMapDebugModeRegistry) getRegistry(ctx context.Context) (*corev1.ConfigMap, error) {
	return c.createRegistryIfNotFound(ctx, false)
}

func (c *configMapDebugModeRegistry) createRegistryIfNotFound(ctx context.Context, ignoreNotFound bool) (*corev1.ConfigMap, error) {
	cm, err := c.configMapInterface.Get(ctx, registryName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) && !ignoreNotFound {
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
	registry, err := c.createRegistryIfNotFound(ctx, true)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, 0, nil
		}
		return false, 0, err
	}

	isEnabled, err = isRegistryEnabled(registry)
	if err != nil {
		return false, 0, err
	}

	disableAtTimestamp, err = getDisableAtTimeStamp(registry)
	if err != nil {
		return false, 0, err
	}

	return
}

func getDisableAtTimeStamp(registry *corev1.ConfigMap) (int64, error) {
	if registry.Data == nil {
		return 0, fmt.Errorf("registry %s is not initialized", registry.Name)
	}

	disableAtTimestampStr, ok := registry.Data[keyDisableAtTimestamp]

	if !ok {
		return 0, nil
	}

	timeDisableAt, err := time.Parse(timestampFormat, disableAtTimestampStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time from disableAtTimestampStr %s: %w", disableAtTimestampStr, err)
	}

	return timeDisableAt.UnixMilli(), nil
}

func isRegistryEnabled(registry *corev1.ConfigMap) (isEnabled bool, err error) {
	if registry.Data == nil {
		return false, fmt.Errorf("registry %s is not initialized", registry.Name)
	}

	isEnabledStr, ok := registry.Data[keyDebugModeEnabled]
	if !ok {
		return false, nil
	}

	isEnabled, err = strconv.ParseBool(isEnabledStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse bool %s: %w", isEnabledStr, err)
	}

	return isEnabled, nil
}

func (c *configMapDebugModeRegistry) BackupDoguLogLevels(ctx context.Context) error {
	registry, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	enabled, err := isRegistryEnabled(registry)
	if err != nil {
		return err
	}

	if !enabled {
		return registryNotEnabledError()
	}

	newRegistry, err := c.doguLogLevelRegistry.MarshalFromCesRegistryToString()
	if err != nil {
		return fmt.Errorf("failed to renew dogu log level registry: %w", err)
	}

	registry.Data[keyDoguLogLevel] = newRegistry

	return c.updateConfigMap(ctx, registry)
}

func registryNotEnabledError() error {
	return fmt.Errorf("registry is not enabled")
}

func (c *configMapDebugModeRegistry) RestoreDoguLogLevels(ctx context.Context) error {
	registry, err := c.getRegistry(ctx)
	if err != nil {
		return err
	}

	enabled, err := isRegistryEnabled(registry)
	if err != nil {
		return err
	}

	if !enabled {
		return registryNotEnabledError()
	}

	doguLogLevelData, ok := registry.Data[keyDoguLogLevel]

	if !ok {
		return fmt.Errorf("missing registry key %s", keyDoguLogLevel)
	}

	err = c.doguLogLevelRegistry.UnMarshalFromStringToCesRegistry(doguLogLevelData)
	if err != nil {
		return fmt.Errorf("failed to restore dogu log level [%s]: %w", doguLogLevelData, err)
	}

	return nil
}
