package logging

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//go:embed testdata/loki-response.json
var lokiResponseTestData []byte

//go:embed testdata/loki-long-response.json
var lokiResponseTestDataLong []byte

//go:embed testdata/loki-empty-response.json
var lokiResponseTestDataEmpty []byte

//go:embed testdata/loki-invalid-response.json
var lokiResponseTestDataInvalid []byte

//go:embed testdata/loki-error-response.json
var lokiResponseTestDataError []byte

//go:embed testdata/loki-non-stream-response.json
var lokiResponseTestDataNoStream []byte

func Test_extractLogLinesFromLokiResponse(t *testing.T) {
	t.Run("should return the response as list of log lines", func(t *testing.T) {
		// given
		lr := &lokiResponse{
			Data: lokiResponseData{
				ResultType: "streams",
				Result: []lokiStreamResult{
					{Values: [][]string{
						// unsorted map!
						{"1655722130600667934", `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
						{"1655722130600667903", `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
						{"1655722130600667919", `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
					}},
				},
			},
		}

		// when
		actual, err := extractLogLinesFromLokiResponse(lr)

		// then
		assert.NoError(t, err)
		expectedLogLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}
		assert.Equal(t, expectedLogLines, actual)
	})

	t.Run("should fail for unparseable timestamp", func(t *testing.T) {
		// given
		lr := &lokiResponse{
			Data: lokiResponseData{
				ResultType: "streams",
				Result: []lokiStreamResult{
					{Values: [][]string{
						// unsorted map!
						{"1655722130600667934", `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
						{"time has run out...", `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
						{"1655722130600667919", `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
					}},
				},
			},
		}

		// when
		_, err := extractLogLinesFromLokiResponse(lr)

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse log timestamp: strconv.ParseInt: parsing \"time has run out...\": invalid syntax")
	})
}

func Test_deduplicateLogLines(t *testing.T) {
	tests := []struct {
		name     string
		logLines []logLine
		want     []logLine
	}{
		{
			name: "remove duplicates",
			logLines: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1235), value: "a"},
				{timestamp: time.Unix(0, 1236), value: "b"},
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1237), value: "c"},
				{timestamp: time.Unix(0, 1237), value: "c"},
				{timestamp: time.Unix(0, 1238), value: "c"},
				{timestamp: time.Unix(0, 1236), value: "b"},
			},
			want: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1235), value: "a"},
				{timestamp: time.Unix(0, 1236), value: "b"},
				{timestamp: time.Unix(0, 1237), value: "c"},
				{timestamp: time.Unix(0, 1238), value: "c"},
			},
		},
		{
			name: "remove no duplicates for different timestamps",
			logLines: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1235), value: "a"},
				{timestamp: time.Unix(0, 1236), value: "a"},
				{timestamp: time.Unix(0, 1237), value: "a"},
			},
			want: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1235), value: "a"},
				{timestamp: time.Unix(0, 1236), value: "a"},
				{timestamp: time.Unix(0, 1237), value: "a"},
			},
		},
		{
			name: "remove no duplicates for different values",
			logLines: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "b"},
				{timestamp: time.Unix(0, 1234), value: "c"},
				{timestamp: time.Unix(0, 1234), value: "aa"},
			},
			want: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "b"},
				{timestamp: time.Unix(0, 1234), value: "c"},
				{timestamp: time.Unix(0, 1234), value: "aa"},
			},
		},
		{
			name: "remove all duplicates",
			logLines: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "a"},
				{timestamp: time.Unix(0, 1234), value: "a"},
			},
			want: []logLine{
				{timestamp: time.Unix(0, 1234), value: "a"},
			},
		},
		{
			name:     "remove none for emtpy slice",
			logLines: []logLine{},
			want:     []logLine{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, deduplicateLogLines(tt.logLines), "deduplicateLogLines(%v)", tt.logLines)
		})
	}
}

