package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewLoginer() *loginer {
	return &loginer{}
}

func (lg *loginer) LoginSever(logSrv AppServer) *loginer {
	lg.logSrv = logSrv
	return lg
}

func (lg *loginer) LoginView(logView AppView) *loginer {
	lg.logView = logView
	return lg
}

func (lg *loginer) LoginHelper(helper DomainHelper) *loginer {
	lg.helper = helper
	return lg
}

func (lg *loginer) LoginStorage(storage AppStorage) *loginer {
	lg.storage = storage
	return lg
}

type loginer struct {
	logSrv  AppServer
	logView AppView
	helper  DomainHelper
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

func (lg *loginer) CheckMasterKey(ctx context.Context, masterPassword string) {
	log.Debugf("checkMasterKey start")
	helloData, err := lg.logSrv.GetHelloData(ctx)
	if err != nil {
		log.Errorf("checkMasterKey err - getHelloData %v", err.Error())
		lg.logView.ShowError(err)
		return
	}

	err = lg.helper.DecryptHello(masterPassword, helloData.HelloEncrypted)
	if err != nil {
		log.Errorf("checkMasterKey err - decryptHello %v", err.Error())
		lg.logView.ShowError(err)
		lg.logView.ShowMasterKeyView(helloData.MasterPasswordHint)
		return
	}

	log.Debugf("hello decrypted success")

	select {
	case <-ctx.Done():
		log.Errorf("checkMasterKey err - context done %v", ctx.Err().Error())
		return
	default:
		log.Debugf("checkMasterKey success")
		lg.storage.SetMasterPassword(masterPassword)
		lg.logView.ShowDataAccessView()
	}
}
