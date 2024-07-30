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

func TestRegistrator_CheckEmail(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailAvailable, nil
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowMsg(gomock.Any()).Times(1)

		reg := app.NewRegistrator().RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("parseEmail_err", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return false
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		reg := app.NewRegistrator().RegHelper(mockHelper).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("checkEmail_err", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockDomainHelper(ctrl)

		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		testErr := errors.New("testErr")
		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailAvailable, testErr
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		reg := app.NewRegistrator().RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.CheckEmail(context.Background(), email)
	})

	t.Run("email_busy", func(t *testing.T) {
		email := "email"
		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, email, str)
			return true
		}).Times(1)

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, str string) (domain.EMailStatus, error) {
			require.Equal(t, email, str)
			return domain.EMailBusy, nil
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
		}).Times(1)

		reg := app.NewRegistrator().RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
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

		mockHelper := NewMockDomainHelper(ctrl)

		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.EMail, str)
			return true
		}).Times(1)

		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return true
		}).Times(1)

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().Registrate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			return nil
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowRegOTPView().Times(1)

		reg := app.NewRegistrator().RegHelper(mockHelper).RegServer(mockSrv).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})

	t.Run("wrong_email", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockDomainHelper(ctrl)

		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.EMail, str)
			return false
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowMsg(gomock.Any())

		reg := app.NewRegistrator().RegHelper(mockHelper).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})

	t.Run("validate_pass_err", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.EMail, str)
			return true
		}).Times(1)
		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return false
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowMsg(gomock.Any())

		reg := app.NewRegistrator().RegHelper(mockHelper).RegView(mockView)
		reg.Registrate(context.Background(), emailData)
	})

	t.Run("registrate_err", func(t *testing.T) {

		emailData := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().ParseEMail(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.EMail, str)
			return true
		}).Times(1)
		mockHelper.EXPECT().CheckAuthPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, emailData.Password, str)
			return true
		}).Times(1)

		mockSrv := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().Registrate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, data *domain.EMailData) error {
			return testErr
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		reg := app.NewRegistrator().RegServer(mockSrv).RegHelper(mockHelper).RegView(mockView)
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

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().PassRegOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return nil
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowRegMasterKeyView().Times(1)

		reg := app.NewRegistrator().RegServer(mockSrv).RegView(mockView)

		reg.PassOTP(context.Background(), otpPass)
	})

	t.Run("err", func(t *testing.T) {

		otpPass := &domain.OTPPass{
			Pass: "otpPass",
		}

		mockSrv := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().PassRegOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pass string) error {
			assert.Equal(t, otpPass.Pass, pass)
			return testErr
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			require.ErrorIs(t, err, testErr)
		}).Times(1)

		reg := app.NewRegistrator().RegServer(mockSrv).RegView(mockView)
		reg.PassOTP(context.Background(), otpPass)
	})

}

func TestRegistrator_InitMasterKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterPassword:     "MasterPassword",
			MasterPasswordHint: "MasterPasswordHint",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().CheckMasterPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		helloEncrypted := "helloEncrypted"
		mockHelper.EXPECT().EncryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(masterPassword, plaintext string) (string, error) {
			assert.Equal(t, masterPassword, keyData.MasterPassword)
			assert.Equal(t, plaintext, rndString)
			return helloEncrypted, nil
		})

		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().InitMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, mData *domain.MasterKeyData) error {

			assert.Equal(t, helloEncrypted, mData.HelloEncrypted)
			assert.Equal(t, keyData.MasterPasswordHint, mData.MasterPasswordHint)
			return nil
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowLoginView().Times(1)

		reg := app.NewRegistrator().RegServer(mockSrv).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("validate_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterPassword:     "MasterPassword",
			MasterPasswordHint: "MasterPasswordHint",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().CheckMasterPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterPassword, str)
			return false
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowMsg(gomock.Any()).Times(1)

		reg := app.NewRegistrator().RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("encrypt_err", func(t *testing.T) {
		keyData := &domain.UnencryptedMasterKeyData{
			MasterPassword:     "MasterPassword",
			MasterPasswordHint: "MasterPasswordHint",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().CheckMasterPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().EncryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(masterPassword, plaintext string) (string, error) {
			assert.Equal(t, masterPassword, keyData.MasterPassword)
			assert.Equal(t, plaintext, rndString)
			return "", testErr
		})

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		reg := app.NewRegistrator().RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})

	t.Run("init_key_err", func(t *testing.T) {

		keyData := &domain.UnencryptedMasterKeyData{
			MasterPassword:     "MasterPassword",
			MasterPasswordHint: "MasterPasswordHint",
		}

		mockHelper := NewMockDomainHelper(ctrl)
		mockHelper.EXPECT().CheckMasterPasswordComplexityLevel(gomock.Any()).DoAndReturn(func(str string) bool {
			require.Equal(t, keyData.MasterPassword, str)
			return true
		}).Times(1)

		rndString := "f1d66fA6Id7gUnQNo4imvp/USizQsg=="
		mockHelper.EXPECT().Random32ByteString().Return(rndString)

		helloEncrypted := "helloEncrypted"
		mockHelper.EXPECT().EncryptHello(gomock.Any(), gomock.Any()).DoAndReturn(func(masterPassword, plaintext string) (string, error) {
			assert.Equal(t, masterPassword, keyData.MasterPassword)
			assert.Equal(t, plaintext, rndString)
			return helloEncrypted, nil
		})

		mockSrv := NewMockAppServer(ctrl)
		testErr := errors.New("testErr")
		mockSrv.EXPECT().InitMasterKey(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, mData *domain.MasterKeyData) error {

			assert.Equal(t, helloEncrypted, mData.HelloEncrypted)
			assert.Equal(t, keyData.MasterPasswordHint, mData.MasterPasswordHint)
			return testErr
		}).Times(1)

		mockView := NewMockAppView(ctrl)
		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, testErr)
		}).Times(1)

		reg := app.NewRegistrator().RegServer(mockSrv).RegView(mockView).RegHelper(mockHelper)
		reg.InitMasterKey(context.Background(), keyData)
	})
}
