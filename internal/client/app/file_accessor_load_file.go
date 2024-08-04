package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

// LoadFile start long server loading operation
func (fl *fileAccessor) LoadFile(ctx context.Context,
	info *domain.FileInfo,
	progerssFn func(send int, all int),
	cancelChan <-chan struct{},
	errorChan chan<- error) {

	log := GetMainLogger()
	action := domain.GetAction(1)
	log.Debugf("%v start", action)

	if err := fl.helper.CheckFileForWrite(info); err != nil {
		err := fmt.Errorf("%w - %v error - load file err %v", domain.ErrClientDataIncorrect, action, err.Error())
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	dir := filepath.Dir(info.Path)
	basename := filepath.Base(info.Path)
	writer, err := fl.helper.CreateStreamFileWriter(dir)
	if err != nil {
		err := fmt.Errorf("%w load file %s err %s", domain.ErrClientInternal, info.Name, err.Error())
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	defer func() {
		err := writer.Rollback(ctx)
		if err != nil {
			err := fmt.Errorf("%w - %v error - writer.Rollback %v", domain.ErrClientDataIncorrect, action, err.Error())
			log.Warn(err.Error())
		}
	}()

	reader, err := fl.appServer.CreateFileReceiver(ctx, info.Name)
	if err != nil {
		err := fmt.Errorf("%w load file %s err %s", domain.ErrClientInternal, info.Name, err.Error())
		log.Warn(err.Error())
		errorChan <- err
		return
	}

	forDecrypChan := make(chan []byte)

	forWriteChan := make(chan []byte)

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
			chunk, err := reader.Next()

			if err != nil {
				if errors.Is(err, io.EOF) {
					isLast = true
				} else {
					errorChan <- err
					close(jobTermiatedCh)
					break Loop
				}
			}
			fileSize := int(reader.FileSize())

			chunkSize := len(chunk)
			readed += chunkSize
			progressCount++
			if progressCount == 10 {
				if progerssFn != nil {
					progerssFn(readed, fileSize)
				}
				progressCount = 0
			}

			if len(chunk) > 0 {
				rdCn = forDecrypChan
			} else {
				close(forDecrypChan)
				break Loop
			}

			select {
			case rdCn <- chunk:
				if isLast {
					close(forDecrypChan)
					break Loop
				}
				rdCn = nil
			case <-jobTermiatedCh:
				log.Debug("read goroutine terminated")
				break Loop
			}
		}
	}()

	// Decrypt operation TODO
	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("decr goroutine complete")
		defer opWg.Done()
		defer close(forWriteChan)

		var decrypted []byte
		var rCh = forDecrypChan
		var wrtCh chan ([]byte)

	Loop:
		for {
			// test stop first
			select {
			case <-jobTermiatedCh:
				log.Debug("decr goroutine terminated")
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
					decrypted = val
					rCh = nil
					wrtCh = forWriteChan
				}
			case wrtCh <- decrypted:
				decrypted = nil
				rCh = forDecrypChan
				wrtCh = nil
			case <-jobTermiatedCh:
				log.Debug("decr goroutine terminated")
				break Loop
			}
		}
	}()
	// Write goroutine
	opWg.Add(1)
	go func() {
		defer GetMainLogger().Debugf("write goroutine complete")
		defer opWg.Done()

		sendCrx, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

	Loop:
		for {
			select {
			case <-jobTermiatedCh:
				log.Debug("send goroutine terminated")
				if err := writer.Rollback(sendCrx); err != nil {
					log.Warnf("sender rollabck error %v", err.Error())
				}
				break Loop
			default:
			}
			select {
			case chunk, ok := <-forWriteChan:
				if !ok {
					// success
					if err := writer.Commit(sendCrx); err != nil {
						log.Warnf("write commit error %v", err.Error())
						errorChan <- err
						close(jobTermiatedCh)
					} else {
						close(jobDoneCh)
					}
					break Loop
				} else {
					if err := writer.WriteChunk(sendCrx, basename, chunk); err != nil {
						log.Warnf("write chun error %v", err.Error())
						errorChan <- err
						if err := writer.Rollback(sendCrx); err != nil {
							log.Warnf("sender rollabck error %v", err.Error())
						}
						close(jobTermiatedCh)
						break Loop
					}
				}
			case <-jobTermiatedCh:
				log.Debug("send goroutine terminated")
				if err := writer.Rollback(sendCrx); err != nil {
					log.Warnf("sender rollabck error %v", err.Error())
				}
				break Loop
			}
		}
	}()

	opWg.Wait()
	log.Debugf("%v complete", action)
}
