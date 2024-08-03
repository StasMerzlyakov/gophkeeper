package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewStreamSender(client proto.FileAccessor_UploadFileClient) *streamSender {
	return &streamSender{
		client: client,
	}
}

var _ domain.StreamFileWriter = (*streamSender)(nil)

type streamSender struct {
	client proto.FileAccessor_UploadFileClient
}

func (ss *streamSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	if len(chunk) == 0 {
		return nil
	}

	return ss.client.Send(&proto.UploadFileRequest{
		Name:        name,
		SizeInBytes: int32(len(chunk)),
		Data:        chunk,
		IsLastChunk: false,
	})
}

func (ss *streamSender) Close(ctx context.Context) error {

	if err := ss.client.Send(&proto.UploadFileRequest{
		Name:        "",
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
