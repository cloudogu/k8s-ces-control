package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	logLevelEnvironmentVariable = "LOG_LEVEL"
	defaultLogLevel             = logrus.WarnLevel

	stagingEnvironmentVariable = "STAGE"
	stageProduction            = "production"
	stageDevelopment           = "development"
)

var currentStage = "development"

// ConfigureApplication performs the default configuration for the control app including configuring the logging and
// current stage of the system.
func ConfigureApplication() error {
	err := configureLogLevel()
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
