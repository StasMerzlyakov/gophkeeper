package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewLoginer() *loginer {
	return &loginer{}
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
	logSrv  LoginServer
	logView LoginView
	helper  LoginHelper
	storage LoginStorage
}

func (lg *loginer) Login(ctx context.Context, data *domain.EMailData) {
	if err := lg.logSrv.Login(ctx, data); err != nil {
		lg.logView.ShowError(err)
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		lg.logView.ShowLogOTPView()
	}
}

func (lg *loginer) PassOTP(ctx context.Context, otpPass *domain.OTPPass) {
	if err := lg.logSrv.PassLoginOTP(ctx, otpPass.Pass); err != nil {
		lg.logView.ShowError(err)
		return
	}

	/*select {
	case <-ctx.Done():
		return
	default: */
	lg.logView.ShowMasterKeyView("")
	//}
}

func (lg *loginer) CheckMasterKey(ctx context.Context, masterKeyPassword string) {
	helloData, err := lg.logSrv.GetHelloData(ctx)
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

	select {
	case <-ctx.Done():
		return
	default:
		lg.storage.SetMasterKey(masterKey)
		lg.logView.ShowDataAccessView()
	}
}
