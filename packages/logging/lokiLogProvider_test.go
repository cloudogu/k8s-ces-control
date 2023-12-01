package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
