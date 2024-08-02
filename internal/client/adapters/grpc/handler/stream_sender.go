package handler

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewStreamSender(name string, client proto.FileAccessor_UploadFileClient) *streamSender {
	return &streamSender{
		name:   name,
		client: client,
	}
}

var _ domain.StreamSender = (*streamSender)(nil)

type streamSender struct {
	name   string
	client proto.FileAccessor_UploadFileClient
}

func (ss *streamSender) Send(chunk []byte) error {

	if len(chunk) == 0 {
		return nil
	}

	return ss.client.Send(&proto.UploadFileRequest{
		Name:        ss.name,
		SizeInBytes: int32(len(chunk)),
		Data:        chunk,
		IsLastChunk: false,
	})
}

func (ss *streamSender) CloseAndRecv() error {

	if err := ss.client.Send(&proto.UploadFileRequest{
		Name:        ss.name,
		SizeInBytes: 0,
		IsLastChunk: true,
	}); err != nil {
		fmt.Println("%w can't send the last chunk", err)
	}

	if _, err := ss.client.CloseAndRecv(); err != nil {
		fmt.Println("%w can't close client", err)
	}
	return nil
}
