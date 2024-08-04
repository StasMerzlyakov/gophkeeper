package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileAccessor() *fileAccessor {
	return &fileAccessor{
		stopCh: make(chan struct{}),
	}
}

type fileAccessor struct {
	appServer  AppServer
	appStorage AppStorage
	helper     DomainHelper
	wg         sync.WaitGroup
	stopCh     chan struct{}
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

func (fl *fileAccessor) Stop(ctx context.Context) {
	close(fl.stopCh)

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh) // stop via close channel
		fl.wg.Wait()
	}()

	select {
	case <-ctx.Done():
		panic("not all goroutines stopped")
	case <-doneCh:
		GetMainLogger().Debugf("all long operation completed")
		return
	}
}

func (fl *fileAccessor) GetFileInfoList(ctx context.Context) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)
	if lst, err := fl.appServer.GetFileInfoList(ctx); err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	} else {
		fl.appStorage.SetFilesInfo(lst)
	}
	log.Debugf("%v complete", action)
	return nil
}

func (fl *fileAccessor) DeleteFile(ctx context.Context, name string) error {
	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)
	if err := fl.appServer.DeleteFileInfo(ctx, name); err != nil {
		err := fmt.Errorf("%w - %v error", err, action)
		log.Warn(err.Error())
		return err
	}
	log.Debugf("%v complete", action)
	return fl.appStorage.DeleteFileInfo(name)
}
