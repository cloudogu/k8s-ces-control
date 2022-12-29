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
	scalingv1 "k8s.io/api/autoscaling/v1"
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

	// TODO create etcd client and read dogu.json for all dogus

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

// RestartDogu restarts the specified dogu
func (s *server) RestartDogu(ctx context.Context, request *pb.DoguAdministrationRequest) (*types.BasicResponse, error) {
	doguName := request.DoguName
	if doguName == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	}

	zeroReplicas := int32(0)
	deployment, err := s.client.AppsV1().Deployments("ecosystem").Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get deployment for dogu %s: %w", doguName, err)
	}
	println(deployment.Spec.Replicas)
	if *deployment.Spec.Replicas == zeroReplicas {
		return &types.BasicResponse{}, s.scaleDeployment(ctx, doguName, 1)
	}

	err = s.scaleDeployment(ctx, doguName, 0)
	if err != nil {
		return nil, err
	}

	deployLabel := fmt.Sprintf("dogu.name=%s", doguName)
	watchInterface, err := s.client.AppsV1().Deployments("ecosystem").
		Watch(ctx, metav1.ListOptions{LabelSelector: deployLabel})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed create watch for deployment wit label %s: %w", deployLabel, err)
	}

	for {
		event := <-watchInterface.ResultChan()
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			return nil, status.Error(codes.Unknown, "watch object is not type of deployment")
		}

		if deployment.Status.Replicas == zeroReplicas {
			return &types.BasicResponse{}, s.scaleDeployment(ctx, doguName, 1)
		} else {
			continue
		}
	}
}

func (s *server) scaleDeployment(ctx context.Context, doguName string, replicas int32) error {
	scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: doguName, Namespace: "ecosystem"}, Spec: scalingv1.ScaleSpec{Replicas: replicas}}
	_, err := s.client.AppsV1().Deployments("ecosystem").UpdateScale(ctx, doguName, scale, metav1.UpdateOptions{})
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to scale deployment to %d: %w", replicas, err)
	}

	return nil
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
