package maintenance

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbMaintenance "github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/generated/types"
)

type debugModeService struct {
	pbMaintenance.UnimplementedDebugModeServer
}

// NewDebugModeService returns an instance of debugModeService.
func NewDebugModeService() *debugModeService {
	return &debugModeService{}
}

// Enable returns an error because the method is unimplemented.
func (s debugModeService) Enable(_ context.Context, _ *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("(fake) Enable maintenance mode...")
	return &types.BasicResponse{}, status.Errorf(codes.Unimplemented, "this service is not yet implemented")
}

// Disable returns an error because the method is unimplemented.
func (s debugModeService) Disable(_ context.Context, _ *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("(fake) Disable maintenance mode...")
	return &types.BasicResponse{}, status.Errorf(codes.Unimplemented, "this service is not yet implemented")
}

// Status return an error because the method is unimplemented.
func (s debugModeService) Status(context.Context, *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	logrus.Debugf("(fake) Get status of mainentance mode...")
	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: false, DisableAtTimestamp: 0}, status.Errorf(codes.Unimplemented, "this service is not yet implemented")
}
