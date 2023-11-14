package debug

import (
	"encoding/json"
	"fmt"
	"github.com/cloudogu/cesapp-lib/registry"
)

const maintenanceKey = "maintenance"

type defaultMaintenanceModeSwitch struct {
	globalConfig registry.ConfigurationContext
}

type maintenanceRegistryObject struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// NewDefaultMaintenanceModeSwitch creates a new instance of defaultMaintenanceModeSwitch.
func NewDefaultMaintenanceModeSwitch(globalConfig registry.ConfigurationContext) *defaultMaintenanceModeSwitch {
	return &defaultMaintenanceModeSwitch{
		globalConfig: globalConfig,
	}
}

// ActivateMaintenanceMode activates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) ActivateMaintenanceMode(title, text string) error {
	value := maintenanceRegistryObject{
		Title: title,
		Text:  text,
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal maintenance globalConfig value object [%+v]: %w", value, err)
	}

	maintenanceJsonStr := string(marshal)
	err = d.globalConfig.Set(maintenanceKey, maintenanceJsonStr)
	if err != nil {
		return fmt.Errorf("failed to set value [%s] with key %s: %w", maintenanceJsonStr, maintenanceKey, err)
	}

	return nil
}

// DeactivateMaintenanceMode deactivates the maintenance mode.
func (d *defaultMaintenanceModeSwitch) DeactivateMaintenanceMode() error {
	err := d.globalConfig.Delete(maintenanceKey)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", maintenanceKey, err)
	}
	return nil
}
