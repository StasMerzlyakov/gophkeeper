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

func TestStorageOperations(t *testing.T) {
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
		err = clear(ctx)
		require.NoError(t, err)
	}()

	err = storage.Ping(ctx)
	require.NoError(t, err)

	err = clear(ctx)
	require.NoError(t, err)

	testEmail := "email@email"
	ok, err := storage.IsEMailAvailable(ctx, testEmail)
	require.NoError(t, err)
	require.True(t, ok)

	regData := &domain.FullRegistrationData{
		EMail:              testEmail,
		PasswordHash:       "PasswordHash",
		PasswordSalt:       "PasswordSalt",
		EncryptedOTPKey:    "EncryptedOTPKey",
		EncryptedMasterKey: "EncryptedMasterKey",
		MasterKeyHint:      "MasterKeyHint",
		HelloEncrypted:     "HelloEncrypted",
	}

	lData, err := storage.GetLoginData(ctx, testEmail)
	require.Nil(t, lData)
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

	err = storage.Registrate(ctx, regData)
	require.NoError(t, err)

	err = storage.Registrate(ctx, regData)
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

	ok, err = storage.IsEMailAvailable(ctx, testEmail)
	require.NoError(t, err)
	require.False(t, ok)

	lData, err = storage.GetLoginData(ctx, testEmail)
	require.NoError(t, err)

	require.Equal(t, regData.EMail, lData.EMail)
	require.Equal(t, regData.EncryptedOTPKey, lData.EncryptedOTPKey)
	require.Equal(t, regData.PasswordHash, lData.PasswordHash)
	require.Equal(t, regData.PasswordSalt, lData.PasswordSalt)
	require.True(t, lData.UserID > 0)

	_, err = storage.GetLoginData(ctx, "testEmail")
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)

	_, err = storage.GetHelloData(ctx)
	require.ErrorIs(t, err, domain.ErrServerInternal)

	ctxWithID := domain.EnrichWithUserID(ctx, lData.UserID)
	helloData, err := storage.GetHelloData(ctxWithID)
	require.NoError(t, err)
	require.Equal(t, regData.HelloEncrypted, helloData.HelloEncrypted)
	require.Equal(t, regData.EncryptedMasterKey, helloData.EncryptedMasterKey)

	require.Equal(t, regData.MasterKeyHint, helloData.MasterKeyPassHint)

	ctxWithID2 := domain.EnrichWithUserID(ctx, -1)
	_, err = storage.GetHelloData(ctxWithID2)
	require.ErrorIs(t, err, domain.ErrServerInternal)
}
