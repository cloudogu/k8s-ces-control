package debug

import (
	"context"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/ces-control-api/generated/types"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"github.com/cloudogu/k8s-registry-lib/config"
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

//nolint:unused
//goland:noinspection GoUnusedType
type doguRegistry interface {
	// GetCurrentOfAll retrieves the specs of all dogus' currently installed versions.
	GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error)
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
	// MarshalFromCesRegistryToString converts the log levels from the ces registry to a string
	MarshalFromCesRegistryToString(ctx context.Context) (string, error)
	// UnMarshalFromStringToCesRegistry writes the log level string to the ces registry.
	UnMarshalFromStringToCesRegistry(ctx context.Context, unmarshal string) error
}

type maintenanceModeSwitch interface {
	// ActivateMaintenanceMode activates the maintenance mode
	ActivateMaintenanceMode(ctx context.Context, title, text string) error
	// DeactivateMaintenanceMode deactivates the maintenance mode.
	DeactivateMaintenanceMode(ctx context.Context) error
}

type doguInterActor interface {
	// StopAllDogus stops all dogus.
	StopAllDogus(ctx context.Context) error
	// StartAllDogus starts all dogus.
	StartAllDogus(ctx context.Context) error
	// SetLogLevelInAllDogus sets the specified log level to all dogus.
	SetLogLevelInAllDogus(ctx context.Context, logLevel string) error
}

type debugModeServer interface {
	// Disable disables the debug mode.
	Disable(context.Context, *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error)
}

type doguConfigRepository interface {
	Get(context.Context, config.SimpleDoguName) (config.DoguConfig, error)
	Update(context.Context, config.DoguConfig) (config.DoguConfig, error)
}

type globalConfigRepository interface {
	Get(ctx context.Context) (config.GlobalConfig, error)
	Update(ctx context.Context, globalConfig config.GlobalConfig) (config.GlobalConfig, error)
}
