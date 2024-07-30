package app_test

import (
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("start_stop", func(t *testing.T) {
		mockSrv := NewMockAppServer(ctrl)
		mockSrv.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
		mockSrv.EXPECT().Stop().Times(1)
		mockSrv.EXPECT().Start().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		appController := app.NewAppController(conf).SetServer(mockSrv)

		appController.Start()
		time.Sleep(3 * time.Second)
		appController.Stop()
	})

	t.Run("start_stop_restore", func(t *testing.T) {

		mockSrv := NewMockAppServer(ctrl)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}

		mockSrv.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
			return nil
		}).AnyTimes()

		mockSrv.EXPECT().Stop().Times(1)
		mockSrv.EXPECT().Start().Times(1)

		data := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EMailData) error {
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return nil
		})

		mockView := NewMockAppView(ctrl)

		mockView.EXPECT().ShowLogOTPView().Times(1)

		appController := app.NewAppController(conf).SetServer(mockSrv).SetInfoView(mockView)

		appController.Start()

		appController.LoginEMail(data)

		mockView.EXPECT().ShowMasterKeyView(gomock.Any()).Do(func(hint string) {
			assert.Empty(t, hint)
		}).Times(1)

		mockSrv.EXPECT().PassLoginOTP(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, pas string) error {
			assert.Equal(t, "pass", pas)
			return nil
		})

		appController.LoginPassOTP(&domain.OTPPass{
			Pass: "pass",
		})

		time.Sleep(1 * time.Second)

		appController.Stop()
	})
}
