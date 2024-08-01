package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileAccessor() *fileAccessor {
	return &fileAccessor{}
}

type fileAccessor struct {
	appServer  AppServer
	appStorage AppStorage
	helper     DomainHelper
}

func (fl *fileAccessor) AppServer(appServer AppServer) *fileAccessor {
	fl.appServer = appServer
	return fl
}

func (fl *fileAccessor) AppStorage(appStorage AppStorage) *fileAccessor {
	fl.appStorage = appStorage
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

	// test by local storage
	if fl.appStorage.IsFileInfoExists(info.Name) {
		return fmt.Errorf("%w fileInfo %s already exists. change name", domain.ErrClientDataIncorrect, info.Name)
	}
	// TODO connect to server

	// add to local storage
	_ = fl.appStorage.AddFileInfo(info)

	return nil
}

func (fl *fileAccessor) GetFileInfoList(ctx context.Context) error {
	// TODO get from server
	return nil
}

func (fl *fileAccessor) DeleteFile(ctx context.Context, name string) error {
	return fl.appStorage.DeleteFileInfo(name)
}

func (fl *fileAccessor) SaveFile(ctx context.Context, info *domain.FileInfo) error {
	if err := fl.helper.CheckFileForWrite(info); err != nil {
		return err
	}

	// TODO store to server

	return fl.appStorage.UpdateFileInfo(info)
}
