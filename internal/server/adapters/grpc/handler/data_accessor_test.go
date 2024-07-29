package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)

		mockService.EXPECT().GetHelloData(gomock.Any()).
			DoAndReturn(func(ctx context.Context) (*domain.HelloData, error) {
				return &domain.HelloData{
					HelloEncrypted:     "hello",
					MasterPasswordHint: "masterKey",
				}, nil
			}).Times(1)

		aService := handler.NewDataAccessor(mockService)

		res, err := aService.Hello(context.Background(), &proto.HelloRequest{})

		require.Nil(t, err)
		require.Equal(t, "hello", res.HelloEncrypted)
		require.Equal(t, "masterKey", res.MasterPasswordHint)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)

		testErr := errors.New("testErr")
		mockService.EXPECT().GetHelloData(gomock.Any()).
			DoAndReturn(func(ctx context.Context) (*domain.HelloData, error) {
				return nil, testErr
			}).Times(1)

		aService := handler.NewDataAccessor(mockService)

		_, err := aService.Hello(context.Background(), &proto.HelloRequest{})

		require.ErrorIs(t, err, testErr)
	})

}

func TestBankCardOps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("get_list_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().GetBankCardList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.EncryptedBankCard, error) {
			return []domain.EncryptedBankCard{
				{
					Number:  "number1",
					Content: "content1",
				},
				{
					Number:  "number2",
					Content: "content2",
				},
			}, nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		resp, err := aService.GetBankCardList(context.Background(), &proto.BankCardListRequest{})
		require.NoError(t, err)
		require.Equal(t, 2, len(resp.Cards))

		assert.Equal(t, "number1", resp.Cards[0].Number)
		assert.Equal(t, "content1", resp.Cards[0].Content)

		assert.Equal(t, "number2", resp.Cards[1].Number)
		assert.Equal(t, "content2", resp.Cards[1].Content)
	})

	t.Run("get_list_err", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		testErr := errors.New("testErr")
		mockService.EXPECT().GetBankCardList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.EncryptedBankCard, error) {
			return nil, testErr
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		resp, err := aService.GetBankCardList(context.Background(), &proto.BankCardListRequest{})
		require.ErrorIs(t, err, testErr)
		require.Nil(t, resp)
	})

	t.Run("create_bank_card_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().CreateBankCard(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
			require.Equal(t, "number1", data.Number)
			require.Equal(t, "content1", data.Content)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.CreateBankCard(context.Background(), &proto.CreateBankCardRequest{
			Number:  "number1",
			Content: "content1",
		})

		require.NoError(t, err)
	})

	t.Run("delete_bank_card_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().DeleteBankCard(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, number string) error {
			require.Equal(t, "number1", number)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.DeleteBankCard(context.Background(), &proto.DeleteBankCardRequest{
			Number: "number1",
		})

		require.NoError(t, err)
	})

	t.Run("update_bank_card_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().UpdateBankCard(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, data *domain.EncryptedBankCard) error {
			require.Equal(t, "number1", data.Number)
			require.Equal(t, "content1", data.Content)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.UpdateBankCard(context.Background(), &proto.UpdateBankCardRequest{
			Number:  "number1",
			Content: "content1",
		})

		require.NoError(t, err)
	})
}

func TestUserPasswordDataOps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("get_list_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().GetUserPasswordDataList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {
			return []domain.EncryptedUserPasswordData{
				{
					Hint:    "hint1",
					Content: "content1",
				},
				{
					Hint:    "hint2",
					Content: "content2",
				},
			}, nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		resp, err := aService.GetUserPasswordDataList(context.Background(), &proto.UserPasswordDataRequest{})
		require.NoError(t, err)
		require.Equal(t, 2, len(resp.Datas))

		assert.Equal(t, "hint1", resp.Datas[0].Hint)
		assert.Equal(t, "content1", resp.Datas[0].Content)

		assert.Equal(t, "hint2", resp.Datas[1].Hint)
		assert.Equal(t, "content2", resp.Datas[1].Content)
	})

	t.Run("get_list_err", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		testErr := errors.New("testErr")
		mockService.EXPECT().GetUserPasswordDataList(gomock.All()).DoAndReturn(func(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {
			return nil, testErr
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		resp, err := aService.GetUserPasswordDataList(context.Background(), &proto.UserPasswordDataRequest{})
		require.ErrorIs(t, err, testErr)
		require.Nil(t, resp)
	})

	t.Run("create_user_data_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().CreateUserPasswordData(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
			require.Equal(t, "hint1", data.Hint)
			require.Equal(t, "content1", data.Content)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.CreateUserPasswordData(context.Background(), &proto.CreateUserPasswordDataRequest{
			Hint:    "hint1",
			Content: "content1",
		})

		require.NoError(t, err)
	})

	t.Run("delete_user_data_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().DeleteUserPasswordData(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, hint string) error {
			require.Equal(t, "hint1", hint)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.DeleteUserPasswordData(context.Background(), &proto.DeleteUserPasswordDataRequest{
			Hint: "hint1",
		})

		require.NoError(t, err)
	})

	t.Run("update_user_data_ok", func(t *testing.T) {
		mockService := NewMockDataAccessor(ctrl)
		mockService.EXPECT().UpdateUserPasswordData(gomock.All(), gomock.All()).DoAndReturn(func(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
			require.Equal(t, "hint1", data.Hint)
			require.Equal(t, "content1", data.Content)
			return nil
		}).Times(1)

		aService := handler.NewDataAccessor(mockService)
		_, err := aService.UpdateUserPasswordData(context.Background(), &proto.UpdateUserPasswordDataRequest{
			Hint:    "hint1",
			Content: "content1",
		})

		require.NoError(t, err)
	})
}
