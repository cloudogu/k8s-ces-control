package config

import (
	"fmt"
	"github.com/bombsimon/logrusr/v2"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/cesapp-lib/core"
	cesregistry "github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

const (
	logLevelEnvironmentVariable = "LOG_LEVEL"
	defaultLogLevel             = logrus.WarnLevel

	stagingEnvironmentVariable = "STAGE"
	stageProduction            = "production"
	stageDevelopment           = "development"

	namespaceEnvironmentVariable = "NAMESPACE"
)

type clusterClient struct {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

var currentStage = stageProduction

// CreateClusterClient creates a new kubernetes.Interface given the locally available cluster configurations.
func CreateClusterClient() (*clusterClient, error) {
	clusterConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load cluster configuration: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	doguClient, err := ecoSystem.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dogu client: %w", err)
	}

	return &clusterClient{EcoSystemV1Alpha1Interface: doguClient, Interface: k8sClient}, nil
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
		currentStage = stageProduction
		return nil
	}

	if stage == stageProduction {
		logrus.Infoln("Using stage [production].")
	} else if stage == stageDevelopment {
		logrus.Warningf("Using stage [development]. This is not recommended for production systems!")
	} else {
		return fmt.Errorf("found invalid value [%s] for environment variable [%s], only the values [production, development] are valid values", stage, stagingEnvironmentVariable)
	}

	currentStage = stage
	return nil
}

// CurrentNamespace contains the namespace from the k8s-ecs-control pod.
var CurrentNamespace = ""

func configureNamespace() error {
	namespace, ok := os.LookupEnv(namespaceEnvironmentVariable)
	if !ok {
		logrus.Errorf("No namespace was set via the environment variable [%s]. A namespace is required.", namespaceEnvironmentVariable)
	}

	CurrentNamespace = namespace
	if CurrentNamespace == "" {
		return fmt.Errorf("found invalid value for namespace [%s]: namespace cannot be empty: set valid value with environment variable [%s]", CurrentNamespace, namespaceEnvironmentVariable)
	}

	logrus.Infof("Using namespace [%s].", CurrentNamespace)

	return nil
}

func configureLogLevel() error {
	logLevel, ok := os.LookupEnv(logLevelEnvironmentVariable)
	if !ok {
		logrus.SetLevel(defaultLogLevel)
		return nil
	}

	logLevelParsed, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level %s to logrus level: %w", logLevel, err)
	}

	logrus.StandardLogger().SetLevel(logLevelParsed)
	logrus.Infof("Using log level: %s", logLevelParsed)

	// create logrus logger that can be styled and formatted
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{})
	logrusLog.SetLevel(logLevelParsed)

	// convert logrus logger to logr logger
	logrusLogrLogger := logrusr.New(logrusLog)
	log.SetLogger(logrusLogrLogger)

	return nil
}

// PrintCloudoguLogo prints the awesome cloudogu logo.
func PrintCloudoguLogo() {
	logrus.Println("                                     ./////,                    ")
	logrus.Println("                                 ./////==//////*                ")
	logrus.Println("                                ////.  ___   ////.              ")
	logrus.Println("                         ,**,. ////  ,////A,  */// ,**,.        ")
	logrus.Println("                    ,/////////////*  */////*  *////////////A    ")
	logrus.Println("                   ////'        \\VA.   '|'   .///'       '///* ")
	logrus.Println("                  *///  .*///*,         |         .*//*,   ///* ")
	logrus.Println("                  (///  (//////)**--_./////_----*//////)   ///) ")
	logrus.Println("                   V///   '°°°°      (/////)      °°°°'   ////  ")
	logrus.Println("                    V/////(////////\\. '°°°' ./////////(///(/'  ")
	logrus.Println("                       'V/(/////////////////////////////V'      ")
}

// GetCesRegistry return the ces registry in the current namespace.
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
