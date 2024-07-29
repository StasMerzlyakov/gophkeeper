package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/usecases"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataAccessor_GetHelloData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)
		encryptedKey := &domain.HelloData{
			HelloEncrypted:     "encryptedKey",
			MasterPasswordHint: "masterPasswordHint",
		}
		mockStorage.EXPECT().GetHelloData(gomock.Any()).Times(1).Return(encryptedKey, nil)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		res, err := da.GetHelloData(context.Background())
		require.NoError(t, err)
		assert.Equal(t, encryptedKey.HelloEncrypted, res.HelloEncrypted)
		assert.Equal(t, encryptedKey.MasterPasswordHint, res.MasterPasswordHint)
	})

	t.Run("err", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)

		testErr := errors.New("testErr")
		mockStorage.EXPECT().GetHelloData(gomock.Any()).Times(1).Return(nil, testErr)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		_, err := da.GetHelloData(context.Background())
		require.ErrorIs(t, err, testErr)

	})

}
