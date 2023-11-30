package logging

import (
	"archive/zip"
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"

	pb "github.com/cloudogu/k8s-ces-control/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/cloudogu/k8s-dogu-operator/api/ecoSystem"
)

const (
	responseMessageMissingDoguname = "Dogu name should not be empty"
)

// TODO: URLs should be made configurable in case parts of the URL must be altered (f. i. "monitoring" may be unavailable
// for organizational reasons)
const lokiGatewareExecutionNamespace = "ecosystem"
const lokiGatewareServiceURL = "http://k8s-loki-gateway." + lokiGatewareExecutionNamespace + ".svc.cluster.local"

type clusterClient interface {
	ecoSystem.EcoSystemV1Alpha1Interface
	kubernetes.Interface
}

type nowClock interface {
	Now() time.Time
}

type realClock struct{}

func (r *realClock) Now() time.Time {
	return time.Now()
}

// NewLoggingService creates a new logging service.
func NewLoggingService(client clusterClient) *loggingService {
	clock := &realClock{}
	return &loggingService{client: client, clock: clock}
}

type loggingService struct {
	pb.UnimplementedDoguLogMessagesServer
	client clusterClient
	clock  nowClock
}

// GetForDogu writes dogu log messages into the stream of the given server.
func (s *loggingService) GetForDogu(request *pb.DoguLogMessageRequest, server pb.DoguLogMessages_GetForDoguServer) error {
	linesCount := int(request.LineCount)
	doguName := request.DoguName
	// delegate to an orderly named method because GetForDogu is misleading but cannot be renamed due to the
	// distributed nature of GRPC definitions
	return writeLogLinesToStream(doguName, linesCount, server)
}

func writeLogLinesToStream(doguName string, linesCount int, server pb.DoguLogMessages_GetForDoguServer) error {
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}
	logrus.Debugf("retrieving %d line(s) of log messages for dogu '%s'", linesCount, doguName)

	logProvider := &lokiLogProvider{
		username: "loki-gateway-user",
		password: "zErGCt9mQVcBbcenFPE3KNYm",
	}

	logLines, err := logProvider.getLogs(doguName, linesCount)
	if err != nil {
		logrus.Errorf("error reading logs: %v", err)
		return createInternalErr(err, codes.InvalidArgument)
	}

	compressedMessagesBytes, err := compressMessages(doguName, logLines)
	if err != nil {
		return err
	}

	err = stream.WriteToStream(compressedMessagesBytes, server)
	if err != nil {
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

func compressMessages(doguName string, logLines []logLine) ([]byte, error) {
	if len(logLines) <= 0 {
		return []byte{}, nil
	}

	compressedMessages := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(compressedMessages)
	logFileMetadata := zip.FileHeader{
		Name:     fmt.Sprintf("%s.log", doguName),
		Modified: time.Now(),
		Method:   zip.Deflate,
	}

	writer, err := zipWriter.CreateHeader(&logFileMetadata)
	if err != nil {
		return nil, err
	}

	totalWrittenBytes := 0
	for _, line := range logLines {
		writtenBytes, err := writer.Write([]byte(line.value + "\n"))
		if err != nil {
			return nil, err
		}
		totalWrittenBytes += writtenBytes
	}

	_ = zipWriter.Close()
	logrus.Debugf("wrote %d byte(s) to archive", totalWrittenBytes)
	return compressedMessages.Bytes(), nil
}

func createInternalErr(err error, code codes.Code) error {
	logrus.Error(err)
	return status.Error(code, err.Error())
}
