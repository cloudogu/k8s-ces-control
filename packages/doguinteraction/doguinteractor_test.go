package doguinteraction

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		doguRegistryMock := newMockDoguRegistry(t)

		// when
		actor := NewDefaultDoguInterActor(repository.DoguConfigRepository{}, clientSetMock, testNamespace, doguRegistryMock)

		// then
		require.NotNil(t, actor)
		assert.NotNil(t, actor.clientSet)
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

		deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", getScale("postgresql", 0), metav1.UpdateOptions{}).Return(nil, nil)
		deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", getScale("postgresql", 1), metav1.UpdateOptions{}).Return(nil, nil)
		deploymentMock.EXPECT().GetScale(testCtx, "postgresql", metav1.GetOptions{}).Return(&scalingv1.Scale{Spec: scalingv1.ScaleSpec{Replicas: 1}}, nil)

		notRolledOutDeploy := getZeroReplicaDeployment()
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(notRolledOutDeploy, nil)
		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		err := sut.RestartDogu(testCtx, "postgresql")

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_RestartDoguWithWait(t *testing.T) {
	t.Run("should return error on empty dogu name", func(t *testing.T) {
		// given
		sut := defaultDoguInterActor{}

		// when
		err := sut.RestartDoguWithWait(testCtx, "", true)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu name must not be empty")
	})

	t.Run("should return error on error getting current scale", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		appsMock := newMockAppsV1Interface(t)
		deployMock := newMockDeploymentInterface(t)
		appsMock.EXPECT().Deployments(testNamespace).Return(deployMock)
		deployMock.EXPECT().GetScale(testCtx, "redmine", metav1.GetOptions{}).Return(nil, assert.AnError)
		clientSetMock.EXPECT().AppsV1().Return(appsMock)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		err := sut.RestartDoguWithWait(testCtx, "redmine", true)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get deployment scale for dogu redmine")
	})

	t.Run("should return error on update scale error", func(t *testing.T) {
		testUpdateScaleError(t, "redmine", 0, 1)
	})

	t.Run("should return error on update scale error on rolled out deployment", func(t *testing.T) {
		testUpdateScaleError(t, "redmine", 1, 0)
	})
}

func testUpdateScaleError(t *testing.T, doguName string, currentReplicas int32, desiredReplicas int32) {
	// given
	clientSetMock := newMockClusterClientSet(t)
	appsMock := newMockAppsV1Interface(t)
	deployMock := newMockDeploymentInterface(t)
	appsMock.EXPECT().Deployments(testNamespace).Return(deployMock)
	deployMock.EXPECT().GetScale(testCtx, doguName, metav1.GetOptions{}).Return(getScale(doguName, currentReplicas), nil)
	deployMock.EXPECT().UpdateScale(testCtx, doguName, getScale(doguName, desiredReplicas), metav1.UpdateOptions{}).Return(nil, assert.AnError)
	clientSetMock.EXPECT().AppsV1().Return(appsMock)

	sut := defaultDoguInterActor{
		clientSet: clientSetMock,
		namespace: testNamespace,
	}

	// when
	err := sut.RestartDoguWithWait(testCtx, doguName, true)

	// then
	require.Error(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf("failed to scale deployment %s to %d", doguName, desiredReplicas))
}

func Test_defaultDoguInterActor_StartDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)

		deploymentMock.EXPECT().UpdateScale(context.TODO(), "postgresql", getScale("postgresql", 1), metav1.UpdateOptions{}).Return(nil, nil)

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
	t.Run("should return error on empty dogu name", func(t *testing.T) {
		// given
		sut := defaultDoguInterActor{}

		// when
		err := sut.StartDoguWithWait(testCtx, "", true)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu name must not be empty")
	})

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

		deploymentMock.EXPECT().UpdateScale(context.TODO(), "postgresql", getScale("postgresql", 1), metav1.UpdateOptions{}).Return(nil, nil)

		notRolledOutDeploy := getDeployment(1, 0, 0, 0, 0)
		rolledOutDeploy := getDeployment(1, 1, 1, 1, 1)
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(notRolledOutDeploy, nil).Once().Run(func(args mock.Arguments) {
			deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(rolledOutDeploy, nil)
		})
		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

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

	t.Run("should return error on empty dogu name", func(t *testing.T) {
		// given
		sut := defaultDoguInterActor{}

		// when
		err := sut.StopDogu(testCtx, "")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "dogu name must not be empty")
	})
}

