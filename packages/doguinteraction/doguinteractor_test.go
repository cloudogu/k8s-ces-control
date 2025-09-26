package doguinteraction

import (
	"context"
	common "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/cloudogu/k8s-registry-lib/config"
	"github.com/cloudogu/k8s-registry-lib/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"testing"
	"time"
)

var testCtx = context.TODO()

func TestNewDefaultDoguInterActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguRestartClientMock := NewMockDoguRestartInterface(t)
		doguRegistryMock := newMockDoguDescriptorGetter(t)

		// when
		actor := NewDefaultDoguInterActor(repository.DoguConfigRepository{}, doguClientMock, doguRestartClientMock, doguRegistryMock)

		// then
		require.NotNil(t, actor)
		assert.NotNil(t, actor.doguClient)
		assert.NotNil(t, actor.doguRestartClient)
		assert.NotNil(t, actor.doguDescriptorGetter)
		assert.NotNil(t, actor.doguConfigRepository)
	})
}

func Test_defaultDoguInterActor_RestartDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "postgresql-"}, Spec: v2.DoguRestartSpec{DoguName: "postgresql"}}

		doguRestartClientMock := NewMockDoguRestartInterface(t)
		doguRestartClientMock.EXPECT().Create(testCtx, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(expectedDoguRestartToCreate, nil)

		sut := defaultDoguInterActor{
			doguRestartClient: doguRestartClientMock,
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

	t.Run("should restart with wait", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "redmine-"}, Spec: v2.DoguRestartSpec{DoguName: "redmine"}}
		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "redmine",
				Stopped: false,
			},
		}

		doguClientMock := NewMockDoguInterface(t)

		doguClientMock.EXPECT().Get(mock.Anything, "redmine", metav1.GetOptions{}).Return(dogu, nil)

		doguRestartClientMock := NewMockDoguRestartInterface(t)
		doguRestartClientMock.EXPECT().Create(testCtx, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(expectedDoguRestartToCreate, nil)

		watcher := watch.NewFake()
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=redmine"}).Return(watcher, nil)

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Action(watch.Modified, expectedDoguRestartToCreate)
		}()

		sut := defaultDoguInterActor{
			doguClient:        doguClientMock,
			doguRestartClient: doguRestartClientMock,
		}

		// when
		err := sut.RestartDoguWithWait(testCtx, "redmine", true)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to restart for error creating restart-cr", func(t *testing.T) {
		// given
		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "redmine-"}, Spec: v2.DoguRestartSpec{DoguName: "redmine"}}

		doguRestartClientMock := NewMockDoguRestartInterface(t)
		doguRestartClientMock.EXPECT().Create(testCtx, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(expectedDoguRestartToCreate, assert.AnError)

		sut := defaultDoguInterActor{
			doguRestartClient: doguRestartClientMock,
		}

		// when
		err := sut.RestartDoguWithWait(testCtx, "redmine", true)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to restart dogu redmine:")
	})

	t.Run("should fail to restart for error while waiting for restart", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "redmine-"}, Spec: v2.DoguRestartSpec{DoguName: "redmine"}}
		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "redmine",
				Stopped: false,
			},
		}

		doguClientMock := NewMockDoguInterface(t)

		doguClientMock.EXPECT().Get(mock.Anything, "redmine", metav1.GetOptions{}).Return(dogu, assert.AnError)

		doguRestartClientMock := NewMockDoguRestartInterface(t)
		doguRestartClientMock.EXPECT().Create(testCtx, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(expectedDoguRestartToCreate, nil)

		watcher := watch.NewFake()
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=redmine"}).Return(watcher, nil)

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Action(watch.Modified, expectedDoguRestartToCreate)
		}()

		sut := defaultDoguInterActor{
			doguClient:        doguClientMock,
			doguRestartClient: doguRestartClientMock,
		}

		// when
		err := sut.RestartDoguWithWait(testCtx, "redmine", true)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error waiting for dogu redmine while restarting: error checking dogu-state while waiting for start/stop: failed to get dogu redmine")
	})
}

func Test_defaultDoguInterActor_startStopDogu(t *testing.T) {
	t.Run("should fail to start/stop for error getting dogu", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		err := sut.startStopDogu(testCtx, "postgresql", true, true)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get dogu postgresql:")
	})

	t.Run("should fail to start/stop for error updating dogu-spec", func(t *testing.T) {
		// given
		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: true,
			},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(testCtx, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		err := sut.startStopDogu(testCtx, "postgresql", false, true)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to start/stop dogu postgresql:")
	})
}

