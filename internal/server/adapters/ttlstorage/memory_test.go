package ttlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/ttlstorage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemStorage1(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	srvConf := &config.ServerConf{
		AuthStageTimeout: 5 * time.Second,
	}

	memStorage := ttlstorage.NewMemStorage(ctx, srvConf)

	regStage1 := domain.RegistrationData{
		EMail:            "test@test.com",
		PasswordHash:     "AABBCC",
		PasswordSalt:     "EE<asd",
		EncryptedOTPPass: "PASS",
		State:            domain.RegistrationStateInit,
	}

	stage1ID := domain.SessionID(uuid.NewString())

	err := memStorage.Create(ctx, stage1ID, regStage1)
	require.NoError(t, err)

	err = memStorage.Create(ctx, stage1ID, regStage1)
	require.ErrorIs(t, err, domain.ErrDublicateKeyViolation)

	val, err := memStorage.Load(ctx, stage1ID)
	require.NoError(t, err)
	require.Equal(t, regStage1, val)

	regStage2 := domain.RegistrationData{
		EMail:            "test2@test.com",
		PasswordHash:     "AABBCCD",
		PasswordSalt:     "EE<asdE",
		EncryptedOTPPass: "PASSSS",
		State:            domain.RegistrationStateAuthPassed,
	}

	stage2ID := domain.SessionID(uuid.NewString())
	val, err = memStorage.Load(ctx, stage2ID)
	require.ErrorIs(t, err, domain.ErrDataNotExists)
	require.Nil(t, val)

	err = memStorage.DeleteAndCreate(ctx, stage1ID, stage2ID, regStage2)
	require.NoError(t, err)

	val, err = memStorage.Load(ctx, stage2ID)
	require.NoError(t, err)
	require.Equal(t, regStage2, val)

	val, err = memStorage.LoadAndDelete(ctx, stage2ID)
	require.NoError(t, err)
	require.Equal(t, regStage2, val)

	val, err = memStorage.LoadAndDelete(ctx, stage2ID)
	require.ErrorIs(t, err, domain.ErrDataNotExists)
	require.Nil(t, val)
}

func TestMemStorage2(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	srvConf := &config.ServerConf{
		AuthStageTimeout: 5 * time.Second,
	}

	memStorage := ttlstorage.NewMemStorage(ctx, srvConf)

	regStage1 := domain.RegistrationData{
		EMail:            "test@test.com",
		PasswordHash:     "AABBCC",
		PasswordSalt:     "EE<asd",
		EncryptedOTPPass: "PASS",
		State:            domain.RegistrationStateInit,
	}

	stage1ID := domain.SessionID(uuid.NewString())

	stage2ID := domain.SessionID(uuid.NewString())

	err := memStorage.DeleteAndCreate(ctx, stage1ID, stage2ID, regStage1)
	require.ErrorIs(t, err, domain.ErrDataNotExists)

	err = memStorage.Create(ctx, stage2ID, regStage1)
	require.NoError(t, err)

	err = memStorage.Create(ctx, stage1ID, regStage1)
	require.NoError(t, err)

	err = memStorage.DeleteAndCreate(ctx, stage1ID, stage2ID, regStage1)
	require.ErrorIs(t, err, domain.ErrDublicateKeyViolation)

}
