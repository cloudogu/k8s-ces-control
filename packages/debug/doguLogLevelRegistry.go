package debug

import (
	"context"
	"errors"
	"fmt"
	dogu2 "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/config"
	"gopkg.in/yaml.v3"
)

const keyDoguConfigLogLevel = "logging/root"

type doguLogLevelYamlRegistryMap struct {
	doguConfigRepository doguConfigRepository
	doguReg              doguDescriptorGetter
	logLevelRegistryMap  map[string]string
}

// NewDoguLogLevelRegistryMap creates an instance of doguLogLevelYamlRegistryMap.
func NewDoguLogLevelRegistryMap(doguConfig doguConfigRepository, doguReg doguDescriptorGetter) *doguLogLevelYamlRegistryMap {
	return &doguLogLevelYamlRegistryMap{
		doguConfigRepository: doguConfig,
		doguReg:              doguReg,
		logLevelRegistryMap:  map[string]string{},
	}
}

// MarshalFromCesRegistryToString marshals the registry to yaml string.
func (d *doguLogLevelYamlRegistryMap) MarshalFromCesRegistryToString(ctx context.Context) (string, error) {
	allDogus, err := d.doguReg.GetCurrentOfAll(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get all dogus: %w", err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig, _ := d.doguConfigRepository.Get(ctx, dogu2.SimpleName(dogu.GetSimpleName()))
		logLevel, exists := doguConfig.Get(keyDoguConfigLogLevel)

		if !exists {
			d.logLevelRegistryMap[dogu.GetSimpleName()] = ""
			continue
		}

		d.logLevelRegistryMap[dogu.GetSimpleName()] = string(logLevel)
	}

	out, err := yaml.Marshal(d.logLevelRegistryMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registry: %w", err)
	}

	return string(out), multiError
}

// UnMarshalFromStringToCesRegistry unmarshal a map as yaml string to ces registry.
func (d *doguLogLevelYamlRegistryMap) UnMarshalFromStringToCesRegistry(ctx context.Context, unmarshal string) error {
	out := map[string]string{}
	err := yaml.Unmarshal([]byte(unmarshal), out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dogu log level from string [%s]: %w", unmarshal, err)
	}
	d.logLevelRegistryMap = out

	return d.restoreToCesRegistry(ctx)
}

// RestoreToCesRegistry writes all log levels to the ces registry.
func (d *doguLogLevelYamlRegistryMap) restoreToCesRegistry(ctx context.Context) error {
	var multiError error
	for dogu, level := range d.logLevelRegistryMap {
		doguConfig, _ := d.doguConfigRepository.Get(ctx, dogu2.SimpleName(dogu))
		// If the dogu had no log level it is defined as an empty string in the registry.
		// In this case we have to delete the entry.
		if level == "" {
			newDoguConfig := doguConfig.Delete(keyDoguConfigLogLevel)
			doguConfig.Config = newDoguConfig
			_, err := d.doguConfigRepository.Update(ctx, doguConfig)
			if err != nil {
				multiError = errors.Join(multiError, getDoguConfigUpdateError(dogu, err))
			}
			continue
		}
		newDoguConfig, err := doguConfig.Set(keyDoguConfigLogLevel, config.Value(level))
		if err != nil {
			multiError = errors.Join(multiError, err)
			continue
		}
		doguConfig.Config = newDoguConfig
		_, err = d.doguConfigRepository.Update(ctx, doguConfig)
		if err != nil {
			multiError = errors.Join(multiError, getDoguConfigUpdateError(dogu, err))
		}
		continue
	}

	return multiError
}

func getDoguConfigUpdateError(dogu string, err error) error {
	return fmt.Errorf("failed to update dogu config for dogu %s: %w", dogu, err)
}
