package stream

import (
	"io"

	pbTypes "github.com/cloudogu/ces-control-api/generated/types"
)

const (
	chunkSize = 64 * 1024 // 64 KiB
)

// GRPCStreamServer is used to stream data.
type GRPCStreamServer interface {
	// Send sends chunked data.
	Send(response *pbTypes.ChunkedDataResponse) error
}

// Writer is used to write data to a stream server.
type Writer func([]byte, GRPCStreamServer) error

// WriteToStream writes data to stream server in chunks.
func WriteToStream(data []byte, server GRPCStreamServer) error {
	resp := &pbTypes.ChunkedDataResponse{}
	for currentByte := 0; currentByte < len(data); currentByte += chunkSize {
		if currentByte+chunkSize > len(data) {
			resp.Data = data[currentByte:]
		} else {
			resp.Data = data[currentByte : currentByte+chunkSize]
		}
		if err := server.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

// WriteReaderToStream reads from an io.Reader and writes data to stream server in chunks.
func WriteReaderToStream(reader io.Reader, server GRPCStreamServer) error {
	buffer := make([]byte, chunkSize)
	resp := &pbTypes.ChunkedDataResponse{}

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			resp.Data = buffer[:n]
			if sendErr := server.Send(resp); sendErr != nil {
				return sendErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
