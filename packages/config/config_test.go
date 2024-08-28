package config

import (
	"os"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"

	doguApiV1 "github.com/cloudogu/k8s-dogu-operator/api/v1"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	logLevelEnv  = "LOG_LEVEL"
	namespaceEnv = "NAMESPACE"
	stageEnv     = "STAGE"
)

func TestConfigureApplication(t *testing.T) {
	t.Run("should set log level, namespace and stage from env vars", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "info")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		previousNamespace, namespaceExists := setEnv(t, namespaceEnv, "ecosystem")
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		previousStage, stageExists := setEnv(t, stageEnv, "development")
		defer cleanupEnv(t, stageEnv, previousStage, stageExists)

		// when
		err := ConfigureApplication()

		// then
		require.NoError(t, err)
		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
		assert.Equal(t, "ecosystem", CurrentNamespace)
		assert.Equal(t, "development", currentStage)
	})
	t.Run("should set log level, namespace and stage from env vars", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "info")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		previousNamespace, namespaceExists := setEnv(t, namespaceEnv, "ecosystem")
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		previousStage, stageExists := setEnv(t, stageEnv, "production")
		defer cleanupEnv(t, stageEnv, previousStage, stageExists)

		// when
		err := ConfigureApplication()

		// then
		require.NoError(t, err)
		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
		assert.Equal(t, "ecosystem", CurrentNamespace)
		assert.Equal(t, "production", currentStage)
	})
	t.Run("should set default log level and stage if not set", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := unsetEnv(t, logLevelEnv)
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		previousStage, stageExists := unsetEnv(t, stageEnv)
		defer cleanupEnv(t, stageEnv, previousStage, stageExists)

		previousNamespace, namespaceExists := setEnv(t, namespaceEnv, "ecosystem")
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		// when
		err := ConfigureApplication()

		// then
		require.NoError(t, err)
		assert.Equal(t, logrus.WarnLevel, logrus.GetLevel())
		assert.Equal(t, "ecosystem", CurrentNamespace)
		assert.Equal(t, "production", currentStage)
	})
	t.Run("should fail when invalid log level is set", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "banana")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		// when
		err := ConfigureApplication()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not parse log level banana to logrus level")
	})
	t.Run("should fail if NAMESPACE env var is not set", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousNamespace, namespaceExists := unsetEnv(t, namespaceEnv)
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "info")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		// when
		err := ConfigureApplication()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "found invalid value for namespace []: namespace cannot be empty: set valid value with environment variable [NAMESPACE]")

		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
	})
	t.Run("should fail if NAMESPACE env var is set to empty string", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "info")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		previousNamespace, namespaceExists := setEnv(t, namespaceEnv, "")
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		// when
		err := ConfigureApplication()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "found invalid value for namespace []: namespace cannot be empty: set valid value with environment variable [NAMESPACE]")

		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
	})
	t.Run("should fail if stage is invalid", func(t *testing.T) {
		// given
		previousStageVar := currentStage
		defer func() { currentStage = previousStageVar }()
		previousNamespaceVar := CurrentNamespace
		defer func() { CurrentNamespace = previousNamespaceVar }()
		previousLogrusLogLevel := logrus.GetLevel()
		defer logrus.SetLevel(previousLogrusLogLevel)

		previousLogLevel, logLevelExists := setEnv(t, logLevelEnv, "info")
		defer cleanupEnv(t, logLevelEnv, previousLogLevel, logLevelExists)

		previousNamespace, namespaceExists := setEnv(t, namespaceEnv, "ecosystem")
		defer cleanupEnv(t, namespaceEnv, previousNamespace, namespaceExists)

		previousStage, stageExists := setEnv(t, stageEnv, "banana")
		defer cleanupEnv(t, stageEnv, previousStage, stageExists)

		// when
		err := ConfigureApplication()

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "found invalid value [banana] for environment variable [STAGE], only the values [production, development] are valid values")

		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
		assert.Equal(t, "ecosystem", CurrentNamespace)
	})
}

func setEnv(t *testing.T, key, value string) (previous string, exists bool) {
	t.Helper()
	previous, exists = os.LookupEnv(key)
	err := os.Setenv(key, value)
	require.NoError(t, err)

	return previous, exists
}

func unsetEnv(t *testing.T, key string) (previous string, exists bool) {
	t.Helper()
	previous, exists = os.LookupEnv(key)
	err := os.Unsetenv(key)
	require.NoError(t, err)

	return previous, exists
}

