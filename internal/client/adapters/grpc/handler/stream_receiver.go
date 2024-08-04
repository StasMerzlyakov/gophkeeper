package handler

import (
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func NewStreamReceiver(stream proto.FileAccessor_LoadFileClient) *streamReceiver {
	return &streamReceiver{
		stream: stream,
	}
}

var _ domain.StreamFileReader = (*streamReceiver)(nil)

type streamReceiver struct {
	stream   proto.FileAccessor_LoadFileClient
	fileSize int64
	once     sync.Once
}

func (sr *streamReceiver) FileSize() int64 {
	return sr.fileSize
}
func (sr *streamReceiver) Next() ([]byte, error) {
	if rep, err := sr.stream.Recv(); err != nil {
		return nil, fmt.Errorf("%w receive err", err)
	} else {
		sr.once.Do(func() {
			sr.fileSize = int64(rep.SizeInBytes)
		})
		return rep.Data[:], nil
	}
}

func (sr *streamReceiver) Close() {
	action := domain.GetAction(1)
	log := app.GetMainLogger()

	log.Debug("%v start", action)

	if err := sr.stream.CloseSend(); err != nil {
		log.Debug("%v close err %v", action, err.Error())
		return
	}
	log.Debug("%v complete", action)
}
