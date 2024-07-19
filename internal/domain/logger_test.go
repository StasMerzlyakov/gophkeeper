package domain_test

import (
	"context"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	lg := domain.GetApplicationLogger()
	lg.Info("action", "err", "msg")

	logger, err := zap.NewProduction()
	require.NoError(t, err)

	sugarLog := logger.Sugar()
	domain.SetApplicationLogger(sugarLog)
	lg.Info("action", "err", "msg")
}

func TestEnrichContextRequestID(t *testing.T) {

	requestUUID := uuid.New()
	reqStr := requestUUID.String()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockLogger(ctrl)

	testLoggerFn := func(msg string, keysAndValues ...any) {
		// Проверяем что что при вызове метода логирования добавляется информация о пользователе и requstId
		requestIDIsChecked := false

		for id, v := range keysAndValues {
			switch v := v.(type) {
			case string:
				if v == domain.LoggerKeyRequestID {
					require.True(t, id+1 < len(keysAndValues), "requestID is not set")
					k := keysAndValues[id+1]
					id, ok := k.(string)
					require.True(t, ok, "requestID is not string")
					require.Equal(t, reqStr, id, "unexpecred requestID value")
					requestIDIsChecked = true
				}
			}
		}
		require.Truef(t, requestIDIsChecked, "requestID is not specified")
	}

	m.EXPECT().Debugw(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	m.EXPECT().Infow(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	m.EXPECT().Errorw(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	ctx := context.Background()

	enrichedCtx := domain.EnrichWithRequestIDLogger(ctx, requestUUID, m)

	log := domain.GetCtxLogger(enrichedCtx)

	log.Errorw("test errorw", "msg", "hello")
	log.Infow("test errorw", "msg", "hello")
	log.Debugw("test errorw", "msg", "hello")
}

func TestEnrichContextUserID(t *testing.T) {

	userID := domain.UserID(uuid.NewString())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockLogger(ctrl)

	testLoggerFn := func(msg string, keysAndValues ...any) {
		// Проверяем что что при вызове метода логирования добавляется информация о пользователе и requstId
		userIDChecked := false

		for id, v := range keysAndValues {
			switch v := v.(type) {
			case string:
				if v == domain.LoggerKeyUserID {
					require.True(t, id+1 < len(keysAndValues), "userID is not set")
					k := keysAndValues[id+1]
					id, ok := k.(string)
					require.True(t, ok, "userID is not string")
					require.Equal(t, string(userID), id, "unexpecred userID value")
					userIDChecked = true
				}
			}
		}
		require.Truef(t, userIDChecked, "userID is not specified")
	}

	m.EXPECT().Debugw(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	m.EXPECT().Infow(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	m.EXPECT().Errorw(gomock.Any(), gomock.Any()).DoAndReturn(testLoggerFn).AnyTimes()

	ctx := context.Background()

	enrichedCtx := domain.EnrichWithUserIDLogger(ctx, userID, m)

	log := domain.GetCtxLogger(enrichedCtx)

	log.Errorw("test errorw", "msg", "hello")
	log.Infow("test errorw", "msg", "hello")
	log.Debugw("test errorw", "msg", "hello")
}

func TestEnrichWithUserID(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		userID := domain.UserID(uuid.NewString())
		extCtx := domain.EnrichWithUserID(ctx, userID)

		userIDCtx, err := domain.GetUserID(extCtx)
		require.NoError(t, err)
		assert.Equal(t, userID, userIDCtx)
	})

	t.Run("ok", func(t *testing.T) {
		_, err := domain.GetUserID(ctx)
		require.ErrorIs(t, err, domain.ErrNotAuthorized)
	})

	t.Run("wrong_data", func(t *testing.T) {
		resultCtx := context.WithValue(ctx, domain.UserIDKey, "asd")
		_, err := domain.GetUserID(resultCtx)
		require.ErrorIs(t, err, domain.ErrServerInternal)
	})

}
