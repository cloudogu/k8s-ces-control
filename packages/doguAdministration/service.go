package doguAdministration

import (
	"context"
	"fmt"
	pb "github.com/cloudogu/ces-control-api/generated/doguAdministration"
	"github.com/cloudogu/ces-control-api/generated/types"
	"github.com/cloudogu/cesapp-lib/core"
	v1bp "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-ces-control/packages/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const responseMessageMissingDoguName = "dogu name is empty"

type logService interface {
	GetLogLevel(context.Context, string) (logging.LogLevel, error)
}

// NewDoguAdministrationServer returns a new administration server instance to start/stop.. etc. Dogus.
func NewDoguAdministrationServer(blueprintLister BlueprintLister, doguDescriptorGetter doguDescriptorGetter, doguInterActor doguInterActor, logService logService) *server {
	return &server{
		blueprintLister:      blueprintLister,
		doguDescriptorGetter: doguDescriptorGetter,
		doguInterActor:       doguInterActor,
		loggingService:       logService,
	}
}

type server struct {
	doguDescriptorGetter doguDescriptorGetter
	pb.UnimplementedDoguAdministrationServer
	blueprintLister BlueprintLister
	doguInterActor  doguInterActor
	loggingService  logService
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
	return status.Errorf(codes.Internal, "%v", fmt.Errorf("failed to %s dogu: %w", verb, err).Error())
}

// GetDoguList returns the list of dogus to administrate (all)
func (s *server) GetDoguList(ctx context.Context, _ *pb.DoguListRequest) (*pb.DoguListResponse, error) {
	doguJsonList, err := s.doguDescriptorGetter.GetCurrentOfAll(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get dogu registry: %w", err)
		logrus.Error(err)
		return nil, err
	}

	return s.createDoguListResponse(ctx, doguJsonList), nil
}

func (s *server) createDoguListResponse(ctx context.Context, dogus []*core.Dogu) *pb.DoguListResponse {
	var result []*pb.Dogu

	for _, dogu := range dogus {

		logLevel, err := s.loggingService.GetLogLevel(ctx, dogu.GetSimpleName())
		if err != nil {
			logrus.Warn(err)
		}

		result = append(result, &pb.Dogu{
			Name:        dogu.GetSimpleName(),
			DisplayName: dogu.DisplayName,
			Version:     dogu.Version,
			Description: dogu.Description,
			Tags:        dogu.Tags,
			LogLevel:    logLevel.String(),
		})
	}

	return &pb.DoguListResponse{
		Dogus: result,
	}
}

func (s *server) GetBlueprintId(ctx context.Context, _ *pb.DoguBlueprinitIdRequest) (*pb.DoguBlueprintIdResponse, error) {
	bpList, err := s.blueprintLister.List(ctx, metav1.ListOptions{})
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
