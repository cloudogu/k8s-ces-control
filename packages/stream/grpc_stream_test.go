package stream

import (
	"errors"
	"io"
	"strings"
	"testing"

	pbTypes "github.com/cloudogu/ces-control-api/generated/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWriteToStream(t *testing.T) {
	t.Run("should write small data in single chunk", func(t *testing.T) {
		data := []byte("small data")
		mockServer := NewMockGRPCStreamServer(t)

		mockServer.EXPECT().Send(mock.MatchedBy(func(resp *pbTypes.ChunkedDataResponse) bool {
			return string(resp.Data) == "small data"
		})).Return(nil).Once()

		err := WriteToStream(data, mockServer)

		require.NoError(t, err)
	})

	t.Run("should write empty data", func(t *testing.T) {
		data := []byte{}
		mockServer := NewMockGRPCStreamServer(t)

		err := WriteToStream(data, mockServer)

		require.NoError(t, err)
		mockServer.AssertNotCalled(t, "Send")
	})

	t.Run("should write large data in multiple chunks", func(t *testing.T) {
		// Create data larger than chunkSize (64 KiB)
		largeData := make([]byte, chunkSize*2+1000)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		mockServer := NewMockGRPCStreamServer(t)
		var capturedChunks [][]byte

		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Run(func(resp *pbTypes.ChunkedDataResponse) {
				// Make a copy of the data since it might be reused
				chunk := make([]byte, len(resp.Data))
				copy(chunk, resp.Data)
				capturedChunks = append(capturedChunks, chunk)
			}).
			Return(nil).Times(3)

		err := WriteToStream(largeData, mockServer)

		require.NoError(t, err)
		assert.Len(t, capturedChunks, 3)
		assert.Equal(t, chunkSize, len(capturedChunks[0]))
		assert.Equal(t, chunkSize, len(capturedChunks[1]))
		assert.Equal(t, 1000, len(capturedChunks[2]))

		// Verify the data is correctly split
		var reconstructed []byte
		for _, chunk := range capturedChunks {
			reconstructed = append(reconstructed, chunk...)
		}
		assert.Equal(t, largeData, reconstructed)
	})

	t.Run("should return error when send fails", func(t *testing.T) {
		data := []byte("test data")
		mockServer := NewMockGRPCStreamServer(t)
		expectedErr := errors.New("send error")

		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Return(expectedErr).Once()

		err := WriteToStream(data, mockServer)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error on second chunk send failure", func(t *testing.T) {
		largeData := make([]byte, chunkSize+100)
		mockServer := NewMockGRPCStreamServer(t)
		expectedErr := errors.New("send error on second chunk")

		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Return(nil).Once()
		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Return(expectedErr).Once()

		err := WriteToStream(largeData, mockServer)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestWriteReaderToStream(t *testing.T) {
	t.Run("should write small data from reader", func(t *testing.T) {
		reader := strings.NewReader("small data")
		mockServer := NewMockGRPCStreamServer(t)

		mockServer.EXPECT().Send(mock.MatchedBy(func(resp *pbTypes.ChunkedDataResponse) bool {
			return string(resp.Data) == "small data"
		})).Return(nil).Once()

		err := WriteReaderToStream(reader, mockServer)

		require.NoError(t, err)
	})

	t.Run("should write empty data from reader", func(t *testing.T) {
		reader := strings.NewReader("")
		mockServer := NewMockGRPCStreamServer(t)

		err := WriteReaderToStream(reader, mockServer)

		require.NoError(t, err)
		mockServer.AssertNotCalled(t, "Send")
	})

	t.Run("should write large data from reader in multiple chunks", func(t *testing.T) {
		// Create data larger than chunkSize
		largeData := make([]byte, chunkSize*2+1000)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}
		reader := strings.NewReader(string(largeData))

		mockServer := NewMockGRPCStreamServer(t)
		var capturedChunks [][]byte

		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Run(func(resp *pbTypes.ChunkedDataResponse) {
				chunk := make([]byte, len(resp.Data))
				copy(chunk, resp.Data)
				capturedChunks = append(capturedChunks, chunk)
			}).
			Return(nil).Times(3)

		err := WriteReaderToStream(reader, mockServer)

		require.NoError(t, err)
		assert.Len(t, capturedChunks, 3)

		// Verify the data is correctly split
		var reconstructed []byte
		for _, chunk := range capturedChunks {
			reconstructed = append(reconstructed, chunk...)
		}
		assert.Equal(t, largeData, reconstructed)
	})

	t.Run("should return error when send fails", func(t *testing.T) {
		reader := strings.NewReader("test data")
		mockServer := NewMockGRPCStreamServer(t)
		expectedErr := errors.New("send error")

		mockServer.EXPECT().Send(mock.AnythingOfType("*types.ChunkedDataResponse")).
			Return(expectedErr).Once()

		err := WriteReaderToStream(reader, mockServer)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should return error when reader fails", func(t *testing.T) {
		mockServer := NewMockGRPCStreamServer(t)
		expectedErr := errors.New("read error")
		errorReader := &errorReader{err: expectedErr}

		err := WriteReaderToStream(errorReader, mockServer)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("should handle reader that returns data and error together", func(t *testing.T) {
		reader := &partialErrorReader{
			data: []byte("partial data"),
			err:  errors.New("read error after data"),
		}
		mockServer := NewMockGRPCStreamServer(t)

		mockServer.EXPECT().Send(mock.MatchedBy(func(resp *pbTypes.ChunkedDataResponse) bool {
			return string(resp.Data) == "partial data"
		})).Return(nil).Once()

		err := WriteReaderToStream(reader, mockServer)

		require.Error(t, err)
		assert.ErrorContains(t, err, "read error after data")
	})
}

// errorReader is a reader that always returns an error
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

// partialErrorReader returns data once then an error
type partialErrorReader struct {
	data []byte
	err  error
	read bool
}

func (p *partialErrorReader) Read(buf []byte) (n int, err error) {
	if p.read {
		return 0, io.EOF
	}
	p.read = true
	n = copy(buf, p.data)
	return n, p.err
}
