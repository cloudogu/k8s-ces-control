package doguinteraction

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var waitTimeout = time.Minute * 10
var waitInterval = time.Second * 5

const (
	doguConfigKeyLogLevel   = "logging/root"
	containerStateCrashLoop = "CrashLoopBackOff"

	stateStarted = "started"
	stateStopped = "stopped"
)

type defaultDoguInterActor struct {
	doguClient           DoguInterface
	doguConfigRepository doguConfigRepository
	doguDescriptorGetter doguDescriptorGetter
	namespace            string
}

// NewDefaultDoguInterActor creates a new instance of defaultDoguInterActor.
func NewDefaultDoguInterActor(doguConfigRepository doguConfigRepository, doguClient DoguInterface, namespace string, doguDescriptorGetter doguDescriptorGetter) *defaultDoguInterActor {
	return &defaultDoguInterActor{doguConfigRepository: doguConfigRepository, doguClient: doguClient, namespace: namespace, doguDescriptorGetter: doguDescriptorGetter}
}

// StartDogu starts the specified dogu.
func (ddi *defaultDoguInterActor) StartDogu(ctx context.Context, doguName string) error {
	return ddi.StartDoguWithWait(ctx, doguName, false)
}

// StartDoguWithWait starts the specified dogu and wait for the deployment rollout if specified.
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

// StopDoguWithWait stops the specified dogu and waits for the deployment rollout if specified.
func (ddi *defaultDoguInterActor) StopDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	return ddi.startStopDogu(ctx, doguName, true, waitForRollout)
}

// RestartDoguWithWait restarts the specified dogu waits for the deployment rollouts if specified.
func (ddi *defaultDoguInterActor) RestartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	// fixme

	return nil

	//scale, err := ddi.clientSet.AppsV1().Deployments(ddi.namespace).GetScale(ctx, doguName, metav1.GetOptions{})
	//if err != nil {
	//	return fmt.Errorf("failed to get deployment scale for dogu %s: %w", doguName, err)
	//}
	//
	//zeroReplicas := int32(0)
	//if scale.Spec.Replicas == zeroReplicas {
	//	return ddi.scaleDeployment(ctx, doguName, 1, waitForRollout)
	//}
	//
	//err = ddi.scaleDeployment(ctx, doguName, 0, true)
	//if err != nil {
	//	return err
	//}
	//
	//return ddi.scaleDeployment(ctx, doguName, 1, waitForRollout)
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

	dogu, err = ddi.doguClient.UpdateSpecWithRetry(ctx, dogu, func(spec v2.DoguSpec) v2.DoguSpec {
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
	timeoutTimer := time.NewTimer(waitTimeout)
	// Use a ticker instead of a kubernetes watch because the watch does not notify on status changes.
	ticker := time.NewTicker(waitInterval)
	for {
		select {
		case <-ticker.C:
			rolledOut, stopWait, err := ddi.doWaitForDoguStartStop(ctx, doguName)
			if err != nil {
				stopWaitChannels(timeoutTimer, ticker)
				return err
			}

			if stopWait || rolledOut {
				stopWaitChannels(timeoutTimer, ticker)
				return nil
			}
		case <-timeoutTimer.C:
			ticker.Stop()
			return fmt.Errorf("failed to wait for deployment %s rollout: timeout reached", doguName)
		}
	}
}

func (ddi *defaultDoguInterActor) doWaitForDoguStartStop(ctx context.Context, doguName string) (rolledOut bool, stopWait bool, err error) {
	logrus.Info(fmt.Sprintf("check rollout status for deployment %s", doguName))
	dogu, getErr := ddi.doguClient.Get(ctx, doguName, metav1.GetOptions{})
	if getErr != nil {
		return false, true, fmt.Errorf("failed to get dogu %s: %w", doguName, getErr)
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
		logrus.Info(fmt.Sprintf("waiting for dogu %q start/stop to finish. Desired state is %s; Current state is %s", dogu.Name, desiredState, currentState))
		return false, false, nil
	}

	logrus.Info(fmt.Sprintf("deployment %q successfully %s", dogu.Name, currentState))
	return true, true, nil
}

func stopWaitChannels(timer *time.Timer, ticker *time.Ticker) {
	timer.Stop()
	ticker.Stop()
}

// SetLogLevelInAllDogus sets the specified log level to all dogus.
func (ddi *defaultDoguInterActor) SetLogLevelInAllDogus(ctx context.Context, logLevel string) error {
	allDogus, err := ddi.doguDescriptorGetter.GetCurrentOfAll(ctx)
	if err != nil {
		return getAllDogusError(err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig, _ := ddi.doguConfigRepository.Get(ctx, config.SimpleDoguName(dogu.GetSimpleName()))
		newConfig, err := doguConfig.Set(doguConfigKeyLogLevel, config.Value(logLevel))
		if err != nil {
			multiError = errors.Join(multiError, err)
		}

		newDoguConfig := config.DoguConfig{
			DoguName: config.SimpleDoguName(dogu.GetSimpleName()),
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
		return getAllDogusError(err)
	}

	sortedDogus := core.SortDogusByInvertedDependency(allDogus)
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
		return getAllDogusError(err)
	}

	sortedDogus := core.SortDogusByDependency(allDogus)
	var multiError error
	for _, dogu := range sortedDogus {
		startErr := ddi.StartDoguWithWait(ctx, dogu.GetSimpleName(), true)
		if startErr != nil {
			multiError = errors.Join(multiError, startErr)
		}
	}

	return multiError
}

func getAllDogusError(err error) error {
	return fmt.Errorf("failed to get all dogus: %w", err)
}
