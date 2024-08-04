package app_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileAccesprLoad(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("load_succes_short_file", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)
		testSender := &testFileSender{
			exptectedName: name,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		size := 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReder := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case <-errorChan:
				t.Error("errorChan is not empty")
			default:
			}

			// testReader
			assert.Equal(t, 0, len(testReder.bytes))
			assert.Equal(t, int32(1), testReder.nextInfok.Load())
			assert.Equal(t, int32(1), testReder.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))
		}
	})

	t.Run("load_succes_big_file", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)
		testSender := &testFileSender{
			exptectedName: name,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		size := 1024 * 4
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)

		testReder := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case <-errorChan:
				t.Error("errorChan is not empty")
			default:
			}

			// testReader
			assert.Equal(t, 0, len(testReder.bytes))
			assert.Equal(t, int32(5), testReder.nextInfok.Load())
			assert.Equal(t, int32(1), testReder.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))
		}
	})

	t.Run("load_succes_big_file_with_tail", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)
		testSender := &testFileSender{
			exptectedName: name,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		size := 1024*4 + 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)

		testReder := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case <-errorChan:
				t.Error("errorChan is not empty")
			default:
			}

			// testReader
			assert.Equal(t, 0, len(testReder.bytes))
			assert.Equal(t, int32(5), testReder.nextInfok.Load())
			assert.Equal(t, int32(1), testReder.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))
		}
	})

	t.Run("load_read_error", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)
		testSender := &testFileSender{
			exptectedName: name,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		testError := errors.New("testError")
		testReder := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, testError)
				// testReader
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("load_write_error", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)
		testError := errors.New("testErr")
		testSender := &testFileErrSender{
			exptectedName: name,
			err:           testError,
		}
		dir := filepath.Dir(fileInfo.Path)

		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		size := 1024*4 + 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReder := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, testError)
				// testReader
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(1), testSender.writeChunk.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("load_read_error_rollabck_err", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)

		rollErr := errors.New("roll_err")
		testSender := &testFileRoolbackErrSender{
			exptectedName: name,
			err:           rollErr,
		}

		dir := filepath.Dir(fileInfo.Path)

		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		testError := errors.New("testError")
		testReder := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, testError)
				// testReader
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("load_succes_commit_err", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		name := filepath.Base(fileInfo.Path)

		testError := errors.New("commit err")
		testSender := &testFileCommitErrSender{
			exptectedName: name,
			err:           testError,
		}
		dir := filepath.Dir(fileInfo.Path)

		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		size := 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReder := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReder, nil
		})

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("load is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, testError)

				// testReader
				assert.Equal(t, 0, len(testReder.bytes))
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(1), testSender.commitInvok.Load())
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("test_app_stop", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		dir := filepath.Dir(fileInfo.Path)
		basename := filepath.Base(fileInfo.Path)
		testSender := &testFileInfiniteSender{
			exptectedName: basename,
		}
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		mockServer := NewMockAppServer(ctrl)

		chunkSize := 1024
		testReader := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
		}).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		time.Sleep(2 * time.Second)
		fa.Stop(ctx)

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientAppStopped)

				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.True(t, testSender.rollbackInvok.Load() > 1)
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("test_app_canceled", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
			assert.Equal(t, fileInfo.Name, inf.Name)
			assert.Equal(t, fileInfo.Path, inf.Path)
			return nil
		}).Times(1)

		dir := filepath.Dir(fileInfo.Path)
		basename := filepath.Base(fileInfo.Path)
		testSender := &testFileInfiniteSender{
			exptectedName: basename,
		}
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return &testChunkDecryptor{}
		})

		mockServer := NewMockAppServer(ctrl)

		chunkSize := 1024
		testReader := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
		}).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.LoadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		time.Sleep(2 * time.Second)
		cancelChan <- struct{}{}

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInteruptoin)

				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.True(t, testSender.rollbackInvok.Load() > 1)
			default:
				t.Error("errorChan not empty")
			}
		}
	})

}

type testChunkDecryptor struct {
	finishInv atomic.Int32
}

func (tcd *testChunkDecryptor) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (tcd *testChunkDecryptor) Finish() error {
	if tcd.finishInv.CompareAndSwap(0, 1) {
		return nil
	}
	return fmt.Errorf("unexpected finish calls")
}

var _ domain.ChunkDecrypter = (*testChunkDecryptor)(nil)
