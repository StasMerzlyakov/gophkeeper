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

func TestCheckEMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("avalable", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)
		mockService.EXPECT().GetEMailStatus(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, email string) (domain.EMailStatus, error) {
				return domain.EMailAvailable, nil
			})

		regHandler := handler.NewRegHandler(mockService)

		res, err := regHandler.CheckEMail(context.Background(), &proto.CheckEMailRequest{
			Email: "email",
		})

		require.NoError(t, err)
		require.Equal(t, proto.CheckEMailResponse_AVAILABLE, res.Status)
	})

	t.Run("avalable", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)
		mockService.EXPECT().GetEMailStatus(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, email string) (domain.EMailStatus, error) {
				return domain.EMailBusy, nil
			})

		regHandler := handler.NewRegHandler(mockService)

		res, err := regHandler.CheckEMail(context.Background(), &proto.CheckEMailRequest{
			Email: "email",
		})

		require.NoError(t, err)
		require.Equal(t, proto.CheckEMailResponse_BUSY, res.Status)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)
		testErr := errors.New("testErr")
		mockService.EXPECT().GetEMailStatus(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, email string) (domain.EMailStatus, error) {
				return domain.EMailBusy, testErr
			})

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.CheckEMail(context.Background(), &proto.CheckEMailRequest{
			Email: "email",
		})

		require.ErrorIs(t, err, testErr)
	})

	t.Run("wrong status", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		mockService.EXPECT().GetEMailStatus(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, email string) (domain.EMailStatus, error) {
				return domain.EMailStatus("3"), nil
			})

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.CheckEMail(context.Background(), &proto.CheckEMailRequest{
			Email: "email",
		})

		require.ErrorIs(t, err, domain.ErrServerInternal)
	})
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		mockService.EXPECT().Registrate(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
				return domain.SessionID("sessionID"), nil
			})

		regHandler := handler.NewRegHandler(mockService)

		resp, err := regHandler.Registrate(context.Background(), &proto.RegistrationRequest{
			Email:    "email",
			Password: "password",
		})

		require.NoError(t, err)
		require.Equal(t, "sessionID", resp.SessionId)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		testErr := errors.New("testErr")
		mockService.EXPECT().Registrate(gomock.Any(), gomock.Any()).Times(1).
			DoAndReturn(func(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
				return domain.SessionID(""), testErr
			})

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.Registrate(context.Background(), &proto.RegistrationRequest{
			Email:    "email",
			Password: "password",
		})

		require.ErrorIs(t, err, testErr)
	})
}

func TestRegPassOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)
		mockService.EXPECT().PassOTP(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error) {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "otpPass", otpPass)
				return domain.SessionID("sessionID"), nil
			}).Times(1)

		regHandler := handler.NewRegHandler(mockService)

		res, err := regHandler.PassOTP(context.Background(), &proto.PassOTPRequest{
			SessionId: "currentID",
			Password:  "otpPass",
		})

		require.Nil(t, err)
		require.Equal(t, "sessionID", res.SessionId)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		testErr := errors.New("testErr")
		mockService.EXPECT().PassOTP(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error) {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "otpPass", otpPass)
				return domain.SessionID(""), testErr
			}).Times(1)

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.PassOTP(context.Background(), &proto.PassOTPRequest{
			SessionId: "currentID",
			Password:  "otpPass",
		})

		require.ErrorIs(t, err, testErr)
	})
}

func TestSetMasterKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		mockService.EXPECT().InitMasterKey(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, mKey *domain.MasterKeyData) error {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "EncryptedMasterKey", mKey.EncryptedMasterKey)
				require.Equal(t, "HelloEncrypted", mKey.HelloEncrypted)
				require.Equal(t, "MasterKeyHint", mKey.MasterKeyHint)
				return nil
			}).Times(1)

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.SetMasterKey(context.Background(), &proto.MasterKeyRequest{
			SessionId:          "currentID",
			EncryptedMasterKey: "EncryptedMasterKey",
			MasterKeyPassHint:  "MasterKeyHint",
			HelloEncrypted:     "HelloEncrypted",
		})

		require.Nil(t, err)
	})

	t.Run("err", func(t *testing.T) {
		mockService := NewMockRegistrator(ctrl)

		testErr := errors.New("testErr")
		mockService.EXPECT().InitMasterKey(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, currentID domain.SessionID, mKey *domain.MasterKeyData) error {
				require.Equal(t, domain.SessionID("currentID"), currentID)
				require.Equal(t, "EncryptedMasterKey", mKey.EncryptedMasterKey)
				require.Equal(t, "HelloEncrypted", mKey.HelloEncrypted)
				require.Equal(t, "MasterKeyHint", mKey.MasterKeyHint)
				return testErr
			}).Times(1)

		regHandler := handler.NewRegHandler(mockService)

		_, err := regHandler.SetMasterKey(context.Background(), &proto.MasterKeyRequest{
			SessionId:          "currentID",
			EncryptedMasterKey: "EncryptedMasterKey",
			MasterKeyPassHint:  "MasterKeyHint",
			HelloEncrypted:     "HelloEncrypted",
		})

		require.ErrorIs(t, err, testErr)
	})
}
