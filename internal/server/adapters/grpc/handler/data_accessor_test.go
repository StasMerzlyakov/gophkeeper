package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/grpc/handler"
	gomock "github.com/golang/mock/gomock"
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
