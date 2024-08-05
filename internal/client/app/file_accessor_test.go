package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileInfoList(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)
		mockStorage := NewMockAppStorage(ctrl)

		lst := []domain.FileInfo{}
		mockServer.EXPECT().GetFileInfoList(gomock.Any()).Return(lst, nil).Times(1)
		mockStorage.EXPECT().SetFilesInfo(gomock.Any()).Do(func(l []domain.FileInfo) {
			assert.Equal(t, lst, l)
		}).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer).AppStorage(mockStorage)
		err := fa.GetFileInfoList(context.Background())
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)

		lst := []domain.FileInfo{}
		testErr := errors.New("testErr")
		mockServer.EXPECT().GetFileInfoList(gomock.Any()).Return(lst, testErr).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer)
		err := fa.GetFileInfoList(context.Background())
		require.ErrorIs(t, err, testErr)
	})

}

func TestDeleteFile(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)

		mockStorage := NewMockAppStorage(ctrl)

		name := "name"
		mockServer.EXPECT().DeleteFileInfo(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, nm string) error {
			assert.Equal(t, name, nm)
			return nil
		}).Times(1)

		mockStorage.EXPECT().DeleteFileInfo(gomock.Any()).DoAndReturn(func(nm string) error {
			assert.Equal(t, name, nm)
			return nil
		}).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer).AppStorage(mockStorage)
		err := fa.DeleteFile(context.Background(), name)
		require.NoError(t, err)
	})

	t.Run("srv_err", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)

		name := "name"
		testErr := errors.New("testErr")
		mockServer.EXPECT().DeleteFileInfo(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, nm string) error {
			assert.Equal(t, name, nm)
			return testErr
		}).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer)
		err := fa.DeleteFile(context.Background(), name)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("internal_err", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)

		mockStorage := NewMockAppStorage(ctrl)

		name := "name"
		mockServer.EXPECT().DeleteFileInfo(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, nm string) error {
			assert.Equal(t, name, nm)
			return nil
		}).Times(1)

		testErr := errors.New("testErr")
		mockStorage.EXPECT().DeleteFileInfo(gomock.Any()).DoAndReturn(func(nm string) error {
			assert.Equal(t, name, nm)
			return testErr
		}).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer).AppStorage(mockStorage)
		err := fa.DeleteFile(context.Background(), name)
		require.ErrorIs(t, err, domain.ErrClientInternal)
	})

	/*t.Run("err", func(t *testing.T) {

		mockServer := NewMockAppServer(ctrl)

		lst := []domain.FileInfo{}
		testErr := errors.New("testErr")
		mockServer.EXPECT().GetFileInfoList(gomock.Any()).Return(lst, testErr).Times(1)

		fa := app.NewFileAccessor().AppServer(mockServer)
		err := fa.GetFileInfoList(context.Background())
		require.ErrorIs(t, err, testErr)
	}) */

}
