package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	pbBackup "github.com/cloudogu/ces-control-api/generated/backup"
	pbDoguAdministration "github.com/cloudogu/ces-control-api/generated/doguAdministration"
	pgHealth "github.com/cloudogu/ces-control-api/generated/health"
	pbLogging "github.com/cloudogu/ces-control-api/generated/logging"
	pbMaintenance "github.com/cloudogu/ces-control-api/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/backup"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	pbDebug "github.com/cloudogu/k8s-ces-control/packages/debug"
	"github.com/cloudogu/k8s-ces-control/packages/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/packages/doguHealth"
	"github.com/cloudogu/k8s-ces-control/packages/doguinteraction"
	"github.com/cloudogu/k8s-ces-control/packages/logging"
	"github.com/cloudogu/k8s-ces-control/packages/supportArchive"
	"github.com/cloudogu/k8s-ces-control/packages/util"
	"github.com/cloudogu/k8s-registry-lib/dogu"
	"github.com/cloudogu/k8s-registry-lib/repository"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var (
	// Version of the application
	Version string
)

func main() {
	err := startCesControl()

	if err != nil {
		logrus.Errorf("%+v\n", err)
		os.Exit(1)
	}

	logrus.Infoln("Gracefully exited k8s-ces-control")
	os.Exit(0)
}

func startCesControl() error {
	err := config.ConfigureApplication()
	if err != nil {
		return err
	}

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

	logrus.Infoln("Starting k8s-ces-control")
	err = app.Run(os.Args)
	return err
}

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "show-stack",
			Usage: "show stacktrace on errors",
		},
	}
}

func printCallWithArguments() {
	logrus.Debugf("Executing command: %s; k8s-ces-control version: %s", os.Args, Version)
}

func configureApplication(_ *cli.Context) error {
	printCallWithArguments()
	return nil
}

func registerServices(client clusterClient, grpcServer grpc.ServiceRegistrar) error {
	lokiLogProvider := logging.NewLokiLogProvider(
		config.CurrentLokiGatewayConfig.Url,
		config.CurrentLokiGatewayConfig.Username,
		config.CurrentLokiGatewayConfig.Password,
	)

	configMapClient := client.CoreV1().ConfigMaps(config.CurrentNamespace)
	doguDescriptorGetter := util.NewDoguGetter(
		dogu.NewDoguVersionRegistry(configMapClient),
		dogu.NewLocalDoguDescriptorRepository(configMapClient),
	)

	globalConfig := repository.NewGlobalConfigRepository(configMapClient)
	doguConfig := repository.NewDoguConfigRepository(configMapClient)

	doguClient := client.Dogus(config.CurrentNamespace)
	doguRestartClient := client.DoguRestarts(config.CurrentNamespace)

	doguInterActor := doguinteraction.NewDefaultDoguInterActor(doguConfig, doguClient, doguRestartClient, doguDescriptorGetter)

	supportArchiveClient := client.SupportArchives(config.CurrentNamespace)

	backupClient := client.Backups(config.CurrentNamespace)
	restoreClient := client.Restores(config.CurrentNamespace)
	backupScheduleClient := client.BackupSchedules(config.CurrentNamespace)

	debugModeClient := client.DebugMode(config.CurrentNamespace)

	loggingService := logging.NewLoggingService(
		lokiLogProvider,
		doguConfig,
		doguInterActor,
		doguDescriptorGetter,
		client.Dogus(config.CurrentNamespace),
	)

	doguAdministrationServer := doguAdministration.NewDoguAdministrationServer(client, doguDescriptorGetter, doguInterActor, loggingService)

	pbLogging.RegisterDoguLogMessagesServer(grpcServer, loggingService)
	pbDoguAdministration.RegisterDoguAdministrationServer(grpcServer, doguAdministrationServer)
	pgHealth.RegisterDoguHealthServer(grpcServer, doguHealth.NewDoguHealthService(client.Dogus(config.CurrentNamespace)))
	debugModeService := pbDebug.NewDebugModeService(debugModeClient, doguInterActor, doguConfig, globalConfig, doguDescriptorGetter, client, config.CurrentNamespace)
	pbMaintenance.RegisterDebugModeServer(grpcServer, debugModeService)
	supportArchiveService := supportArchive.NewSupportArchiveService(supportArchiveClient, &http.Client{})
	pbMaintenance.RegisterSupportArchiveServer(grpcServer, supportArchiveService)
	watcher := pbDebug.NewDefaultConfigMapRegistryWatcher(configMapClient, debugModeService)
	watcher.StartWatch(context.Background())
	backupService := backup.NewBackupService(backupClient, restoreClient, backupScheduleClient)
	pbBackup.RegisterBackupManagementServer(grpcServer, backupService)
	// health endpoint used to determine the healthiness of the app
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	return nil
}

func registerServerForServiceDiscovery(grpcServer *grpc.Server) {
	reflection.Register(grpcServer)
}

func startServerCommand() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Aliases:   []string{"s"},
		Usage:     fmt.Sprintf("starts the service and listens on port %s", port),
		ArgsUsage: "",
		Flags:     []cli.Flag{},
		Action:    startServerAction,
	}
}

func startServerAction(_ *cli.Context) error {
	config.PrintCloudoguLogo()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	client, err := config.CreateClusterClient()
	if err != nil {
		return fmt.Errorf("failed to create cluster client")
	}

	grpcServer := grpc.NewServer()
	err = registerServices(client, grpcServer)
	if err != nil {
		logrus.Fatalf("failed to register services: %s", err.Error())
		return err
	}

	if config.IsDevelopmentStage() {
		logrus.Debugln("Register k8s-ces-control to be used with the service discovery")
		registerServerForServiceDiscovery(grpcServer)
	}

	logrus.Infof("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
