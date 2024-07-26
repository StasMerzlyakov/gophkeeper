package handler_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/grpc/handler"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/StasMerzlyakov/gophkeeper/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gp "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestHandler(t *testing.T) {
	caFile := filepath.Join(TestDataDirectory, "test-ca-cert.pem")

	conf := &config.ClientConf{
		CACert:            caFile,
		ServerAddress:     fmt.Sprintf("localhost:%d", srvPort),
		InterationTimeout: 3 * time.Second,
	}

	hnd, err := handler.NewHandler(conf)
	require.NoError(t, err)
	defer hnd.Stop()

	t.Run("checkEMail_ok", func(t *testing.T) {
		email := "email"

		rgHandler.checkEMailFn = func(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error) {
			require.Equal(t, email, req.Email)
			return &proto.CheckEMailResponse{
				Status: proto.CheckEMailResponse_AVAILABLE,
			}, nil
		}
		ctx := context.Background()
		res, err := hnd.CheckEMail(ctx, email)
		require.NoError(t, err)

		assert.Equal(t, domain.EMailAvailable, res)
	})

	t.Run("checkEMail_busy", func(t *testing.T) {
		email := "email"

		rgHandler.checkEMailFn = func(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error) {
			require.Equal(t, email, req.Email)
			return &proto.CheckEMailResponse{
				Status: proto.CheckEMailResponse_BUSY,
			}, nil
		}
		ctx := context.Background()
		res, err := hnd.CheckEMail(ctx, email)
		require.NoError(t, err)

		assert.Equal(t, domain.EMailBusy, res)
	})

	t.Run("checkEMail_err", func(t *testing.T) {
		email := "email"
		rgHandler.checkEMailFn = func(ctx context.Context, req *proto.CheckEMailRequest) (*proto.CheckEMailResponse, error) {
			require.Equal(t, email, req.Email)
			return nil, status.Error(gp.Internal, "")
		}
		ctx := context.Background()
		_, err := hnd.CheckEMail(ctx, email)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, gp.Internal, st.Code())
	})

	t.Run("registrate_ok", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "email",
			Password: "password",
		}
		sessionID := "sessionID"

		rgHandler.registrateFn = func(ctx context.Context, req *proto.RegistrationRequest) (*proto.RegistrationResponse, error) {
			assert.Equal(t, data.EMail, req.Email)
			assert.Equal(t, data.Password, req.Password)
			return &proto.RegistrationResponse{
				SessionId: sessionID,
			}, nil
		}
		ctx := context.Background()
		err := hnd.Registrate(ctx, data)
		require.NoError(t, err)
		assert.Equal(t, sessionID, hnd.SessionID())
	})

	t.Run("passRegOTP_ok", func(t *testing.T) {
		otpPass := "optPass"
		currentSessID := "sessionID"

		sessionID := "sessionID"

		hnd.SetSessionID(currentSessID)

		rgHandler.passOTPFn = func(ctx context.Context, req *proto.PassOTPRequest) (*proto.PassOTPResponse, error) {
			assert.Equal(t, otpPass, req.Password)
			assert.Equal(t, currentSessID, req.SessionId)
			return &proto.PassOTPResponse{
				SessionId: sessionID,
			}, nil
		}

		ctx := context.Background()
		err := hnd.PassRegOTP(ctx, otpPass)
		require.NoError(t, err)
		assert.Equal(t, sessionID, hnd.SessionID())
	})

	t.Run("initMasterKey_ok", func(t *testing.T) {
		mKey := &domain.MasterKeyData{
			EncryptedMasterKey: "EncryptedMasterKey",
			MasterKeyHint:      "MasterKeyHint",
			HelloEncrypted:     "HelloEncrypted",
		}

		currentSessID := "sessionID"

		hnd.SetSessionID(currentSessID)

		rgHandler.setMasterKeyFn = func(ctx context.Context, req *proto.MasterKeyRequest) (*proto.MasterKeyResponse, error) {
			assert.Equal(t, req.SessionId, currentSessID)
			assert.Equal(t, req.EncryptedMasterKey, mKey.EncryptedMasterKey)
			assert.Equal(t, req.MasterKeyPassHint, mKey.MasterKeyHint)
			assert.Equal(t, req.HelloEncrypted, mKey.HelloEncrypted)
			return &proto.MasterKeyResponse{}, nil
		}

		ctx := context.Background()
		err := hnd.InitMasterKey(ctx, mKey)
		require.NoError(t, err)
		assert.Equal(t, "", hnd.SessionID())
		assert.Equal(t, "", hnd.JWTToken())
	})

	t.Run("login_ok", func(t *testing.T) {
		data := &domain.EMailData{
			EMail:    "email",
			Password: "password",
		}

		sessionId := "sessionID"

		hnd.SetSessionID("")

		athService.loginFn = func(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
			assert.Equal(t, data.EMail, req.Email)
			assert.Equal(t, data.Password, req.Password)

			return &proto.LoginResponse{
				SessionId: sessionId,
			}, nil
		}

		ctx := context.Background()
		err := hnd.Login(ctx, data)
		require.NoError(t, err)
		assert.Equal(t, sessionId, hnd.SessionID())
		assert.Equal(t, "", hnd.JWTToken())
	})

	t.Run("passRegOTP_ok", func(t *testing.T) {
		otpPass := "optPass"
		currentSessID := "sessionID"

		hnd.SetSessionID(currentSessID)
		jwtTolen := "token"

		athService.passOTPFn = func(ctx context.Context, req *proto.PassOTPRequest) (*proto.AuthResponse, error) {
			assert.Equal(t, otpPass, req.Password)
			assert.Equal(t, currentSessID, req.SessionId)
			return &proto.AuthResponse{
				Token: jwtTolen,
			}, nil
		}

		ctx := context.Background()
		err := hnd.PassLoginOTP(ctx, otpPass)
		require.NoError(t, err)
		assert.Equal(t, "", hnd.SessionID())
		assert.Equal(t, jwtTolen, hnd.JWTToken())
	})

	t.Run("dtAccessor_ok", func(t *testing.T) {

		jwtTolen := "token"

		resp := &proto.HelloResponse{
			HelloEncrypted:     "HelloEncrypted",
			EncryptedMasterKey: "EncryptedMasterKey",
			MasterKeyPassHint:  "MasterKeyPassHint",
		}
		dtAccessor.helloFn = func(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			require.True(t, ok)
			values := md[domain.AuthorizationMetadataTokenName]
			require.Equal(t, 1, len(values))
			assert.Equal(t, jwtTolen, values[0])
			return resp, nil
		}

		ctx := context.Background()
		hnd.SetJWTToken(jwtTolen)
		data, err := hnd.GetHelloData(ctx)
		require.NoError(t, err)
		assert.Equal(t, "", hnd.SessionID())
		assert.Equal(t, jwtTolen, hnd.JWTToken())

		assert.Equal(t, resp.HelloEncrypted, data.HelloEncrypted)
		assert.Equal(t, resp.EncryptedMasterKey, data.EncryptedMasterKey)
		assert.Equal(t, resp.MasterKeyPassHint, data.MasterKeyPassHint)
	})
}
