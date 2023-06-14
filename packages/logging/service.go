package logging

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
const lokiGatewareExecutionNamespace = "monitoring"
const lokiGatewareServiceURL = "http://loki-gateway." + lokiGatewareExecutionNamespace + ".svc.cluster.local:80"

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
	return writeLogLinesToStream(s.client, doguName, linesCount, s.clock, server)
}

func writeLogLinesToStream(client clusterClient, doguName string, linesCount int, clock nowClock, server pb.DoguLogMessages_GetForDoguServer) error {
	if doguName == "" {
		return status.Error(codes.InvalidArgument, responseMessageMissingDoguname)
	}
	logrus.Debugf("retrieving %d line(s) of log messages for dogu '%s'", linesCount, doguName)

	logFileData, err := readLogs(client, clock, doguName, linesCount)
	if err != nil {
		return createInternalErr(err, codes.InvalidArgument)
	}

	compressedMessagesBytes, err := compressMessages(doguName, logFileData)
	if err != nil {
		return err
	}

	err = stream.WriteToStream(compressedMessagesBytes, server)
	if err != nil {
		return createInternalErr(err, codes.Internal)
	}

	return nil
}

// buildLokiQueryUrl returns a Loki query over a range of time using the given pod regexp part and the maximum number of
// results being returned.
func buildLokiQueryUrl(podRegexp string, limit int, clock nowClock) (string, error) {
	baseUrl, err := url.Parse(lokiGatewareServiceURL)
	if err != nil {
		return "", err
	}

	baseUrl = baseUrl.JoinPath("/loki/api/v1/query_range")
	params := url.Values{}
	queryParam := fmt.Sprintf("{pod=~\"%s.*\"}", podRegexp)
	params.Add("query", queryParam)
	params.Add("direction", "backward")
	baseUrl.RawQuery = params.Encode()
	if limit != 0 {
		baseUrl.RawQuery += fmt.Sprintf("&limit=%d", limit)
	}
	startDate := clock.Now().Add(-(time.Hour * 24 * 7))
	baseUrl.RawQuery += fmt.Sprintf("&start=%d", startDate.UnixNano())

	return baseUrl.String(), nil
}

func doLokiHttpQuery(client clusterClient, lokiUrl string) (*http.Response, error) {
	secret, err := client.CoreV1().Secrets(lokiGatewareExecutionNamespace).Get(context.Background(), "loki-credentials",
		v1.GetOptions{})
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("failed to fetch loki secret: %w", err), codes.Canceled)
	}

	req, err := http.NewRequest(http.MethodGet, lokiUrl, nil)
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("failed to create request with url [%s]: %w", lokiUrl, err), codes.Canceled)
	}
	req.SetBasicAuth(string(secret.Data["username"]), string(secret.Data["password"]))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("failed to execute request with url [%s]: %w", lokiUrl, err), codes.Canceled)
	}

	return resp, nil
}

func readLogs(client clusterClient, clock nowClock, name string, count int) ([]byte, error) {
	lokiUrl, err := buildLokiQueryUrl(name, count, clock)
	if err != nil {
		return nil, err
	}

	resp, err := doLokiHttpQuery(client, lokiUrl)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	isHttpStatusErrorish := resp.StatusCode >= http.StatusMultipleChoices
	if isHttpStatusErrorish {
		return nil, createInternalErr(fmt.Errorf("loki http error: status: %s, code: %d", resp.Status, resp.StatusCode), codes.Canceled)
	}

	scanner := bufio.NewScanner(resp.Body)
	var reqBytes []byte
	for scanner.Scan() {
		reqBytes = append(reqBytes, scanner.Bytes()...)
	}

	lokiResp := &LokiResponse{}
	err = json.Unmarshal(reqBytes, lokiResp)
	if err != nil {
		return nil, createInternalErr(fmt.Errorf("failed to unmarshal response: %w", err), codes.Canceled)
	}

	if lokiResp.Status != "success" {
		return nil, createInternalErr(fmt.Errorf("loki response status is not successfull"), codes.Canceled)
	}

	if lokiResp.Data.ResultType != "streams" {
		return nil, createInternalErr(fmt.Errorf("loki response data aren't streams"), codes.Canceled)
	}

	if len(lokiResp.Data.Result) == 0 {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}
	for _, s := range extractRawLogsFromLokiResponseData(lokiResp.Data) {
		buf.WriteString(fmt.Sprintf("%s\n", s))
	}

	return buf.Bytes(), nil
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

func extractRawLogsFromLokiResponseData(lokiResponseData LokiResponseData) []string {
	var unsortedLogs = make(map[string]string)
	streams := lokiResponseData.Result
	for _, lokiStream := range streams {
		for _, value := range lokiStream.Values {
			unsortedLogs[value[0]] = value[1]
		}
	}

	keys := make([]string, 0, len(unsortedLogs))
	for k := range unsortedLogs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sortedLogs []string
	for _, k := range keys {
		sortedLogs = append(sortedLogs, unsortedLogs[k])
	}

	return sortedLogs
}

func createInternalErr(err error, code codes.Code) error {
	logrus.Error(err)
	return status.Error(code, err.Error())
}

// LokiResponse represents the root structure of a query response.
type LokiResponse struct {
	Status string           `json:"status"`
	Data   LokiResponseData `json:"data"`
}

// LokiResponseData contains log stream results and metadata ResultType. ResultType could be "stream" oder "vector".
type LokiResponseData struct {
	ResultType string             `json:"resultType"`
	Result     []LokiStreamResult `json:"result"`
}

// LokiStreamResult the stream and the log values.
type LokiStreamResult struct {
	Stream LokiStream `json:"stream"`
	Values [][]string `json:"values"`
}

// LokiStream contains metadata for the stream result.
type LokiStream struct {
	Container string `json:"container"`
	Filename  string `json:"filename"`
	Job       string `json:"job"`
	Namespace string `json:"namespace"`
	NodeName  string `json:"node_name"`
	Pod       string `json:"pod"`
	Stream    string `json:"stream"`
	App       string `json:"app"`
}
