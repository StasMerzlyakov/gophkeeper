package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/usecases"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthentification_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		sessionID := domain.SessionID("sessionID")

		mockStflStorage := NewMockStateFullStorage(ctrl)
		loginData := &domain.LoginData{
			UserID:          1,
			EncryptedOTPKey: "encryptedOTPKey",
			EMail:           "test@email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(loginData, nil)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(sessionID)

		mockHelper.EXPECT().ValidateAccountPass(gomock.Eq(data.Password),
			gomock.Eq(loginData.PasswordHash),
			gomock.Eq(loginData.PasswordSalt)).Times(1).
			Return(true, nil)

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().Create(gomock.Any(), gomock.Eq(sessionID), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, sID domain.SessionID, input any) error {
				lData, ok := input.(domain.LoginData)
				require.True(t, ok)
				assert.Equal(t, loginData.EMail, lData.EMail)
				assert.Equal(t, loginData.EncryptedOTPKey, lData.EncryptedOTPKey)
				assert.Equal(t, loginData.UserID, lData.UserID)
				assert.Equal(t, loginData.PasswordHash, lData.PasswordHash)
				assert.Equal(t, loginData.PasswordSalt, lData.PasswordSalt)
				return nil
			})

		auth := usecases.NewAuth(nil).RegistrationHelper(mockHelper).StateFullStorage(mockStflStorage).TemporaryStorage(mockTempStorage)
		sID, err := auth.Login(context.Background(), data)
		require.NoError(t, err)
		assert.Equal(t, sessionID, sID)
	})

	t.Run("get_data_err", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStflStorage := NewMockStateFullStorage(ctrl)

		testErr := errors.New("testErr")
		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(nil, testErr)

		auth := usecases.NewAuth(nil).StateFullStorage(mockStflStorage)
		_, err := auth.Login(context.Background(), data)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("check_pass_err", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStflStorage := NewMockStateFullStorage(ctrl)
		loginData := &domain.LoginData{
			UserID:          1,
			EncryptedOTPKey: "encryptedOTPKey",
			EMail:           "test@email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(loginData, nil)

		mockHelper := NewMockRegistrationHelper(ctrl)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().ValidateAccountPass(gomock.Eq(data.Password),
			gomock.Eq(loginData.PasswordHash),
			gomock.Eq(loginData.PasswordSalt)).Times(1).
			Return(false, testErr)

		auth := usecases.NewAuth(nil).StateFullStorage(mockStflStorage).RegistrationHelper(mockHelper)
		_, err := auth.Login(context.Background(), data)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("pass_wrong", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStflStorage := NewMockStateFullStorage(ctrl)
		loginData := &domain.LoginData{
			UserID:          1,
			EncryptedOTPKey: "encryptedOTPKey",
			EMail:           "test@email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(loginData, nil)

		mockHelper := NewMockRegistrationHelper(ctrl)

		mockHelper.EXPECT().ValidateAccountPass(gomock.Eq(data.Password),
			gomock.Eq(loginData.PasswordHash),
			gomock.Eq(loginData.PasswordSalt)).Times(1).
			Return(false, nil)

		auth := usecases.NewAuth(nil).StateFullStorage(mockStflStorage).RegistrationHelper(mockHelper)
		_, err := auth.Login(context.Background(), data)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("create_err", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		sessionID := domain.SessionID("sessionID")

		mockStflStorage := NewMockStateFullStorage(ctrl)
		loginData := &domain.LoginData{
			UserID:          1,
			EncryptedOTPKey: "encryptedOTPKey",
			EMail:           "test@email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(loginData, nil)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(sessionID)

		mockHelper.EXPECT().ValidateAccountPass(gomock.Eq(data.Password),
			gomock.Eq(loginData.PasswordHash),
			gomock.Eq(loginData.PasswordSalt)).Times(1).
			Return(true, nil)

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		testErr := errors.New("testErr")
		mockTempStorage.EXPECT().Create(gomock.Any(), gomock.Eq(sessionID), gomock.Any()).Times(1).Return(testErr)

		auth := usecases.NewAuth(nil).StateFullStorage(mockStflStorage).RegistrationHelper(mockHelper).TemporaryStorage(mockTempStorage)

		_, err := auth.Login(context.Background(), data)
		require.ErrorIs(t, err, testErr)
	})
}

func TestCheckOTP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		userID := domain.UserID(1)

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.LoginData{
			UserID:          userID,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(true, nil)

		jwtTok := domain.JWTToken("jwtTok")
		mockHelper.EXPECT().CreateJWTToken(gomock.Eq(userID)).Times(1).Return(jwtTok, nil)

		auth := usecases.NewAuth(conf).RegistrationHelper(mockHelper).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		jTok, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.NoError(t, err)
		assert.Equal(t, jwtTok, jTok)
	})

	t.Run("load_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}

		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		testEx := errors.New("testEx")
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			nil, testEx,
		)

		auth := usecases.NewAuth(conf).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.ErrorIs(t, err, testEx)
	})

	t.Run("wrong_data", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}

		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{}

		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		auth := usecases.NewAuth(conf).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("validate_pass_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		userID := domain.UserID(1)

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.LoginData{
			UserID:          userID,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(false, testErr)

		auth := usecases.NewAuth(conf).TemporaryStorage(mockTempStorage).RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("otp_pass_fail", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		userID := domain.UserID(1)

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.LoginData{
			UserID:          userID,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(false, nil)

		auth := usecases.NewAuth(conf).TemporaryStorage(mockTempStorage).RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.Error(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("jwt_create_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		userID := domain.UserID(1)

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.LoginData{
			UserID:          userID,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(true, nil)

		testErr := errors.New("testErr")
		mockHelper.EXPECT().CreateJWTToken(gomock.Eq(userID)).Times(1).Return(domain.JWTToken(""), testErr)

		auth := usecases.NewAuth(conf).TemporaryStorage(mockTempStorage).RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := auth.CheckOTP(ctx, currentID, otpPass)
		require.ErrorIs(t, err, testErr)
	})
}
