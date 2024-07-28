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

		mockSrv := NewMockServer(ctrl)
		mockSrv.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
		mockSrv.EXPECT().Stop().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		appController := app.NewAppController(conf).SetServer(mockSrv)

		appController.Start()
		time.Sleep(3 * time.Second)
		appController.Stop()
	})

	t.Run("start_stop_restore", func(t *testing.T) {

		mockSrv := NewMockServer(ctrl)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}

		mockSrv.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
			//time.Sleep(2 * conf.InterationTimeout) // timeout
			return nil
		}).Times(1)

		mockSrv.EXPECT().Stop().Times(1)

		data := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EMailData) error {
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			time.Sleep(2 * conf.InterationTimeout) // timeout
			return nil
		})

		mockView := NewMockInfoView(ctrl)

		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, domain.ErrClientServerTimeout)
		}).Times(1)

		mockView.EXPECT().ShowMsg(gomock.Any()).Do(func(msg string) {
		}).Times(1)

		appController := app.NewAppController(conf).SetServer(mockSrv).SetInfoView(mockView)

		appController.Start()

		appController.LoginEMail(data)
		assert.Equal(t, domain.ClientStatusOffline, appController.GetStatus())

		time.Sleep(4 * time.Second)
		assert.Equal(t, domain.ClientStatusOnline, appController.GetStatus())

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
		appController.Stop()
	})

	t.Run("start_stop_offline", func(t *testing.T) {

		mockSrv := NewMockServer(ctrl)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}

		mockSrv.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
			time.Sleep(2 * conf.InterationTimeout) // timeout
			return nil
		}).MinTimes(1).MaxTimes(3)

		mockSrv.EXPECT().Stop().Times(1)

		data := &domain.EMailData{
			EMail:    "email",
			Password: "pass",
		}

		mockSrv.EXPECT().Login(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EMailData) error {
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			time.Sleep(2 * conf.InterationTimeout) // timeout
			return nil
		})

		mockView := NewMockInfoView(ctrl)

		mockView.EXPECT().ShowError(gomock.Any()).Do(func(err error) {
			assert.ErrorIs(t, err, domain.ErrClientServerTimeout)
		}).Times(1)

		mockView.EXPECT().ShowMsg(gomock.Any()).Do(func(msg string) {
			assert.Equal(t, "server is offline", msg)
		}).Times(1)

		appController := app.NewAppController(conf).SetServer(mockSrv).SetInfoView(mockView)

		appController.Start()

		appController.LoginEMail(data)
		assert.Equal(t, domain.ClientStatusOffline, appController.GetStatus())

		time.Sleep(4 * time.Second)
		assert.Equal(t, domain.ClientStatusOffline, appController.GetStatus())

		appController.LoginPassOTP(&domain.OTPPass{
			Pass: "pass",
		})
		appController.Stop()
	})

}
