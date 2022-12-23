package doguAdministration

import (
	"context"
	"github.com/cloudogu/cesapp-lib/core"
	pb "github.com/cloudogu/k8s-ces-control/generated/doguAdministration"
	"github.com/cloudogu/k8s-ces-control/generated/types"
)

func NewDoguAdministrationServer() *server {
	return &server{}
}

type server struct {
	pb.UnimplementedDoguAdministrationServer
}

// GetDoguList returns the list of dogus to administrate (all)
func (s *server) GetDoguList(_ context.Context, _ *pb.DoguListRequest) (*pb.DoguListResponse, error) {
	//dogus, err := s.administrator.getDoguList()
	//if err != nil {
	//	log.Error(err)
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//if len(dogus) < 1 {
	//	return &pb.DoguListResponse{}, nil
	//}
	//
	//return createDoguListResponse(dogus), nil

	fakeDoguList := []*pb.Dogu{
		{
			Name:        "test1",
			DisplayName: "Fake Dogu 1",
			Version:     "Fake Version",
			Description: "Fake Dogu Description",
			Tags:        nil,
		},
	}
	return &pb.DoguListResponse{
		Dogus: fakeDoguList,
	}, nil
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

func createDoguListResponse(dogus []*core.Dogu) *pb.DoguListResponse {
	var result []*pb.Dogu
	for _, dogu := range dogus {
		result = append(result, &pb.Dogu{
			Name:        dogu.GetSimpleName(),
			DisplayName: dogu.DisplayName,
			Version:     dogu.Version,
			Description: dogu.Description,
			Tags:        dogu.Tags,
		})
	}

	return &pb.DoguListResponse{
		Dogus: result,
	}
}
