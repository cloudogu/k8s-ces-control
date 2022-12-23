package doguAdministration

import (
	"context"
	"fmt"
	pb "github.com/cloudogu/k8s-ces-control/generated/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/generated/types"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
func (s *server) StartDogu(_ context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	//doguName := request.DoguName
	//if doguName == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	//}
	//message, err := s.administrator.startDogu(doguName)
	//log.Info(message)
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}

	return &types.BasicResponse{}, nil
}

// StopDogu stops the specified dogu
func (s *server) StopDogu(_ context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	//doguName := request.DoguName
	//if doguName == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	//}
	//message, err := s.administrator.stopDogu(doguName)
	//log.Info(message)
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	return &types.BasicResponse{}, nil
}

// RestartDogu restarts the specified dogu
func (s *server) RestartDogu(_ context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	//doguName := request.DoguName
	//if doguName == "" {
	//	return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	//}
	//messageStopDogu, err := s.administrator.stopDogu(doguName)
	//log.Info(messageStopDogu)
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//messsageStartDogu, err := s.administrator.startDogu(doguName)
	//log.Info(messsageStartDogu)
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
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
