package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier)
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

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier)
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

		otpPass := "otpPass"

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass, pass)
			return nil
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier)
		loginer.PassOTP(context.Background(), otpPass)
	})

	t.Run("err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)

		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		otpPass := "otpPass"
		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass, pass)
			return testErr
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier)
		loginer.PassOTP(context.Background(), otpPass)
	})
}

func TestLoginer_CheckMasterKey(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		mockVier.EXPECT().ShowDataAccessView().Times(1)

		masterKeyPassword := "masterKeyPassword"

		helloData := &domain.HelloData{
			EncryptedMasterKey: "encryptedMasterKey",
			HelloEncrypted:     "helloEncrypted",
			MasterKeyPassHint:  "masterKeyHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)
		masterKey := "masterKey"

		mockHlp.EXPECT().DecryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) (string, error) {
			assert.Equal(t, masterKeyPassword, secretKey)
			assert.Equal(t, helloData.EncryptedMasterKey, ciphertext)
			return masterKey, nil
		}).Times(1)

		helloDecrypted := "hello"
		mockHlp.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(ciphertext string, passphrase string) ([]byte, error) {
			assert.Equal(t, masterKey, passphrase)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return []byte(helloDecrypted), nil
		}).Times(1)

		mockHlp.EXPECT().CheckHello(gomock.Any()).DoAndReturn(func(helloDecrypted string) (bool, error) {
			assert.Equal(t, helloDecrypted, helloDecrypted)
			return true, nil
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
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

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})

	t.Run("decryptMasterKey_err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		masterKeyPassword := "masterKeyPassword"

		helloData := &domain.HelloData{
			EncryptedMasterKey: "encryptedMasterKey",
			HelloEncrypted:     "helloEncrypted",
			MasterKeyPassHint:  "masterKeyHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)

		mockHlp.EXPECT().DecryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) (string, error) {
			assert.Equal(t, masterKeyPassword, secretKey)
			assert.Equal(t, helloData.EncryptedMasterKey, ciphertext)
			return "", testErr
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})

	t.Run("decryptData_err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		masterKeyPassword := "masterKeyPassword"

		helloData := &domain.HelloData{
			EncryptedMasterKey: "encryptedMasterKey",
			HelloEncrypted:     "helloEncrypted",
			MasterKeyPassHint:  "masterKeyHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)

		masterKey := "masterKey"

		mockHlp.EXPECT().DecryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) (string, error) {
			assert.Equal(t, masterKeyPassword, secretKey)
			assert.Equal(t, helloData.EncryptedMasterKey, ciphertext)
			return masterKey, nil
		}).Times(1)

		mockHlp.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(ciphertext string, passphrase string) ([]byte, error) {
			assert.Equal(t, masterKey, passphrase)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return nil, testErr
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})

	t.Run("checkHello_err", func(t *testing.T) {
		mockVier := NewMockLoginView(ctrl)
		testErr := errors.New("testErr")
		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		masterKeyPassword := "masterKeyPassword"

		helloData := &domain.HelloData{
			EncryptedMasterKey: "encryptedMasterKey",
			HelloEncrypted:     "helloEncrypted",
			MasterKeyPassHint:  "masterKeyHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)
		masterKey := "masterKey"

		mockHlp.EXPECT().DecryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) (string, error) {
			assert.Equal(t, masterKeyPassword, secretKey)
			assert.Equal(t, helloData.EncryptedMasterKey, ciphertext)
			return masterKey, nil
		}).Times(1)

		helloDecrypted := "hello"
		mockHlp.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(ciphertext string, passphrase string) ([]byte, error) {
			assert.Equal(t, masterKey, passphrase)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return []byte(helloDecrypted), nil
		}).Times(1)

		mockHlp.EXPECT().CheckHello(gomock.Any()).DoAndReturn(func(helloDecrypted string) (bool, error) {
			assert.Equal(t, helloDecrypted, helloDecrypted)
			return false, testErr
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})

	t.Run("checkHello_false", func(t *testing.T) {

		masterKeyPassword := "masterKeyPassword"

		helloData := &domain.HelloData{
			EncryptedMasterKey: "encryptedMasterKey",
			HelloEncrypted:     "helloEncrypted",
			MasterKeyPassHint:  "masterKeyHint",
		}

		mockSrv := NewMockLoginServer(ctrl)
		mockSrv.EXPECT().GetHelloData(gomock.Any()).Return(helloData, nil).Times(1)

		mockHlp := NewMockLoginHelper(ctrl)
		masterKey := "masterKey"

		mockHlp.EXPECT().DecryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, ciphertext string) (string, error) {
			assert.Equal(t, masterKeyPassword, secretKey)
			assert.Equal(t, helloData.EncryptedMasterKey, ciphertext)
			return masterKey, nil
		}).Times(1)

		helloDecrypted := "hello"
		mockHlp.EXPECT().DecryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(ciphertext string, passphrase string) ([]byte, error) {
			assert.Equal(t, masterKey, passphrase)
			assert.Equal(t, helloData.HelloEncrypted, ciphertext)
			return []byte(helloDecrypted), nil
		}).Times(1)

		mockHlp.EXPECT().CheckHello(gomock.Any()).DoAndReturn(func(helloDecrypted string) (bool, error) {
			assert.Equal(t, helloDecrypted, helloDecrypted)
			return false, nil
		}).Times(1)

		mockVier := NewMockLoginView(ctrl)

		mockVier.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
		}).Times(1)

		mockVier.EXPECT().ShowMasterKeyView(gomock.Any()).Do(func(msg string) {
			assert.Equal(t, helloData.MasterKeyPassHint, msg)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		loginer := app.NewLoginer(conf).LoginSever(mockSrv).LoginView(mockVier).LoginHelper(mockHlp)
		loginer.CheckMasterKey(context.Background(), masterKeyPassword)
	})
}
