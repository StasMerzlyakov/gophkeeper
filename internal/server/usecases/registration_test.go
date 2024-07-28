package usecases_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/usecases"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestDataDirectory = "../../../testdata/"

func TestRegistrator_Regisration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qrFilePath := filepath.Join(TestDataDirectory, "QR.png")
	qrFile, err := os.Open(qrFilePath)
	require.NoError(t, err)

	defer func() {
		err := qrFile.Close()
		require.NoError(t, err)
	}()

	qr, err := io.ReadAll(qrFile)
	require.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		passHash := "hashed_password"
		passSalt := "passSalt"
		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(&domain.HashData{
			Hash: passHash,
			Salt: passSalt,
		}, nil)

		qrKey := "qrKey"
		mockHelper.EXPECT().GenerateQR(gomock.Eq(data.EMail)).Times(1).Return(qrKey, qr, nil)

		encryptedOTPKey := "encryptedOTPKey"
		mockHelper.EXPECT().EncryptOTPKey(gomock.Eq(qrKey)).Times(1).Return(encryptedOTPKey, nil)

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, sID domain.SessionID, input any) error {
				regData, ok := input.(domain.RegistrationData)
				require.True(t, ok)
				assert.Equal(t, sessionID, sID)
				assert.Equal(t, data.EMail, regData.EMail)
				assert.Equal(t, passHash, regData.PasswordHash)
				assert.Equal(t, passSalt, regData.PasswordSalt)
				assert.Equal(t, encryptedOTPKey, regData.EncryptedOTPKey)
				assert.Equal(t, domain.RegistrationStateInit, regData.State)
				return nil
			})

		mockSender := NewMockEMailSender(ctrl)
		mockSender.EXPECT().Send(gomock.Any(), gomock.Eq(data.EMail), gomock.Eq(qr)).Times(1)

		registrator := usecases.NewRegistrator(conf).
			TemporaryStorage(mockTempStorage).
			EMailSender(mockSender).
			StateFullStorage(mockStorage).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		sID, err := registrator.Registrate(ctx, data)
		assert.NoError(t, err)
		assert.Equal(t, sessionID, sID)
	})

	t.Run("email_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		testErr := errors.New("testErr")

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(false, testErr)

		registrator := usecases.NewRegistrator(conf).
			StateFullStorage(mockStorage).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("email_busy", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(false, nil)

		registrator := usecases.NewRegistrator(conf).
			StateFullStorage(mockStorage).RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("check_email_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockHelper := NewMockRegistrationHelper(ctrl)

		testErr := errors.New("CheckEMilErr")

		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return false, testErr
		})

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("hash_passport_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		testErr := errors.New("HashPasswordErr")

		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(nil, testErr)

		registrator := usecases.NewRegistrator(conf).
			StateFullStorage(mockStorage).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("generte_qr_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		passHash := "hashed_password"
		passSalt := "passSalt"
		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(&domain.HashData{
			Hash: passHash,
			Salt: passSalt,
		}, nil)

		testErr := errors.New("generate qr err")
		mockHelper.EXPECT().GenerateQR(gomock.Eq(data.EMail)).Times(1).Return("", nil, testErr)

		registrator := usecases.NewRegistrator(conf).
			StateFullStorage(mockStorage).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("encrypt_data_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		passHash := "hashed_password"
		passSalt := "passSalt"
		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(&domain.HashData{
			Hash: passHash,
			Salt: passSalt,
		}, nil)

		qrKey := "qrKey"
		mockHelper.EXPECT().GenerateQR(gomock.Eq(data.EMail)).Times(1).Return(qrKey, qr, nil)

		testErr := errors.New("encrypt data err")
		mockHelper.EXPECT().EncryptOTPKey(gomock.Eq(qrKey)).Times(1).Return("", testErr)

		registrator := usecases.NewRegistrator(conf).
			StateFullStorage(mockStorage).
			RegistrationHelper(mockHelper)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("create_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		passHash := "hashed_password"
		passSalt := "passSalt"
		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(&domain.HashData{
			Hash: passHash,
			Salt: passSalt,
		}, nil)

		qrKey := "qrKey"
		mockHelper.EXPECT().GenerateQR(gomock.Eq(data.EMail)).Times(1).Return(qrKey, qr, nil)

		encryptedOTPKey := "encryptedOTPKey"
		mockHelper.EXPECT().EncryptOTPKey(gomock.Eq(qrKey)).Times(1).Return(encryptedOTPKey, nil)

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		testErr := errors.New("create data err")
		mockTempStorage.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, sID domain.SessionID, input any) error {
				regData, ok := input.(domain.RegistrationData)
				require.True(t, ok)
				assert.Equal(t, sessionID, sID)
				assert.Equal(t, data.EMail, regData.EMail)
				assert.Equal(t, passHash, regData.PasswordHash)
				assert.Equal(t, passSalt, regData.PasswordSalt)
				assert.Equal(t, encryptedOTPKey, regData.EncryptedOTPKey)
				assert.Equal(t, domain.RegistrationStateInit, regData.State)
				return testErr
			})

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			StateFullStorage(mockStorage).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})

	t.Run("send_err", func(t *testing.T) {
		conf := &config.ServerConf{
			DomainName:          "localhost",
			ServerEncryptionKey: "secret_key",
		}

		data := &domain.EMailData{
			EMail:    "test@email",
			Password: "test_pass",
		}
		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(true, nil)

		sessionID := domain.SessionID(uuid.NewString())
		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(domain.SessionID(sessionID))
		mockHelper.EXPECT().CheckEMailData(gomock.Any()).Times(1).DoAndReturn(func(dt *domain.EMailData) (bool, error) {
			require.NotNil(t, dt)
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return true, nil
		})

		passHash := "hashed_password"
		passSalt := "passSalt"
		mockHelper.EXPECT().HashPassword(gomock.Eq(data.Password)).Times(1).Return(&domain.HashData{
			Hash: passHash,
			Salt: passSalt,
		}, nil)

		qrKey := "qrKey"
		mockHelper.EXPECT().GenerateQR(gomock.Eq(data.EMail)).Times(1).Return(qrKey, qr, nil)

		encryptedOTPKey := "encryptedOTPKey"
		mockHelper.EXPECT().EncryptOTPKey(gomock.Eq(qrKey)).Times(1).Return(encryptedOTPKey, nil)

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, sID domain.SessionID, input any) error {
				regData, ok := input.(domain.RegistrationData)
				require.True(t, ok)
				assert.Equal(t, sessionID, sID)
				assert.Equal(t, data.EMail, regData.EMail)
				assert.Equal(t, passHash, regData.PasswordHash)
				assert.Equal(t, passSalt, regData.PasswordSalt)
				assert.Equal(t, encryptedOTPKey, regData.EncryptedOTPKey)
				assert.Equal(t, domain.RegistrationStateInit, regData.State)
				return nil
			})

		mockSender := NewMockEMailSender(ctrl)
		testErr := errors.New("send qr err")
		mockSender.EXPECT().Send(gomock.Any(), gomock.Eq(data.EMail), gomock.Eq(qr)).Times(1).Return(testErr)

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			StateFullStorage(mockStorage).
			TemporaryStorage(mockTempStorage).
			EMailSender(mockSender)

		ctx := context.Background()
		_, err := registrator.Registrate(ctx, data)
		assert.ErrorIs(t, err, testErr)
	})
}