func Test_defaultDoguInterActor_StartAllDogus(t *testing.T) {
	t.Run("should return error on get all dogus error", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)
		sut := defaultDoguInterActor{
			doguRegistry: doguRegistryMock,
		}

		// when
		err := sut.StartAllDogus(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})

	t.Run("should start dogus in correct order", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		oldWaitInterval := waitInterval
		waitInterval = time.Second * 1
		defer func() { waitInterval = oldWaitInterval }()

		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql"}}}}, nil)

		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		podMock := newMockPodInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		postgresPod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "postgresql"}}}}
		postgresqlPodList := &corev1.PodList{Items: []corev1.Pod{postgresPod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=postgresql"}).Return(postgresqlPodList, nil)
		redminePod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "redmine"}}}}
		redminePodList := &corev1.PodList{Items: []corev1.Pod{redminePod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

		// This ensures the start order.
		deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", getScale("postgresql", 1), metav1.UpdateOptions{}).Return(nil, nil).Run(func(ctx context.Context, deploymentName string, scale *scalingv1.Scale, opts metav1.UpdateOptions) {
			deploymentMock.EXPECT().UpdateScale(testCtx, "redmine", getScale("redmine", 1), metav1.UpdateOptions{}).Return(nil, nil)
		})

		sut := defaultDoguInterActor{
			doguRegistry:         doguRegistryMock,
			doguConfigRepository: repository.DoguConfigRepository{},
			clientSet:            clientSetMock,
			namespace:            testNamespace,
		}

		notRolledOutDeploy := getZeroReplicaDeployment()
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(notRolledOutDeploy, nil)
		deploymentMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(notRolledOutDeploy, nil)

		// when
		err := sut.StartAllDogus(testCtx)

		// then
		require.NoError(t, err)
	})
}

func getZeroReplicaDeployment() *v1.Deployment {
	return getDeployment(0, 0, 0, 0, 0)
}

func getDeployment(specReplicas int32, statusReplicas int32, readyReplicas int32, updatedReplicas int32, availableReplicas int32) *v1.Deployment {
	specReplica := specReplicas
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.DeploymentSpec{Replicas: &specReplica},
		Status:     v1.DeploymentStatus{Replicas: statusReplicas, ReadyReplicas: readyReplicas, UpdatedReplicas: updatedReplicas, AvailableReplicas: availableReplicas},
	}
}

func getScale(deployName string, replicas int32) *scalingv1.Scale {
	return &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: deployName, Namespace: testNamespace}, Spec: scalingv1.ScaleSpec{Replicas: replicas}}
}

