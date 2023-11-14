package debug

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

const keyDoguConfigLogLevel = "logging/root"

type doguLogLevelYamlRegistryMap struct {
	registry map[string]string
}

// NewDoguLogLevelRegistryMap creates an instance of doguLogLevelYamlRegistryMap.
func NewDoguLogLevelRegistryMap() *doguLogLevelYamlRegistryMap {
	return &doguLogLevelYamlRegistryMap{
		registry: map[string]string{},
	}
}

// MarshalFromCesRegistryToString marshals the registry to yaml string.
func (d *doguLogLevelYamlRegistryMap) MarshalFromCesRegistryToString(cesRegistry cesRegistry) (string, error) {
	d.registry = map[string]string{}

	allDogus, err := cesRegistry.DoguRegistry().GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to get all dogus: %w", err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig := cesRegistry.DoguConfig(dogu.GetSimpleName())
		exists, existsErr := doguConfig.Exists(keyDoguConfigLogLevel)
		if existsErr != nil {
			multiError = errors.Join(multiError, existsErr)
		}

		if !exists {
			d.registry[dogu.GetSimpleName()] = ""
			continue
		}

		logLevel, getErr := doguConfig.Get(keyDoguConfigLogLevel)
		if getErr != nil {
			multiError = errors.Join(multiError, getErr)
			continue
		}

		d.registry[dogu.GetSimpleName()] = logLevel
	}

	out, err := yaml.Marshal(d.registry)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registry: %w", err)
	}

	return string(out), multiError
}

// UnMarshalFromStringToCesRegistry unmarshal a map as yaml string to ces registry.
func (d *doguLogLevelYamlRegistryMap) UnMarshalFromStringToCesRegistry(cesRegistry cesRegistry, unmarshal string) error {
	out := map[string]string{}
	err := yaml.Unmarshal([]byte(unmarshal), out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dogu log level from string [%s]: %w", unmarshal, err)
	}
	d.registry = out

	return d.restoreToCesRegistry(cesRegistry)
}

// RestoreToCesRegistry writes all log levels to the ces registry.
func (d *doguLogLevelYamlRegistryMap) restoreToCesRegistry(cesRegistry cesRegistry) error {
	var multiError error
	for dogu, level := range d.registry {
		doguConfig := cesRegistry.DoguConfig(dogu)
		// If the dogu had no log level it is defined as an empty string in the registry.
		// In this case we have to delete the entry.
		if level == "" {
			deleteErr := doguConfig.Delete(keyDoguConfigLogLevel)
			if deleteErr != nil {
				multiError = errors.Join(multiError, deleteErr)
			}
			continue
		}

		err := doguConfig.Set(keyDoguConfigLogLevel, level)
		if err != nil {
			multiError = errors.Join(multiError, err)
			continue
		}
	}

	return multiError
}
