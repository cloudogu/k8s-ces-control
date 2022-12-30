package main

import (
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	cesregistry "github.com/cloudogu/cesapp-lib/registry"
	pbDoguAdministration "github.com/cloudogu/k8s-ces-control/generated/doguAdministration"
	pgHealth "github.com/cloudogu/k8s-ces-control/generated/health"
	pbLogging "github.com/cloudogu/k8s-ces-control/generated/logging"
	pbMaintenance "github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/account"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/cloudogu/k8s-ces-control/packages/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/packages/doguHealth"
	"github.com/cloudogu/k8s-ces-control/packages/logging"
	"github.com/cloudogu/k8s-ces-control/packages/maintenance"
	"github.com/cloudogu/k8s-ces-control/packages/ssl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
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
		createServiceAccountCommand(),
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

func registerServices(grpcServer *grpc.Server) error {
	cesReg, err := cesregistry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", config.CurrentNamespace)},
	})
	if err != nil {
		return fmt.Errorf("failed to create CES registry: %w", err)
	}

	client, err := config.CreateClusterClient()
	if err != nil {
		return fmt.Errorf("failed to create cluster client")
	}

	pbLogging.RegisterDoguLogMessagesServer(grpcServer, logging.NewLoggingService(client))
	pbDoguAdministration.RegisterDoguAdministrationServer(grpcServer, doguAdministration.NewDoguAdministrationServer(client, cesReg))
	pgHealth.RegisterDoguHealthServer(grpcServer, doguHealth.NewDoguHealthService(client))
	pbMaintenance.RegisterDebugModeServer(grpcServer, maintenance.NewDebugModeService())

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
		Usage:     "starts the service and listens on port 50051",
		ArgsUsage: "",
		Flags:     []cli.Flag{},
		Action:    startServerAction,
	}
}

func startServerAction(_ *cli.Context) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	creds, err := readTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	err = registerServices(grpcServer)
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

func readTLSCredentials() (credentials.TransportCredentials, error) {
	cesReg, err := getCesRegistry()
	if err != nil {
		return nil, err
	}

	reader := ssl.NewManager(cesReg.GlobalConfig())
	return reader.GetCertificateCredentials()
}

func getCesRegistry() (cesregistry.Registry, error) {
	cesReg, err := cesregistry.New(core.Registry{
		Type:      "etcd",
		Endpoints: []string{fmt.Sprintf("http://etcd.%s.svc.cluster.local:4001", config.CurrentNamespace)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CES registry: %w", err)
	}

	return cesReg, nil
}

func createServiceAccountCommand() *cli.Command {
	managerCreator := func(serviceName string) (account.ServiceAccountManager, error) {
		cesRegistry, err := getCesRegistry()
		if err != nil {
			return account.ServiceAccountManager{}, err
		}
		manager, err := account.NewServiceAccountManager(serviceName, cesRegistry)
		if err != nil {
			return account.ServiceAccountManager{}, err
		}
		return manager, nil
	}
	return &cli.Command{
		Name:      "service-account-create",
		Usage:     "creates a service account for the given service",
		ArgsUsage: "SERVICE_NAME",
		Action:    getServiceAccountAction(managerCreator),
	}
}

type serviceAccountManagerCreator func(serviceName string) (account.ServiceAccountManager, error)

func getServiceAccountAction(getManager serviceAccountManagerCreator) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		err := validateArgsCount(ctx, 1)
		if err != nil {
			return createServiceAccountErr("", "create", err)
		}
		serviceName := ctx.Args().First()
		accountManager, err := getManager(serviceName)
		if err != nil {
			return createServiceAccountErr(serviceName, "create", err)
		}
		result, err := accountManager.Create(ctx.Context)
		if err != nil {
			return createServiceAccountErr(serviceName, "create", err)
		}
		fmt.Println(result)
		return nil
	}
}

func createServiceAccountErr(serviceName, action string, err error) error {
	return fmt.Errorf("failed to %s service account for service %s: %w", action, serviceName, err)
}

func validateArgsCount(ctx *cli.Context, requiredCount int) error {
	actualArgsCount := ctx.Args().Len()
	if actualArgsCount < requiredCount {
		return errors.Errorf("The command '%s' requires at least %d argument(s)", ctx.Command.Name, requiredCount)
	}
	return nil
}
