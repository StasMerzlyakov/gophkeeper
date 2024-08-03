package handler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewFileAccessor(accessor FileAccessor) *fileAccessor {
	return &fileAccessor{
		accessor: accessor,
	}
}

type fileAccessor struct {
	proto.UnimplementedFileAccessorServer
	accessor FileAccessor
}

func (fa *fileAccessor) GetFileInfoList(ctx context.Context, empty *empty.Empty) (*proto.GetFileInfoListResponse, error) {
	action := domain.GetAction(1)

	dd, err := fa.accessor.GetFileInfoList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}

	resp := &proto.GetFileInfoListResponse{}

	for _, inf := range dd {
		resp.FileInfo = append(resp.FileInfo, &proto.FileInfo{
			Name: inf.Name,
		})
	}
	return resp, nil
}

func (fa *fileAccessor) DeleteFileInfo(ctx context.Context, req *proto.DeleteFileInfoRequest) (*empty.Empty, error) {
	action := domain.GetAction(1)

	err := fa.accessor.DeleteFileInfo(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("%v err - %w", action, err)
	}
	return nil, nil
}

func (fa *fileAccessor) UploadFile(stream proto.FileAccessor_UploadFileServer) (err error) {

	action := domain.GetAction(1)

	ctx := stream.Context()
	log := domain.GetCtxLogger(ctx)
	log.Debugw(action, "msg", "start")

	var storage domain.StreamFileWriter

	storage, err = fa.accessor.CreateStreamFileWriter(ctx)
	if err != nil {
		return fmt.Errorf("%v err - %w", action, err)
	}

	defer func() {
		if err != nil {
			log.Debugw(action, "msg", "start rollback")
			if rlErr := storage.Rollback(ctx); rlErr != nil {
				log.Errorw(action, "err", fmt.Sprintf("storage.Rollback err %s", rlErr.Error()))
			}
		} else {
			log.Debugw(action, "msg", "start commit")
			if cmErr := storage.Commit(ctx); cmErr != nil {
				log.Errorw(action, "err", fmt.Sprintf("storage.Commit err %s", cmErr.Error()))
			}
			if err := stream.SendAndClose(nil); err != nil {
				log.Errorw(action, "err", fmt.Sprintf("SendAndClose err %s", err.Error()))
			}
		}
	}()
	received := 0
	var req *proto.UploadFileRequest

	for {
		req, err = stream.Recv()

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = fmt.Errorf("%v recv err - %w", action, err)
			return err
		}

		if req.Cancel {
			// client cancel operation
			err = fmt.Errorf("%v client operation cancelation", action)
			return err
		}

		received += len(req.Data)

		log.Debugw(action, "msg", fmt.Sprintf("received %d, actual %d", req.SizeInBytes, len(req.Data)))

		err = storage.WriteChunk(ctx, req.Name, req.Data)
		if err != nil {
			err = fmt.Errorf("%v write chunk err - %w", action, err)
			return err
		}
	}

	log.Debugw(action, "msg", fmt.Sprintf("received %d", received))
	return nil
}

func (fa *fileAccessor) LoadFile(lr *proto.LoadFileRequest, ls proto.FileAccessor_LoadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method LoadFile not implemented")
}

//rpc DeleteFileInfo(DeleteFileInfoRequest) returns (google.protobuf.Empty);
