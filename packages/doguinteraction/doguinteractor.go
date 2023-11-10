package doguinteraction

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	scalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var waitTimeout = time.Minute * 10
var waitInterval = time.Second * 5

const (
	doguConfigKeyLogLevel   = "logging/root"
	containerStateCrashLoop = "CrashLoopBackOff"
)

type defaultDoguInterActor struct {
	clientSet clusterClientSet
	registry  cesRegistry
	namespace string
}

// NewDefaultDoguInterActor creates a new instance of defaultDoguInterActor.
func NewDefaultDoguInterActor(clientSet clusterClientSet, namespace string) *defaultDoguInterActor {
	return &defaultDoguInterActor{clientSet: clientSet, namespace: namespace}
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

	return ddi.scaleDeployment(ctx, doguName, 1, waitForRollout)
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

	return ddi.scaleDeployment(ctx, doguName, 0, waitForRollout)
}

// RestartDoguWithWait restarts the specified dogu waits for the deployment rollouts if specified.
func (ddi *defaultDoguInterActor) RestartDoguWithWait(ctx context.Context, doguName string, waitForRollout bool) error {
	if doguName == "" {
		return emptyDoguNameError()
	}

	scale, err := ddi.clientSet.AppsV1().Deployments(ddi.namespace).GetScale(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment scale for dogu %s: %w", doguName, err)
	}

	zeroReplicas := int32(0)
	if scale.Spec.Replicas == zeroReplicas {
		return ddi.scaleDeployment(ctx, doguName, 1, waitForRollout)
	}

	err = ddi.scaleDeployment(ctx, doguName, 0, true)
	if err != nil {
		return err
	}

	return ddi.scaleDeployment(ctx, doguName, 1, waitForRollout)
}

// RestartDogu restarts the specified dogu.
func (ddi *defaultDoguInterActor) RestartDogu(ctx context.Context, doguName string) error {
	return ddi.RestartDoguWithWait(ctx, doguName, false)
}

func emptyDoguNameError() error {
	return fmt.Errorf("dogu name must not be empty")
}

func (ddi *defaultDoguInterActor) scaleDeployment(ctx context.Context, doguName string, replicas int32, waitForRollout bool) error {
	scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: doguName, Namespace: ddi.namespace}, Spec: scalingv1.ScaleSpec{Replicas: replicas}}
	_, err := ddi.clientSet.AppsV1().Deployments(ddi.namespace).UpdateScale(ctx, doguName, scale, metav1.UpdateOptions{})
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to scale deployment %s to %d: %s", doguName, replicas, err.Error())
	}

	if waitForRollout {
		return ddi.waitForDeploymentRollout(ctx, doguName)
	}

	return nil
}

func (ddi *defaultDoguInterActor) waitForDeploymentRollout(ctx context.Context, doguName string) error {
	timeoutTimer := time.NewTimer(waitTimeout)
	// Use a ticker instead of a kubernetes watch because the watch does not notify on status changes.
	ticker := time.NewTicker(waitInterval)
	for {
		select {
		case <-ticker.C:
			rolledOut, stopWait, err := ddi.doWaitForDeploymentRollout(ctx, doguName)
			if err != nil {
				logrus.Error(err)
				stopWaitChannels(timeoutTimer, ticker)
			}

			if stopWait || rolledOut {
				stopWaitChannels(timeoutTimer, ticker)
			}

			return nil
		case <-timeoutTimer.C:
			ticker.Stop()
			return fmt.Errorf("failed to wait for deployment %s rollout: timeout reached", doguName)
		}
	}
}

func (ddi *defaultDoguInterActor) doWaitForDeploymentRollout(ctx context.Context, doguName string) (rolledOut bool, stopWait bool, err error) {
	logrus.Info(fmt.Sprintf("check rollout status for deployment %s", doguName))
	deployment, getErr := ddi.clientSet.AppsV1().Deployments(ddi.namespace).Get(ctx, doguName, metav1.GetOptions{})
	if getErr != nil {
		return false, true, fmt.Errorf("failed to get deployment %s: %w", doguName, getErr)
	}

	isInCrashLoop, err := ddi.isDoguContainerInCrashLoop(ctx, doguName)
	if err != nil || isInCrashLoop {
		return false, true, err
	}

	if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		logrus.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated", deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas))
		return false, false, nil
	}
	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		logrus.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d old replicas are pending termination", deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas))
		return false, false, nil
	}
	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		logrus.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d of %d updated replicas are available", deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas))
		return false, false, nil
	}
	logrus.Info(fmt.Sprintf("deployment %q successfully rolled out", deployment.Name))
	return true, true, nil
}

func stopWaitChannels(timer *time.Timer, ticker *time.Ticker) {
	timer.Stop()
	ticker.Stop()
}

func (ddi *defaultDoguInterActor) isDoguContainerInCrashLoop(ctx context.Context, doguName string) (bool, error) {
	list, getErr := ddi.clientSet.CoreV1().Pods(ddi.namespace).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("dogu.name=%s", doguName)})
	if getErr != nil {
		return false, fmt.Errorf("failed to get pods of deployment %s", doguName)
	}

	for _, pod := range list.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name != doguName {
				continue
			}

			containerWaitState := containerStatus.State.Waiting

			if containerWaitState != nil && containerWaitState.Reason == containerStateCrashLoop {
				logrus.Error(fmt.Errorf("some containers are in a crash loop"), fmt.Sprintf("skip waiting rollout for deployment %s", doguName))
				return true, nil
			}
		}
	}

	return false, nil
}

// SetLogLevelInAllDogus sets the specified log level to all dogus.
func (ddi *defaultDoguInterActor) SetLogLevelInAllDogus(logLevel string) error {
	allDogus, err := ddi.registry.DoguRegistry().GetAll()
	if err != nil {
		return getAllDogusError(err)
	}

	var multiError error
	for _, dogu := range allDogus {
		doguConfig := ddi.registry.DoguConfig(dogu.GetSimpleName())
		setErr := doguConfig.Set(doguConfigKeyLogLevel, logLevel)
		if setErr != nil {
			multiError = errors.Join(multiError, setErr)
		}
	}

	return multiError
}

// StopAllDogus stops all dogus in the correct dependency order.
func (ddi *defaultDoguInterActor) StopAllDogus(ctx context.Context) error {
	allDogus, err := ddi.registry.DoguRegistry().GetAll()
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
	allDogus, err := ddi.registry.DoguRegistry().GetAll()
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
