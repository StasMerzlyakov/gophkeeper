package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAppController(conf *config.ClientConf) *appController {
	return &appController{
		conf:        conf,
		status:      domain.ClientStatusOnline,
		loginer:     NewLoginer(),
		pinger:      NewPinger(),
		registrator: NewRegistrator(),
	}
}

func (ac *appController) Start() {
	ac.wg.Add(1)
	ac.exitChan = make(chan struct{}, 1)
	go func() {
		defer ac.wg.Done()
		for {
			select {
			case <-time.After(2 * ac.conf.InterationTimeout):
				ac.CheckServerStatus()
			case <-ac.exitChan:
				return
			}
		}
	}()

}

func (ac *appController) CheckServerStatus() {

	ac.invokeFn(func(ctx context.Context) {
		err := ac.pinger.Ping(ctx)
		if err == nil {
			select {
			case <-ctx.Done():
				// timeout
				return
			default:
				ac.status = domain.ClientStatusOnline
				ac.infoView.ShowMsg("Server is online")
			}
		}
	}, domain.ClientStatusOffline, "", false)
}

func (ac *appController) Stop() {
	close(ac.exitChan)
	ac.server.Stop()
	ac.wg.Wait()
}

func (ac *appController) SetInfoView(view InfoView) *appController {
	ac.infoView = view
	ac.loginer.LoginView(view)
	ac.registrator.RegView(view)
	return ac
}

func (ac *appController) SetServer(server Server) *appController {
	ac.pinger.SetPinger(server)
	ac.loginer.LoginSever(server)
	ac.registrator.RegServer(server)
	ac.server = server
	return ac
}

type appController struct {
	conf        *config.ClientConf
	status      domain.ClientStatus
	infoView    InfoView
	loginer     *loginer
	pinger      *pinger
	registrator *registrator
	server      Server
	exitChan    chan struct{}
	wg          sync.WaitGroup
}

func (ac *appController) GetStatus() domain.ClientStatus {
	return ac.status
}

func (ac *appController) invokeFn(fn func(ctx context.Context), runStatus domain.ClientStatus, msg string, showErr bool) {
	if ac.status != runStatus {
		if msg != "" {
			ac.infoView.ShowMsg(msg)
		}
	} else {
		timedCtx, cancelCtxFn := context.WithTimeout(context.Background(), ac.conf.InterationTimeout)
		defer cancelCtxFn()
		ac.wg.Add(1)
		doneCh := make(chan struct{}, 1)
		go func() {
			defer ac.wg.Done()
			fn(timedCtx)
			select {
			case <-timedCtx.Done():
				return // timeout
			default:
				doneCh <- struct{}{}
			}
		}()

		select {
		case <-timedCtx.Done():
			ac.status = domain.ClientStatusOffline
			if showErr {
				ac.infoView.ShowError(fmt.Errorf("%w server timeout", domain.ErrClientServerTimeout))
			}
		case <-ac.exitChan:
			return
		case <-doneCh:
			return
		}
	}
}

func (ac *appController) LoginEMail(data *domain.EMailData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.Login(ctx, data)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) LoginPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.PassOTP(ctx, otpPass)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) LoginCheckMasterKey(masterKeyPassword string) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.CheckMasterKey(ctx, masterKeyPassword)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegEMail(data *domain.EMailData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.Registrate(ctx, data)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.PassOTP(ctx, otpPass)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.InitMasterKey(ctx, mKey)
	}, domain.ClientStatusOnline, "server is offline", true)
}