func cleanupEnv(t *testing.T, key, previousValue string, exists bool) {
	t.Helper()
	if exists {
		err := os.Setenv(key, previousValue)
		require.NoError(t, err)
	} else {
		err := os.Unsetenv(key)
		require.NoError(t, err)
	}
}

func TestIsDevelopmentStage(t *testing.T) {
	t.Run("should be true for stage development", func(t *testing.T) {
		// given
		previousStage := currentStage
		defer func() { currentStage = previousStage }()
		currentStage = "development"

		// when
		actual := IsDevelopmentStage()

		// then
		assert.True(t, actual)
	})
	t.Run("should be true for stage production", func(t *testing.T) {
		// given
		previousStage := currentStage
		defer func() { currentStage = previousStage }()
		currentStage = "production"

		// when
		actual := IsDevelopmentStage()

		// then
		assert.False(t, actual)
	})
}

func TestPrintCloudoguLogo(t *testing.T) {
	t.Run("should print logo", func(t *testing.T) {
		// given
		defer logrus.SetOutput(os.Stderr)
		buf := new(strings.Builder)
		logrus.SetOutput(buf)

		// when
		PrintCloudoguLogo()

		// then
		lines := strings.Split(buf.String(), "\n")
		assert.Len(t, lines, 12)
		assert.Contains(t, lines[0], "                                     ./////,                    ")
		assert.Contains(t, lines[1], "                                 ./////==//////*                ")
		assert.Contains(t, lines[2], "                                ////.  ___   ////.              ")
		assert.Contains(t, lines[3], "                         ,**,. ////  ,////A,  */// ,**,.        ")
		assert.Contains(t, lines[4], "                    ,/////////////*  */////*  *////////////A    ")
		assert.Contains(t, lines[5], "                   ////'        \\\\VA.   '|'   .///'       '///* ")
		assert.Contains(t, lines[6], "                  *///  .*///*,         |         .*//*,   ///* ")
		assert.Contains(t, lines[7], "                  (///  (//////)**--_./////_----*//////)   ///) ")
		assert.Contains(t, lines[8], "                   V///   '°°°°      (/////)      °°°°'   ////  ")
		assert.Contains(t, lines[9], "                    V/////(////////\\\\. '°°°' ./////////(///(/'  ")
		assert.Contains(t, lines[10], "                       'V/(/////////////////////////////V'      ")
	})
}

func TestCreateClusterClient(t *testing.T) {
	t.Run("should fail to load cluster config", func(t *testing.T) {
		// given
		originalConfigFunc := ctrl.GetConfig
		defer func() { ctrl.GetConfig = originalConfigFunc }()

		failingConfigFunc := func() (*rest.Config, error) {
			return nil, assert.AnError
		}
		ctrl.GetConfig = failingConfigFunc

		// when
		actual, err := CreateClusterClient()

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to load cluster configuration")
	})
	t.Run("should fail to create kubernetes client", func(t *testing.T) {
		// given
		originalConfigFunc := ctrl.GetConfig
		defer func() { ctrl.GetConfig = originalConfigFunc }()

		incorrectConfigFunc := func() (*rest.Config, error) {
			return &rest.Config{
				ExecProvider: &api.ExecConfig{},
				AuthProvider: &api.AuthProviderConfig{},
			}, nil
		}
		ctrl.GetConfig = incorrectConfigFunc

		// when
		actual, err := CreateClusterClient()

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "failed to create kubernetes client")
	})
	t.Run("should fail to create kubernetes client", func(t *testing.T) {
		// given
		originalConfigFunc := ctrl.GetConfig
		defer func() { ctrl.GetConfig = originalConfigFunc }()

		dummyConfigFunc := func() (*rest.Config, error) {
			return &rest.Config{}, nil
		}
		ctrl.GetConfig = dummyConfigFunc

		originalAddSchemeFunc := doguApiV1.AddToScheme
		defer func() { doguApiV1.AddToScheme = originalAddSchemeFunc }()

		failingAddSchemeFunc := func(s *runtime.Scheme) error {
			return assert.AnError
		}
		doguApiV1.AddToScheme = failingAddSchemeFunc

		// when
		actual, err := CreateClusterClient()

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create dogu client")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		originalConfigFunc := ctrl.GetConfig
		defer func() { ctrl.GetConfig = originalConfigFunc }()

		dummyConfigFunc := func() (*rest.Config, error) {
			return &rest.Config{}, nil
		}
		ctrl.GetConfig = dummyConfigFunc

		// when
		actual, err := CreateClusterClient()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, actual.EcoSystemV1Alpha1Interface)
		assert.NotNil(t, actual.Interface)
		assert.NotNil(t, actual.BlueprintLister)
	})
}
