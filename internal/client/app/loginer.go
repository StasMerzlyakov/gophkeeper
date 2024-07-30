package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewLoginer() *loginer {
	return &loginer{}
}

func (lg *loginer) LoginSever(logSrv AppServer) *loginer {
	lg.logSrv = logSrv
	return lg
}

func (lg *loginer) LoginHelper(helper DomainHelper) *loginer {
	lg.helper = helper
	return lg
}

type loginer struct {
	logSrv AppServer
	helper DomainHelper
}

func (lg *loginer) Login(ctx context.Context, data *domain.EMailData) error {
	log := GetMainLogger()

	log.Debugf("login %v start", data.EMail)

	if err := lg.logSrv.Login(ctx, data); err != nil {
		err := fmt.Errorf("%w - login error", err)
		log.Warn(err.Error())
		return err
	}
	log.Debugf("login %v success", data.EMail)
	return nil
}

func (lg *loginer) PassOTP(ctx context.Context, otpPass *domain.OTPPass) error {
	log := GetMainLogger()
	log.Debugf("passOTP start")
	if err := lg.logSrv.PassLoginOTP(ctx, otpPass.Pass); err != nil {
		err := fmt.Errorf("%w - passOTP error", err)
		log.Warn(err.Error())
		return err
	}
	log.Debugf("passOTP success")
	return nil
}

// CheckMasterKey return nil, "" if ok;  error and hint if error and hint available
func (lg *loginer) CheckMasterKey(ctx context.Context, masterPassword string) (error, string) {
	log.Debug("checkMasterKey start")
	helloData, err := lg.logSrv.GetHelloData(ctx)
	if err != nil {
		err := fmt.Errorf("%w - checkMasterKey err", err)
		log.Warn(err.Error())
		return err, ""
	}

	err = lg.helper.DecryptHello(masterPassword, helloData.HelloEncrypted)
	if err != nil {
		err := fmt.Errorf("%w checkMasterKey err", err)
		log.Warn(err.Error())
		return err, helloData.MasterPasswordHint
	}

	log.Debug("checkMasterKey success")
	return nil, ""
}
