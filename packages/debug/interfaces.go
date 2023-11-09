package debug

import (
	"context"
	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type clusterClientSet interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type coreV1Interface interface {
	v1.CoreV1Interface
}

type configMapInterface interface {
	v1.ConfigMapInterface
}

type cesRegistry interface {
	registry.Registry
}

//nolint:unused
//goland:noinspection GoUnusedType
type configurationContext interface {
	registry.ConfigurationContext
}

//nolint:unused
//goland:noinspection GoUnusedType
type doguRegistry interface {
	registry.DoguRegistry
}

type debugModeRegistry interface {
	// Enable enables the debug mode in the registry.
	Enable(ctx context.Context, timer int32) error
	// Disable disables the debug mode registry.
	Disable(ctx context.Context) error
	// Status returns a boolean if the mode is enabled or disabled and if the status is enabled the timestamp where the
	// mode should be automatically disabled. If the mode is disabled the timestamp will be 0.
	Status(ctx context.Context) (isEnabled bool, DisableAtTimestamp int64, err error)
	// BackupDoguLogLevels saves all current log levels from all dogus.
	BackupDoguLogLevels(ctx context.Context) error
	// RestoreDoguLogLevels restores all backuped log levels from dogus.
	RestoreDoguLogLevels(ctx context.Context) error
}

type doguLogLevelRegistry interface {
	MarshalFromCesRegistryToString() (string, error)
	UnMarshalFromStringToCesRegistry(unmarshal string) error
}

type maintenanceModeSwitch interface {
	// ActivateMaintenanceMode activates the maintenance mode
	ActivateMaintenanceMode(title, text string) error
	// DeactivateMaintenanceMode deactivates the maintenance mode.
	DeactivateMaintenanceMode() error
}

type doguInterActor interface {
	StopAllDogus(ctx context.Context) error
	StartAllDogus(ctx context.Context) error
	SetLogLevelInAllDogus(logLevel string) error
}
