package app_test

import (
	"context"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPi(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPinger := NewMockPinger(ctrl)
	mockPinger.EXPECT().Ping(gomock.Any()).Return(nil).Times(1)

	p := app.NewPinger()
	p.SetPinger(mockPinger)

	err := p.Ping(context.Background())
	require.NoError(t, err)
}
