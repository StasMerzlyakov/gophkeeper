package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewLoginer(conf *config.ClientConf) *loginer {
	return &loginer{
		conf: conf,
	}
}

func (lg *loginer) LoginSever(logSrv LoginServer) *loginer {
	lg.logSrv = logSrv
	return lg
}

func (lg *loginer) LoginView(logView LoginView) *loginer {
	lg.logView = logView
	return lg
}

func (lg *loginer) LoginHelper(helper LoginHelper) *loginer {
	lg.helper = helper
	return lg
}

func (lg *loginer) LoginStorage(storage LoginStorage) *loginer {
	lg.storage = storage
	return lg
}

type loginer struct {
	conf    *config.ClientConf
	logSrv  LoginServer
	logView LoginView
	helper  LoginHelper
	storage LoginStorage
}

func (lg *loginer) Login(ctx context.Context, data *domain.EMailData) {
	timedCtx, fn := context.WithTimeout(ctx, lg.conf.InterationTimeout)
	defer fn()

	if err := lg.logSrv.Login(timedCtx, data); err != nil {
		lg.logView.ShowError(err)
		return
	}

	lg.logView.ShowLogOTPView()
}

func (lg *loginer) PassOTP(ctx context.Context, otpPass *domain.OTPPass) {
	timedCtx, fn := context.WithTimeout(ctx, lg.conf.InterationTimeout)
	defer fn()

	if err := lg.logSrv.PassLoginOTP(timedCtx, otpPass.Pass); err != nil {
		lg.logView.ShowError(err)
		return
	}

	lg.logView.ShowMasterKeyView("")

}

func (lg *loginer) CheckMasterKey(ctx context.Context, masterKeyPassword string) {
	timedCtx, fn := context.WithTimeout(ctx, lg.conf.InterationTimeout)
	defer fn()

	helloData, err := lg.logSrv.GetHelloData(timedCtx)
	if err != nil {
		lg.logView.ShowError(err)
		return
	}

	masterKey, err := lg.helper.DecryptMasterKey(masterKeyPassword, helloData.EncryptedMasterKey)
	if err != nil {
		lg.logView.ShowError(err)
		return
	}

	helloDecrypted, err := lg.helper.DecryptShortData(helloData.HelloEncrypted, masterKey)
	if err != nil {
		lg.logView.ShowError(err)
		return
	}

	ok, err := lg.helper.CheckHello(string(helloDecrypted))
	if err != nil {
		lg.logView.ShowError(err)
		return
	}

	if !ok {
		lg.logView.ShowError(fmt.Errorf("%w check master key passord", domain.ErrAuthDataIncorrect))
		lg.logView.ShowMasterKeyView(helloData.MasterKeyPassHint)
		return
	}

	lg.storage.SetMasterKey(masterKey)
	lg.logView.ShowDataAccessView()
}
