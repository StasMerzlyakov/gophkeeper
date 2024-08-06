package handler

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
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

	log := app.GetMainLogger()
	if err := ss.client.Send(&proto.UploadFileRequest{
		Name:        name,
		SizeInBytes: int32(len(chunk)),
		Data:        chunk,
		Cancel:      false,
	}); err != nil {
		log.Errorf("send error %s ", err.Error())
		return err
	} else {
		log.Debugf("send %d bytes sucess", len(chunk))
	}

	return nil
}

func (ss *streamSender) Commit(ctx context.Context) error {
	if _, err := ss.client.CloseAndRecv(); err != nil {
		return err
	}
	return nil
}

func (ss *streamSender) Rollback(ctx context.Context) error {
	// How to send error to server?
	//return ss.client.CloseSend()
	if err := ss.client.Send(&proto.UploadFileRequest{
		Cancel: true,
	}); err != nil {
		return err
	}
	return ss.client.CloseSend()
}
