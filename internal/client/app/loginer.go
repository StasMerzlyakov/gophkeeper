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

func (lg *loginer) LoginStorage(storage AppStorage) *loginer {
	lg.storage = storage
	return lg
}

type loginer struct {
	logSrv  LoginServer
	logView LoginView
	helper  LoginHelper
	storage AppStorage
}

func (lg *loginer) Login(ctx context.Context, data *domain.EMailData) {
	log := GetMainLogger()

	log.Debugf("login %v start", data.EMail)

	if err := lg.logSrv.Login(ctx, data); err != nil {
		lg.logView.ShowError(err)
		return
	}

	select {
	case <-ctx.Done():
		log.Errorf("context done %v", ctx.Err().Error())
		return
	default:
		log.Debugf("login %v succes", data.EMail)
		lg.logView.ShowLogOTPView()
	}
}

func (lg *loginer) PassOTP(ctx context.Context, otpPass *domain.OTPPass) {
	log := GetMainLogger()
	log.Debugf("passOTP start")
	if err := lg.logSrv.PassLoginOTP(ctx, otpPass.Pass); err != nil {
		lg.logView.ShowError(err)
		return
	}

	select {
	case <-ctx.Done():
		log.Errorf("context done %v", ctx.Err().Error())
		return
	default:
		log.Debugf("passOTP success")
		lg.logView.ShowMasterKeyView("")
	}
}

func (lg *loginer) CheckMasterKey(ctx context.Context, masterKeyPassword string) {
	log.Debugf("checkMasterKey start")
	helloData, err := lg.logSrv.GetHelloData(ctx)
	if err != nil {
		log.Errorf("checkMasterKey err - getHelloData %v", err.Error())
		lg.logView.ShowError(err)
		return
	}

	masterKey, err := lg.helper.DecryptMasterKey(masterKeyPassword, helloData.EncryptedMasterKey)
	if err != nil {
		log.Errorf("checkMasterKey err - decryptMasterKey %v", err.Error())
		lg.logView.ShowError(err)
		return
	}

	log.Debugf("masterKey decrypted")

	helloDecrypted, err := lg.helper.DecryptShortData(helloData.HelloEncrypted, masterKey)
	if err != nil {
		log.Errorf("checkMasterKey err - decryptShortData %v", err.Error())
		lg.logView.ShowError(err)
		return
	}

	log.Debugf("hello decrypted")

	ok, err := lg.helper.CheckHello(string(helloDecrypted))
	if err != nil {
		log.Errorf("checkMasterKey err - checkHello %v", err.Error())
		lg.logView.ShowError(err)
		return
	}

	if !ok {
		lg.logView.ShowError(fmt.Errorf("%w check master key passord", domain.ErrAuthDataIncorrect))
		lg.logView.ShowMasterKeyView(helloData.MasterKeyPassHint)
		return
	}

	log.Debugf("checkHello success")

	select {
	case <-ctx.Done():
		log.Errorf("context done %v", ctx.Err().Error())
		return
	default:
		log.Debugf("checkMasterKey success")
		lg.storage.SetMasterKey(masterKey)
		lg.logView.ShowDataAccessView()
	}
}
