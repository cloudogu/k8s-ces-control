package doguAdministration

import (
	"context"
	"fmt"
	pb "github.com/cloudogu/k8s-ces-control/generated/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/generated/types"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const responseMessageMissingDoguname = "dogu name is empty"

func NewDoguAdministrationServer() (*server, error) {
	client, err := config.CreateClusterClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster client")
	}

	return &server{client: client}, nil
}

type server struct {
	pb.UnimplementedDoguAdministrationServer
	client config.ClusterClient
}

// GetDoguList returns the list of dogus to administrate (all)
func (s *server) GetDoguList(ctx context.Context, _ *pb.DoguListRequest) (*pb.DoguListResponse, error) {
	list, err := s.client.EcoSystemApi.Dogus("ecosystem").List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if len(list.Items) < 1 {
		return &pb.DoguListResponse{}, nil
	}

	return createDoguListResponse(list), nil
}

// StartDogu starts the specified dogu
func (s *server) StartDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	}

	err := s.scaleDeployment(ctx, doguName, 1)

	return &types.BasicResponse{}, err
}

// StopDogu stops the specified dogu
func (s *server) StopDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	}

	err := s.scaleDeployment(ctx, doguName, 0)

	return &types.BasicResponse{}, err
}

func (s *server) scaleDeployment(ctx context.Context, doguName string, replicas int) error {
	deployment, err := s.getDeployment(ctx, doguName)
	if err != nil {
		return err
	}

	replicas32 := int32(replicas)
	deployment.Spec.Replicas = &replicas32

	err = s.updateDeployment(ctx, deployment, doguName)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) getDeployment(ctx context.Context, doguName string) (*appsv1.Deployment, error) {
	deployment, err := s.client.AppsV1().Deployments("ecosystem").Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get deployment for dogu %s: %w", doguName, err)
	}

	return deployment, nil
}

func (s *server) updateDeployment(ctx context.Context, deployment *appsv1.Deployment, doguName string) error {
	_, err := s.client.AppsV1().Deployments("ecosystem").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to update deployment for dogu %s: %w", doguName, err)
	}

	return nil
}

// RestartDogu restarts the specified dogu
func (s *server) RestartDogu(_ context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	// doguName := request.DoguName
	// if doguName == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	// }
	// messageStopDogu, err := s.administrator.stopDogu(doguName)
	// log.Info(messageStopDogu)
	// if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	// }
	// messsageStartDogu, err := s.administrator.startDogu(doguName)
	// log.Info(messsageStartDogu)
	// if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	// }
	return &types.BasicResponse{}, nil
}

func createDoguListResponse(dogus *v1.DoguList) *pb.DoguListResponse {
	var result []*pb.Dogu
	for _, dogu := range dogus.Items {
		result = append(result, &pb.Dogu{
			Name:        dogu.Name,
			DisplayName: dogu.Spec.Name,
			Version:     dogu.Spec.Version,
			Description: dogu.Spec.Name,
			Tags:        []string{dogu.Spec.Name},
		})
	}

	return &pb.DoguListResponse{
		Dogus: result,
	}
}
