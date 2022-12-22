package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const (
	port           = ":50051"
	environmentDev = "development"
)

var (
	// Version of the application
	Version string
)

func main() {
	app := cli.NewApp()
	app.Name = "k8s-ces-control"
	app.Usage = "Control you EcoSystem with ease!"
	app.Version = Version
	app.Flags = createGlobalFlags()
	app.Before = configureApplication
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		startServerCommand(),
	}

	err := app.Run(os.Args)
	os.Exit(checkMainError(err))
}

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "define log level",
			Value: "warn",
		},
		&cli.BoolFlag{
			Name:  "show-stack",
			Usage: "show stacktrace on errors",
		},
	}
}

func configureLogging(logLevel string) error {
	logLevelParsed, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("could not parse log level %s to logrus level: %w", logLevel, err)
	}

	logrus.StandardLogger().SetLevel(logLevelParsed)
	return nil
}

func isPrintStack() bool {
	for _, arg := range os.Args {
		if arg == "--show-stack" {
			return true
		}
	}
	return false
}

func checkMainError(err error) int {
	if err != nil {
		if isPrintStack() {
			logrus.Errorf("%+v\n", err)
		} else {
			logrus.Errorf("%+s\n", err)
		}
		return 1
	}
	return 0
}

func printCallWithArguments() {
	logrus.Debugf("Executing command: %s; cesappd version: %s", os.Args, Version)
}

func configureApplication(cliCtx *cli.Context) error {
	logLevel := cliCtx.String("log-level")
	err := configureLogging(logLevel)
	if err != nil {
		return err
	}

	printCallWithArguments()
	return nil
}

func registerServices(grpcServer *grpc.Server) error {
	//pbLogging.RegisterDoguLogMessagesServer(grpcServer, loggingService.NewLoggingService())
	//adminService, err := doguAdministration.NewDoguAdministrationService()
	//if err != nil {
	//	return fmt.Errorf("failed to create new dogu administration service: %w", err)
	//}

	//pbDoguAdministration.RegisterDoguAdministrationServer(grpcServer, adminService)
	//healthService, err := health.NewDoguHealthService()
	//if err != nil {
	//	return fmt.Errorf("failed to create new dogu health service: %w", err)
	//}

	//pbHealth.RegisterDoguHealthServer(grpcServer, healthService)
	//backupService, err := backup.NewBackupManagementService()
	//if err != nil {
	//	return fmt.Errorf("failed to create new backup service: %w", err)
	//}

	//pbBackup.RegisterBackupManagementServer(grpcServer, backupService)
	//maintenanceService, err := maintenance.NewDebugModeService()
	//if err != nil {
	//	return fmt.Errorf("failed to create new maintenance service: %w", err)
	//}

	//pbMaintenance.RegisterDebugModeServer(grpcServer, maintenanceService)
	//supportArchiveService, err := maintenance.NewSupportArchiveService()
	//if err != nil {
	//	return fmt.Errorf("failed to create new support archive service: %w", err)
	//}

	//pbMaintenance.RegisterSupportArchiveServer(grpcServer, supportArchiveService)
	return nil
}

func registerServerForServiceDiscovery(grpcServer *grpc.Server) {
	reflection.Register(grpcServer)
}

func startServerCommand() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Aliases:   []string{"s"},
		Usage:     "starts the service and listens on port 50051",
		ArgsUsage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Usage: "define execution environment",
				Value: "production",
			},
		},
		Action: startServerAction,
	}
}

func startServerAction(ctx *cli.Context) error {
	executionEnvironment := ctx.String("env")
	logrus.Infof("starting cesappd in %s environment...", executionEnvironment)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	//creds, err := credentials.NewServerTLSFromFile("/etc/ces/cesappd/server.crt", "/etc/ces/cesappd/server.key")
	//if err != nil {
	//	logrus.Fatalf("Failed to setup TLS: %v", err)
	//}

	grpcServer := grpc.NewServer()
	err = registerServices(grpcServer)
	if err != nil {
		logrus.Fatalf("failed to register services: %w", err)
		return err
	}

	if executionEnvironment == environmentDev {
		logrus.Debugln("register service discovery")
		registerServerForServiceDiscovery(grpcServer)
	}

	logrus.Infof("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
