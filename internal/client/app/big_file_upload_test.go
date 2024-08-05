package app_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestDataDirectory = "../../../testdata"

func TestFileUploadingNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// encrypt big file

	testFileName := "bigfile.mp4"
	testFilePath := filepath.Join(TestDataDirectory, testFileName)

	// create temp file name
	f, err := os.CreateTemp("", "big-file-encrypted-")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	require.NoError(t, os.Remove(f.Name()))

	tempFileDirectory := filepath.Dir(f.Name())

	tempFileBasename := filepath.Base(f.Name())

	readFileInfo := &domain.FileInfo{
		Name: tempFileBasename,
		Path: testFilePath,
	}

	_ = os.Remove(f.Name()) // remove if exists
	defer func() {
		_ = os.Remove(f.Name()) // after test
	}()

	masterKey := "masterKey"

	uploadMockHelper := NewMockDomainHelper(ctrl)

	uploadMockHelper.EXPECT().CreateStreamFileReader(gomock.Any()).DoAndReturn(func(fInf *domain.FileInfo) (domain.StreamFileReader, error) {
		assert.Equal(t, readFileInfo.Name, fInf.Name)
		assert.Equal(t, readFileInfo.Path, fInf.Path)
		return domain.CreateStreamFileReader(fInf.Path)
	}).Times(1)

	uploadMockHelper.EXPECT().CheckFileForRead(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
		return domain.CheckFileForRead(inf)
	}).Times(1)

	uploadMockHelper.EXPECT().CreateChunkEncrypter(gomock.Any()).DoAndReturn(func(pass string) (domain.ChunkEncrypter, error) {
		assert.Equal(t, masterKey, pass)
		return domain.NewChunkEncrypter(pass)
	}).Times(1)

	uploadMockStorage := NewMockAppStorage(ctrl)
	uploadMockStorage.EXPECT().IsFileInfoExists(gomock.Any()).Return(false).Times(1)
	uploadMockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

	uploadMockServer := NewMockAppServer(ctrl)
	uploadMockServer.EXPECT().CreateFileSender(gomock.Any()).DoAndReturn(func(ctx context.Context) (domain.StreamFileWriter, error) {
		return domain.CreateStreamFileWriter(tempFileDirectory)
	}).Times(1)

	baseTestFileName := filepath.Base(f.Name())

	defer func() {
		_ = os.Remove(filepath.Join(tempFileDirectory, domain.TempFileNamePrefix+baseTestFileName))
	}()

	uploader := app.NewFileAccessor().DomainHelper(uploadMockHelper).AppStorage(uploadMockStorage).AppServer(uploadMockServer)

	ctx, doneFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer doneFn()

	cancelChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)

	doneCh := make(chan struct{}, 1)
	go func() {
		defer func() { doneCh <- struct{}{} }()
		uploader.UploadFile(ctx, readFileInfo, nil, cancelChan, errorChan)
	}()

	select {
	case <-ctx.Done():
		t.Error("upload is not complete")
	case <-doneCh:
		// test errorChan is empty
		select {
		case err := <-errorChan:
			require.NoError(t, err)
		default:
		}
	}

	// test upload function
	require.NoError(t, checkUploadedFile(f.Name(), masterKey))

	t.Log("encryption success")
	fmt.Println("-------------------------")

	// read comlete
	decryptFileInfo := &domain.FileInfo{
		Name: testFileName,
		Path: f.Name(),
	}

	mockHeler := NewMockDomainHelper(ctrl)
	mockHeler.EXPECT().CheckFileForWrite(gomock.Any()).DoAndReturn(func(inf *domain.FileInfo) error {
		assert.Equal(t, decryptFileInfo.Name, inf.Name)
		assert.Equal(t, decryptFileInfo.Path, inf.Path)
		return nil
	}).Times(1)

	mockHeler.EXPECT().CreateStreamFileWriter(gomock.Any()).DoAndReturn(func(name string) (domain.StreamFileWriter, error) {
		assert.Equal(t, tempFileDirectory, name)
		return &testFileInfiniteSender{
			exptectedName: baseTestFileName,
		}, nil

	}).Times(1)

	mockStorage := NewMockAppStorage(ctrl)
	mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)
	mockHeler.EXPECT().CreateChunkDecrypter(gomock.Any()).DoAndReturn(func(pass string) domain.ChunkDecrypter {
		assert.Equal(t, masterKey, pass)
		return domain.NewChunkDecrypter(masterKey)
	})

	mockServer := NewMockAppServer(ctrl)

	mockServer.EXPECT().CreateFileReceiver(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string) (domain.StreamFileReader, error) {
		assert.Equal(t, decryptFileInfo.Name, name)
		return domain.CreateStreamFileReader(f.Name())
	}).Times(1)

	fa := app.NewFileAccessor().DomainHelper(mockHeler).AppStorage(mockStorage).AppServer(mockServer)

	ctx, doneFn = context.WithTimeout(context.Background(), 10*time.Second)
	defer doneFn()

	cancelChan = make(chan struct{}, 1)
	errorChan = make(chan error, 1)

	doneCh = make(chan struct{}, 1)
	go func() {
		defer func() { doneCh <- struct{}{} }()
		fa.LoadFile(ctx, decryptFileInfo, nil, cancelChan, errorChan)
	}()

	select {
	case <-ctx.Done():
		t.Error("load is not complete")
	case <-doneCh:
		// test errorChan is empty
		select {
		case err := <-errorChan:
			require.NoError(t, err)
		default:
		}
	}

}

func checkUploadedFile(filePath string, masterKey string) error {
	decryptor := domain.NewChunkDecrypter(masterKey)
	chunkReader, err := domain.NewStreamFileReader(filePath, 4*domain.KiB)
	if err != nil {
		return err
	}

	size := 0
	for {
		chunk, err := chunkReader.Next()
		if err == nil {
			size += len(chunk)
			_, err := decryptor.WriteChunk(chunk)
			if err != nil {
				return err
			}

		} else {
			if err == io.EOF {
				if len(chunk) != 0 {
					return fmt.Errorf("unexpected chunk len")
				}

				chunkReader.Close()
				if err = decryptor.Finish(); err != nil {
					return err
				}
				break
			}
			return fmt.Errorf("unexpected")
		}
	}
	return nil
}
