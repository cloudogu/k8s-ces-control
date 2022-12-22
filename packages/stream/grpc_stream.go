package stream

import (
	pbTypes "github.com/cloudogu/k8s-ces-control/generated/types"
)

const (
	chunkSize = 64 * 1024 // 64 KiB
)

type GRPCStreamServer interface {
	Send(response *pbTypes.ChunkedDataResponse) error
}

type Writer func([]byte, GRPCStreamServer) error

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
