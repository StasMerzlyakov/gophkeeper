package handler_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	gomock "github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type bufSaver struct {
	name       string
	buf        bytes.Buffer
	closeCount atomic.Int32
}

func (bf *bufSaver) Send(ctx context.Context, name string, chunk []byte) error {
	if bf.name != name {
		return fmt.Errorf("expected %s actual %s", bf.name, name)
	}
	_, err := bf.buf.Write(chunk)
	return err
}

func (bf *bufSaver) CloseAndRecv(ctx context.Context) error {
	val := bf.closeCount.Add(1)
	if val == 1 {
		return nil
	}
	return fmt.Errorf("unexpected CloseAndRecv invokation %d", val)
}

var _ domain.StreamFileWriter = (*bufSaver)(nil)

type fileServ struct {
	name string
	grpc.ServerStream
	buff       []byte
	chunkSize  int
	closeCount atomic.Int32
	ctx        context.Context
}

func (fSrv *fileServ) SendAndClose(*empty.Empty) error {
	val := fSrv.closeCount.Add(1)
	if val == 1 {
		return nil
	}
	return fmt.Errorf("unexpected CloseAndRecv invokation %d", val)
}

func (fSrv *fileServ) Context() context.Context {
	return fSrv.ctx
}

func (fSrv *fileServ) Recv() (*proto.UploadFileRequest, error) {
	if len(fSrv.buff) == 0 {
		return nil, fmt.Errorf("unexpected call")
	}

	var res []byte
	if len(fSrv.buff) > fSrv.chunkSize {
		res, fSrv.buff = fSrv.buff[:fSrv.chunkSize], fSrv.buff[fSrv.chunkSize:]
		return &proto.UploadFileRequest{
			Name:        fSrv.name,
			SizeInBytes: int32(len(res)),
			Data:        res,
			IsLastChunk: false,
		}, nil
	} else {
		res, fSrv.buff = fSrv.buff[:], nil
		return &proto.UploadFileRequest{
			Name:        fSrv.name,
			SizeInBytes: int32(len(res)),
			Data:        res,
			IsLastChunk: true,
		}, nil
	}
}

var _ proto.FileAccessor_UploadFileServer = (*fileServ)(nil)

func TestFileAccessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("upload_ok", func(t *testing.T) {
		mockFileAccessor := NewMockFileAccessor(ctrl)

		fileName := "fileName"
		bufSaver := &bufSaver{
			name: fileName,
		}
		mockFileAccessor.EXPECT().CreateStreamSaver().Return(bufSaver, nil).Times(1)

		chunkSize := 1024
		iter := 10
		tailSize := 512

		fullSize := iter*chunkSize + tailSize
		toSend := make([]byte, fullSize)

		n, err := rand.Read(toSend)
		require.NoError(t, err)
		require.Equal(t, fullSize, n)

		aService := handler.NewFileAccessor(mockFileAccessor)
		flSrv := &fileServ{
			name:      fileName,
			buff:      toSend,
			chunkSize: chunkSize,
			ctx:       context.Background(),
		}
		err = aService.UploadFile(flSrv)
		require.NoError(t, err)

		require.True(t, bytes.Equal(toSend, bufSaver.buf.Bytes()))

		require.Equal(t, int32(1), bufSaver.closeCount.Load())
		require.Equal(t, int32(1), flSrv.closeCount.Load())

	})

	t.Run("delete_file_ok", func(t *testing.T) {
		mockService := NewMockFileAccessor(ctrl)
		mockService.EXPECT().DeleteFileInfo(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, name string) error {
			require.Equal(t, "fileName", name)
			return nil
		}).Times(1)

		aService := handler.NewFileAccessor(mockService)
		_, err := aService.DeleteFileInfo(context.Background(), &proto.DeleteFileInfoRequest{
			Name: "fileName",
		})

		require.NoError(t, err)
	})

	t.Run("get_filelist_ok", func(t *testing.T) {
		mockService := NewMockFileAccessor(ctrl)
		mockService.EXPECT().GetFileInfoList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.FileInfo, error) {
			return []domain.FileInfo{
				{
					Name: "number1",
				},
				{
					Name: "number2",
				},
			}, nil
		}).Times(1)

		aService := handler.NewFileAccessor(mockService)
		resp, err := aService.GetFileInfoList(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, 2, len(resp.FileInfo))

		assert.Equal(t, "number1", resp.FileInfo[0].Name)

		assert.Equal(t, "number2", resp.FileInfo[1].Name)
	})

	t.Run("get_filelist_err", func(t *testing.T) {
		mockService := NewMockFileAccessor(ctrl)
		testErr := errors.New("testErr")
		mockService.EXPECT().GetFileInfoList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.FileInfo, error) {
			return nil, testErr
		}).Times(1)

		aService := handler.NewFileAccessor(mockService)
		resp, err := aService.GetFileInfoList(context.Background(), nil)
		require.ErrorIs(t, err, testErr)
		require.Nil(t, resp)
	})
}
