package logging

import (
	"cmp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

const defaultQueryLimit = 1000
const maxQueryLimit = 5000 // the max query limit for Loki

type nowClock interface {
	Now() time.Time
}

type realClock struct{}

func (r *realClock) Now() time.Time {
	return time.Now()
}

type LokiLogProvider struct {
	gatewayUrl string
	username   string
	password   string
	clock      nowClock
	httpClient *http.Client
}

func NewLokiLogProvider(gatewayUrl string, username string, password string) *LokiLogProvider {
	return &LokiLogProvider{
		gatewayUrl: gatewayUrl,
		username:   username,
		password:   password,
		clock:      &realClock{},
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// getLogs queries logs from loki until the given lineCount is reached, or no more logs are available.
// Since loki requires a max time-window of 30 days and has a limit of max 5000 log-lines per request,
// a query can result in multiple request for the loki backend.
// The time-window is adjusted accordingly to get all logs even if they are older than 30 days.
func (llp *LokiLogProvider) getLogs(doguName string, linesCount int) ([]logLine, error) {
	endDate := llp.clock.Now()
	startDate := createQueryStartDateFromEndDate(endDate)

	result := make([]logLine, 0)
	for {
		limit := calculateQueryLimit(linesCount, len(result))

		logLines, err := llp.queryLogsFromLoki(doguName, startDate, endDate, "", limit)
		if err != nil {
			return nil, fmt.Errorf("failed to query logs from loki: %w", err)
		}

		if len(logLines) <= 0 {
			// no new lines to read => nothing is left
			break
		}

		// prepend logs
		result = append(logLines, result...)

		if len(logLines) < limit {
			// the query returned fewer lines than requested => nothing is left
			break
		}

		if linesCount > 0 && len(result) >= linesCount {
			// we reached the maximum lines count
			break
		}

		// there are still logs to read -> start with the newest log timestamp from the last response
		endDate = logLines[0].timestamp
		startDate = createQueryStartDateFromEndDate(endDate)
	}

	// because multiple logs can happen at the exact same timestamp and the query is batched over time,
	// it is possible that consecutive batches contain the same logLines. These need to be removed.
	result = deduplicateLogLines(result)

	// truncate log-lines if we got more than requested
	if linesCount > 0 && len(result) > linesCount {
		result = result[:linesCount]
	}

	logrus.Debugf("finished loading logs for linecount; got %d logLines", len(result))

	return result, nil
}

// queryLogs queries logs from loki for the given time-window (startDate and endDate).
// The allowed time-window is max 30 days (limited by loki).
// Since loki also has a limit of max 5000 log-lines per request, a query can result in multiple request for the loki backend.
func (llp *LokiLogProvider) queryLogs(doguName string, startDate time.Time, endDate time.Time, filter string) ([]logLine, error) {
	if endDate.IsZero() {
		endDate = llp.clock.Now()
	}

	if startDate.IsZero() {
		startDate = createQueryStartDateFromEndDate(endDate)
	}

	result := make([]logLine, 0)
	limit := defaultQueryLimit
	for {

		logLines, err := llp.queryLogsFromLoki(doguName, startDate, endDate, filter, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to query logs from loki: %w", err)
		}

		if len(logLines) <= 0 {
			// no new lines to read => nothing is left
			break
		}

		// prepend logs
		result = append(logLines, result...)

		if len(logLines) < limit {
			// the query returned fewer lines than requested => nothing is left
			break
		}

		// there are still logs to read -> start with the newest log timestamp from the last response
		endDate = logLines[0].timestamp
	}

	// because multiple logs can happen at the exact same timestamp and the query is batched over time,
	// it is possible that consecutive batches contain the same logLines. These need to be removed.
	result = deduplicateLogLines(result)

	logrus.Debugf("finished querying logs; got %d logLines", len(result))

	return result, nil
}

func (llp *LokiLogProvider) queryLogsFromLoki(doguName string, startDate time.Time, endDate time.Time, filter string, limit int) ([]logLine, error) {
	if limit <= 0 {
		limit = defaultQueryLimit
	}

	if limit > maxQueryLimit {
		return nil, fmt.Errorf("the given limit of %d exceeds the maximum limit of %d", limit, maxQueryLimit)
	}

	query := fmt.Sprintf("{pod=~\"%s.*\"}", doguName)

	if len(strings.TrimSpace(filter)) != 0 {
		query = fmt.Sprintf("%s |= \"%s\" ", query, filter)
	}

	logrus.Debugf("running loki query for '%s' from %s to %s with limit %d", doguName, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), limit)
	lokiQueryUrl, err := buildLokiQueryUrl(llp.gatewayUrl, query, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to build loki-query: %w", err)
	}

	lokiResp, err := llp.doLokiHttpQuery(lokiQueryUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to execute loki-query: %w", err)
	}

	logLines, err := extractLogLinesFromLokiResponse(lokiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to extract logs from loki response: %w", err)
	}

	return logLines, nil
}

func calculateQueryLimit(linesCount int, resultCount int) int {
	if linesCount <= 0 {
		return defaultQueryLimit
	}

	remainingCount := linesCount - resultCount

	if remainingCount < 0 {
		return 0
	}

	if remainingCount < defaultQueryLimit {
		return remainingCount
	}
	return defaultQueryLimit
}

// createQueryStartDateFromEndDate calculates the start date for a loki query based on the given end date.
// Since loki query run backwards (in time) the start date is calculated based on the end date.
// The start date is set 30 days before the given end date, because that is the maximum range that loki allows.
func createQueryStartDateFromEndDate(endDate time.Time) time.Time {
	return endDate.Add(-24 * 30 * time.Hour)
}

// buildLokiQueryUrl returns a Loki query over a range of time using the given query regexp part, a start date, an end date and the maximum number of
// results being returned.
func buildLokiQueryUrl(lokiBaseUrl string, query string, startDate time.Time, endDate time.Time, limit int) (string, error) {
	baseUrl, err := url.Parse(lokiBaseUrl)
	if err != nil {
		return "", err
	}

	baseUrl = baseUrl.JoinPath("/loki/api/v1/query_range")

	params := baseUrl.Query()
	params.Set("query", query)
	params.Set("direction", "backward")
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("start", fmt.Sprintf("%d", startDate.UnixNano()))
	params.Set("end", fmt.Sprintf("%d", endDate.UnixNano()))

	baseUrl.RawQuery = params.Encode()

	return baseUrl.String(), nil
}

func (llp *LokiLogProvider) doLokiHttpQuery(lokiUrl string) (*lokiResponse, error) {
	logrus.Debugf("running loki query with URL: %s", lokiUrl)
	req, err := http.NewRequest(http.MethodGet, lokiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request with url [%s]: %w", lokiUrl, err)
	}
	req.SetBasicAuth(llp.username, llp.password)
	resp, err := llp.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request with url [%s]: %w", lokiUrl, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Errorf("error while closing response body of loki query: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil || len(responseData) == 0 {
			responseData = []byte(fmt.Sprintf("faild to read error response: %v", err))
		}

		return nil, fmt.Errorf("loki http error: status: %s, code: %d; response-body: %s", resp.Status, resp.StatusCode, responseData)
	}

	return parseLokiResponse(resp.Body)
}

func parseLokiResponse(lokiResult io.Reader) (*lokiResponse, error) {
	lokiResp := &lokiResponse{}
	err := json.NewDecoder(lokiResult).Decode(lokiResp)
	if err != nil {
		return lokiResp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if lokiResp.Status != "success" {
		return lokiResp, fmt.Errorf("loki response status is not successfull; status is %s", lokiResp.Status)
	}

	if lokiResp.Data.ResultType != "streams" {
		return lokiResp, fmt.Errorf("loki response data aren't streams; resultType is %s", lokiResp.Data.ResultType)
	}

	return lokiResp, nil
}

func extractLogLinesFromLokiResponse(lokiResponse *lokiResponse) ([]logLine, error) {
	var logLines = make([]logLine, 0)
	for _, lokiStream := range lokiResponse.Data.Result {
		for _, value := range lokiStream.Values {
			nanos, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse log timestamp: %w", err)
			}
			logLines = append(logLines, logLine{
				timestamp: time.Unix(0, nanos),
				value:     value[1],
			})
		}
	}

	slices.SortFunc(logLines, func(a, b logLine) int {
		return cmp.Compare(a.timestamp.UnixNano(), b.timestamp.UnixNano())
	})

	return logLines, nil
}

func deduplicateLogLines(logLines []logLine) []logLine {
	result := make([]logLine, 0)
	uniqueMap := make(map[string]bool)
	for _, ll := range logLines {
		uniqueKey := fmt.Sprintf("%d_%s", ll.timestamp.UnixNano(), ll.value)
		_, exists := uniqueMap[uniqueKey]
		if !exists {
			uniqueMap[uniqueKey] = true
			result = append(result, ll)
		}
	}
	logrus.Debugf("removed %d duplicates from log result", len(logLines)-len(result))

	return result
}

// lokiResponse represents the root structure of a query response.
type lokiResponse struct {
	Status string           `json:"status"`
	Data   lokiResponseData `json:"data"`
}

// lokiResponseData contains log stream results and metadata ResultType. ResultType could be "stream" oder "vector".
type lokiResponseData struct {
	// ResultType contains the type of the response data. May be one "stream", "matrix", "vector".
	ResultType string             `json:"resultType"`
	Result     []lokiStreamResult `json:"result"`
}

// lokiStreamResult the stream and the log values.
type lokiStreamResult struct {
	Stream lokiStream `json:"stream"`
	// Values contains the logs as slices of which the first field consists of
	// a timestamp as epoch second and the second the log line as JSON
	Values [][]string `json:"values"`
}

// lokiStream contains label metadata for the stream result.
type lokiStream struct {
	Container string `json:"container"`
	Filename  string `json:"filename"`
	Job       string `json:"job"`
	Namespace string `json:"namespace"`
	NodeName  string `json:"node_name"`
	Pod       string `json:"pod"`
	Stream    string `json:"stream"`
	App       string `json:"app"`
}