func Test_calculateQueryLimit(t *testing.T) {
	tests := []struct {
		name        string
		linesCount  int
		resultCount int
		want        int
	}{
		{"default query limit for 0 lines count", 0, 0, defaultQueryLimit},
		{"default query limit for remaining count higher than default", defaultQueryLimit + 100, 10, defaultQueryLimit},
		{"remaining count lower than default", defaultQueryLimit + 100, defaultQueryLimit + 10, 90},
		{"remaining count is 0", 100, 100, 0},
		{"remaining count is negative", 100, 200, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, calculateQueryLimit(tt.linesCount, tt.resultCount), "calculateQueryLimit(%v, %v)", tt.linesCount, tt.resultCount)
		})
	}
}

func Test_buildLokiQueryUrl(t *testing.T) {
	tests := []struct {
		name        string
		lokiBaseUrl string
		query       string
		startDate   time.Time
		endDate     time.Time
		limit       int
		want        string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "should set all parameters",
			lokiBaseUrl: "http://loki:8001",
			query:       "{pod=~\"test.*\"}",
			startDate:   time.Unix(0, 1701697208000046548),
			endDate:     time.Unix(0, 1696419608000094634),
			limit:       500,
			want:        "http://loki:8001/loki/api/v1/query_range?direction=backward&end=1696419608000094634&limit=500&query=%7Bpod%3D~%22test.%2A%22%7D&start=1701697208000046548",
			wantErr:     assert.NoError,
		},
		{
			name:        "should fail for wrong url",
			lokiBaseUrl: "t:/\\\foo/bar",
			query:       "{pod=~\"test.*\"}",
			startDate:   time.Unix(0, 1701697208000046548),
			endDate:     time.Unix(0, 1696419608000094634),
			limit:       500,
			want:        "",
			wantErr:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildLokiQueryUrl(tt.lokiBaseUrl, tt.query, tt.startDate, tt.endDate, tt.limit)
			if !tt.wantErr(t, err, fmt.Sprintf("buildLokiQueryUrl(%v, %v, %v, %v, %v)", tt.lokiBaseUrl, tt.query, tt.startDate, tt.endDate, tt.limit)) {
				return
			}
			assert.Equalf(t, tt.want, got, "buildLokiQueryUrl(%v, %v, %v, %v, %v)", tt.lokiBaseUrl, tt.query, tt.startDate, tt.endDate, tt.limit)
		})
	}
}

func TestLokiLogProvider_getLogs(t *testing.T) {
	t.Run("should get logs", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=1000&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestData)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		actual, err := sut.getLogs("test", 0)

		// then
		require.NoError(t, err)
		expectedLogLines := []logLine{
			{timestamp: time.Unix(0, 1569266492548155000), value: "bar"},
			{timestamp: time.Unix(0, 1569266497240578000), value: "foo"},
		}
		assert.Equal(t, expectedLogLines, actual)
	})

	t.Run("should get logs with limit", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=11&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestDataLong)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		actual, err := sut.getLogs("test", 11)

		// then
		require.NoError(t, err)
		assert.Len(t, actual, 11)
	})

	t.Run("should truncate logs with limit", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=10&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestDataLong)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		actual, err := sut.getLogs("test", 10)

		// then
		require.NoError(t, err)
		assert.Len(t, actual, 10)
	})

	t.Run("should get no logs with empty response", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=10&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestDataEmpty)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		actual, err := sut.getLogs("test", 10)

		// then
		require.NoError(t, err)
		assert.Len(t, actual, 0)
	})

	t.Run("should fail on building query", func(t *testing.T) {
		// given
		sut := &LokiLogProvider{
			gatewayUrl: "t:/\\\foo/bar",
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		_, err := sut.getLogs("test", 0)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to build loki-query")
	})

	t.Run("should fail to do query", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=1000&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write([]byte("foobar"))
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		_, err := sut.getLogs("test", 0)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "faild to execute loki-query: failed to unmarshal response")
	})

	t.Run("should fail to do extract log lines", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/loki/api/v1/query_range?direction=backward&end=1655722130600667903&limit=1000&query=%7Bpod%3D~%22test.%2A%22%7D&start=1653130130600667903", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestDataInvalid)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			gatewayUrl: svr.URL,
			username:   "admin",
			password:   "admin123",
			clock:      &testClock{time.Unix(0, 1655722130600667903)},
		}

		// when
		_, err := sut.getLogs("test", 0)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to extract logs from loki response: failed to parse log timestamp")
	})

}

