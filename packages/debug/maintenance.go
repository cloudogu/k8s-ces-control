package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

const maintenanceKey = "maintenance"

type defaultMaintenanceModeSwitch struct {
	globalConfig repository.GlobalConfigRepository
}

type maintenanceRegistryObject struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// NewDefaultMaintenanceModeSwitch creates a new instance of defaultMaintenanceModeSwitch.
func NewDefaultMaintenanceModeSwitch(globalConfigRepository repository.GlobalConfigRepository) *defaultMaintenanceModeSwitch {
	return &defaultMaintenanceModeSwitch{
		globalConfig: globalConfigRepository,
	}
}

// ActivateMaintenanceMode activates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) ActivateMaintenanceMode(ctx context.Context, title, text string) error {
	gConfig, err := d.globalConfig.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get globalConfig: %w", err)
	}
	value := maintenanceRegistryObject{
		Title: title,
		Text:  text,
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance globalConfig value object [%+v]: %w", value, err)
	}

	maintenanceJsonStr := string(marshal)
	newGlobalConfig, err := gConfig.Set(maintenanceKey, config.Value(maintenanceJsonStr))
	if err != nil {
		return err
	}
	_, err = d.globalConfig.Update(ctx, config.GlobalConfig{newGlobalConfig})
	if err != nil {
		return fmt.Errorf("failed to set value [%s] with key %s: %w", maintenanceJsonStr, maintenanceKey, err)
	}

	return nil
}

// DeactivateMaintenanceMode deactivates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) DeactivateMaintenanceMode(ctx context.Context) error {
	gConfig, err := d.globalConfig.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get globalConfig: %w", err)
	}
	newGlobalConfig := gConfig.Delete(maintenanceKey)

	_, err = d.globalConfig.Update(ctx, config.GlobalConfig{newGlobalConfig})
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", maintenanceKey, err)
	}

	_, err = d.globalConfig.Update(ctx, config.GlobalConfig{newGlobalConfig})
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", maintenanceKey, err)
	}
	return nil
}
