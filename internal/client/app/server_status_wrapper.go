package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewStatusWrapper(clientConf *config.ClientConf, server AppServer) *serverStatusWrapper {
	return &serverStatusWrapper{
		server: server,
		conf:   clientConf,
		status: domain.ClientStatusOffline,
	}
}

// serverStatusWrapper wrapp server communication object and hold connection status.
type serverStatusWrapper struct {
	server   AppServer
	conf     *config.ClientConf
	status   domain.ClientStatus
	stopChan chan struct{}
	wg       sync.WaitGroup
}

var _ AppServer = (*serverStatusWrapper)(nil)

func (aw *serverStatusWrapper) Start() {
	aw.wg.Add(1)
	aw.stopChan = make(chan struct{}, 1)
	go func() {
		defer aw.wg.Done()
		_ = aw.invokeFn(context.Background(), aw.ping) // Start ping immediatly
		for {
			select {
			case <-time.After(2 * aw.conf.InterationTimeout):
				_ = aw.invokeFn(context.Background(), aw.ping)
			case <-aw.stopChan:
				return
			}
		}
	}()
}

func (aw *serverStatusWrapper) GetStatus() domain.ClientStatus {
	return aw.status
}

func (aw *serverStatusWrapper) ping(ctx context.Context) error {
	log := GetMainLogger()
	log.Debug("ping start")
	if err := aw.server.Ping(ctx); err != nil {
		aw.status = domain.ClientStatusOffline
		log.Warn("ping err - server is not available")
		return err
	} else {
		log.Debug("server is online")
		aw.status = domain.ClientStatusOnline
		return nil
	}
}

func (aw *serverStatusWrapper) invokeOnlineFn(ctx context.Context, fn func(ctx context.Context) error) error {
	if aw.status != domain.ClientStatusOnline {
		return fmt.Errorf("%w server is offline", domain.ErrServerIsNotResponding)
	}
	return aw.invokeFn(ctx, fn)
}

func (aw *serverStatusWrapper) invokeFn(ctx context.Context, fn func(ctx context.Context) error) error {
	timedCtx, cancelFn := context.WithTimeout(ctx, aw.conf.InterationTimeout)
	defer cancelFn()

	aw.wg.Add(1)
	resultCh := make(chan error, 1)
	go func() {
		defer aw.wg.Done()
		err := fn(timedCtx)
		resultCh <- err
	}()

	select {
	case err := <-resultCh:
		return err
	case <-aw.stopChan:
		// Applicaiton is stopped.
		return nil
	case <-timedCtx.Done():
		// timeout
		aw.status = domain.ClientStatusOffline
		log := GetMainLogger()
		log.Warn("timeout - server is not available")
		return fmt.Errorf("%w server is offline", domain.ErrServerIsNotResponding)
	}
}

func (aw *serverStatusWrapper) Stop() {
	if aw.stopChan != nil {
		aw.stopChan <- struct{}{}
	}
	aw.wg.Wait()
	if aw.server != nil {
		aw.server.Stop()
	}
}

func (aw *serverStatusWrapper) Ping(ctx context.Context) error {
	return aw.server.Ping(ctx)
}

func (aw *serverStatusWrapper) CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error) {
	var emailStatus domain.EMailStatus = domain.EMailBusy
	var err error
	fn := func(ctx context.Context) error {
		emailStatus, err = aw.server.CheckEMail(ctx, email)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return emailStatus, retErr
}

func (aw *serverStatusWrapper) Registrate(ctx context.Context, data *domain.EMailData) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.Registrate(ctx, data)
	})
}

func (aw *serverStatusWrapper) PassRegOTP(ctx context.Context, otpPass string) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.PassRegOTP(ctx, otpPass)
	})
}

func (aw *serverStatusWrapper) InitMasterKey(ctx context.Context, mKey *domain.MasterKeyData) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.InitMasterKey(ctx, mKey)
	})
}

func (aw *serverStatusWrapper) Login(ctx context.Context, data *domain.EMailData) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.Login(ctx, data)
	})
}

func (aw *serverStatusWrapper) PassLoginOTP(ctx context.Context, otpPass string) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.PassLoginOTP(ctx, otpPass)
	})
}

func (aw *serverStatusWrapper) GetHelloData(ctx context.Context) (*domain.HelloData, error) {
	var data *domain.HelloData
	var err error
	fn := func(ctx context.Context) error {
		data, err = aw.server.GetHelloData(ctx)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return data, retErr
}

func (aw *serverStatusWrapper) GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error) {
	var data []domain.EncryptedBankCard
	var err error
	fn := func(ctx context.Context) error {
		data, err = aw.server.GetBankCardList(ctx)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return data, retErr
}
func (aw *serverStatusWrapper) CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.CreateBankCard(ctx, bnkCard)
	})
}

func (aw *serverStatusWrapper) UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.UpdateBankCard(ctx, bnkCard)
	})
}

func (aw *serverStatusWrapper) DeleteBankCard(ctx context.Context, number string) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.DeleteBankCard(ctx, number)
	})
}

func (aw *serverStatusWrapper) GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error) {
	var data []domain.EncryptedUserPasswordData
	var err error
	fn := func(ctx context.Context) error {
		data, err = aw.server.GetUserPasswordDataList(ctx)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return data, retErr
}

func (aw *serverStatusWrapper) CreateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.CreateUserPasswordData(ctx, data)
	})
}
func (aw *serverStatusWrapper) UpdateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.UpdateUserPasswordData(ctx, data)
	})
}
func (aw *serverStatusWrapper) DeleteUserPasswordData(ctx context.Context, hint string) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.DeleteUserPasswordData(ctx, hint)
	})
}

func (aw *serverStatusWrapper) GetFileInfoList(ctx context.Context) ([]domain.FileInfo, error) {
	var data []domain.FileInfo
	var err error
	fn := func(ctx context.Context) error {
		data, err = aw.server.GetFileInfoList(ctx)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return data, retErr
}

func (aw *serverStatusWrapper) DeleteFileInfo(ctx context.Context, name string) error {
	return aw.invokeOnlineFn(ctx, func(ctx context.Context) error {
		return aw.server.DeleteFileInfo(ctx, name)
	})
}

func (aw *serverStatusWrapper) SendFile(ctx context.Context) (domain.StreamFileWriter, error) {

	var data domain.StreamFileWriter
	var err error
	fn := func(ctx context.Context) error {
		data, err = aw.server.SendFile(ctx)
		return err
	}
	retErr := aw.invokeOnlineFn(ctx, fn)
	return data, retErr
}
