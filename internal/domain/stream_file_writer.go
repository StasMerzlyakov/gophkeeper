package domain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const TempFileNamePrefix = "_"

func NewStreamFileWriter(folderPath string) (*streamFileWriter, error) {

	return &streamFileWriter{
		folderPath: folderPath,
	}, nil
}

var _ StreamFileWriter = (*streamFileWriter)(nil)

type streamFileWriter struct {
	folderPath   string
	fileName     string
	once         sync.Once
	tempFilePath string
	file         *os.File
}

func (sw *streamFileWriter) WriteChunk(ctx context.Context, name string, chunk []byte) error {

	var onceErr error

	sw.once.Do(func() {
		sw.fileName = name

		if err := os.MkdirAll(sw.folderPath, os.ModePerm); err != nil {
			onceErr = fmt.Errorf("%w can't create bucket dir %s", ErrServerInternal, sw.folderPath)
			return
		}

		fullPath := filepath.Join(sw.folderPath, TempFileNamePrefix+name)

		if _, err := os.Stat(fullPath); err == nil {
			// file exists
			onceErr = fmt.Errorf("%w file already exists", ErrClientDataIncorrect)
			return
		} else if !errors.Is(err, os.ErrNotExist) {
			onceErr = fmt.Errorf("%w test filePaht %s error", ErrServerInternal, fullPath)
			return
		}

		f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			onceErr = fmt.Errorf("%w open file %s error %s", ErrServerInternal, fullPath, err.Error())
			return
		}
		sw.file = f
		sw.tempFilePath = fullPath
	})
	if onceErr != nil {
		return onceErr
	}

	if _, err := sw.file.Write(chunk); err != nil {
		return fmt.Errorf("%w write chunk of %s err %s", ErrServerInternal, name, err.Error())
	}
	return nil
}

func (sw *streamFileWriter) Commit(ctx context.Context) error {

	if err := sw.file.Close(); err != nil {
		return fmt.Errorf("%w close file %s err %s", ErrServerInternal, sw.file.Name(), err.Error())
	}

	fileDir := filepath.Dir(sw.tempFilePath)
	rightFilePath := filepath.Join(fileDir, sw.fileName)

	if err := os.Rename(sw.tempFilePath, rightFilePath); err != nil {
		return fmt.Errorf("%w rename file %s err %s", ErrServerInternal, sw.file.Name(), err.Error())
	}

	return nil
}

func (sw *streamFileWriter) Rollback(ctx context.Context) error {

	if err := sw.file.Close(); err != nil {
		return fmt.Errorf("%w close file %s err %s", ErrServerInternal, sw.file.Name(), err.Error())
	}

	if err := os.RemoveAll(sw.tempFilePath); err != nil {
		return fmt.Errorf("%w remove temp file %s err %s", ErrServerInternal, sw.file.Name(), err.Error())
	}

	return nil
}
