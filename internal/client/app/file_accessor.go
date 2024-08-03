package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileAccessor() *fileAccessor {
	return &fileAccessor{
		stopCh: make(chan struct{}),
	}
}

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
)

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

// UploadFile start long server uploading operation.
func (fl *fileAccessor) UploadFile(ctx context.Context,
	info *domain.FileInfo,
	resultHandler func(err error),
	progerssFn func(send int, all int)) (cancelFn func(), err error) {

	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := fl.helper.CheckFileForRead(info); err != nil {
		err := fmt.Errorf("%w - %v error - upload file err %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return nil, err
	}

	// test by local storage
	if fl.appStorage.IsFileInfoExists(info.Name) {
		err := fmt.Errorf("%w fileInfo %s already exists. change name", domain.ErrClientDataIncorrect, info.Name)
		log.Warn(err.Error())
		return nil, err
	}

	reader, err := fl.helper.CreateFileStreamer(info)
	if err != nil {
		err := fmt.Errorf("%w upload file %s err %s", domain.ErrClientInternal, info.Name, err.Error())
		log.Warn(err.Error())
		return nil, err
	}

	forEncryptChan := make(chan []byte) // chan for file readed chunck

	forSendChan := make(chan []byte) // chan for prepared chunks (encryption)

	// Reading operation
	fileSize := int(reader.FileSize())
	progerssFn(0, fileSize)

	var opWg sync.WaitGroup

	errorChan := make(chan error)
	doneCh := make(chan struct{})
	cancelCh := make(chan struct{})

	var opErr error

	cancelFn = func() {
		errorChan <- fmt.Errorf("opeartion was cancelsed")
		log.Debug("opeartion was cancelsed")
	}

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("wait goroutine complete")
		defer opWg.Done()
		select {
		case opErr = <-errorChan: // error occures
			close(cancelCh)
			log.Warn(fmt.Sprintf("opeartion %s error %s", action, err.Error()))
		case <-fl.stopCh:
		case <-doneCh: // success
		}
	}()

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("read goroutine complete")
		defer opWg.Done()
		defer reader.Close()
		defer close(forEncryptChan)
		readed := 0
		var rdCn chan []byte
		var isLast bool
	Loop:
		for {
			chunk, err := reader.Next()
			chunkSize := len(chunk)
			readed += chunkSize
			progerssFn(readed, fileSize)

			if err != nil {
				if errors.Is(err, io.EOF) {
					isLast = true
				} else {
					errorChan <- err
				}
			}
			if len(chunk) > 0 {
				rdCn = forEncryptChan
			} else {
				break Loop
			}

			select {
			case rdCn <- chunk:

				if isLast {
					break Loop
				}
				rdCn = nil
			case <-fl.stopCh:
				// complete operation and finish
				break Loop
			case <-cancelCh:
				// any error occurs
				break Loop
			}
		}

	}()

	// Encryption operation
	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("encr goroutine complete")
		defer opWg.Done()
		defer close(forSendChan)

		var encrypted []byte
		var rCh = forEncryptChan
		var wrtCh chan ([]byte)

	Loop:
		for {
			select {
			case val, ok := <-rCh:
				if !ok {
					// success
					break Loop
				} else {
					// TODO - do encryption
					encrypted = val
					time.Sleep(500 * time.Millisecond)
					rCh = nil
					wrtCh = forSendChan
				}
			case wrtCh <- encrypted:
				encrypted = nil
				rCh = forEncryptChan
				wrtCh = nil
			case <-cancelCh:
				// any error occurs
				break Loop
			case <-fl.stopCh:
				// complete operation and finish
				break Loop
			}
		}
	}()

	// Send goroutine
	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("send goroutine complete")
		defer opWg.Done()
		defer close(doneCh)

		sender, err := fl.appServer.SendFile(ctx)
		if err != nil {
			err := fmt.Errorf("%w upload file %s err", err, info.Name)
			log.Warn(err.Error())
			errorChan <- err
			return
		}

	Loop:
		for {
			select {
			case chunk, ok := <-forSendChan:
				if !ok {
					// success
					err := sender.Commit(ctx)
					if err != nil {
						GetMainLogger().Debugf("sender closed err %s", err.Error())
						errorChan <- err
					} else {
						GetMainLogger().Debugf("sender closed")
					}
					break Loop
				} else {
					if err := sender.WriteChunk(ctx, info.Name, chunk); err != nil {
						if errors.Is(err, io.EOF) {
							// go grpc call ClientStream.SendMsg is not blocked
							// success
							GetMainLogger().Debugf("message is sent")
							break Loop
						}
						GetMainLogger().Infof("chunk sending err %s", err.Error())
						errorChan <- err
						break Loop
					}
					GetMainLogger().Debugf("chunk is sent")
				}
			case <-cancelCh:
				// any error occurs
				break Loop
			case <-fl.stopCh:
				// complete operation and finish
				break Loop
			}
		}
	}()

	opWg.Wait()
	resultHandler(opErr)
	log.Debugf("%v complete", action)
	return cancelFn, nil
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
