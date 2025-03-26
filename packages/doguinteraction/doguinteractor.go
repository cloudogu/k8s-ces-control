package doguinteraction

import (
	"context"
	"errors"
	"fmt"
	dogu2 "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

var waitTimeout = time.Minute * 10

const (
	doguConfigKeyLogLevel = "logging/root"

	stateStarted = "started"
	stateStopped = "stopped"
)

type defaultDoguInterActor struct {
	doguClient           DoguInterface
	doguRestartClient    DoguRestartInterface
	doguConfigRepository doguConfigRepository
	doguDescriptorGetter doguDescriptorGetter
}

// NewDefaultDoguInterActor creates a new instance of defaultDoguInterActor.
func NewDefaultDoguInterActor(doguConfigRepository doguConfigRepository, doguClient DoguInterface, doguRestartClient DoguRestartInterface, doguDescriptorGetter doguDescriptorGetter) *defaultDoguInterActor {
	return &defaultDoguInterActor{
		doguConfigRepository: doguConfigRepository,
		doguClient:           doguClient,
		doguRestartClient:    doguRestartClient,
		doguDescriptorGetter: doguDescriptorGetter,
	}
}

// StartDogu starts the specified dogu.
func (ddi *defaultDoguInterActor) StartDogu(ctx context.Context, doguName string) error {
	return ddi.StartDoguWithWait(ctx, doguName, false)
}

// StartDoguWithWait starts the specified dogu and waits until started if specified.
func (ddi *defaultDoguInterActor) StartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	return ddi.startStopDogu(ctx, doguName, false, waitForRollout)
}

// StopDogu stops the specified dogu.
func (ddi *defaultDoguInterActor) StopDogu(ctx context.Context, doguName string) error {
	return ddi.StopDoguWithWait(ctx, doguName, false)
}

// StopDoguWithWait stops the specified dogu and waits until stopped if specified.
func (ddi *defaultDoguInterActor) StopDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	return ddi.startStopDogu(ctx, doguName, true, waitForRollout)
}

// RestartDoguWithWait restarts the specified dogu waits until restarted if specified.
func (ddi *defaultDoguInterActor) RestartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	doguRestart := &v2.DoguRestart{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", doguName),
		},
		Spec: v2.DoguRestartSpec{
			DoguName: doguName,
		},
	}
	if _, err := ddi.doguRestartClient.Create(ctx, doguRestart, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to restart dogu %s: %w", doguName, err)
	}

	if waitForRollout {
		if err := ddi.waitForDoguStartStop(ctx, doguName); err != nil {
			return fmt.Errorf("error waiting for dogu %s while restarting: %w", doguName, err)
		}
	}

	return nil
}

// RestartDogu restarts the specified dogu.
func (ddi *defaultDoguInterActor) RestartDogu(ctx context.Context, doguName string) error {
	return ddi.RestartDoguWithWait(ctx, doguName, false)
}

func emptyDoguNameError() error {
	return fmt.Errorf("dogu name must not be empty")
}

func (ddi *defaultDoguInterActor) startStopDogu(ctx context.Context, doguName string, shouldStop bool, waitForRollout bool) error {
	dogu, err := ddi.doguClient.Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get dogu %s: %w", doguName, err)
	}

	dogu.Spec.Stopped = shouldStop

	_, err = ddi.doguClient.UpdateSpecWithRetry(ctx, dogu, func(spec v2.DoguSpec) v2.DoguSpec {
		spec.Stopped = shouldStop
		return spec
	}, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to start/stop dogu %s: %w", doguName, err)
	}

	if waitForRollout {
		return ddi.waitForDoguStartStop(ctx, doguName)
	}

	return nil
}