func TestRegistrator_GetEMailStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	email := "test@email"

	t.Run("available", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(email)).Times(1).Return(true, nil)

		registrator := usecases.NewRegistrator(nil).
			StateFullStorage(mockStorage)

		ctx := context.Background()
		st, err := registrator.GetEMailStatus(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, domain.EMailAvailable, st)
	})

	t.Run("busy", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(email)).Times(1).Return(false, nil)

		registrator := usecases.NewRegistrator(nil).
			StateFullStorage(mockStorage)

		ctx := context.Background()
		st, err := registrator.GetEMailStatus(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, domain.EMailBusy, st)
	})

	t.Run("err", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)
		testErr := errors.New("test err")
		mockStorage.EXPECT().IsEMailAvailable(gomock.Any(), gomock.Eq(email)).Times(1).Return(false, testErr)

		registrator := usecases.NewRegistrator(nil).
			StateFullStorage(mockStorage)

		ctx := context.Background()
		st, err := registrator.GetEMailStatus(ctx, email)
		assert.ErrorIs(t, err, testErr)
		assert.Equal(t, domain.EMailBusy, st)
	})
}

func TestRegistrator_PassOTP(t *testing.T) {
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

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(true, nil)

		newSessionID := domain.SessionID("newSessionID")
		mockHelper.EXPECT().NewSessionID().Return(newSessionID)

		mockTempStorage.EXPECT().DeleteAndCreate(gomock.Any(), gomock.Eq(currentID), gomock.Eq(newSessionID), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context,
				oldSessionID domain.SessionID,
				sessionID domain.SessionID,
				input any,
			) error {
				regData, ok := input.(domain.RegistrationData)
				require.True(t, ok)
				assert.Equal(t, currentID, oldSessionID)
				assert.Equal(t, newSessionID, sessionID)
				assert.Equal(t, data.EMail, regData.EMail)
				assert.Equal(t, data.PasswordHash, regData.PasswordHash)
				assert.Equal(t, data.PasswordSalt, regData.PasswordSalt)
				assert.Equal(t, data.EncryptedOTPKey, regData.EncryptedOTPKey)
				assert.Equal(t, domain.RegistrationStateAuth, regData.State)
				return nil
			})

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		nId, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.NoError(t, err)
		assert.Equal(t, newSessionID, nId)
	})

	t.Run("load_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		testErr := errors.New("testErr")
		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, testErr,
		)

		registrator := usecases.NewRegistrator(conf).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("wrong_type_loaded", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}

		encryptedOTPKey := "asdasdasd!iasd"
		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateAuth,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		registrator := usecases.NewRegistrator(conf).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("decrypt_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"

		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		testErr := errors.New("test error")

		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			"", testErr,
		)

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("validate_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		testErr := errors.New("test_err")
		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(false, testErr)

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("validate_pass_wrong", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(false, nil)

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, domain.ErrAuthDataIncorrect)
	})

	t.Run("delete_and_create_err", func(t *testing.T) {
		currentID := domain.SessionID("currentID")
		conf := &config.ServerConf{
			ServerEncryptionKey: "serverKey",
		}
		encryptedOTPKey := "asdasdasd!iasd"
		decryptedOTPKey := "otp_key"
		otpPass := "123345"

		mockTempStorage := NewMockTemporaryStorage(ctrl)

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage.EXPECT().Load(gomock.Any(), gomock.Eq(currentID)).Times(1).Return(
			data, nil,
		)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().DecryptOTPKey(gomock.Eq(encryptedOTPKey)).Times(1).Return(
			decryptedOTPKey, nil,
		)

		mockHelper.EXPECT().ValidateOTPCode(gomock.Eq(decryptedOTPKey), gomock.Eq(otpPass)).Times(1).Return(true, nil)

		newSessionID := domain.SessionID("newSessionID")
		mockHelper.EXPECT().NewSessionID().Return(newSessionID)

		testErr := errors.New("test_err")
		mockTempStorage.EXPECT().DeleteAndCreate(gomock.Any(), gomock.Eq(currentID), gomock.Eq(newSessionID), gomock.Any()).Times(1).
			Return(testErr)

		registrator := usecases.NewRegistrator(conf).
			RegistrationHelper(mockHelper).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		_, err := registrator.PassOTP(ctx, domain.SessionID(currentID), otpPass)
		require.ErrorIs(t, err, testErr)
	})

}

