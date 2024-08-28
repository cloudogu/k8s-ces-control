package debug

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-registry-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/repository"
)

type LogLevelRepository struct {
	repository repository.DoguConfigRepository
}

func NewLogLevelRepository(configRepository repository.DoguConfigRepository) *LogLevelRepository {
	return &LogLevelRepository{repository: configRepository}
}

type DoguRegistry struct {
	versionRegistry      dogu.DoguVersionRegistry
	descriptorRepository dogu.LocalDoguDescriptorRepository
}

func NewDoguRegistry(versionRegistry dogu.DoguVersionRegistry, descriptorRepository dogu.LocalDoguDescriptorRepository) *DoguRegistry {
	return &DoguRegistry{
		versionRegistry:      versionRegistry,
		descriptorRepository: descriptorRepository,
	}
}

func (r *DoguRegistry) GetCurrent(ctx context.Context, simpleDoguName string) (*core.Dogu, error) {
	current, err := r.versionRegistry.GetCurrent(ctx, dogu.SimpleDoguName(simpleDoguName))
	if err != nil {
		return nil, err
	}
	get, err := r.descriptorRepository.Get(ctx, current)
	if err != nil {
		return nil, err
	}
	return get, nil
}
func (r *DoguRegistry) GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error) {
	allCurrentDoguVersions, err := r.versionRegistry.GetCurrentOfAll(ctx)
	if err != nil {
		return nil, err
	}
	allCurrentDogus := []*core.Dogu{}
	for _, doguVersion := range allCurrentDoguVersions {
		dogu, err := r.descriptorRepository.Get(ctx, doguVersion)
		if err != nil {
			return nil, err
		}
		allCurrentDogus = append(allCurrentDogus, dogu)
	}
	return allCurrentDogus, nil
}
