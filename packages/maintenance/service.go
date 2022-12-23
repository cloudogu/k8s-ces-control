package maintenance

import (
	"context"
	pbMaintenance "github.com/cloudogu/k8s-ces-control/generated/maintenance"
	"github.com/cloudogu/k8s-ces-control/generated/types"
	"github.com/sirupsen/logrus"
)

const (
	debugModeDisableAtKey = "debug/disable_at"
)

type debugModeService struct {
	pbMaintenance.UnimplementedDebugModeServer
}

func NewDebugModeService() *debugModeService {
	return &debugModeService{}
}

func (s debugModeService) Enable(_ context.Context, request *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("(fake) Enable maintenance mode...")
	return &types.BasicResponse{}, nil
}

func (s debugModeService) Disable(_ context.Context, request *pbMaintenance.ToggleDebugModeRequest) (*types.BasicResponse, error) {
	logrus.Info("(fake) Disable maintenance mode...")
	return &types.BasicResponse{}, nil
}

func (s debugModeService) Status(context.Context, *types.BasicRequest) (*pbMaintenance.DebugModeStatusResponse, error) {
	logrus.Debugf("(fake) Get status of mainentance mode...")
	return &pbMaintenance.DebugModeStatusResponse{IsEnabled: false, DisableAtTimestamp: 0}, nil
}
