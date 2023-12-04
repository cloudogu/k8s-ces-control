package logging

import (
	"archive/zip"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"
)

func TestNewLoggingService(t *testing.T) {
	t.Run("should create query clock", func(t *testing.T) {
		// given
		llp := &LokiLogProvider{}

		// when
		sut := NewLoggingService(llp)

		// then
		require.NotNil(t, sut)
		assert.NotNil(t, sut.logProvider)
		assert.Equal(t, llp, sut.logProvider)
	})
}

//func Test_writeLogLinesToStream(t *testing.T) {
//	t.Run("should write logs to stream", func(t *testing.T) {
//		// given
//		mockedLogProvider := newMockLogProvider(t)
//
//		logLines := []logLine{
//			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
//			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
//			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
//		}
//
//		// when
//		actualMyDoguLogZip, err := writeLogLinesToStream(mockedLogProvider, "my-dogu", 222, )
//
//		// then
//		require.NoError(t, err)
//		assert.NotNil(t, actualMyDoguLogZip)
//
//		// test text equality
//		zipreader, err := zip.NewReader(bytes.NewReader(actualMyDoguLogZip), int64(len(actualMyDoguLogZip)))
//		require.NoError(t, err)
//		assert.Len(t, zipreader.File, 1)
//
//		for _, zipfile := range zipreader.File {
//			assert.Equal(t, "my-dogu.log", zipfile.Name)
//
//			fc, err := zipfile.Open()
//			require.NoError(t, err)
//			defer func() { _ = fc.Close() }()
//			actualFileContent, err := io.ReadAll(fc)
//			require.NoError(t, err)
//
//			assert.Len(t, actualFileContent, 333)
//		}
//	})
//}

func Test_compressMessages(t *testing.T) {
	t.Run("should compress log lines", func(t *testing.T) {
		// given
		logLines := []logLine{
			{timestamp: time.Unix(0, 1655722130600667903), value: `{"log":"Mon Jun 20 10:48:50 UTC 2022 -- Logging1\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667919), value: `{"log":"Mon Jun 20 10:48:51 UTC 2022 -- Logging2\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
			{timestamp: time.Unix(0, 1655722130600667934), value: `{"log":"Mon Jun 20 10:48:52 UTC 2022 -- Logging3\n","stream":"stdout","time":"2022-06-20T10:48:50.432098057Z"}`},
		}

		// when
		actualMyDoguLogZip, err := compressMessages("my-dogu", logLines)

		// then
		require.NoError(t, err)
		assert.NotNil(t, actualMyDoguLogZip)

		// test text equality
		zipreader, err := zip.NewReader(bytes.NewReader(actualMyDoguLogZip), int64(len(actualMyDoguLogZip)))
		require.NoError(t, err)
		assert.Len(t, zipreader.File, 1)

		for _, zipfile := range zipreader.File {
			assert.Equal(t, "my-dogu.log", zipfile.Name)

			fc, err := zipfile.Open()
			require.NoError(t, err)
			defer func() { _ = fc.Close() }()
			actualFileContent, err := io.ReadAll(fc)
			require.NoError(t, err)

			assert.Len(t, actualFileContent, 333)
		}
	})
}
