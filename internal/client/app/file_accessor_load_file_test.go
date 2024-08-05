package app_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
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

		testDecrypter := &testChunkDecrypter{}

		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		size := 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
				assert.Equal(t, fileInfo.Name, name)
				return testReader, nil
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
			assert.Equal(t, 0, len(testReader.bytes))
			assert.Equal(t, int32(1), testReader.nextInfok.Load())
			assert.Equal(t, int32(1), testReader.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))

			// testDecryptor
			assert.Equal(t, int32(1), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		size := 1024 * 4
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)

		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
			assert.Equal(t, 0, len(testReader.bytes))
			assert.Equal(t, int32(5), testReader.nextInfok.Load())
			assert.Equal(t, int32(1), testReader.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))

			// testDecryptor
			assert.Equal(t, int32(1), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		size := 1024*4 + 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)

		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
			assert.Equal(t, 0, len(testReader.bytes))
			assert.Equal(t, int32(5), testReader.nextInfok.Load())
			assert.Equal(t, int32(1), testReader.closeInfok.Load())

			// testSender
			assert.Equal(t, size, len(testSender.buf.Bytes()))
			assert.True(t, bytes.Equal(buf, testSender.buf.Bytes()))

			// testDecryptor
			assert.Equal(t, int32(1), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		testError := errors.New("testError")
		testReader := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrServerIsNotResponding)
				// testReader
				assert.Equal(t, int32(1), testReader.nextInfok.Load())
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		size := 1024*4 + 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrClientInternal)
				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(1), testSender.writeChunk.Load())

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		testError := errors.New("testError")
		testReader := &testFileStreamerErr{
			size: 512,
			err:  testError,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrServerIsNotResponding)
				// testReader
				assert.Equal(t, int32(1), testReader.nextInfok.Load())
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				assert.True(t, testSender.rollbackInvok.Load() > 0)
				assert.Equal(t, int32(0), testSender.commitInvok.Load())
				assert.Equal(t, int32(0), testSender.writeChunk.Load())

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		size := 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
			assert.Equal(t, fileInfo.Name, name)
			return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrClientInternal)

				// testReader
				assert.Equal(t, 0, len(testReader.bytes))
				assert.Equal(t, int32(1), testReader.nextInfok.Load())
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.Equal(t, int32(1), testSender.rollbackInvok.Load())
				assert.Equal(t, int32(1), testSender.commitInvok.Load())

				// testDecryptor
				assert.Equal(t, int32(1), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
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
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientAppStopped)

				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.True(t, testSender.rollbackInvok.Load() > 1)

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
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

		testDecrypter := &testChunkDecrypter{}
		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
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
			t.Error("load is not complete")
		case <-doneCh:
			// test errorChan is empty
			select {
			case err := <-errorChan:
				assert.ErrorIs(t, err, domain.ErrClientInteruptoin)

				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.True(t, testSender.rollbackInvok.Load() > 1)

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("decrypt_write_err", func(t *testing.T) {

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

		basename := filepath.Base(fileInfo.Path)
		testSender := &testFileInfiniteSender{
			exptectedName: basename,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		testErr := errors.New("error")
		testDecrypter := &testWriteErrChunkDecrypter{
			err: testErr,
		}

		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		chunkSize := 1024
		testReader := &testFileInfinitStreamer{
			chunkSize: chunkSize,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
				assert.Equal(t, fileInfo.Name, name)
				return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrClientInternal)

				// testReader
				assert.Equal(t, int32(1), testReader.closeInfok.Load())

				// testSender
				assert.True(t, testSender.rollbackInvok.Load() > 1)

				// testDecryptor
				assert.Equal(t, int32(0), testDecrypter.finishInv.Load())
			default:
				t.Error("errorChan not empty")
			}
		}
	})

	t.Run("decrypt_finish_err", func(t *testing.T) {

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

		basename := filepath.Base(fileInfo.Path)
		testSender := &testFileInfiniteSender{
			exptectedName: basename,
		}
		dir := filepath.Dir(fileInfo.Path)
		mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
			assert.Equal(t, dir, name)
			return testSender, nil

		}).Times(1)

		testErr := errors.New("error")
		testDecrypter := &testFinishErrChunkDecrypter{
			err: testErr,
		}

		mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
			assert.Equal(t, "masterPass", pass)
			return testDecrypter
		})

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return("masterPass").Times(1)

		size := 512
		chunkSize := 1024
		buf := make([]byte, size)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		testReader := &testFileStreamer{
			chunkSize: chunkSize,
			size:      size,
			bytes:     buf,
		}

		mockServer := NewMockAppServer(ctrl)

		mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
				assert.Equal(t, fileInfo.Name, name)
				return testReader, nil
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
				assert.ErrorIs(t, err, domain.ErrClientDataIsNotRestored)

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

type testChunkDecrypter struct {
	finishInv atomic.Int32
}

func (tcd *testChunkDecrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (tcd *testChunkDecrypter) Finish() error {
	tcd.finishInv.Add(1)
	return nil
}

var _ domain.ChunkDecrypter = (*testChunkDecrypter)(nil)

type testWriteErrChunkDecrypter struct {
	err       error
	finishInv atomic.Int32
}

func (tcd *testWriteErrChunkDecrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, tcd.err
}

func (tcd *testWriteErrChunkDecrypter) Finish() error {
	tcd.finishInv.Add(1)
	return nil
}

var _ domain.ChunkDecrypter = (*testWriteErrChunkDecrypter)(nil)

type testFinishErrChunkDecrypter struct {
	err error
}

func (tcd *testFinishErrChunkDecrypter) WriteChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (tcd *testFinishErrChunkDecrypter) Finish() error {
	return tcd.err
}

var _ domain.ChunkDecrypter = (*testWriteErrChunkDecrypter)(nil)
