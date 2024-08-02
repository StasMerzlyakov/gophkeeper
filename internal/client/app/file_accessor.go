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
	if err := fl.helper.CheckFileForRead(info); err != nil {
		return nil, fmt.Errorf("%w upload file err", err)
	}

	// test by local storage
	if fl.appStorage.IsFileInfoExists(info.Name) {
		return nil, fmt.Errorf("%w fileInfo %s already exists. change name", domain.ErrClientDataIncorrect, info.Name)
	}

	reader, err := fl.helper.CreateFileStreamer(info)
	if err != nil {
		return nil, fmt.Errorf("%w can't upload file %s", err, info.Name)
	}

	sendCtx, sendCancelFn := context.WithCancel(ctx)
	sender, err := fl.appServer.SendFile(sendCtx, info.Name)
	if err != nil {
		sendCancelFn()
		return nil, fmt.Errorf("%w can't upload file %s", err, info.Name)
	}

	readChan := make(chan []byte) // chan for file readed chunck

	preparedChan := make(chan []byte) // chan for prepared chunks (encryption)

	// Reading operation
	fileSize := int(reader.FileSize())
	progerssFn(0, fileSize)

	var opWg sync.WaitGroup

	errorChan := make(chan error)
	doneCh := make(chan struct{})

	var opErr error

	cancelFn = func() {
		errorChan <- fmt.Errorf("opeartion was cancelsed")
	}

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("whait goroutine complete")
		defer opWg.Done()
		select {
		case opErr = <-errorChan: // error occures
			close(doneCh)
		case <-fl.stopCh:
		}
	}()

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("read goroutine complete")
		defer opWg.Done()
		defer reader.Close()
		sended := 0
		var rdCn chan []byte
		var isLast bool
	Loop:
		for {
			chunk, err := reader.Next()
			chunkSize := len(chunk)
			sended += chunkSize
			progerssFn(sended, fileSize)

			if err != nil {
				if errors.Is(err, io.EOF) {
					isLast = true
				} else {
					errorChan <- err
				}
			}
			if len(chunk) > 0 {
				rdCn = readChan
			}

			select {
			case rdCn <- chunk:
				sended += chunkSize
				progerssFn(sended, fileSize)
				if isLast {
					close(readChan)
					break Loop
				}
				rdCn = nil
			case <-fl.stopCh:
				// complete operation and finish
				break Loop
			case <-doneCh:
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

		var encrypted []byte
		var rCh = readChan
		var wrtCh chan ([]byte)

	Loop:
		for {
			select {
			case val, ok := <-rCh:
				if !ok {
					// success
					close(preparedChan)
					break Loop
				} else {
					// TODO - do encryption
					encrypted = val
					time.Sleep(500 * time.Millisecond)
					rCh = nil
					wrtCh = preparedChan
				}
			case wrtCh <- encrypted:
				encrypted = nil
				rCh = readChan
				wrtCh = nil
			case <-doneCh:
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
	Loop:
		for {
			select {
			case chunk, ok := <-preparedChan:
				if !ok {
					// success
					err := sender.CloseAndRecv()
					if err != nil {
						errorChan <- err
					} else {
						close(doneCh)
					}
					break Loop
				} else {
					if err := sender.Send(chunk); err != nil {
						errorChan <- err
						break Loop
					}
				}
			case <-doneCh:
				// any error occurs
				break Loop
			case <-fl.stopCh:
				// complete operation and finish
				break Loop
			}
		}
	}()

	// Whait result gorotine
	fl.wg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("wait result goroutine complete")
		defer fl.wg.Done()
		defer sendCancelFn()
		opWg.Wait()
		if opErr == nil {
			fl.appStorage.AddFileInfo(info)
		}
		resultHandler(opErr)
	}()

	return cancelFn, nil
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
