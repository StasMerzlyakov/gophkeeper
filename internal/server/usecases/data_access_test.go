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

func TestDataAccessor_GetBankCardList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		list := []domain.EncryptedBankCard{
			{
				Number:  "Number1",
				Content: "Content1",
			},
			{
				Number:  "Number2",
				Content: "Content1",
			},
		}

		mockStorage.EXPECT().GetBankCardList(gomock.Any()).Times(1).Return(list, nil)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		res, err := da.GetBankCardList(context.Background())
		require.NoError(t, err)
		require.Equal(t, 2, len(res))
		require.Equal(t, "Number1", res[0].Number)
		require.Equal(t, "Number2", res[1].Number)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		list := []domain.EncryptedBankCard{
			{
				Number:  "Number1",
				Content: "Content1",
			},
			{
				Number:  "Number2",
				Content: "Content1",
			},
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().GetBankCardList(gomock.Any()).Times(1).Return(list, testErr)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		res, err := da.GetBankCardList(context.Background())
		require.Nil(t, res)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_GetUserPasswordDataList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		list := []domain.EncryptedUserPasswordData{
			{
				Hint:    "Hint1",
				Content: "Content1",
			},
			{
				Hint:    "Hint2",
				Content: "Content1",
			},
		}

		mockStorage.EXPECT().GetUserPasswordDataList(gomock.Any()).Times(1).Return(list, nil)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		res, err := da.GetUserPasswordDataList(context.Background())
		require.NoError(t, err)
		require.Equal(t, 2, len(res))
		require.Equal(t, "Hint1", res[0].Hint)
		require.Equal(t, "Hint2", res[1].Hint)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		list := []domain.EncryptedUserPasswordData{
			{
				Hint:    "Hint1",
				Content: "Content1",
			},
			{
				Hint:    "Hint2",
				Content: "Content1",
			},
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().GetUserPasswordDataList(gomock.Any()).Times(1).Return(list, testErr)
		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		res, err := da.GetUserPasswordDataList(context.Background())
		require.Nil(t, res)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_CreateBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedBankCard{
			Number:  "Number1",
			Content: "Content1",
		}

		mockStorage.EXPECT().CreateBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
				assert.Equal(t, dt.Number, data.Number)
				assert.Equal(t, dt.Content, data.Content)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.CreateBankCard(context.Background(), dt)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedBankCard{
			Number:  "Number1",
			Content: "Content1",
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().CreateBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
				assert.Equal(t, dt.Number, data.Number)
				assert.Equal(t, dt.Content, data.Content)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.CreateBankCard(context.Background(), dt)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_CreateUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedUserPasswordData{
			Hint:    "Hint1",
			Content: "Content1",
		}

		mockStorage.EXPECT().CreateUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
				assert.Equal(t, dt.Hint, data.Hint)
				assert.Equal(t, dt.Content, data.Content)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.CreateUserPasswordData(context.Background(), dt)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedUserPasswordData{
			Hint:    "Hint1",
			Content: "Content1",
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().CreateUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
				assert.Equal(t, dt.Hint, data.Hint)
				assert.Equal(t, dt.Content, data.Content)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.CreateUserPasswordData(context.Background(), dt)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_UpdateBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedBankCard{
			Number:  "Number1",
			Content: "Content1",
		}

		mockStorage.EXPECT().UpdateBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
				assert.Equal(t, dt.Number, data.Number)
				assert.Equal(t, dt.Content, data.Content)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.UpdateBankCard(context.Background(), dt)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedBankCard{
			Number:  "Number1",
			Content: "Content1",
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().UpdateBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
				assert.Equal(t, dt.Number, data.Number)
				assert.Equal(t, dt.Content, data.Content)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.UpdateBankCard(context.Background(), dt)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_UpdateUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedUserPasswordData{
			Hint:    "Hint1",
			Content: "Content1",
		}

		mockStorage.EXPECT().UpdateUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
				assert.Equal(t, dt.Hint, data.Hint)
				assert.Equal(t, dt.Content, data.Content)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.UpdateUserPasswordData(context.Background(), dt)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)
		dt := &domain.EncryptedUserPasswordData{
			Hint:    "Hint1",
			Content: "Content1",
		}

		testErr := errors.New("testErr")
		mockStorage.EXPECT().UpdateUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
				assert.Equal(t, dt.Hint, data.Hint)
				assert.Equal(t, dt.Content, data.Content)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.UpdateUserPasswordData(context.Background(), dt)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_DeleteBankCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)

		number := "Number1"

		mockStorage.EXPECT().DeleteBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, num string) error {
				assert.Equal(t, number, num)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.DeleteBankCard(context.Background(), number)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)

		number := "Number1"

		testErr := errors.New("testErr")
		mockStorage.EXPECT().DeleteBankCard(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, num string) error {
				assert.Equal(t, number, num)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.DeleteBankCard(context.Background(), number)
		require.ErrorIs(t, err, testErr)
	})
}

func TestDataAccessor_DeleteUserPasswordData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockStorage := NewMockStateFullStorage(ctrl)

		hint := "Hint1"

		mockStorage.EXPECT().DeleteUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, hnt string) error {
				assert.Equal(t, hint, hnt)
				return nil
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)

		err := da.DeleteUserPasswordData(context.Background(), hint)
		require.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {

		mockStorage := NewMockStateFullStorage(ctrl)

		hint := "Hint1"

		testErr := errors.New("testErr")
		mockStorage.EXPECT().DeleteUserPasswordData(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, hnt string) error {
				assert.Equal(t, hint, hnt)
				return testErr
			}).Times(1)

		da := usecases.NewDataAccessor(nil).StateFullStorage(mockStorage)
		err := da.DeleteUserPasswordData(context.Background(), hint)
		require.ErrorIs(t, err, testErr)
	})
}