func (ddi *defaultDoguInterActor) waitForDoguStartStop(ctx context.Context, doguName string) error {
	timeoutCtx, cancelTimeoutCtx := context.WithTimeoutCause(ctx, waitTimeout, fmt.Errorf("timout (%v) reached waiting for dogu %s", waitTimeout, doguName))
	defer cancelTimeoutCtx()

	watchOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", doguName),
	}
	watcher, err := ddi.doguClient.Watch(timeoutCtx, watchOptions)
	if err != nil {
		return fmt.Errorf("error starting watch for dogu %s: %w", doguName, err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			return fmt.Errorf("error in watch while waiting for start/stop: %v", event.Object)
		default:
			isInDesiredState, cErr := ddi.checkIfDoguInDesiredStopState(timeoutCtx, doguName)
			if cErr != nil {
				return fmt.Errorf("error checking dogu-state while waiting for start/stop: %w", cErr)
			}

			if isInDesiredState {
				return nil
			}
		}
	}

	return fmt.Errorf("watch for dogu %s stopped: %v", doguName, context.Cause(timeoutCtx))
}

func (ddi *defaultDoguInterActor) checkIfDoguInDesiredStopState(ctx context.Context, doguName string) (isInDesiredState bool, err error) {
	dogu, getErr := ddi.doguClient.Get(ctx, doguName, metav1.GetOptions{})
	if getErr != nil {
		return false, fmt.Errorf("failed to get dogu %s: %w", doguName, getErr)
	}

	desiredState := stateStarted
	if dogu.Spec.Stopped {
		desiredState = stateStopped
	}
	currentState := stateStarted
	if dogu.Status.Stopped {
		currentState = stateStopped
	}

	if dogu.Spec.Stopped != dogu.Status.Stopped {
		logrus.Debug(fmt.Sprintf("dogu %q has NOT reached desired start/stop state. Desired state is %s; Current state is %s", doguName, desiredState, currentState))
		return false, nil
	}

	logrus.Debug(fmt.Sprintf("dogu %q has reached desired start/stop state: %s", doguName, desiredState))
	return true, nil
}

// SetLogLevelInAllDogus sets the specified log level to all dogus.
func (ddi *defaultDoguInterActor) SetLogLevelInAllDogus(ctx context.Context, logLevel string) error {
	allDogus, err := ddi.doguDescriptorGetter.GetCurrentOfAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting all dogus while setting log-level: %w", err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig, _ := ddi.doguConfigRepository.Get(ctx, dogu2.SimpleName(dogu.GetSimpleName()))
		newConfig, err := doguConfig.Set(doguConfigKeyLogLevel, config.Value(logLevel))
		if err != nil {
			multiError = errors.Join(multiError, err)
		}

		newDoguConfig := config.DoguConfig{
			DoguName: dogu2.SimpleName(dogu.GetSimpleName()),
			Config:   newConfig,
		}

		_, err = ddi.doguConfigRepository.Update(ctx, newDoguConfig)
		if err != nil {
			multiError = errors.Join(multiError, err)
		}
	}

	return multiError
}

// StopAllDogus stops all dogus in the correct dependency order.
func (ddi *defaultDoguInterActor) StopAllDogus(ctx context.Context) error {
	allDogus, err := ddi.doguDescriptorGetter.GetCurrentOfAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting all dogus while stopping: %w", err)
	}

	sortedDogus, err := core.SortDogusByInvertedDependencyWithError(allDogus)
	if err != nil {
		return fmt.Errorf("error sorting dogus while stopping: %w", err)
	}

	var multiError error
	for _, dogu := range sortedDogus {
		stopErr := ddi.StopDoguWithWait(ctx, dogu.GetSimpleName(), true)
		if stopErr != nil {
			multiError = errors.Join(multiError, stopErr)
		}
	}

	return multiError
}

// StartAllDogus starts all dogus in the correct dependency order.
func (ddi *defaultDoguInterActor) StartAllDogus(ctx context.Context) error {
	allDogus, err := ddi.doguDescriptorGetter.GetCurrentOfAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting all dogus while starting: %w", err)
	}

	sortedDogus, err := core.SortDogusByDependencyWithError(allDogus)
	if err != nil {
		return fmt.Errorf("error sorting dogus while starting: %w", err)
	}

	var multiError error
	for _, dogu := range sortedDogus {
		startErr := ddi.StartDoguWithWait(ctx, dogu.GetSimpleName(), true)
		if startErr != nil {
			multiError = errors.Join(multiError, startErr)
		}
	}

	return multiError
}
