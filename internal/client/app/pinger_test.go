package app_test

import (
	"context"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	gomock "github.com/golang/mock/gomock"
)

func TestPi(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPinger := NewMockPinger(ctrl)
	mockPinger.EXPECT().Ping(gomock.Any()).DoAndReturn(nil).Times(1)

	p := app.NewPinger()
	p.SetPinger(mockPinger)

	p.Ping(context.Background())
}
