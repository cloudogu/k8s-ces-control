package config

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	cesregistry "github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	logLevelEnvironmentVariable = "LOG_LEVEL"
	defaultLogLevel             = logrus.WarnLevel

	stagingEnvironmentVariable = "STAGE"
	stageProduction            = "production"
	stageDevelopment           = "development"

	namespaceEnvironmentVariable = "NAMESPACE"
)

type ClusterClient struct {
	EcoSystemApi ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

var currentStage = "development"

// CreateClusterClient creates a new kubernetes.Interface given the locally available cluster configurations.
func CreateClusterClient() (ClusterClient, error) {
	clusterConfig, err := ctrl.GetConfig()
	if err != nil {
		return ClusterClient{}, fmt.Errorf("failed to load cluster configuration: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return ClusterClient{}, fmt.Errorf("failed to create kubernetes client")
	}

	doguClient, err := ecoSystem.NewForConfig(clusterConfig)
	if err != nil {
		return ClusterClient{}, fmt.Errorf("failed to create dogu client")
	}

	return ClusterClient{EcoSystemApi: doguClient, Interface: k8sClient}, nil
}

// ConfigureApplication performs the default configuration for the control app including configuring the logging and
// current stage of the system.
func ConfigureApplication() error {
	err := configureLogLevel()
	if err != nil {
		return err
	}

	err = configureNamespace()
	if err != nil {
		return err
	}

	err = configureCurrentStage()
	if err != nil {
		return err
	}

	return nil
}

// IsDevelopmentStage return true whether the current stage is set to development.
func IsDevelopmentStage() bool {
	return currentStage == stageDevelopment
}

func configureCurrentStage() error {
	stage, ok := os.LookupEnv(stagingEnvironmentVariable)
	if !ok {
		logrus.Printf("No stage was set via the environment variable [%s]. Using default stage [production].", stagingEnvironmentVariable)
		return nil
	}

	currentStage = stage
	if stage == stageProduction {
		logrus.Println("Using stage [production].")
	} else if stage == stageDevelopment {
		logrus.Warningf("Using stage [development]. This is not recommended for production systems!")
	} else {
		return fmt.Errorf("found invalid value [%s] for environment variable [%s], only the values [production, development] are valid values", stage, stagingEnvironmentVariable)
	}

	return nil
}

var CurrentNamespace = ""

func configureNamespace() error {
	namespace, ok := os.LookupEnv(namespaceEnvironmentVariable)
	if !ok {
		logrus.Printf("No namespace was set via the environment variable [%s]. A namespace is required.", namespaceEnvironmentVariable)
		return nil
	}

	CurrentNamespace = namespace
	if CurrentNamespace == "" {
		return fmt.Errorf("found invalid value for namespace [%s]: namespace cannot be empty: set valid value with environment variable [%s]", CurrentNamespace, namespaceEnvironmentVariable)
	}

	logrus.Printf("Using namespace [%s].", CurrentNamespace)
	return nil
}

func configureLogLevel() error {
	printCloudoguLogo()

	logLevel, ok := os.LookupEnv(logLevelEnvironmentVariable)
	if !ok {
		logrus.SetLevel(defaultLogLevel)
		return nil
	}

	logLevelParsed, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level %s to logrus level: %w", logLevel, err)
	}

	logrus.Infof("Using log level: %s", logLevelParsed)
	logrus.StandardLogger().SetLevel(logLevelParsed)
	return nil
}

func printCloudoguLogo() {
	logrus.Println("                                     ./////,                    ")
	logrus.Println("                                 ./////==//////*                ")
	logrus.Println("                                ////.  ___   ////.              ")
	logrus.Println("                         ,**,. ////  ,////A,  */// ,**,.        ")
	logrus.Println("                    ,/////////////*  */////*  *////////////A    ")
	logrus.Println("                   ////'        \\VA.   '|'   .///'       '///*  ")
	logrus.Println("                  *///  .*///*,         |         .*//*,   ///* ")
	logrus.Println("                  (///  (//////)**--_./////_----*//////)   ///) ")
	logrus.Println("                   V///   '°°°°      (/////)      °°°°'   ////  ")
	logrus.Println("                    V/////(////////\\. '°°°' ./////////(///(/'   ")
	logrus.Println("                       'V/(/////////////////////////////V'      ")
}

func GetCesRegistry() (cesregistry.Registry, error) {
	cesReg, err := cesregistry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", CurrentNamespace)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CES registry: %w", err)
	}

	return cesReg, nil
}
