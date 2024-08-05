package app_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileAccesprUpload(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("upload_succes_short_file", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

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

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
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
			assert.Equal(t, int32(1), testSender.commitInvok.Load())
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))

			assert.Equal(t, int32(1), testEncrypter.finishInv.Load())
		}
	})

	t.Run("upload_succes_big_file", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

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

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
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

			// testEncrypter
			assert.Equal(t, int32(1), testEncrypter.finishInv.Load())
		}
	})

	t.Run("upload_succes_big_file_with_tail", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

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

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
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

			// testEncrypter
			assert.Equal(t, int32(1), testEncrypter.finishInv.Load())
		}

	})

	t.Run("upload_read_error_wrong_file_size", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		testReder := &testFileStreamerErr{
			size: -1,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInternal)
				// testReader
				assert.Equal(t, int32(0), testReder.nextInfok.Load())
				assert.Equal(t, int32(0), testReder.closeInfok.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("upload_read_error", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		testError := errors.New("testError")
		testReder := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInternal)
				// testReader
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("upload_write_error", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

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

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testError := errors.New("testErr")
		testSender := &testFileErrSender{
			exptectedName: fileInfo.Name,
			err:           testError,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrServerIsNotResponding)
				// testReader
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(1), testSender.writeChunk.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("upload_read_error_rollabck_err", func(t *testing.T) {
		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		testError := errors.New("testError")
		testReder := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		rollErr := errors.New("roll_err")
		testSender := &testFileRoolbackErrSender{
			exptectedName: fileInfo.Name,
			err:           rollErr,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInternal)
				// testReader
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})

	t.Run("upload_succes_commit_err", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
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

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testError := errors.New("commit err")
		testSender := &testFileCommitErrSender{
			exptectedName: fileInfo.Name,
			err:           testError,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrServerIsNotResponding)

				// testReader
				assert.Equal(t, 0, len(testReder.bytes))
				assert.Equal(t, int32(1), testReder.nextInfok.Load())
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(0), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(1), testSender.commitInvok.Load())

				// testEncrypter
				assert.Equal(t, int32(1), testEncrypter.finishInv.Load())
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
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		chunkSize := 1024
		testReder := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileInfiniteSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
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
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
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
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		chunkSize := 1024
		testReder := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testEncrypter := &testChunkEncrypter{}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileInfiniteSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		time.Sleep(2 * time.Second)
		cancelChan <- struct{}{}

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInteruptoin)

				// testReader
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("test_app_encrypt_write_chunck_err", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

		chunkSize := 1024
		testReder := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testErr := errors.New("testErr")
		testEncrypter := &testWriteChunkErrEncrypter{
			err: testErr,
		}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileInfiniteSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInternal)

				// testReader
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())

				// testEncrypter
				assert.Equal(t, int32(0), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("test_app_encrypt_finish_err", func(t *testing.T) {

		fileInfo := &domain.FileInfo{
			Name: "autotest",
			Path: "./build/atuotest",
		}

		mockHeler := NewMockDomainHelper(ctrl)
		mockHeler.EXPECT().CheckFileForRead(gomock.Any()).Return(nil).Times(1)

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
		mockHeler.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, fInf.Name)
			assert.Equal(t, fileInfo.Path, fInf.Path)
			return testReder, nil
		})

		testErr := errors.New("testErr")
		testEncrypter := &testFinishErrEncrypter{
			err: testErr,
		}
		mockHeler.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
			assert.Equal(t, "masterPass", pass)
			return testEncrypter, nil
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		mockServer := NewMockAppServer(ctrl)
		testSender := &testFileInfiniteSender{
			exptectedName: fileInfo.Name,
		}
		mockServer.EXPECT().CreateFileSender(gomock.Any()).Return(testSender, nil).Times(1)

		fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

		ctx, doneFn := context.WithTimeout(context.Background(), 3*time.Second)
		defer doneFn()

		cancelChan := make(chan struct{}, 1)
		errorChan := make(chan error, 1)

		doneCh := make(chan struct{}, 1)
		go func() {
			defer func() { doneCh <- struct{}{} }()
			fa.UploadFile(ctx, fileInfo, nil, cancelChan, errorChan)
		}()

		select {
		case <-ctx.Done():
			t.Error("upload is not complete")
		case <-doneCh:
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInternal)

				// testReader
				assert.Equal(t, int32(1), testReder.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())

				// testEncrypter
				assert.Equal(t, int32(1), testEncrypter.finishInv.Load())
			default:
				t.Error("errorChan is empty")
			}
		}
	})
}

