package postgres_test

import (
	"context"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/server/adapters/storage/postgres"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("debug_fn", func(t *testing.T) {

		mockLoger := NewMockLogger(ctrl)
		mockLoger.EXPECT().Debugw(gomock.Any(), gomock.Any()).Times(2).
			Do(func(msg string, keysAndValues ...any) {
				assert.Equal(t, "hello", msg)
				assert.Equal(t, "1", keysAndValues[0])
				assert.Equal(t, 2, keysAndValues[1])
				assert.Equal(t, "2", keysAndValues[2])
				assert.Equal(t, "3", keysAndValues[3])
			})

		adapter := postgres.NewLogAdapter(mockLoger)

		ctx := context.Background()

		adapter.Log(ctx, tracelog.LogLevelDebug, "hello", map[string]any{
			"1": 2,
			"2": "3",
		})

		adapter.Log(ctx, tracelog.LogLevelTrace, "hello", map[string]any{
			"1": 2,
			"2": "3",
		})
	})

	t.Run("debug_fn", func(t *testing.T) {

		mockLoger := NewMockLogger(ctrl)
		mockLoger.EXPECT().Infow(gomock.Any(), gomock.Any()).Times(1).
			Do(func(msg string, keysAndValues ...any) {
				assert.Equal(t, "hello", msg)
				assert.Equal(t, "1", keysAndValues[0])
				assert.Equal(t, 2, keysAndValues[1])
				assert.Equal(t, "2", keysAndValues[2])
				assert.Equal(t, "3", keysAndValues[3])
			})

		adapter := postgres.NewLogAdapter(mockLoger)

		ctx := context.Background()

		adapter.Log(ctx, tracelog.LogLevelInfo, "hello", map[string]any{
			"1": 2,
			"2": "3",
		})
	})

	t.Run("debug_fn", func(t *testing.T) {

		mockLoger := NewMockLogger(ctrl)
		mockLoger.EXPECT().Errorw(gomock.Any(), gomock.Any()).Times(1).
			Do(func(msg string, keysAndValues ...any) {
				assert.Equal(t, "hello", msg)
				assert.Equal(t, "1", keysAndValues[0])
				assert.Equal(t, 2, keysAndValues[1])
				assert.Equal(t, "2", keysAndValues[2])
				assert.Equal(t, "3", keysAndValues[3])
			})

		adapter := postgres.NewLogAdapter(mockLoger)

		ctx := context.Background()

		adapter.Log(ctx, tracelog.LogLevelError, "hello", map[string]any{
			"1": 2,
			"2": "3",
		})
	})
}
