package logging

import (
	pb "github.com/cloudogu/k8s-ces-control/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	responseMessageMissingDoguname = "Dogu name should not be empty"
)

func NewLoggingService() pb.DoguLogMessagesServer {
	return &loggingService{}
}

type loggingService struct {
	pb.UnimplementedDoguLogMessagesServer
}

func (s *loggingService) GetForDogu(request *pb.DoguLogMessageRequest, server pb.DoguLogMessages_GetForDoguServer) error {
	linesCount := int(request.LineCount)
	doguName := request.DoguName
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}
	logrus.Debugf("retrieving %d line(s) of log messages for dogu '%s'", linesCount, doguName)

	compressedMessages := "this is a nice test"
	err := stream.WriteToStream([]byte(compressedMessages), server)
	if err != nil {
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

func createInternalErr(err error, code codes.Code) error {
	logrus.Error(err)
	return status.Error(code, err.Error())
}