type testFileSender struct {
	buf           bytes.Buffer
	commitInvok   atomic.Int32
	rollbackInvok atomic.Int32
	writeChunk    atomic.Int32
	exptectedName string
}

var _ domain.StreamFileWriter = (*testFileSender)(nil)

func (tf *testFileSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	tf.writeChunk.Add(1)
	if tf.exptectedName != name {
		return fmt.Errorf("unexpected name")
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk i nil")
	}
	_, err := tf.buf.Write(chunk)
	return err

}
func (tf *testFileSender) Commit(ctx context.Context) error {
	tf.commitInvok.Add(1)
	return nil
}
func (tf *testFileSender) Rollback(ctx context.Context) error {
	tf.rollbackInvok.Add(1)
	return nil
}

type testFileInfiniteSender struct {
	commitInvok   atomic.Int32
	rollbackInvok atomic.Int32
	writeChunk    atomic.Int32
	exptectedName string
}

var _ domain.StreamFileWriter = (*testFileInfiniteSender)(nil)

func (tf *testFileInfiniteSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	tf.writeChunk.Add(1)
	if tf.exptectedName != name {
		return fmt.Errorf("unexpected name")
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk i nil")
	}
	return nil

}
func (tf *testFileInfiniteSender) Commit(ctx context.Context) error {
	tf.commitInvok.Add(1)
	return nil
}
func (tf *testFileInfiniteSender) Rollback(ctx context.Context) error {
	tf.rollbackInvok.Add(1)
	return nil
}

type testFileErrSender struct {
	err           error
	commitInvok   atomic.Int32
	rollbackInvok atomic.Int32
	writeChunk    atomic.Int32
	exptectedName string
}

var _ domain.StreamFileWriter = (*testFileErrSender)(nil)

func (tf *testFileErrSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	tf.writeChunk.Add(1)
	if tf.exptectedName != name {
		return fmt.Errorf("unexpected name")
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk i nil")
	}
	return tf.err

}
func (tf *testFileErrSender) Commit(ctx context.Context) error {
	tf.commitInvok.Add(1)
	return nil
}
func (tf *testFileErrSender) Rollback(ctx context.Context) error {
	tf.rollbackInvok.Add(1)
	return nil
}

type testFileRoolbackErrSender struct {
	err           error
	commitInvok   atomic.Int32
	rollbackInvok atomic.Int32
	writeChunk    atomic.Int32
	exptectedName string
}

var _ domain.StreamFileWriter = (*testFileRoolbackErrSender)(nil)

func (tf *testFileRoolbackErrSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	tf.writeChunk.Add(1)
	if tf.exptectedName != name {
		return fmt.Errorf("unexpected name")
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk i nil")
	}
	return nil

}
func (tf *testFileRoolbackErrSender) Commit(ctx context.Context) error {
	tf.commitInvok.Add(1)
	return nil
}
func (tf *testFileRoolbackErrSender) Rollback(ctx context.Context) error {
	tf.rollbackInvok.Add(1)
	return tf.err
}

type testFileCommitErrSender struct {
	err           error
	commitInvok   atomic.Int32
	rollbackInvok atomic.Int32
	writeChunk    atomic.Int32
	exptectedName string
}

var _ domain.StreamFileWriter = (*testFileCommitErrSender)(nil)