func Test_defaultDoguInterActor_StartDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: true,
			},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(testCtx, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(dogu, nil)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
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

	t.Run("error starting watch", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 3
		defer func() { waitTimeout = oldWaitTimeout }()

		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(mock.Anything, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=postgresql"}).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		err := sut.StartDoguWithWait(testCtx, "postgresql", true)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error starting watch for dogu postgresql")
	})

	t.Run("error in watch for dogu", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 3
		defer func() { waitTimeout = oldWaitTimeout }()

		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(mock.Anything, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(dogu, nil)

		watcher := watch.NewFakeWithChanSize(5, false)
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=postgresql"}).Return(watcher, nil)

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Action(watch.Error, expectedUpdateDogu)
		}()

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		err := sut.StartDoguWithWait(testCtx, "postgresql", true)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "error in watch while waiting for start/stop:")
	})

	t.Run("error wait-timeout reached", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = 0
		defer func() { waitTimeout = oldWaitTimeout }()

		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
			Status: v2.DoguStatus{Stopped: true},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(mock.Anything, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(dogu, nil)

		watcher := watch.NewFakeWithChanSize(5, false)

		var callCtx context.Context

		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=postgresql"}).Return(watcher, nil).Run(
			func(watchCtx context.Context, opts metav1.ListOptions) {
				callCtx = watchCtx
				time.Sleep(2 * time.Second)
				watcher.Stop()
			})

		go func() {
			time.Sleep(1 * time.Second)
			watcher.Action(watch.Modified, expectedUpdateDogu)
		}()

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		err := sut.StartDoguWithWait(testCtx, "postgresql", true)

		expectedTimeoutErr := context.Cause(callCtx).Error()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, expectedTimeoutErr)
	})

}

func Test_defaultDoguInterActor_StopDogu(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		dogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: false,
			},
		}

		expectedUpdateDogu := &v2.Dogu{
			Spec: v2.DoguSpec{
				Name:    "postgresql",
				Stopped: true,
			},
		}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(dogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(testCtx, expectedUpdateDogu, mock.Anything, metav1.UpdateOptions{}).Return(dogu, nil)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
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
		doguRegistryMock := newMockDoguDescriptorGetter(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)
		sut := defaultDoguInterActor{
			doguDescriptorGetter: doguRegistryMock,
		}

		// when
		err := sut.StartAllDogus(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error getting all dogus while starting:")
	})

	t.Run("should start dogus in correct order", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		doguRegistryMock := newMockDoguDescriptorGetter(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql", Type: core.DependencyTypeDogu}}}}, nil)

		postgresqlDogu := &v2.Dogu{Spec: v2.DoguSpec{Name: "postgresql", Stopped: true}}
		redmineDogu := &v2.Dogu{Spec: v2.DoguSpec{Name: "redmine", Stopped: true}}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(mock.Anything, "postgresql", metav1.GetOptions{}).Return(postgresqlDogu, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, postgresqlDogu, mock.Anything, metav1.UpdateOptions{}).Return(postgresqlDogu, nil).Run(func(ctx context.Context, dogu *v2.Dogu, modifySpecFn func(v2.DoguSpec) v2.DoguSpec, opts metav1.UpdateOptions) {
			// This ensures the start order.
			doguClientMock.EXPECT().Get(mock.Anything, "redmine", metav1.GetOptions{}).Return(redmineDogu, nil)
			doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, redmineDogu, mock.Anything, metav1.UpdateOptions{}).Return(redmineDogu, nil)
			redmineWatcher := watch.NewFake()
			doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=redmine"}).Return(redmineWatcher, nil)
			go func() {
				time.Sleep(2 * time.Second)
				redmineWatcher.Action(watch.Modified, redmineDogu)
			}()
		})

		postgresqlWatcher := watch.NewFake()
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=postgresql"}).Return(postgresqlWatcher, nil)

		go func() {
			time.Sleep(1 * time.Second)
			postgresqlWatcher.Action(watch.Modified, postgresqlDogu)
		}()

		sut := defaultDoguInterActor{
			doguDescriptorGetter: doguRegistryMock,
			doguConfigRepository: repository.DoguConfigRepository{},
			doguClient:           doguClientMock,
		}

		// when
		err := sut.StartAllDogus(testCtx)

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_StopAllDogus(t *testing.T) {
	t.Run("should return error on get all dogus error", func(t *testing.T) {
		// given

		doguRegistryMock := newMockDoguDescriptorGetter(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)
		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguDescriptorGetter: doguRegistryMock,
		}

		// when
		err := sut.StopAllDogus(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error getting all dogus while stopping:")
	})

	t.Run("should stop dogus in correct order", func(t *testing.T) {
		// given
		oldWaitTimeout := waitTimeout
		waitTimeout = time.Second * 10
		defer func() { waitTimeout = oldWaitTimeout }()

		doguRegistryMock := newMockDoguDescriptorGetter(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return([]*core.Dogu{{Name: "official/postgresql"}, {Name: "official/redmine", Dependencies: []core.Dependency{{Name: "postgresql", Type: core.DependencyTypeDogu}}}}, nil)

		postgresqlDoguStopped := &v2.Dogu{Spec: v2.DoguSpec{Name: "postgresql", Stopped: true}, Status: v2.DoguStatus{Stopped: true}}
		redmineDoguStopped := &v2.Dogu{Spec: v2.DoguSpec{Name: "redmine", Stopped: true}, Status: v2.DoguStatus{Stopped: true}}

		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(mock.Anything, "redmine", metav1.GetOptions{}).Return(redmineDoguStopped, nil)
		doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, redmineDoguStopped, mock.Anything, metav1.UpdateOptions{}).Return(redmineDoguStopped, nil).Run(func(ctx context.Context, dogu *v2.Dogu, modifySpecFn func(v2.DoguSpec) v2.DoguSpec, opts metav1.UpdateOptions) {
			// This ensures the stop order.
			doguClientMock.EXPECT().Get(mock.Anything, "postgresql", metav1.GetOptions{}).Return(postgresqlDoguStopped, nil)
			doguClientMock.EXPECT().UpdateSpecWithRetry(mock.Anything, postgresqlDoguStopped, mock.Anything, metav1.UpdateOptions{}).Return(postgresqlDoguStopped, nil)
			postgresqlWatcher := watch.NewFake()
			doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=postgresql"}).Return(postgresqlWatcher, nil)
			go func() {
				time.Sleep(1 * time.Second)
				postgresqlWatcher.Action(watch.Modified, redmineDoguStopped)
			}()
		})

		redmineWatcher := watch.NewFake()
		doguClientMock.EXPECT().Watch(mock.Anything, metav1.ListOptions{FieldSelector: "metadata.name=redmine"}).Return(redmineWatcher, nil)
		go func() {
			time.Sleep(1 * time.Second)
			redmineWatcher.Action(watch.Modified, redmineDoguStopped)
		}()

		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguDescriptorGetter: doguRegistryMock,
			doguClient:           doguClientMock,
		}

		// when
		err := sut.StopAllDogus(testCtx)

		// then
		require.NoError(t, err)
	})
}

