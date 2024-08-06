package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserPasswordDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		passDataList := []domain.EncryptedUserPasswordData{
			{
				Hint:    "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetUserPasswordDataList(gomock.Any()).Return(passDataList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, passDataList[0].Content, cnt)
			return `{"hint":"100", "login":"login","password":"pass"}`, nil
		}).Times(1)

		mockStorage.EXPECT().SetUserPasswordDatas(gomock.Any()).Do(func(crds []domain.UserPasswordData) {
			require.Equal(t, 1, len(crds))
			crd := crds[0]
			assert.Equal(t, "100", crd.Hint)
			assert.Equal(t, "login", crd.Login)
			assert.Equal(t, "pass", crd.Passwrod)
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetUserPasswordDataList(context.Background())
		require.NoError(t, err)
	})

	t.Run("get_user_pass_login_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		testErr := errors.New("testErr")
		mockServer.EXPECT().GetUserPasswordDataList(gomock.Any()).Return(nil, testErr).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage)
		err := da.GetUserPasswordDataList(context.Background())
		require.ErrorIs(t, err, testErr)
	})

	t.Run("decrypt_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)
		passDataList := []domain.EncryptedUserPasswordData{
			{
				Hint:    "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetUserPasswordDataList(gomock.Any()).Return(passDataList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		testErr := errors.New("testErr")
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, passDataList[0].Content, cnt)
			return `{"hint":"100", "login":"login","password":"pass"}`, testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetUserPasswordDataList(context.Background())
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})
	t.Run("json_err", func(t *testing.T) {

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		mockServer := NewMockAppServer(ctrl)

		passDataList := []domain.EncryptedUserPasswordData{
			{
				Hint:    "100",
				Content: "content",
			},
		}
		mockServer.EXPECT().GetUserPasswordDataList(gomock.Any()).Return(passDataList, nil).Times(1)

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, cnt string) (string, error) {
			assert.Equal(t, masterKey, key)
			assert.Equal(t, passDataList[0].Content, cnt)
			return `-`, nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.GetUserPasswordDataList(context.Background())
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

}

func TestAddUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().CreateUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EncryptedUserPasswordData) error {
			require.NotNil(t, dt)
			assert.Equal(t, passData.Hint, dt.Hint)
			assert.Equal(t, encryptedCnt, dt.Content)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddUserPasswordData(context.Background(), passData)
		require.NoError(t, err)
	})

	t.Run("check_card_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}
		testErr := errors.New("testErr")
		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().DomainHelper(mockHelper)
		err := da.AddUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("emcrypt_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return "", testErr
		}).Times(1)

		da := app.NewDataAccessor().AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, domain.ErrClientInternal)
	})

	t.Run("server_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().CreateUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EncryptedUserPasswordData) error {
			require.NotNil(t, dt)
			assert.Equal(t, passData.Hint, dt.Hint)
			assert.Equal(t, encryptedCnt, dt.Content)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.AddUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, testErr)
	})

}

func TestUpdateUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().UpdateUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EncryptedUserPasswordData) error {
			require.NotNil(t, dt)
			assert.Equal(t, passData.Hint, dt.Hint)
			assert.Equal(t, encryptedCnt, dt.Content)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.UpdateUserPasswordData(context.Background(), passData)
		require.NoError(t, err)
	})

	t.Run("check_card_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}
		testErr := errors.New("testErr")
		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().DomainHelper(mockHelper)
		err := da.UpdateUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("emcrypt_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return "", testErr
		}).Times(1)

		da := app.NewDataAccessor().AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.UpdateUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, domain.ErrClientInternal)
	})

	t.Run("server_err", func(t *testing.T) {

		mockHelper := NewMockDomainHelper(ctrl)

		passData := &domain.UserPasswordData{
			Hint:     "Hint",
			Login:    "Login",
			Passwrod: "Passwod",
		}

		mockHelper.EXPECT().CheckUserPasswordData(gomock.Any()).DoAndReturn(func(crd *domain.UserPasswordData) error {
			assert.Same(t, passData, crd)
			return nil
		}).Times(1)

		masterKey := "masterKey"
		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().GetMasterPassword().Return(masterKey).Times(1)

		encryptedCnt := "encrypted"
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(key string, content string) (string, error) {

			assert.Equal(t, masterKey, key)
			cnt, err := json.Marshal(passData)
			require.NoError(t, err)
			assert.JSONEq(t, string(cnt), content)
			return encryptedCnt, nil
		}).Times(1)

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().UpdateUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EncryptedUserPasswordData) error {
			require.NotNil(t, dt)
			assert.Equal(t, passData.Hint, dt.Hint)
			assert.Equal(t, encryptedCnt, dt.Content)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer).AppStorage(mockStorage).DomainHelper(mockHelper)
		err := da.UpdateUserPasswordData(context.Background(), passData)
		require.ErrorIs(t, err, testErr)
	})

}

func TestDeleteUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		hint := "yandex"

		mockServer := NewMockAppServer(ctrl)
		mockServer.EXPECT().DeleteUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, hnt string) error {
			assert.Equal(t, hint, hnt)
			return nil
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer)
		err := da.DeleteUserPasswordData(context.Background(), hint)
		require.NoError(t, err)
	})

	t.Run("server_err", func(t *testing.T) {

		hint := "yandex"

		mockServer := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockServer.EXPECT().DeleteUserPasswordData(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, hnt string) error {
			assert.Equal(t, hint, hnt)
			return testErr
		}).Times(1)

		da := app.NewDataAccessor().AppSever(mockServer)
		err := da.DeleteUserPasswordData(context.Background(), hint)
		require.ErrorIs(t, err, testErr)
	})
}
