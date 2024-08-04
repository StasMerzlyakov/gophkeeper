package app

import (
	"context"
	"errors"
	"fmt"
	"io"
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

// UploadFile start long server uploading operation.
func (fl *fileAccessor) UploadFile(ctx context.Context,
	info *domain.FileInfo,
	progerssFn func(send int, all int),
	cancelChan <-chan struct{},
	errorChan chan<- error) {

	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := fl.helper.CheckFileForRead(info); err != nil {
		err := fmt.Errorf("%w - %v error - upload file err %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	// test by local storage
	if fl.appStorage.IsFileInfoExists(info.Name) {
		err := fmt.Errorf("%w fileInfo %s already exists. change name", domain.ErrClientDataIncorrect, info.Name)
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	reader, err := fl.helper.CreateFileStreamer(info)
	if err != nil {
		err := fmt.Errorf("%w upload file %s err %s", domain.ErrClientInternal, info.Name, err.Error())
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	forEncryptChan := make(chan []byte) // chan for file readed chunck

	forSendChan := make(chan []byte) // chan for prepared chunks (encryption)

	// Reading operation
	fileSize := int(reader.FileSize())
	if fileSize <= 0 {
		errorChan <- fmt.Errorf("%w unexpected file size", domain.ErrClientInternal)
		return
	}
	if progerssFn != nil {
		progerssFn(0, fileSize)
	}

	var opWg sync.WaitGroup

	jobDoneCh := make(chan struct{})
	jobTermiatedCh := make(chan struct{})

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("wait goroutine complete")
		defer opWg.Done()
		select {
		case <-cancelChan:
			errorChan <- fmt.Errorf("%w opeartion was canceled", domain.ErrClientInteruptoin)
			close(jobTermiatedCh)
			log.Debugf("%v operation cancled", action)
		case <-fl.stopCh:
			errorChan <- fmt.Errorf("%w opeartion was stopped", domain.ErrClientAppStopped)
			close(jobTermiatedCh)
			log.Debugf("%v app stop", action)
		case <-jobDoneCh: // success
			log.Debugf("%v done", action)
		case <-jobTermiatedCh: // error occures
			log.Debugf("%v termianted", action)
		}
	}()

	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("read goroutine complete")
		defer opWg.Done()
		defer reader.Close()
		readed := 0
		var rdCn chan []byte
		var isLast bool
		var progressCount int
	Loop:
		for {
			// test stop first
			select {
			case <-jobTermiatedCh:
				log.Debug("read goroutine terminated")
				break Loop
			default:
			}
			chunk, err := reader.Next()
			chunkSize := len(chunk)
			readed += chunkSize
			progressCount++
			if progressCount == 10 {
				if progerssFn != nil {
					progerssFn(readed, fileSize)
				}
				progressCount = 0
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					isLast = true
				} else {
					errorChan <- err
					close(jobTermiatedCh)
					break Loop
				}
			}
			if len(chunk) > 0 {
				rdCn = forEncryptChan
			} else {
				close(forEncryptChan)
				break Loop
			}

			select {
			case rdCn <- chunk:
				if isLast {
					close(forEncryptChan)
					break Loop
				}
				rdCn = nil
			case <-jobTermiatedCh:
				log.Debug("read goroutine terminated")
				break Loop
			}
		}

	}()

	// Encryption operation TODO
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
			// test stop first
			select {
			case <-jobTermiatedCh:
				log.Debug("encr goroutine terminated")
				break Loop
			default:
			}

			select {
			case val, ok := <-rCh:
				if !ok {
					// success
					break Loop
				} else {
					// TODO - do encryption
					encrypted = val
					rCh = nil
					wrtCh = forSendChan
				}
			case wrtCh <- encrypted:
				encrypted = nil
				rCh = forEncryptChan
				wrtCh = nil
			case <-jobTermiatedCh:
				log.Debug("encr goroutine terminated")
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

		sender, err := fl.appServer.CreateFileSender(sendCrx)
		if err != nil {
			err := fmt.Errorf("%w upload file %s err", err, info.Name)
			log.Warn(err.Error())
			errorChan <- err
			return
		}

	Loop:
		for {
			select {
			case <-jobTermiatedCh:
				log.Debug("send goroutine terminated")
				if err := sender.Rollback(sendCrx); err != nil {
					log.Warnf("sender rollabck error %v", err.Error())
				}
				break Loop
			default:
			}
			select {
			case chunk, ok := <-forSendChan:
				if !ok {
					// success
					if err := sender.Commit(sendCrx); err != nil {
						log.Warnf("write commit error %v", err.Error())
						errorChan <- err
						close(jobTermiatedCh)
					} else {
						close(jobDoneCh)
					}
					break Loop
				} else {
					if err := sender.WriteChunk(sendCrx, info.Name, chunk); err != nil {
						log.Warnf("write chun error %v", err.Error())
						errorChan <- err
						if err := sender.Rollback(sendCrx); err != nil {
							log.Warnf("sender rollabck error %v", err.Error())
						}
						close(jobTermiatedCh)
						break Loop
					}
				}
			case <-jobTermiatedCh:
				log.Debug("send goroutine terminated")
				if err := sender.Rollback(sendCrx); err != nil {
					log.Warnf("sender rollabck error %v", err.Error())
				}
				break Loop
			}
		}
	}()

	opWg.Wait()
	log.Debugf("%v complete", action)
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
