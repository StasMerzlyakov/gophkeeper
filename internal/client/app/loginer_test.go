package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginer_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		loginData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}
		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			assert.Equal(t, loginData.EMail, data.EMail)
			assert.Equal(t, loginData.Password, data.Password)
			return nil
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv)
		err := loginer.Login(context.Background(), loginData)
		assert.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		testErr := errors.New("testErr")
		loginData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}
		mockSrv := NewMockAppServer(ctrl)

		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			assert.Equal(t, loginData.EMail, data.EMail)
			assert.Equal(t, loginData.Password, data.Password)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv)
		err := loginer.Login(context.Background(), loginData)
		assert.ErrorIs(t, err, testErr)
	})
}

func TestLoginer_PassOTP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return nil
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv)
		err := loginer.PassOTP(context.Background(), otpPass)
		assert.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		testErr := errors.New("testErr")

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}
		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv)
		err := loginer.PassOTP(context.Background(), otpPass)
		assert.ErrorIs(t, err, testErr)
	})

}

func TestLoginer_CheckMasterKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		masterPassword := "masterPassword"

		helloData := &domain.HelloData{
			HelloEncrypted:     "helloEncrypted",
			MasterPasswordHint: "masterPasswordHint",
		}

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockDomainHelper(ctrl)

		mockHlp.EXPECT().DecryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) error {
			assert.Equal(t, masterPassword, secretKey)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return nil
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginHelper(mockHlp)
		err, hint := loginer.CheckMasterKey(context.Background(), masterPassword)
		assert.NoError(t, err)
		assert.Empty(t, hint)
	})

	t.Run("getHelloData_err", func(t *testing.T) {

		testErr := errors.New("testErr")

		masterKeyPassword := "masterKeyPassword"

		mockSrv := NewMockAppServer(ctrl)

		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(nil, testErr).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv)
		err, hint := loginer.CheckMasterKey(context.Background(), masterKeyPassword)
		assert.ErrorIs(t, err, testErr)
		assert.Empty(t, hint)
	})

	t.Run("decryptMasterKey_err", func(t *testing.T) {

		testErr := errors.New("testErr")
		masterPassword := "masterPassword"

		helloData := &domain.HelloData{
			HelloEncrypted:     "helloEncrypted",
			MasterPasswordHint: "masterPasswordHint",
		}

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockDomainHelper(ctrl)

		mockHlp.EXPECT().DecryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) error {
			assert.Equal(t, masterPassword, secretKey)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return testErr
		}).Times(1)

		loginer := app.NewLoginer().LoginSever(mockSrv).LoginHelper(mockHlp)
		err, hint := loginer.CheckMasterKey(context.Background(), masterPassword)
		assert.ErrorIs(t, err, testErr)
		assert.Equal(t, helloData.MasterPasswordHint, hint)
	})
}
