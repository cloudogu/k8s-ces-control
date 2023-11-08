package doguinteraction

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	scalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var waitTimeout = time.Minute * 10

type defaultDoguInterActor struct {
	clientSet clusterClientSet
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

// RestartDogu restarts the specified dogu
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
		return status.Errorf(codes.Unknown, "failed to scale deployment to %d: %s", replicas, err.Error())
	}

	if waitForRollout {
		return ddi.waitForDeploymentRollout(ctx, doguName)
	}

	return nil
}

func (ddi *defaultDoguInterActor) waitForDeploymentRollout(ctx context.Context, doguName string) error {
	logger := log.FromContext(ctx)
	deployLabel := fmt.Sprintf("dogu.name=%s", doguName)
	timeoutSeconds := int64(waitTimeout.Seconds())
	timer := time.NewTimer(waitTimeout)
	watch, err := ddi.clientSet.AppsV1().Deployments(ddi.namespace).Watch(ctx, metav1.ListOptions{LabelSelector: deployLabel, TimeoutSeconds: &timeoutSeconds})
	if err != nil {
		return fmt.Errorf("failed create watch for deployment wit label %s: %s", deployLabel, err)
	}

	for {
		select {
		case event := <-watch.ResultChan():
			deployment, ok := event.Object.(*appsv1.Deployment)
			if !ok {
				logger.Error(fmt.Errorf("watch object %+v is not type of deployment", event.Object), fmt.Sprintf("failed to watch deployment %s", doguName))
				continue
			}

			if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
				logger.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated", deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas))
				continue
			}
			if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
				logger.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d old replicas are pending termination", deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas))
				continue
			}
			if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
				logger.Info(fmt.Sprintf("waiting for deployment %q rollout to finish: %d of %d updated replicas are available", deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas))
				continue
			}
			logger.Info(fmt.Sprintf("deployment %q successfully rolled out\n", deployment.Name))
			return nil
		case <-timer.C:
			return fmt.Errorf("failed to wait for deployment %s rollout: timeout reached", doguName)
		}
	}
}
