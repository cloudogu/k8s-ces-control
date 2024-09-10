package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudogu/k8s-registry-lib/config"
)

const maintenanceKey = "maintenance"

type defaultMaintenanceModeSwitch struct {
	globalConfigRepo globalConfigRepository
}

type maintenanceRegistryObject struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// NewDefaultMaintenanceModeSwitch creates a new instance of defaultMaintenanceModeSwitch.
func NewDefaultMaintenanceModeSwitch(globalConfigRepository globalConfigRepository) *defaultMaintenanceModeSwitch {
	return &defaultMaintenanceModeSwitch{
		globalConfigRepo: globalConfigRepository,
	}
}

// ActivateMaintenanceMode activates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) ActivateMaintenanceMode(ctx context.Context, title, text string) error {
	value := maintenanceRegistryObject{
		Title: title,
		Text:  text,
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance globalConfigRepo value object [%+v]: %w", value, err)
	}

	maintenanceJsonStr := string(marshal)

	globalConfig, err := d.globalConfigRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get global config: %w", err)
	}

	updatedGlobalConfig, err := globalConfig.Set(maintenanceKey, config.Value(maintenanceJsonStr))
	if err != nil {
		return fmt.Errorf("failed to set global config key %q value %q: %w", maintenanceKey, maintenanceJsonStr, err)
	}

	_, err = d.globalConfigRepo.Update(ctx, config.GlobalConfig{Config: updatedGlobalConfig})
	if err != nil {
		return fmt.Errorf("failed update global config for key %q value %s: %w", maintenanceKey, maintenanceJsonStr, err)
	}

	return nil
}

// DeactivateMaintenanceMode deactivates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) DeactivateMaintenanceMode(ctx context.Context) error {
	gConfig, err := d.globalConfigRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get globalConfig: %w", err)
	}

	newGlobalConfig := gConfig.Delete(maintenanceKey)

	_, err = d.globalConfigRepo.Update(ctx, config.GlobalConfig{Config: newGlobalConfig})
	if err != nil {
		return fmt.Errorf("failed to update global config for key %q: %w", maintenanceKey, err)
	}

	return nil
}