func Test_defaultDoguInterActor_StopAllDogus(t *testing.T) {
	t.Run("should return error on get all dogus error", func(t *testing.T) {
		// given

		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)
		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguRegistry:         doguRegistryMock,
		}

		// when
		err := sut.StopAllDogus(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})

	t.Run("should stop dogus in correct order", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		oldWaitInterval := waitInterval
		waitInterval = time.Second * 1
		defer func() { waitInterval = oldWaitInterval }()

		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql"}}}}, nil)

		clientSetMock := newMockClusterClientSet(t)
		appsV1Mock := newMockAppsV1Interface(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		deploymentMock := newMockDeploymentInterface(t)
		podMock := newMockPodInterface(t)
		appsV1Mock.EXPECT().Deployments(testNamespace).Return(deploymentMock)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		postgresPod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "postgresql"}}}}
		postgresqlPodList := &corev1.PodList{Items: []corev1.Pod{postgresPod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=postgresql"}).Return(postgresqlPodList, nil)
		redminePod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "redmine"}}}}
		redminePodList := &corev1.PodList{Items: []corev1.Pod{redminePod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

		// This ensures the start order.
		deploymentMock.EXPECT().UpdateScale(testCtx, "redmine", getScale("redmine", 0), metav1.UpdateOptions{}).Return(nil, nil).Run(func(ctx context.Context, deploymentName string, scale *scalingv1.Scale, opts metav1.UpdateOptions) {
			deploymentMock.EXPECT().UpdateScale(testCtx, "postgresql", getScale("postgresql", 0), metav1.UpdateOptions{}).Return(nil, nil)
		})

		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguRegistry:         doguRegistryMock,
			clientSet:            clientSetMock,
			namespace:            testNamespace,
		}

		notRolledOutDeploy := getZeroReplicaDeployment()
		deploymentMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(notRolledOutDeploy, nil)
		deploymentMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(notRolledOutDeploy, nil)

		// when
		err := sut.StopAllDogus(testCtx)

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_SetLogLevelInAllDogus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(
			[]*core.Dogu{
				{Name: "official/postgresql"},
				{Name: "official/redmine"},
			},
			nil,
		)

		doguConfigRepositoryMock := newMockDoguConfigRepository(t)
		doguConfigA := config.CreateDoguConfig("postgresql", config.Entries{})
		doguConfigB := config.CreateDoguConfig("redmine", config.Entries{})
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("postgresql")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("redmine")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).RunAndReturn(func(ctx context.Context, doguConfig config.DoguConfig) (config.DoguConfig, error) {
			get, b := doguConfig.Get("logging/root")
			require.True(t, b)
			assert.Equal(t, "DEBUG", get.String())
			return doguConfig, nil
		}).Times(2)

		sut := defaultDoguInterActor{
			doguConfigRepository: doguConfigRepositoryMock,
			doguRegistry:         doguRegistryMock,
			namespace:            testNamespace,
		}

		// when
		err := sut.SetLogLevelInAllDogus(testCtx, "DEBUG")

		// then
		require.NoError(t, err)
	})

	t.Run("should return errors on error updating log levels", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(
			[]*core.Dogu{
				{Name: "official/postgresql"},
				{Name: "official/redmine"},
			},
			nil,
		)

		doguConfigRepositoryMock := newMockDoguConfigRepository(t)
		doguConfigA := config.CreateDoguConfig("postgresql", config.Entries{})
		doguConfigB := config.CreateDoguConfig("redmine", config.Entries{})
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("postgresql")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), config.SimpleDoguName("redmine")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).RunAndReturn(func(ctx context.Context, doguConfig config.DoguConfig) (config.DoguConfig, error) {
			get, b := doguConfig.Get("logging/root")
			require.True(t, b)
			assert.Equal(t, "DEBUG", get.String())
			return doguConfig, assert.AnError
		}).Times(2)

		sut := defaultDoguInterActor{
			doguConfigRepository: doguConfigRepositoryMock,
			doguRegistry:         doguRegistryMock,
			namespace:            testNamespace,
		}

		// when
		err := sut.SetLogLevelInAllDogus(testCtx, "DEBUG")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "assert.AnError general error for testing\nassert.AnError general error for testing")
	})

	t.Run("should return errors on error getting current dogus", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguRegistry(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguRegistry:         doguRegistryMock,
			namespace:            testNamespace,
		}

		// when
		err := sut.SetLogLevelInAllDogus(testCtx, "DEBUG")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get all dogus")
	})
}

func Test_defaultDoguInterActor_isDoguContainerInCrashLoop(t *testing.T) {
	t.Run("should ignore other container status than dogu ones", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		podMock := newMockPodInterface(t)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		redminePod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "OtherContainer"}, {Name: "redmine"}}}}
		redminePodList := &corev1.PodList{Items: []corev1.Pod{redminePod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		loop, err := sut.isDoguContainerInCrashLoop(testCtx, "redmine")

		// then
		require.NoError(t, err)
		assert.False(t, loop)
	})

	t.Run("should return true if container is in restart loop", func(t *testing.T) {
		// given
		clientSetMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)
		podMock := newMockPodInterface(t)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		redminePod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "redmine", State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"}}}}}}
		redminePodList := &corev1.PodList{Items: []corev1.Pod{redminePod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

		sut := defaultDoguInterActor{
			clientSet: clientSetMock,
			namespace: testNamespace,
		}

		// when
		loop, err := sut.isDoguContainerInCrashLoop(testCtx, "redmine")

		// then
		require.NoError(t, err)
		assert.True(t, loop)
	})
}

