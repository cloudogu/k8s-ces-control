package doguinteraction

import (
	"context"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/apps/v1"
	scalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

const testNamespace = "ecosystem"

var testCtx = context.TODO()

func TestNewDefaultDoguInterActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)

		// when
		actor := NewDefaultDoguInterActor(clientSetMock, testNamespace)

		// then
		require.NotNil(t, actor)
	})
}

func Test_defaultDoguInterActor_RestartDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		oldWaitInterval := waitInterval
		waitInterval = time.Second * 1
		defer func() { waitInterval = oldWaitInterval }()

		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		podMock := newMockPodInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		runningPostgresPod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "postgresql"}}}}
		podList := &corev1.PodList{Items: []corev1.Pod{runningPostgresPod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=postgresql"}).Return(podList, nil)

		scaleDown := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: 0}}
		scaleUp := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: 1}}
		deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", scaleDown, metav1.UpdateOptions{}).Return(nil, nil)
		deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", scaleUp, metav1.UpdateOptions{}).Return(nil, nil)
		deploymentMock.EXPECT().GetScale(testCtx, "postgresql", metav1.GetOptions{}).Return(&scalingv1.Scale{Spec: scalingv1.ScaleSpec{Replicas: 1}}, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}
		zeroReplicas := int32(0)
		rolledOutDeployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       v1.DeploymentSpec{Replicas: &zeroReplicas},
			Status:     v1.DeploymentStatus{Replicas: zeroReplicas, ReadyReplicas: zeroReplicas, UpdatedReplicas: zeroReplicas, AvailableReplicas: zeroReplicas},
		}
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(rolledOutDeployment, nil)

		// when
		err := sut.RestartDogu(testCtx, "postgresql")

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_RestartDoguWithWait(t *testing.T) {

}

func Test_defaultDoguInterActor_StartDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)

		scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: 1}}
		deploymentMock.EXPECT().UpdateScale(context.TODO(), "postgresql", scale, metav1.UpdateOptions{}).Return(nil, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		err := sut.StartDogu(testCtx, "postgresql")

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_StartDoguWithWait(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		oldWaitInterval := waitInterval
		waitInterval = time.Second * 1
		defer func() { waitInterval = oldWaitInterval }()

		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		podMock := newMockPodInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		runningPostgresPod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "postgresql"}}}}
		podList := &corev1.PodList{Items: []corev1.Pod{runningPostgresPod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=postgresql"}).Return(podList, nil)

		scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: 1}}
		deploymentMock.EXPECT().UpdateScale(context.TODO(), "postgresql", scale, metav1.UpdateOptions{}).Return(nil, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}
		zeroReplicas := int32(0)
		rolledOutDeployment := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       v1.DeploymentSpec{Replicas: &zeroReplicas},
			Status:     v1.DeploymentStatus{Replicas: zeroReplicas, ReadyReplicas: zeroReplicas, UpdatedReplicas: zeroReplicas, AvailableReplicas: zeroReplicas},
		}
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(rolledOutDeployment, nil)

		// when
		err := sut.StartDoguWithWait(testCtx, "postgresql", true)

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_StopDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)

		scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: 0}}
		deploymentMock.EXPECT().UpdateScale(context.TODO(), "postgresql", scale, metav1.UpdateOptions{}).Return(nil, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		err := sut.StopDogu(testCtx, "postgresql")

		// then
		require.NoError(t, err)
	})
}
