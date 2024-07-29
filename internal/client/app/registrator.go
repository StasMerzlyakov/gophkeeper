package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewRegistrator() *registrator {
	return &registrator{}
}

func (reg *registrator) RegServer(srv RegServer) *registrator {
	reg.srv = srv
	return reg
}

func (reg *registrator) RegView(view RegView) *registrator {
	reg.view = view
	return reg
}

func (reg *registrator) RegHelper(helper RegHelper) *registrator {
	reg.helper = helper
	return reg
}

type registrator struct {
	srv    RegServer
	view   RegView
	helper RegHelper
}

func (reg *registrator) CheckEmail(ctx context.Context, email string) {
	if !reg.helper.ParseEMail(email) {
		reg.view.ShowError(fmt.Errorf("%w - wrong email format", domain.ErrClientDataIncorrect))
		return
	}

	if status, err := reg.srv.CheckEMail(ctx, email); err != nil {
		reg.view.ShowError(err)

	} else {
		if status == domain.EMailBusy {
			reg.view.ShowError(fmt.Errorf("%w - email is busy", domain.ErrClientDataIncorrect))
		} else {
			reg.view.ShowMsg("email avaliable")
		}
	}
}

func (reg *registrator) Registrate(ctx context.Context, data *domain.EMailData) {
	log := GetMainLogger()

	log.Debugf("regstration %v start", data.EMail)
	if !reg.helper.ParseEMail(data.EMail) {
		err := fmt.Sprintf("registration %v err wrong email format", data.EMail)
		log.Errorf(err)
		reg.view.ShowMsg(err)
		return
	}

	if !reg.helper.CheckAuthPasswordComplexityLevel(data.Password) {
		err := "registration err - password too slow"
		log.Errorf(err)
		reg.view.ShowMsg(err)
		return
	}

	if err := reg.srv.Registrate(ctx, data); err != nil {
		log.Errorf("registration err - %v", err.Error())
		reg.view.ShowError(err)
		return
	}

	select {
	case <-ctx.Done():
		GetMainLogger().Errorf("registration err - context id done - %v", ctx.Err().Error())
		return
	default:
		GetMainLogger().Infof("registration of %v started", data.EMail)
		reg.view.ShowRegOTPView()
	}

}

func (reg *registrator) PassOTP(ctx context.Context, otpPass *domain.OTPPass) {
	log := GetMainLogger()

	log.Debugf("passOTP start")

	if err := reg.srv.PassRegOTP(ctx, otpPass.Pass); err != nil {
		err = fmt.Errorf("%w - passOTP err", err)
		log.Error(err)
		reg.view.ShowError(err)
		return
	}
	select {
	case <-ctx.Done():
		log.Errorf("context done %v", ctx.Err().Error())
		return
	default:
		log.Debugf("passOTP success")
		reg.view.ShowRegMasterKeyView()
	}

}

func (reg *registrator) InitMasterKey(ctx context.Context, mKey *domain.UnencryptedMasterKeyData) {
	log := GetMainLogger()

	if !reg.helper.CheckMasterPasswordComplexityLevel(mKey.MasterPassword) {
		err := fmt.Errorf("%w - master key too slow", domain.ErrClientDataIncorrect)
		log.Info(err.Error())
		reg.view.ShowMsg(err.Error())
		return
	}

	helloStr := reg.helper.Random32ByteString()
	helloEncrypted, err := reg.helper.EncryptHello(mKey.MasterPassword, helloStr)
	if err != nil {
		err = fmt.Errorf("%w - initMasterKey err", err)
		log.Error(err)
		reg.view.ShowError(err)
		return
	}

	mData := &domain.MasterKeyData{
		MasterPasswordHint: mKey.MasterPasswordHint,
		HelloEncrypted:     helloEncrypted,
	}

	if err := reg.srv.InitMasterKey(ctx, mData); err != nil {
		err = fmt.Errorf("%w - server err", err)
		log.Error(err)
		reg.view.ShowError(err)
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		reg.view.ShowLoginView()
	}
}
