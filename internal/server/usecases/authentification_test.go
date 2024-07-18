package usecases_test

import (
	"context"
	"testing"

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
			EncryptedOTPKey: "encryptedOTPKey",
			EMail:           "test@email",
			PasswordHash:    "Hash",
			PasswordSalt:    "Salt",
		}

		mockStflStorage.EXPECT().GetLoginData(gomock.Any(), gomock.Eq(data.EMail)).Times(1).Return(loginData, nil)

		mockHelper := NewMockRegistrationHelper(ctrl)
		mockHelper.EXPECT().NewSessionID().Times(1).Return(sessionID)

		mockHelper.EXPECT().CheckPassword(gomock.Eq(data.Password),
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
				assert.Equal(t, loginData.PasswordHash, lData.PasswordHash)
				assert.Equal(t, loginData.PasswordSalt, lData.PasswordSalt)
				return nil
			})

		auth := usecases.NewAuth(nil, mockStflStorage, mockTempStorage, nil, mockHelper)
		sID, err := auth.Login(context.Background(), data)
		require.NoError(t, err)
		assert.Equal(t, sessionID, sID)
	})
}
