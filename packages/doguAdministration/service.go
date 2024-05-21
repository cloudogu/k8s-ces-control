package doguAdministration

import (
	"context"
	"fmt"
	pb "github.com/cloudogu/ces-control-api/generated/doguAdministration"
	"github.com/cloudogu/ces-control-api/generated/types"
	"github.com/cloudogu/cesapp-lib/core"
	v1bp "github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/kubernetes/blueprintcr/v1"
	"github.com/cloudogu/k8s-ces-control/packages/doguinteraction"
	"github.com/cloudogu/k8s-ces-control/packages/logging"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const responseMessageMissingDoguName = "dogu name is empty"

// NewDoguAdministrationServer returns a new administration server instance to start/stop.. etc. Dogus.
func NewDoguAdministrationServer(client clusterClient, reg cesRegistry, namespace string, loggingService *logging.LoggingService) *server {
	return &server{client: client,
		doguRegistry:   reg.DoguRegistry(),
		doguInterActor: doguinteraction.NewDefaultDoguInterActor(client, namespace, reg),
		ns:             namespace,
		loggingService: *loggingService,
	}
}

type server struct {
	doguRegistry doguRegistry
	pb.UnimplementedDoguAdministrationServer
	client         clusterClient
	doguInterActor doguInterActor
	ns             string
	loggingService logging.LoggingService
}

// StartDogu starts the specified dogu
func (s *server) StartDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguName)
	}

	err := s.doguInterActor.StartDogu(ctx, doguName)
	if err != nil {
		return &types.BasicResponse{}, getGRPCInternalDoguActionError("start", err)
	}

	return &types.BasicResponse{}, err
}

// StopDogu stops the specified dogu
func (s *server) StopDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguName)
	}

	err := s.doguInterActor.StopDogu(ctx, doguName)
	if err != nil {
		return &types.BasicResponse{}, getGRPCInternalDoguActionError("stop", err)
	}

	return &types.BasicResponse{}, err
}

// RestartDogu restarts the specified dogu
func (s *server) RestartDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguName)
	}

	err := s.doguInterActor.RestartDogu(ctx, doguName)
	if err != nil {
		return &types.BasicResponse{}, getGRPCInternalDoguActionError("restart", err)
	}

	return &types.BasicResponse{}, nil
}

func getGRPCInternalDoguActionError(verb string, err error) error {
	return status.Errorf(codes.Internal, fmt.Errorf("failed to %s dogu: %w", verb, err).Error())
}

// GetDoguList returns the list of dogus to administrate (all)
func (s *server) GetDoguList(ctx context.Context, _ *pb.DoguListRequest) (*pb.DoguListResponse, error) {
	list, err := s.client.Dogus(s.ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if len(list.Items) < 1 {
		return &pb.DoguListResponse{}, nil
	}

	doguJsonList, err := s.getDoguJsonList(list.Items)
	if err != nil {
		logrus.Error(fmt.Errorf("failed to get dogus from etcd"))
		return nil, err
	}

	return s.createDoguListResponse(doguJsonList), nil
}

func (s *server) getDoguJsonList(doguListItems []v1.Dogu) (dogus []*core.Dogu, multiErr error) {
	for _, doguListItem := range doguListItems {
		dogu, err := s.doguRegistry.Get(doguListItem.GetName())
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}

		dogus = append(dogus, dogu)
	}

	return dogus, multiErr
}

func (s *server) createDoguListResponse(dogus []*core.Dogu) *pb.DoguListResponse {
	var result []*pb.Dogu

	for _, dogu := range dogus {

		var logLevel, err = s.loggingService.GetLogLevel(dogu.GetSimpleName())
		if err != nil {
			logLevel, err = findDefaultLogLevel(dogu)
			if err != nil {
				logrus.Warn(err)
			}
		}

		result = append(result, &pb.Dogu{
			Name:        dogu.GetSimpleName(),
			DisplayName: dogu.DisplayName,
			Version:     dogu.Version,
			Description: dogu.Description,
			Tags:        dogu.Tags,
			LogLevel:    logLevel,
		})
	}

	return &pb.DoguListResponse{
		Dogus: result,
	}
}

func (s *server) GetBlueprintId(ctx context.Context, _ *pb.DoguBlueprinitIdRequest) (*pb.DoguBlueprintIdResponse, error) {
	bpList, err := s.client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get blueprint list")
	}

	if len(bpList.Items) == 0 {
		return nil, status.Errorf(codes.NotFound, "could not found blueprintID")
	}

	currentBlueprintID := getCurrentBlueprintID(bpList.Items)

	return &pb.DoguBlueprintIdResponse{BlueprintId: currentBlueprintID}, nil
}

func getCurrentBlueprintID(blueprintList []v1bp.Blueprint) string {
	var latestBluePrint = blueprintList[0]

	for _, bp := range blueprintList {
		if bp.CreationTimestamp.Time.After(latestBluePrint.CreationTimestamp.Time) {
			latestBluePrint = bp
		}
	}

	return latestBluePrint.Name
}

func findDefaultLogLevel(dogu *core.Dogu) (string, error) {
	for i := range dogu.Configuration {
		if dogu.Configuration[i].Name == "logging/root" {
			return dogu.Configuration[i].Default, nil
		}
	}
	return "UNKNOWN", fmt.Errorf("could not find default log level for Dogu %s", dogu.Name)
}
