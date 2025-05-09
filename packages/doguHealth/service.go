package doguHealth

import (
	"context"
	pbHealth "github.com/cloudogu/ces-control-api/generated/health"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const checkTypeContainer = "container"
const responseMessageMissingDoguname = "dogu name is empty"

// NewDoguHealthService return a new health server to retrieve health information from Dogus.
func NewDoguHealthService(client doguClient) *server {
	return &server{client: client}
}

type server struct {
	pbHealth.UnimplementedDoguHealthServer
	client doguClient
}

// GetByName retrieves the health information about a given dogu if it is installed.
func (s *server) GetByName(ctx context.Context, request *pbHealth.DoguHealthRequest) (*pbHealth.DoguHealthResponse, error) {
	logrus.Debugf("Check healthy state of dogu [%s]", request.DoguName)
	if request.GetDoguName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, responseMessageMissingDoguname)
	}

	return s.getDoguHealthResponse(ctx, request.DoguName)
}

// GetByNames retrieves the health information about the given dogus if they are installed.
func (s *server) GetByNames(ctx context.Context, request *pbHealth.DoguHealthListRequest) (*pbHealth.DoguHealthMapResponse, error) {
	logrus.Debugf("Check healthy state of dogus [%s]", request.Dogus)
	return s.getDoguListHealthResponse(ctx, request.Dogus)
}

// GetAll retrieves health information about all installed dogus.
func (s *server) GetAll(ctx context.Context, _ *pbHealth.DoguHealthAllRequest) (*pbHealth.DoguHealthMapResponse, error) {
	logrus.Debugf("Check healthy state of all dogus")
	doguList, err := s.client.List(ctx, metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var dogus []string
	for _, dogu := range doguList.Items {
		dogus = append(dogus, dogu.Name)
	}

	return s.getDoguListHealthResponse(ctx, dogus)
}

func (s *server) getDoguListHealthResponse(ctx context.Context, doguNameList []string) (*pbHealth.DoguHealthMapResponse, error) {
	doguHealthList := map[string]*pbHealth.DoguHealthResponse{}

	var multiErr error
	allHealthy := true
	for _, dogu := range doguNameList {
		doguHealthResponse, err := s.getDoguHealthResponse(ctx, dogu)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)
		}

		if !doguHealthResponse.Healthy {
			allHealthy = false
		}

		doguHealthList[dogu] = doguHealthResponse
	}

	return &pbHealth.DoguHealthMapResponse{
		AllHealthy: allHealthy,
		Results:    doguHealthList,
	}, multiErr
}

func (s *server) getDoguHealthResponse(ctx context.Context, doguName string) (*pbHealth.DoguHealthResponse, error) {
	dogu, err := s.client.Get(ctx, doguName, metav1.GetOptions{})
	if err != nil {
		errResponse := &pbHealth.DoguHealthResponse{
			FullName:    doguName,
			ShortName:   doguName,
			DisplayName: doguName,
			Healthy:     false,
			Results:     []*pbHealth.DoguHealthCheck{},
		}
		return errResponse, err
	}

	isHealthy := dogu.Status.Health == v2.AvailableHealthStatus
	response := &pbHealth.DoguHealthResponse{
		FullName:    doguName,
		ShortName:   doguName,
		DisplayName: doguName,
		Healthy:     isHealthy,
		Results:     []*pbHealth.DoguHealthCheck{},
	}

	// When dogu.spec.stopped=false it indicates that the dogu is not stopped.
	// The admin dogu notices that there is a Dogu which is not healthy but not  stopped and assumes: Dogu is starting
	containerStatusResult := &pbHealth.DoguHealthCheck{
		Type:     checkTypeContainer,
		Success:  !dogu.Spec.Stopped,
		Message:  "Check whether a dogu is not stopped (ready or not).",
		Critical: false,
	}

	response.Results = append(response.Results, containerStatusResult)

	return response, nil
}
