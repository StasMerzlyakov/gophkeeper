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

func TestLoginer_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		mockVier.EXPECT().ShowLogOTPView().Times(1)

		loginData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}
		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			assert.Equal(t, loginData.EMail, data.EMail)
			assert.Equal(t, loginData.Password, data.Password)
			return nil
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier)
		loginer.Login(context.Background(), loginData)
	})

	t.Run("err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)

		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		loginData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}
		mockSrv := NewMockLoginServer(ctrl)

		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			assert.Equal(t, loginData.EMail, data.EMail)
			assert.Equal(t, loginData.Password, data.Password)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier)
		loginer.Login(context.Background(), loginData)
	})
}

func TestLoginer_PassOTP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		mockVier.EXPECT().ShowMasterKeyView(gomock.Any()).Do(func(msg string) {
			assert.Equal(t, msg, "")
		}).Times(1)

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return nil
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier)
		loginer.PassOTP(context.Background(), otpPass)
	})

	t.Run("err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)

		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}
		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier)
		loginer.PassOTP(context.Background(), otpPass)
	})
}

func TestLoginer_CheckMasterKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		mockVier.EXPECT().ShowDataAccessView().Times(1)

		masterPassword := "masterPassword"

		helloData := &domain.HelloData{
			HelloEncrypted:     "helloEncrypted",
			MasterPasswordHint: "masterPasswordHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)

		mockHlp.EXPECT().DecryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) error {
			assert.Equal(t, masterPassword, secretKey)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return nil
		}).Times(1)

		mockStorage := NewMockAppStorage(ctrl)
		mockStorage.EXPECT().SetMasterPassword(gomock.Any()).Do(func(mKey string) {
			require.Equal(t, masterPassword, mKey)
		})

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp).LoginStorage(mockStorage)
		loginer.CheckMasterKey(context.Background(), masterPassword)
	})

	t.Run("getHelloData_err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		masterKeyPassword := "masterKeyPassword"

		mockSrv := NewMockLoginServer(ctrl)

		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(nil, testErr).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})

	t.Run("decryptMasterKey_err", func(t *testing.T) {

		mockVier := NewMockLoginView(ctrl)
		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		masterPassword := "masterPassword"

		helloData := &domain.HelloData{
			HelloEncrypted:     "helloEncrypted",
			MasterPasswordHint: "masterPasswordHint",
		}

		mockVier.EXPECT().ShowMasterKeyView(gomock.Any()).Do(func(hint string) {
			assert.Equal(t, helloData.MasterPasswordHint, hint)
		}).Times(1)

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)

		mockHlp.EXPECT().DecryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) error {
			assert.Equal(t, masterPassword, secretKey)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterPassword)
	})
}
