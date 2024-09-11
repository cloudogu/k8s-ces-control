package debug

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-registry-lib/dogu"
)

type doguGetter struct {
	versionRegistry dogu.DoguVersionRegistry
	doguRepository  dogu.LocalDoguDescriptorRepository
}

func NewDoguGetter(versionRegistry dogu.DoguVersionRegistry, doguRepository dogu.LocalDoguDescriptorRepository) *doguGetter {
	return &doguGetter{
		versionRegistry: versionRegistry,
		doguRepository:  doguRepository,
	}
}

func (r *doguGetter) GetCurrent(ctx context.Context, simpleDoguName string) (*core.Dogu, error) {
	current, err := r.versionRegistry.GetCurrent(ctx, dogu.SimpleDoguName(simpleDoguName))
	if err != nil {
		return nil, fmt.Errorf("failed to get current version for dogu %s: %w", simpleDoguName, err)
	}
	get, err := r.doguRepository.Get(ctx, current)
	if err != nil {
		return nil, fmt.Errorf("failed to get current dogu %s: %w", simpleDoguName, err)
	}
	return get, nil
}
func (r *doguGetter) GetCurrentOfAll(ctx context.Context) ([]*core.Dogu, error) {
	allCurrentDoguVersions, err := r.versionRegistry.GetCurrentOfAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all current dogu versions: %w", err)
	}
	all, err := r.doguRepository.GetAll(ctx, allCurrentDoguVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to get all dogus: %w", err)
	}

	var allCurrentDogus []*core.Dogu
	for _, doguSpec := range all {
		allCurrentDogus = append(allCurrentDogus, doguSpec)
	}
	return allCurrentDogus, nil
}
