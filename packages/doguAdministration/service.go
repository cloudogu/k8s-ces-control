package doguAdministration

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
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

	cesregistry "github.com/cloudogu/cesapp-lib/registry"
	"github.com/hashicorp/go-multierror"
)

const responseMessageMissingDoguname = "dogu name is empty"

func NewDoguAdministrationServer(client config.ClusterClient, reg cesregistry.Registry) *server {
	return &server{client: client, doguRegistry: reg.DoguRegistry()}
}

type server struct {
	doguRegistry cesregistry.DoguRegistry
	pb.UnimplementedDoguAdministrationServer
	client config.ClusterClient
}

// GetDoguList returns the list of dogus to administrate (all)
func (s *server) GetDoguList(ctx context.Context, _ *pb.DoguListRequest) (*pb.DoguListResponse, error) {
	list, err := s.client.EcoSystemApi.Dogus(config.CurrentNamespace).List(ctx, metav1.ListOptions{})
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

	return createDoguListResponse(doguJsonList), nil
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
	deployment, err := s.client.AppsV1().Deployments(config.CurrentNamespace).Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get deployment for dogu %s: %s", doguName, err.Error())
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
	watchInterface, err := s.client.AppsV1().Deployments(config.CurrentNamespace).
		Watch(ctx, metav1.ListOptions{LabelSelector: deployLabel})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed create watch for deployment wit label %s: %s", deployLabel, err.Error())
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
	scale := &scalingv1.Scale{ObjectMeta: metav1.ObjectMeta{Name: doguName, Namespace: config.CurrentNamespace}, Spec: scalingv1.ScaleSpec{Replicas: replicas}}
	_, err := s.client.AppsV1().Deployments(config.CurrentNamespace).UpdateScale(ctx, doguName, scale, metav1.UpdateOptions{})
	if err != nil {
		return status.Errorf(codes.Unknown, "failed to scale deployment to %d: %s", replicas, err.Error())
	}

	return nil
}

func (s *server) getDoguJsonList(doguListItems []v1.Dogu) (dogus []*core.Dogu, multiErr error) {
	for _, doguListItem := range doguListItems {
		dogu, err := s.doguRegistry.Get(doguListItem.GetName())
		if err != nil {
			multiErr = multierror.Append(err, err)
		}

		dogus = append(dogus, dogu)
	}

	return dogus, multiErr
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
