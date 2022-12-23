package doguHealth

import (
	"context"
	pbHealth "github.com/cloudogu/k8s-ces-control/generated/health"
	"github.com/sirupsen/logrus"
)

func NewDoguHealthService() pbHealth.DoguHealthServer {
	//doguHealthMonitor, err := newDoguHealthMonitor()
	//if err != nil {
	//	return nil, fmt.Errorf("failed to create new dogu health monitor: %w", err)
	//}
	//return &server{doguHealthMonitor: &doguHealthMonitor}, nil
	return &server{}
}

type server struct {
	pbHealth.UnimplementedDoguHealthServer
	//doguHealthMonitor *doguHealthMonitor
}

// GetByName retrieves the health information about a given dogu if it is installed.
func (s *server) GetByName(_ context.Context, request *pbHealth.DoguHealthRequest) (*pbHealth.DoguHealthResponse, error) {
	//healthCheckResults, err := s.doguHealthMonitor.CheckOne(request.DoguName)
	//if err != nil {
	//	return nil, createInternalErr(err)
	//}
	//healthData := processHealthCheckResults(*healthCheckResults)
	//response := &pbHealth.DoguHealthResponse{
	//	FullName:    healthCheckResults.fullName,
	//	ShortName:   healthCheckResults.shortName,
	//	DisplayName: healthCheckResults.displayName,
	//	Healthy:     healthData.Healthy,
	//	Results:     healthData.Results,
	//}
	logrus.Debugf("Check healthy state of dogu [%s]", request.DoguName)
	response := &pbHealth.DoguHealthResponse{
		FullName:    "Fake 1",
		ShortName:   "Fake 1",
		DisplayName: "Fake 1",
		Healthy:     true,
		Results: []*pbHealth.DoguHealthCheck{
			{
				Type:     "Test",
				Success:  true,
				Message:  "WOOOW",
				Critical: false,
			},
		},
	}
	return response, nil
}

func (s *server) GetByNames(_ context.Context, request *pbHealth.DoguHealthListRequest) (*pbHealth.DoguHealthMapResponse, error) {
	//healthCheckResults, err := s.doguHealthMonitor.CheckMany(request.Dogus)
	//if err != nil {
	//	return nil, createInternalErr(err)
	//}
	//doguHealthList, allDogusAreHealthy := processCheckResults(healthCheckResults)

	logrus.Debugf("Check healthy state of dogus [%s]", request.Dogus)
	doguHealthList := map[string]*pbHealth.DoguHealthResponse{
		"test1": {
			FullName:    "test1",
			ShortName:   "t1",
			DisplayName: "Fake Dogu 1",
			Healthy:     true,
			Results:     nil,
		},
	}
	return &pbHealth.DoguHealthMapResponse{
		AllHealthy: true,
		Results:    doguHealthList,
	}, nil
}

// GetAll retrieves health information about all installed dogus.
func (s *server) GetAll(_ context.Context, _ *pbHealth.DoguHealthAllRequest) (*pbHealth.DoguHealthMapResponse, error) {
	logrus.Debugf("Check healthy state of all dogus")
	doguHealthList := map[string]*pbHealth.DoguHealthResponse{
		"": {
			FullName:    "test1",
			ShortName:   "t1",
			DisplayName: "Fake Dogu 1",
			Healthy:     true,
			Results:     nil,
		},
	}
	return &pbHealth.DoguHealthMapResponse{
		AllHealthy: true,
		Results:    doguHealthList,
	}, nil
}
