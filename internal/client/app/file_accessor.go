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
	progerssFn func(send int, all int),
	cancelChan <-chan struct{}) error {

	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := fl.helper.CheckFileForRead(info); err != nil {
		err := fmt.Errorf("%w - %v error - upload file err %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		return err
	}

	// test by local storage
	if fl.appStorage.IsFileInfoExists(info.Name) {
		err := fmt.Errorf("%w fileInfo %s already exists. change name", domain.ErrClientDataIncorrect, info.Name)
		log.Warn(err.Error())
		return err
	}

	reader, err := fl.helper.CreateFileStreamer(info)
	if err != nil {
		err := fmt.Errorf("%w upload file %s err %s", domain.ErrClientInternal, info.Name, err.Error())
		log.Warn(err.Error())
		return err
	}

	forEncryptChan := make(chan []byte) // chan for file readed chunck

	forSendChan := make(chan []byte) // chan for prepared chunks (encryption)

	// Reading operation
	fileSize := int(reader.FileSize())
	progerssFn(0, fileSize)

	var opWg sync.WaitGroup

	errorChan := make(chan error, 1)
	doneCh := make(chan struct{})
	stopCh := make(chan struct{})

	var opErr error

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("wait goroutine complete")
		defer opWg.Done()
		select {
		case opErr = <-errorChan: // error occures
			log.Warn(fmt.Sprintf("opeartion %s error %s", action, opErr.Error()))
			close(stopCh)
		case <-cancelChan:
			log.Debug("opeartion was cancelsed")
			opErr = fmt.Errorf("opeartion was cancelsed")
			close(stopCh)
		case <-fl.stopCh:
			log.Warn("opeartion 1")
		case <-doneCh: // success
			log.Warn("opeartion 2")
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
			case <-stopCh:
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
			case <-stopCh:
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

		sendCrx, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		sender, err := fl.appServer.SendFile(sendCrx)
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
					err := sender.Commit(sendCrx)
					if err != nil {
						GetMainLogger().Debugf("sender closed err %s", err.Error())
						errorChan <- err
					} else {
						GetMainLogger().Debugf("sender closed")
						defer close(doneCh)
					}
					break Loop
				} else {
					if err := sender.WriteChunk(sendCrx, info.Name, chunk); err != nil {
						errorChan <- err
						GetMainLogger().Infof("chunk sending err %s", err.Error())
						break Loop
					}
					GetMainLogger().Debugf("chunk is sent")
				}
			case <-stopCh:
				// any error occurs
				sender.Rollback(sendCrx)
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
	return nil
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
