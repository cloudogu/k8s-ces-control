package debug

import (
	"fmt"
	hashicorperror "github.com/hashicorp/go-multierror"
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

// MarshalToString marshals the registry to yaml string.
func (d *doguLogLevelYamlRegistryMap) MarshalToString() (string, error) {
	d.registry = map[string]string{}

	allDogus, err := d.cesRegistry.DoguRegistry().GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to get all dogus: %w", err)
	}

	// TODO errors.Join
	var multierror *hashicorperror.Error
	for _, dogu := range allDogus {
		doguConfig := d.cesRegistry.DoguConfig(dogu.GetSimpleName())
		logLevel, getErr := doguConfig.Get(keyDoguConfigLogLevel)
		if getErr != nil {
			multierror = hashicorperror.Append(multierror, getErr)
			continue
		}

		d.registry[dogu.GetSimpleName()] = logLevel
	}

	out, err := yaml.Marshal(d.registry)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registry: %w", err)
	}

	return string(out), multierror.ErrorOrNil()
}

// UnMarshalFromString unmarshal a map as yaml string to doguLogLevelYamlRegistryMap.
func (d *doguLogLevelYamlRegistryMap) UnMarshalFromString(unmarshal string) (*doguLogLevelYamlRegistryMap, error) {
	out := map[string]string{}
	err := yaml.Unmarshal([]byte(unmarshal), out)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dogu log level from string [%s]: %w", unmarshal, err)
	}
	d.registry = out

	return d, nil
}

// RestoreToCesRegistry writes all log levels to the ces registry.
func (d *doguLogLevelYamlRegistryMap) RestoreToCesRegistry() error {
	var multiError *hashicorperror.Error
	for dogu, level := range d.registry {
		doguConfig := d.cesRegistry.DoguConfig(dogu)
		err := doguConfig.Set(keyDoguConfigLogLevel, level)
		if err != nil {
			return hashicorperror.Append(multiError, err)
		}
	}

	return multiError.ErrorOrNil()
}
