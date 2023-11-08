package debug

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

const keyDoguConfigLogLevel = "logging/root"

type doguLogLevelYamlRegistryMap struct {
	cesRegistry cesRegistry
	registry    map[string]string
}

// NewDoguLogLevelRegistryMap creates an instance of doguLogLevelYamlRegistryMap.
func NewDoguLogLevelRegistryMap(registry cesRegistry) *doguLogLevelYamlRegistryMap {
	return &doguLogLevelYamlRegistryMap{
		cesRegistry: registry,
		registry:    map[string]string{},
	}
}

// MarshalFromCesRegistryToString marshals the registry to yaml string.
func (d *doguLogLevelYamlRegistryMap) MarshalFromCesRegistryToString() (string, error) {
	d.registry = map[string]string{}

	allDogus, err := d.cesRegistry.DoguRegistry().GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to get all dogus: %w", err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig := d.cesRegistry.DoguConfig(dogu.GetSimpleName())
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
func (d *doguLogLevelYamlRegistryMap) UnMarshalFromStringToCesRegistry(unmarshal string) error {
	out := map[string]string{}
	err := yaml.Unmarshal([]byte(unmarshal), out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dogu log level from string [%s]: %w", unmarshal, err)
	}
	d.registry = out

	return d.restoreToCesRegistry()
}

// RestoreToCesRegistry writes all log levels to the ces registry.
func (d *doguLogLevelYamlRegistryMap) restoreToCesRegistry() error {
	var multiError error
	for dogu, level := range d.registry {
		doguConfig := d.cesRegistry.DoguConfig(dogu)
		err := doguConfig.Set(keyDoguConfigLogLevel, level)
		if err != nil {
			multiError = errors.Join(multiError, err)
			continue
		}
	}

	return multiError
}
