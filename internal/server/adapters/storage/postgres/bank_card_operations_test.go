package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestBankCardOperations(t *testing.T) {
	ctx, cancelFN := context.WithCancel(context.Background())

	defer cancelFN()

	connString, err := postgresContainer.ConnectionString(ctx)

	require.NoError(t, err)

	storage := postgres.NewStorage(ctx, &config.ServerConf{
		MaxConns:        5,
		DatabaseURI:     connString,
		MaxConnLifetime: 2 * time.Minute,
		MaxConnIdleTime: 2 * time.Minute,
	})

	defer func() {
		storage.Close()
		err = clear(ctx)
		require.NoError(t, err)
	}()

	err = storage.Ping(ctx)
	require.NoError(t, err)

	err = clear(ctx)
	require.NoError(t, err)

	testEmail := "email@email"
	testEmail2 := "email2@email"
	ok, err := storage.IsEMailAvailable(ctx, testEmail)
	require.NoError(t, err)
	require.True(t, ok)

	regData := &domain.FullRegistrationData{
		EMail:              testEmail,
		PasswordHash:       "PasswordHash",
		PasswordSalt:       "PasswordSalt",
		EncryptedOTPKey:    "EncryptedOTPKey",
		MasterPasswordHint: "MasterPasswordHint",
		HelloEncrypted:     "HelloEncrypted",
	}

	err = storage.Registrate(ctx, regData)
	require.NoError(t, err)

	lData, err := storage.GetLoginData(ctx, testEmail)
	require.NoError(t, err)
	require.NotNil(t, lData)

	regData2 := &domain.FullRegistrationData{
		EMail:              testEmail2,
		PasswordHash:       "PasswordHash",
		PasswordSalt:       "PasswordSalt",
		EncryptedOTPKey:    "EncryptedOTPKey",
		MasterPasswordHint: "MasterPasswordHint",
		HelloEncrypted:     "HelloEncrypted",
	}

	err = storage.Registrate(ctx, regData2)
	require.NoError(t, err)

	lData2, err := storage.GetLoginData(ctx, testEmail2)
	require.NoError(t, err)
	require.NotNil(t, lData)

	userID2 := lData2.UserID

	userID := lData.UserID

	userIdCtx := domain.EnrichWithUserID(ctx, userID)
	userId2Ctx := domain.EnrichWithUserID(ctx, userID2)

	_, err = storage.GetBankCardList(ctx)
	require.ErrorIs(t, err, domain.ErrServerInternal)

	resp, err := storage.GetBankCardList(userIdCtx)
	require.NoError(t, err)
	require.Equal(t, 0, len(resp))

	err = storage.CreateBankCard(userIdCtx, &domain.EncryptedBankCard{
		Number:  "number1",
		Content: "content1",
	})

	require.NoError(t, err)

	err = storage.CreateBankCard(userId2Ctx, &domain.EncryptedBankCard{
		Number:  "number1",
		Content: "content1",
	})

	require.NoError(t, err)

	err = storage.CreateBankCard(userId2Ctx, &domain.EncryptedBankCard{
		Number:  "number1",
		Content: "content1",
	})

	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

	resp, err = storage.GetBankCardList(userIdCtx)
	require.NoError(t, err)
	require.Equal(t, 1, len(resp))
	require.Equal(t, "content1", resp[0].Content)

	err = storage.UpdateBankCard(userIdCtx, &domain.EncryptedBankCard{
		Number:  "number1",
		Content: "content2",
	})

	err = storage.DeleteBankCard(userId2Ctx, "number1")
	require.NoError(t, err)

	err = storage.DeleteBankCard(userId2Ctx, "number2")
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

	resp, err = storage.GetBankCardList(userIdCtx)
	require.NoError(t, err)
	require.Equal(t, 1, len(resp))
	require.Equal(t, "content2", resp[0].Content)

	err = storage.DeleteBankCard(userIdCtx, "number1")
	require.NoError(t, err)

	resp, err = storage.GetBankCardList(userIdCtx)
	require.NoError(t, err)
	require.Equal(t, 0, len(resp))

}