func Test_defaultDoguInterActor_SetLogLevelInAllDogus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguDescriptorGetter(t)
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
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), common.SimpleName("postgresql")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), common.SimpleName("redmine")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).RunAndReturn(func(ctx context.Context, doguConfig config.DoguConfig) (config.DoguConfig, error) {
			get, b := doguConfig.Get("logging/root")
			require.True(t, b)
			assert.Equal(t, "DEBUG", get.String())
			return doguConfig, nil
		}).Times(2)

		sut := defaultDoguInterActor{
			doguConfigRepository: doguConfigRepositoryMock,
			doguDescriptorGetter: doguRegistryMock,
		}

		// when
		err := sut.SetLogLevelInAllDogus(testCtx, "DEBUG")

		// then
		require.NoError(t, err)
	})

	t.Run("should return errors on error updating log levels", func(t *testing.T) {
		// given
		doguRegistryMock := newMockDoguDescriptorGetter(t)
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
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), common.SimpleName("postgresql")).Return(doguConfigA, nil)
		doguConfigRepositoryMock.EXPECT().Get(context.TODO(), common.SimpleName("redmine")).Return(doguConfigB, nil)
		doguConfigRepositoryMock.EXPECT().Update(context.TODO(), mock.Anything).RunAndReturn(func(ctx context.Context, doguConfig config.DoguConfig) (config.DoguConfig, error) {
			get, b := doguConfig.Get("logging/root")
			require.True(t, b)
			assert.Equal(t, "DEBUG", get.String())
			return doguConfig, assert.AnError
		}).Times(2)

		sut := defaultDoguInterActor{
			doguConfigRepository: doguConfigRepositoryMock,
			doguDescriptorGetter: doguRegistryMock,
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
		doguRegistryMock := newMockDoguDescriptorGetter(t)
		doguRegistryMock.EXPECT().GetCurrentOfAll(testCtx).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguConfigRepository: repository.DoguConfigRepository{},
			doguDescriptorGetter: doguRegistryMock,
		}

		// when
		err := sut.SetLogLevelInAllDogus(testCtx, "DEBUG")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error getting all dogus while setting log-level:")
	})
}

func Test_defaultDoguInterActor_checkIfDoguInDesiredStopState(t *testing.T) {
	t.Run("should return error on dogu get error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(nil, assert.AnError)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		_, err := sut.checkIfDoguInDesiredStopState(testCtx, "redmine")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get dogu redmine")
	})

	t.Run("should stop if dogu is in desired state", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(&v2.Dogu{
			Spec:   v2.DoguSpec{Name: "redmine", Stopped: true},
			Status: v2.DoguStatus{Stopped: true},
		}, nil)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		isInDesiredState, err := sut.checkIfDoguInDesiredStopState(testCtx, "redmine")

		// then
		require.NoError(t, err)
		assert.True(t, isInDesiredState)
	})

	t.Run("should not stop if dogu is NOT in desired state", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().Get(testCtx, "redmine", metav1.GetOptions{}).Return(&v2.Dogu{
			Spec:   v2.DoguSpec{Name: "redmine", Stopped: true},
			Status: v2.DoguStatus{Stopped: false},
		}, nil)

		sut := defaultDoguInterActor{
			doguClient: doguClientMock,
		}

		// when
		isInDesiredState, err := sut.checkIfDoguInDesiredStopState(testCtx, "redmine")

		// then
		require.NoError(t, err)
		assert.False(t, isInDesiredState)
	})
}
