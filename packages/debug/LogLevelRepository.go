package debug

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-registry-lib/dogu"
)

type DoguRegistry struct {
	versionRegistry dogu.DoguVersionRegistry
	dogu.LocalDoguDescriptorRepository
}

func NewDoguRegistry(versionRegistry dogu.DoguVersionRegistry) *DoguRegistry {
	return &DoguRegistry{
		versionRegistry: versionRegistry,
	}
}

func (r *DoguRegistry) GetCurrent(ctx context.Context, simpleDoguName string) (*core.Dogu, error) {
	current, err := r.versionRegistry.GetCurrent(ctx, dogu.SimpleDoguName(simpleDoguName))
	if err != nil {
		return nil, err
	}
	get, err := r.LocalDoguDescriptorRepository.Get(ctx, current)
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
	var allCurrentDogus []*core.Dogu
	for _, doguVersion := range allCurrentDoguVersions {
		d, err := r.LocalDoguDescriptorRepository.Get(ctx, doguVersion)
		if err != nil {
			return nil, err
		}
		allCurrentDogus = append(allCurrentDogus, d)
	}
	return allCurrentDogus, nil
}
