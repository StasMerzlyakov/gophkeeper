package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileAccessor() *fileAccessor {
	return &fileAccessor{}
}

type fileAccessor struct {
	appServer AppServer
	helper    DomainHelper
}

func (fl *fileAccessor) AppSever(appServer AppServer) *fileAccessor {
	fl.appServer = appServer
	return fl
}

func (fl *fileAccessor) DomainHelper(helper DomainHelper) *fileAccessor {
	fl.helper = helper
	return fl
}

func (fl *fileAccessor) UploadFile(ctx context.Context, info *domain.FileInfo) error {
	if err := fl.helper.CheckFileForRead(info); err != nil {
		return err
	}
	return nil
}
