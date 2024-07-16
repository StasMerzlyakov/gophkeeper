package domain_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
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
