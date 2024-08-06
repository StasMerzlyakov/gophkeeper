package handler

import (
	"fmt"
	"io"
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

	action := domain.GetAction(1)
	log := app.GetMainLogger()

	if rep, err := sr.stream.Recv(); err != nil {

		if err == io.EOF {
			return nil, nil
		}

		log.Errorf("%v receive err %v ???", action, sr.stream.Context().Err())
		return nil, fmt.Errorf("%w receive err", err)
	} else {
		sr.once.Do(func() {
			sr.fileSize = int64(rep.SizeInBytes)
		})
		log.Debugf("%v receive %v bytes", action, len(rep.Data))
		return rep.Data[:], nil
	}
}

func (sr *streamReceiver) Close() {
	/*action := domain.GetAction(1)
	log := app.GetMainLogger()

	log.Debugf("%v start", action)

	if err := sr.stream.CloseSend(); err != nil {
		log.Debugf("%v close err %v", action, err.Error())
		return
	}
	log.Debugf("%v complete", action) */
}