func Test_defaultDoguInterActor_doWaitForDeploymentRollout(t *testing.T) {
	t.Run("should return error on deployment get error", func(t *testing.T) {
		// given
		clientMock := newMockClusterClientSet(t)
		appsMock := newMockAppsV1Interface(t)
		deployMock := newMockDeploymentInterface(t)
		deployMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(nil, assert.AnError)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		appsMock.EXPECT().Deployments(testNamespace).Return(deployMock)
		sut := defaultDoguInterActor{clientSet: clientMock, namespace: testNamespace}

		// when
		_, _, err := sut.doWaitForDeploymentRollout(testCtx, "redmine")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get deployment redmine")
	})

	t.Run("should stop if container is in crash loop", func(t *testing.T) {
		// given
		deployMock := newMockDeploymentInterface(t)
		deployMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(nil, nil)
		clientMock := newMockClusterClientSet(t)
		coreV1Mock := newMockCoreV1Interface(t)
		clientMock.EXPECT().CoreV1().Return(coreV1Mock)
		appsMock := newMockAppsV1Interface(t)
		clientMock.EXPECT().AppsV1().Return(appsMock)
		appsMock.EXPECT().Deployments(testNamespace).Return(deployMock)

		podMock := newMockPodInterface(t)
		coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
		redminePod := corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Name: "redmine", State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"}}}}}}
		redminePodList := &corev1.PodList{Items: []corev1.Pod{redminePod}}
		podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

		sut := defaultDoguInterActor{clientSet: clientMock, namespace: testNamespace}

		// when
		_, shouldStop, err := sut.doWaitForDeploymentRollout(testCtx, "redmine")

		// then
		require.NoError(t, err)
		assert.True(t, shouldStop)
	})

	t.Run("should not stop if deployment replicas are not updated in spec", func(t *testing.T) {
		// given
		notAllSpecsUpdatedDeployment := getDeployment(1, 0, 0, 0, 0)
		testPendingDeploymentRollout(t, notAllSpecsUpdatedDeployment)
	})

	t.Run("should not stop if not all replicas are updated", func(t *testing.T) {
		// given
		notAllSpecsUpdatedDeployment := getDeployment(1, 2, 0, 1, 0)
		testPendingDeploymentRollout(t, notAllSpecsUpdatedDeployment)
	})

	t.Run("should not stop if not replicas are available", func(t *testing.T) {
		// given
		notAllSpecsUpdatedDeployment := getDeployment(2, 2, 0, 2, 1)
		testPendingDeploymentRollout(t, notAllSpecsUpdatedDeployment)
	})
}

func testPendingDeploymentRollout(t *testing.T, deployment *v1.Deployment) {
	deployMock := newMockDeploymentInterface(t)
	deployMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(deployment, nil)
	clientMock := newMockClusterClientSet(t)
	coreV1Mock := newMockCoreV1Interface(t)
	clientMock.EXPECT().CoreV1().Return(coreV1Mock)
	appsMock := newMockAppsV1Interface(t)
	clientMock.EXPECT().AppsV1().Return(appsMock)
	appsMock.EXPECT().Deployments(testNamespace).Return(deployMock)

	podMock := newMockPodInterface(t)
	coreV1Mock.EXPECT().Pods(testNamespace).Return(podMock)
	redminePodList := &corev1.PodList{Items: []corev1.Pod{}}
	podMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "dogu.name=redmine"}).Return(redminePodList, nil)

	sut := defaultDoguInterActor{clientSet: clientMock, namespace: testNamespace}

	// when
	_, shouldStop, err := sut.doWaitForDeploymentRollout(testCtx, "redmine")

	// then
	require.NoError(t, err)
	assert.False(t, shouldStop)
}

func Test_defaultDoguInterActor_waitForDeploymentRollout(t *testing.T) {
	t.Run("should return error on timeout", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second
		defer func() { waitTimeout = oldWaitTimeout }()

		oldWaitInterval := waitInterval
		waitInterval = time.Second * 2
		defer func() { waitInterval = oldWaitInterval }()

		sut := defaultDoguInterActor{}

		// when
		err := sut.waitForDeploymentRollout(testCtx, "redmine")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to wait for deployment redmine rollout: timeout reached")
	})
}
