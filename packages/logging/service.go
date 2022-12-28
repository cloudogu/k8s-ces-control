package logging

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	pb "github.com/cloudogu/k8s-ces-control/generated/logging"
	"github.com/cloudogu/k8s-ces-control/packages/config"
	"github.com/cloudogu/k8s-ces-control/packages/stream"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/url"
	"time"
)

const (
	responseMessageMissingDoguname = "Dogu name should not be empty"
)

func NewLoggingService() (*loggingService, error) {
	client, err := config.CreateClusterClient()
	if err != nil {
		return nil, err
	}
	return &loggingService{client: client}, nil
}

type loggingService struct {
	client config.ClusterClient
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

func buildLokiQueryUrl(name string, count int) (string, error) {
	baseUrl, err := url.Parse("http://loki-gateway.monitoring.svc.cluster.local:80")
	if err != nil {
		return "", err
	}

	baseUrl.Path += "loki/api/v1/query"
	params := url.Values{}
	queryParam := fmt.Sprintf("{pod=~\"%s.*\"}", name)
	params.Add("query", queryParam)
	params.Add("direction", "forward")
	baseUrl.RawQuery = params.Encode()
	if count != 0 {
		baseUrl.RawQuery += fmt.Sprintf("&limit=%d", count)
	}

	return baseUrl.String(), nil
}

func (s *loggingService) doLokiHttpQuery(lokiUrl string) (*http.Response, error) {
	secret, err := s.client.CoreV1().Secrets("monitoring").Get(context.Background(), "loki-credentials",
		v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch loki secret: %w", err)
	}

	req, err := http.NewRequest("GET", lokiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with url [%s]: %w", lokiUrl, err)
	}
	req.SetBasicAuth(string(secret.Data["username"]), string(secret.Data["password"]))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request with url [%s]: %w", lokiUrl, err)
	}

	return resp, nil
}

// TODO return all logs?
func (s *loggingService) readLogs(name string, count int) ([]byte, error) {
	lokiUrl, err := buildLokiQueryUrl(name, count)
	if err != nil {
		return nil, err
	}

	resp, err := s.doLokiHttpQuery(lokiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("loki http error: status: %s, code: %d", resp.Status, resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var reqBytes []byte
	for scanner.Scan() {
		reqBytes = append(reqBytes, scanner.Bytes()...)
	}

	return reqBytes, nil
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
