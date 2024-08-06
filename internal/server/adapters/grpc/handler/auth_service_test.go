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

func TestLogin(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockAuthService(ctrl)

		mockService.EXPECT().Login(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
				require.Equal(t, "email", data.EMail)
				require.Equal(t, "password", data.Password)
				return domain.SessionID("sessionID"), nil
			}).Times(1)

		aService := handler.NewAuthService(mockService)

		res, err := aService.Login(context.Background(), &proto.LoginRequest{
			Email:    "email",
			Password: "password",
		})

		require.Nil(t, err)
		require.Equal(t, "sessionID", res.SessionId)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockAuthService(ctrl)

		testErr := errors.New("error")
		mockService.EXPECT().Login(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
				require.Equal(t, "email", data.EMail)
				require.Equal(t, "password", data.Password)
				return "", testErr
			}).Times(1)

		aService := handler.NewAuthService(mockService)

		_, err := aService.Login(context.Background(), &proto.LoginRequest{
			Email:    "email",
			Password: "password",
		})

		require.ErrorIs(t, err, testErr)
	})
}

func TestAuthPassOTP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockAuthService(ctrl)

		mockService.EXPECT().CheckOTP(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error) {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "otpPass", otpPass)
				return domain.JWTToken("token"), nil
			}).Times(1)

		aService := handler.NewAuthService(mockService)

		res, err := aService.PassOTP(context.Background(), &proto.PassOTPRequest{
			SessionId: "currentID",
			Password:  "otpPass",
		})

		require.Nil(t, err)
		require.Equal(t, "token", res.Token)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockAuthService(ctrl)

		testErr := errors.New("testErr")

		mockService.EXPECT().CheckOTP(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error) {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "otpPass", otpPass)
				return domain.JWTToken(""), testErr
			}).Times(1)

		aService := handler.NewAuthService(mockService)

		_, err := aService.PassOTP(context.Background(), &proto.PassOTPRequest{
			SessionId: "currentID",
			Password:  "otpPass",
		})

		require.ErrorIs(t, err, testErr)
	})
}