func TestRegistration_InitMasterKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		currentID := domain.SessionID("currentID")

		encryptedMasterKey := "EncryptedMasterKey"
		masterKeyHint := "MasterKeyHint"
		helloEncrypted := "HelloEncrypted"
		encryptedOTPKey := "asdasdasd!iasd"

		keyData := &domain.MasterKeyData{
			EncryptedMasterKey: encryptedMasterKey,
			MasterKeyHint:      masterKeyHint,
			HelloEncrypted:     helloEncrypted,
		}

		data := domain.RegistrationData{
			State:           domain.RegistrationStateAuth,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).DoAndReturn(func(ctx context.Context,
			id domain.SessionID,
		) (any, error) {
			return data, nil
		})

		mockStflStorage := NewMockStateFullStorage(ctrl)

		mockStflStorage.EXPECT().Registrate(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, fullData *domain.FullRegistrationData) error {
				assert.Equal(t, data.EncryptedOTPKey, fullData.EncryptedOTPKey)
				assert.Equal(t, data.EMail, fullData.EMail)
				assert.Equal(t, data.PasswordHash, fullData.PasswordHash)
				assert.Equal(t, data.PasswordSalt, fullData.PasswordSalt)
				assert.Equal(t, keyData.EncryptedMasterKey, fullData.EncryptedMasterKey)
				assert.Equal(t, keyData.MasterKeyHint, fullData.MasterKeyHint)
				assert.Equal(t, keyData.HelloEncrypted, fullData.HelloEncrypted)
				return nil
			})

		registrator := usecases.NewRegistrator(nil).
			StateFullStorage(mockStflStorage).
			TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		err := registrator.InitMasterKey(ctx, currentID, keyData)
		require.NoError(t, err)
	})

	t.Run("load_and_del_err", func(t *testing.T) {

		currentID := domain.SessionID("currentID")

		encryptedMasterKey := "EncryptedMasterKey"
		masterKeyHint := "MasterKeyHint"
		helloEncrypted := "HelloEncrypted"

		keyData := &domain.MasterKeyData{
			EncryptedMasterKey: encryptedMasterKey,
			MasterKeyHint:      masterKeyHint,
			HelloEncrypted:     helloEncrypted,
		}

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		testErr := errors.New("test_error")
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).DoAndReturn(func(ctx context.Context,
			id domain.SessionID,
		) (any, error) {
			return nil, testErr
		})

		registrator := usecases.NewRegistrator(nil).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		err := registrator.InitMasterKey(ctx, currentID, keyData)
		require.ErrorIs(t, err, testErr)
	})

	t.Run("wrong_data", func(t *testing.T) {

		currentID := domain.SessionID("currentID")

		encryptedMasterKey := "EncryptedMasterKey"
		masterKeyHint := "MasterKeyHint"
		helloEncrypted := "HelloEncrypted"

		keyData := &domain.MasterKeyData{
			EncryptedMasterKey: encryptedMasterKey,
			MasterKeyHint:      masterKeyHint,
			HelloEncrypted:     helloEncrypted,
		}

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).DoAndReturn(func(ctx context.Context,
			id domain.SessionID,
		) (any, error) {
			return "asdasd", nil
		})

		registrator := usecases.NewRegistrator(nil).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		err := registrator.InitMasterKey(ctx, currentID, keyData)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("wrong_state", func(t *testing.T) {

		currentID := domain.SessionID("currentID")

		encryptedMasterKey := "EncryptedMasterKey"
		masterKeyHint := "MasterKeyHint"
		helloEncrypted := "HelloEncrypted"
		encryptedOTPKey := "asdasdasd!iasd"

		keyData := &domain.MasterKeyData{
			EncryptedMasterKey: encryptedMasterKey,
			MasterKeyHint:      masterKeyHint,
			HelloEncrypted:     helloEncrypted,
		}

		data := domain.RegistrationData{
			State:           domain.RegistrationStateInit,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).DoAndReturn(func(ctx context.Context,
			id domain.SessionID,
		) (any, error) {
			return data, nil
		})

		registrator := usecases.NewRegistrator(nil).TemporaryStorage(mockTempStorage)

		ctx := context.Background()
		err := registrator.InitMasterKey(ctx, currentID, keyData)
		require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	})

	t.Run("registrate_err", func(t *testing.T) {

		currentID := domain.SessionID("currentID")

		encryptedMasterKey := "EncryptedMasterKey"
		masterKeyHint := "MasterKeyHint"
		helloEncrypted := "HelloEncrypted"
		encryptedOTPKey := "asdasdasd!iasd"

		keyData := &domain.MasterKeyData{
			EncryptedMasterKey: encryptedMasterKey,
			MasterKeyHint:      masterKeyHint,
			HelloEncrypted:     helloEncrypted,
		}

		data := domain.RegistrationData{
			State:           domain.RegistrationStateAuth,
			EncryptedOTPKey: encryptedOTPKey,
			EMail:           "Email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockTempStorage := NewMockTemporaryStorage(ctrl)
		mockTempStorage.EXPECT().LoadAndDelete(gomock.Any(), gomock.Eq(currentID)).Times(1).DoAndReturn(func(ctx context.Context,
			id domain.SessionID,
		) (any, error) {
			return data, nil
		})

		mockStflStorage := NewMockStateFullStorage(ctrl)
		testErr := errors.New("testErr")
		mockStflStorage.EXPECT().Registrate(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, fullData *domain.FullRegistrationData) error {
				return testErr
			})

		mockHelper := NewMockRegistrationHelper(ctrl)

		registrator := usecases.NewRegistrator(nil).TemporaryStorage(mockTempStorage).StateFullStorage(mockStflStorage).RegistrationHelper(mockHelper)

		ctx := context.Background()
		err := registrator.InitMasterKey(ctx, currentID, keyData)
		require.ErrorIs(t, err, testErr)
	})
}