func TestLokiLogProvider_doLokiHttpQuery(t *testing.T) {
	t.Run("should do query", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			_, err := w.Write(lokiResponseTestData)
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			username: "admin",
			password: "admin123",
		}

		// when
		actual, err := sut.doLokiHttpQuery(svr.URL + "/test")

		// then
		require.NoError(t, err)
		expectedResponse := &lokiResponse{
			Status: "success",
			Data: lokiResponseData{
				ResultType: "streams",
				Result: []lokiStreamResult{
					{
						Stream: lokiStream{
							Filename: "/var/log/myproject.log",
							Job:      "varlogs",
						},
						Values: [][]string{
							{"1569266497240578000", "foo"},
							{"1569266492548155000", "bar"},
						},
					},
				},
			},
		}
		assert.Equal(t, expectedResponse, actual)
	})

	t.Run("should fail to create request", func(t *testing.T) {
		// given
		sut := &LokiLogProvider{}

		// when
		_, err := sut.doLokiHttpQuery(":\\\\//ff\f")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create request with url")
	})

	t.Run("should fail to execute request", func(t *testing.T) {
		// given
		sut := &LokiLogProvider{
			username: "admin",
			password: "admin123",
		}

		// when
		_, err := sut.doLokiHttpQuery("/test")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to execute request with url [/test]")
	})

	t.Run("should fail for error response", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("500 - Something bad happened!"))
			require.NoError(t, err)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			username: "admin",
			password: "admin123",
		}

		// when
		_, err := sut.doLokiHttpQuery(svr.URL + "/test")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "loki http error: status: 500 Internal Server Error, code: 500; response-body: 500 - Something bad happened!")
	})
	t.Run("should fail for error response without body", func(t *testing.T) {
		// given
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test", r.URL.String())
			username, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "admin", username)
			assert.Equal(t, "admin123", password)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer svr.Close()

		sut := &LokiLogProvider{
			username: "admin",
			password: "admin123",
		}

		// when
		_, err := sut.doLokiHttpQuery(svr.URL + "/test")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "loki http error: status: 500 Internal Server Error, code: 500; response-body: faild to read error response")
	})
}

func TestLokiLogProvider_parseLokiResponse(t *testing.T) {
	t.Run("should parse response", func(t *testing.T) {
		// given

		// when
		actual, err := parseLokiResponse(lokiResponseTestData)

		// then
		require.NoError(t, err)
		expectedResponse := &lokiResponse{
			Status: "success",
			Data: lokiResponseData{
				ResultType: "streams",
				Result: []lokiStreamResult{
					{
						Stream: lokiStream{
							Filename: "/var/log/myproject.log",
							Job:      "varlogs",
						},
						Values: [][]string{
							{"1569266497240578000", "foo"},
							{"1569266492548155000", "bar"},
						},
					},
				},
			},
		}
		assert.Equal(t, expectedResponse, actual)
	})

	t.Run("should fail on error response", func(t *testing.T) {
		// given

		// when
		_, err := parseLokiResponse(lokiResponseTestDataError)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "loki response status is not successfull; status is error")
	})

	t.Run("should fail on non stream response", func(t *testing.T) {
		// given

		// when
		_, err := parseLokiResponse(lokiResponseTestDataNoStream)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "loki response data aren't streams; resultType is not a stream")
	})
}

func TestNewLokiLogProvider(t *testing.T) {
	t.Run("should create LokiLogProvider", func(t *testing.T) {
		// given

		// when
		llp := NewLokiLogProvider("gatewayUrl", "user", "password")

		// then
		require.NotNil(t, llp)
		assert.Equal(t, "gatewayUrl", llp.gatewayUrl)
		assert.Equal(t, "user", llp.username)
		assert.Equal(t, "password", llp.password)
		assert.IsType(t, &realClock{}, llp.clock)
	})
}

func Test_realClock_Now(t *testing.T) {
	sut := new(realClock)
	actual := sut.Now()
	assert.IsType(t, actual, time.Now())
}

type testClock struct {
	time time.Time
}

func (tc *testClock) Now() time.Time {
	return tc.time
}
