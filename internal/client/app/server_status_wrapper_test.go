package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServerStatusWrapper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("stop_is_ok", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).Return(nil).Times(0)
		mockServ.EXPECT().Stop().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Stop()
	})

	t.Run("start_stop", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).Return(nil).Times(1)
		mockServ.EXPECT().Stop().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 2 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Start()
		time.Sleep(1 * time.Second)
		assert.Equal(t, domain.ClientStatusOnline, wrapper.GetStatus())
		wrapper.Stop()
	})

	t.Run("start_long ping", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
			time.Sleep(2 * time.Second) // long operation
			return nil
		}).Return(nil).Times(1)
		mockServ.EXPECT().Stop().Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Start()
		time.Sleep(1 * time.Second)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Stop()
	})
}

func TestWrapperRegistrationFn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("checkEmail", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
		mockServ.EXPECT().Stop().Times(1)

		mockServ.EXPECT().CheckEMail(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, email string) (domain.EMailStatus, error) {
			assert.Equal(t, "email", email)
			return domain.EMailAvailable, nil
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Start()
		time.Sleep(1 * time.Second)
		assert.Equal(t, domain.ClientStatusOnline, wrapper.GetStatus())
		resp, err := wrapper.CheckEMail(context.Background(), "email")
		assert.NoError(t, err)
		assert.Equal(t, domain.EMailAvailable, resp)

		wrapper.Stop()
	})

	t.Run("registrate", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
		mockServ.EXPECT().Stop().Times(1)

		data := &domain.EMailData{
			EMail:    "email",
			Password: "password",
		}

		mockServ.EXPECT().Registrate(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, dt *domain.EMailData) error {
			assert.Equal(t, data.EMail, dt.EMail)
			assert.Equal(t, data.Password, dt.Password)
			return nil
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Start()
		time.Sleep(1 * time.Second)
		assert.Equal(t, domain.ClientStatusOnline, wrapper.GetStatus())
		err := wrapper.Registrate(context.Background(), data)
		assert.NoError(t, err)

		wrapper.Stop()
	})

	t.Run("getHelloData", func(t *testing.T) {
		mockServ := NewMockAppServer(ctrl)
		mockServ.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
		mockServ.EXPECT().Stop().Times(1)

		retData := &domain.HelloData{
			HelloEncrypted:     "ASDvbasda1",
			MasterPasswordHint: "Hint",
		}
		mockServ.EXPECT().GetHelloData(gomock.Any()).DoAndReturn(func(ctx context.Context) (*domain.HelloData, error) {

			return retData, nil
		}).Times(1)

		conf := &config.ClientConf{
			InterationTimeout: 1 * time.Second,
		}
		wrapper := app.NewStatusWrapper(conf, mockServ)
		assert.Equal(t, domain.ClientStatusOffline, wrapper.GetStatus())
		wrapper.Start()
		time.Sleep(1 * time.Second)
		assert.Equal(t, domain.ClientStatusOnline, wrapper.GetStatus())
		resp, err := wrapper.GetHelloData(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, retData.HelloEncrypted, resp.HelloEncrypted)
		assert.Equal(t, retData.MasterPasswordHint, resp.MasterPasswordHint)

		wrapper.Stop()
	})

}
