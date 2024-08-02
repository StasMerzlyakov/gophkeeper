package handler

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
)

func (h *handler) GetFileInfoList(ctx context.Context) ([]domain.FileInfo, error) {
	resp, err := h.fileAccessor.GetFileInfoList(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: getFileInfoList err ", err)
	}

	var lst []domain.FileInfo
	for _, inf := range resp.FileInfo {
		lst = append(lst, domain.FileInfo{
			Name: inf.Name,
		})
	}

	return lst, nil
}

func (h *handler) DeleteFileInfo(ctx context.Context, name string) error {
	_, err := h.fileAccessor.DeleteFileInfo(ctx, &proto.DeleteFileInfoRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("%w: DeleteFileInfoList err ", err)
	}

	return nil
}

func (h *handler) SendFile(ctx context.Context, name string) (domain.StreamSender, error) {

	if stream, err := h.fileAccessor.UploadFile(ctx); err != nil {
		return nil, fmt.Errorf("%w SendFile err", err)
	} else {
		return NewStreamSender(name, stream), nil
	}
}