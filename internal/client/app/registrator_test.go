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
	"github.com/stretchr/testify/require"
)

func TestRegistrator_CheckEmail(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailAvailable, nil
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowMsg(gomock.Any()).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("parseEmail_err", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return false
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("checkEmail_err", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockRegHelper(ctrl)

		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		testErr := errors.New("testErr")
		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailAvailable, testErr
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("email_busy", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailBusy, nil
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})
}

func TestRegistrator_Registrate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return true
		}).Times(1)

		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().Registrate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			return nil
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowRegOTPView().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}

		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})

	t.Run("validate_pass_err", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return false
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}

		reg := app.NewRegistrator(conf).RegHelper(mockHelper).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})

	t.Run("registrate_err", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return true
		}).Times(1)

		mockSrv := NewMockRegServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().Registrate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			return testErr
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}

		reg := app.NewRegistrator(conf).RegServer(mockSrv).RegHelper(mockHelper).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})
}

func TestRegistrator_PassOTP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}

		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().PassRegOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return nil
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowInitMasterKeyView().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}

		reg := app.NewRegistrator(conf).RegServer(mockSrv).RegView(mockView)

		reg.PassOTP(context.Background(), otpPass)
	})

	t.Run("err", func(t *testing.T) {

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}

		mockSrv := NewMockRegServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().PassRegOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return testErr
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}

		reg := app.NewRegistrator(conf).RegServer(mockSrv).RegView(mockView)
		reg.PassOTP(context.Background(), otpPass)
	})

}

func TestRegistrator_InitMasterKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		encryptedMasterKey := "encryptedMasterKey"
		mockHelper.EXPECT().EncryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, plaintext string) (string, error) {
			assert.Equal(t, secretKey, keyData.MasterKeyPassword)
			assert.Equal(t, plaintext, rndString)
			return encryptedMasterKey, nil
		})

		helloStr := "helloStr"
		mockHelper.EXPECT().GenerateHello().Return(helloStr, nil)

		helloEncrypted := "helloEncrypted"

		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(hello []byte, masteKey string) (string, error) {
			assert.Equal(t, helloStr, string(hello))
			assert.Equal(t, masteKey, rndString)
			return helloEncrypted, nil
		})

		mockSrv := NewMockRegServer(ctrl)
		mockSrv.EXPECT().InitMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, mData *domain.MasterKeyData) error {

			assert.Equal(t, encryptedMasterKey, mData.EncryptedMasterKey)
			assert.Equal(t, keyData.MasterKeyHint, mData.MasterKeyHint)
			assert.Equal(t, helloEncrypted, mData.HelloEncrypted)

			return nil
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowLoginView().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegServer(mockSrv).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("validate_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return false
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("encrypt_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, plaintext string) (string, error) {
			assert.Equal(t, secretKey, keyData.MasterKeyPassword)
			assert.Equal(t, plaintext, rndString)
			return "", testErr
		})

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("generate_hello_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		encryptedMasterKey := "encryptedMasterKey"
		mockHelper.EXPECT().EncryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, plaintext string) (string, error) {
			assert.Equal(t, secretKey, keyData.MasterKeyPassword)
			assert.Equal(t, plaintext, rndString)
			return encryptedMasterKey, nil
		})

		testErr := errors.New("testErr")
		mockHelper.EXPECT().GenerateHello().Return("", testErr)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("encrypt_aes256_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		encryptedMasterKey := "encryptedMasterKey"
		mockHelper.EXPECT().EncryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, plaintext string) (string, error) {
			assert.Equal(t, secretKey, keyData.MasterKeyPassword)
			assert.Equal(t, plaintext, rndString)
			return encryptedMasterKey, nil
		})

		helloStr := "helloStr"
		mockHelper.EXPECT().GenerateHello().Return(helloStr, nil)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(hello []byte, masteKey string) (string, error) {
			assert.Equal(t, helloStr, string(hello))
			assert.Equal(t, masteKey, rndString)
			return "", testErr
		})

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("init_key_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterKeyPassword: "MasterKeyPassword",
			MasterKeyHint:     "MasterKeyHint",
		}

		mockHelper := NewMockRegHelper(ctrl)
		mockHelper.EXPECT().CheckMasterKeyPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterKeyPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		encryptedMasterKey := "encryptedMasterKey"
		mockHelper.EXPECT().EncryptMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(secretKey string, plaintext string) (string, error) {
			assert.Equal(t, secretKey, keyData.MasterKeyPassword)
			assert.Equal(t, plaintext, rndString)
			return encryptedMasterKey, nil
		})

		helloStr := "helloStr"
		mockHelper.EXPECT().GenerateHello().Return(helloStr, nil)

		helloEncrypted := "helloEncrypted"

		mockHelper.EXPECT().EncryptShortData(gomock.Any(), gomock.Any()).DoAndReturn(func(hello []byte, masteKey string) (string, error) {
			assert.Equal(t, helloStr, string(hello))
			assert.Equal(t, masteKey, rndString)
			return helloEncrypted, nil
		})

		mockSrv := NewMockRegServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().InitMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, mData *domain.MasterKeyData) error {

			assert.Equal(t, encryptedMasterKey, mData.EncryptedMasterKey)
			assert.Equal(t, keyData.MasterKeyHint, mData.MasterKeyHint)
			assert.Equal(t, helloEncrypted, mData.HelloEncrypted)

			return testErr
		}).Times(1)

		mockView := NewMockRegView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		reg := app.NewRegistrator(conf).RegServer(mockSrv).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})
}
