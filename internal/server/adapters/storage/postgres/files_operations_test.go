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

func TestFileOperations(t *testing.T) {
	ctx, cancelFN := context.WithCancel(context.Background())

	defer cancelFN()

	connString, err := postgresContainer.ConnectionString(ctx)

	require.NoError(t, err)

	storage := postgres.NewStorage(ctx, &config.ServerConf{
		MaxConns:        5,
		DatabaseDN:      connString,
		MaxConnLifetime: 2 * time.Minute,
		MaxConnIdleTime: 2 * time.Minute,
	})

	defer func() {
		storage.Close()
		err = clear(ctx)
		require.NoError(t, err)
	}()

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
		MasterPasswordHint: "MasterPasswordHint",
		HelloEncrypted:     "HelloEncrypted",
	}

	err = storage.Registrate(ctx, regData)
	require.NoError(t, err)

	err = storage.Registrate(ctx, regData)
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
	lData, err := storage.GetLoginData(ctx, testEmail)
	require.NoError(t, err)

	ctxWithID := domain.EnrichWithUserID(ctx, lData.UserID)
	bucket, err := storage.GetUserFilesBucket(ctxWithID)
	require.NoError(t, err)
	require.NotEmpty(t, bucket)

	ctxWithID2 := domain.EnrichWithUserID(ctx, -1)
	_, err = storage.GetUserFilesBucket(ctxWithID2)
	require.ErrorIs(t, err, domain.ErrClientDataIncorrect)
}