func (tf *testFileCommitErrSender) WriteChunk(ctx context.Context, name string, chunk []byte) error {
	tf.writeChunk.Add(1)
	if tf.exptectedName != name {
		return fmt.Errorf("unexpected name")
	}
	if len(chunk) == 0 {
		return fmt.Errorf("chunk i nil")
	}
	return nil

}
func (tf *testFileCommitErrSender) Commit(ctx context.Context) error {
	tf.commitInvok.Add(1)
	return tf.err
}
func (tf *testFileCommitErrSender) Rollback(ctx context.Context) error {
	tf.rollbackInvok.Add(1)
	return nil
}

type testFileStreamer struct {
	bytes         []byte
	size          int
	chunkSize     int
	fileSizeInvok atomic.Int32
	nextInfok     atomic.Int32
	closeInfok    atomic.Int32
}

var _ domain.StreamFileReader = (*testFileStreamer)(nil)

func (tf *testFileStreamer) FileSize() int64 {
	tf.fileSizeInvok.Add(1)
	return int64(tf.size)
}
func (tf *testFileStreamer) Next() ([]byte, error) {
	tf.nextInfok.Add(1)
	ln := len(tf.bytes)
	var ret []byte
	if ln < tf.chunkSize {

		ret, tf.bytes = tf.bytes[:ln], nil
		return ret, io.EOF
	} else {
		ret, tf.bytes = tf.bytes[:tf.chunkSize], tf.bytes[tf.chunkSize:]
		return ret, nil
	}
}

func (tf *testFileStreamer) Close() {
	tf.closeInfok.Add(1)
}

type testFileInfinitStreamer struct {
	chunkSize     int
	fileSizeInvok atomic.Int32
	nextInfok     atomic.Int32
	closeInfok    atomic.Int32
}

var _ domain.StreamFileReader = (*testFileInfinitStreamer)(nil)

func (tf *testFileInfinitStreamer) FileSize() int64 {
	tf.fileSizeInvok.Add(1)
	return int64(20 * domain.GiB)
}
func (tf *testFileInfinitStreamer) Next() ([]byte, error) {
	tf.nextInfok.Add(1)

	res := make([]byte, tf.chunkSize)
	if _, err := rand.Read(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (tf *testFileInfinitStreamer) Close() {
	tf.closeInfok.Add(1)
}

type testFileStreamerErr struct {
	err           error
	size          int
	fileSizeInvok atomic.Int32
	nextInfok     atomic.Int32
	closeInfok    atomic.Int32
}

var _ domain.StreamFileReader = (*testFileStreamerErr)(nil)

func (tf *testFileStreamerErr) FileSize() int64 {
	tf.fileSizeInvok.Add(1)
	return int64(tf.size)
}
func (tf *testFileStreamerErr) Next() ([]byte, error) {
	tf.nextInfok.Add(1)
	return nil, tf.err
}

func (tf *testFileStreamerErr) Close() {
	tf.closeInfok.Add(1)
}

type testChunkEncrypter struct {
	finishInv atomic.Int32
}

func (tcd *testChunkEncrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (tcd *testChunkEncrypter) Finish() ([]byte, error) {
	if tcd.finishInv.CompareAndSwap(0, 1) {
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected finish calls")
}

var _ domain.ChunkEncrypter = (*testChunkEncrypter)(nil)

type testWriteChunkErrEncrypter struct {
	finishInv atomic.Int32
	err       error
}

func (tcd *testWriteChunkErrEncrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, tcd.err
}

func (tcd *testWriteChunkErrEncrypter) Finish() ([]byte, error) {
	tcd.finishInv.Add(1)
	return nil, nil
}

var _ domain.ChunkEncrypter = (*testWriteChunkErrEncrypter)(nil)

type testFinishErrEncrypter struct {
	finishInv atomic.Int32
	err       error
}

func (tcd *testFinishErrEncrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (tcd *testFinishErrEncrypter) Finish() ([]byte, error) {
	tcd.finishInv.Add(1)
	return nil, tcd.err
}

var _ domain.ChunkEncrypter = (*testFinishErrEncrypter)(nil)
