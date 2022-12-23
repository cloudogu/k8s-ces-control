package logging

import (
	"archive/zip"
	"bytes"
	"fmt"
	pb "github.com/cloudogu/k8s-ces-control/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
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

	logFileData, err := s.readLogs(doguName, linesCount)
	if err != nil {
		return createInternalErr(err, codes.InvalidArgument)
	}

	compressedMessagesBytes, err := compressMessages(request.DoguName, logFileData)
	if err != nil {
		return err
	}

	err = stream.WriteToStream(compressedMessagesBytes, server)
	if err != nil {
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

func (s *loggingService) readLogs(name string, count int) ([]byte, error) {
	logFileData := name + `
 OMG Im starting...
calculating...
damn...	
`
	if count > 0 {
		logFileData += "\nonly some lines"
		return []byte(logFileData), nil
	}
	return []byte(logFileData), nil
}

func compressMessages(doguName string, logLines []byte) ([]byte, error) {
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
	writtenBytes, err := writer.Write(logLines)
	if err != nil {
		return nil, err
	}
	_ = zipWriter.Close()
	logrus.Debugf("wrote %d byte(s) to archive", writtenBytes)
	return compressedMessages.Bytes(), nil
}

func createInternalErr(err error, code codes.Code) error {
	logrus.Error(err)
	return status.Error(code, err.Error())
}
