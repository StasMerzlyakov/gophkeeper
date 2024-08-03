package handler

import (
	"context"
	"fmt"

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

func (fa *fileAccessor) UploadFile(fs proto.FileAccessor_UploadFileServer) error {
	action := domain.GetAction(1)

	storage, err := fa.accessor.CreateStreamSaver()
	if err != nil {
		return fmt.Errorf("%v err - %w", action, err)
	}

	for {
		req, recErr := fs.Recv()

		if recErr != nil {
			return fmt.Errorf("%v recv err - %w", action, err)
		}

		if len(req.Data) > 0 {
			err := storage.Send(fs.Context(), req.Name, req.Data)
			if err != nil {
				return fmt.Errorf("%v recv err - %w", action, err)
			}
		}

		if req.IsLastChunk {
			err := storage.CloseAndRecv(fs.Context())
			if err != nil {
				return fmt.Errorf("%v recv close - %w", action, err)
			}
			if err := fs.SendAndClose(nil); err != nil {
				// client dead?
				return fmt.Errorf("%v recv sendAdnClose err - %w", action, err)
			}
			return nil
		}
	}
}

func (fa *fileAccessor) LoadFile(lr *proto.LoadFileRequest, ls proto.FileAccessor_LoadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method LoadFile not implemented")
}

//rpc DeleteFileInfo(DeleteFileInfoRequest) returns (google.protobuf.Empty);
