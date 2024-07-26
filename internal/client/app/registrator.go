package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewRegistrator(conf *config.ClientConf) *registrator {
	return &registrator{
		conf: conf,
	}
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
	conf   *config.ClientConf
	srv    RegServer
	view   RegView
	helper RegHelper
}

func (reg *registrator) CheckEmail(ctx context.Context, email string) {
	if !reg.helper.ParseEMail(email) {
		reg.view.ShowError(fmt.Errorf("%w - wrong email format", domain.ErrClientDataIncorrect))
		return
	}

	timedCtx, fn := context.WithTimeout(ctx, reg.conf.InterationTimeout)
	defer fn()

	if status, err := reg.srv.CheckEMail(timedCtx, email); err != nil {
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
	timedCtx, fn := context.WithTimeout(ctx, reg.conf.InterationTimeout)
	defer fn()

	if !reg.helper.CheckAuthPasswordComplexityLevel(data.Password) {
		reg.view.ShowError(fmt.Errorf("%w - password too slow", domain.ErrClientDataIncorrect))
		return
	}

	if err := reg.srv.Registrate(timedCtx, data); err != nil {
		reg.view.ShowError(err)
		return
	}
	reg.view.ShowRegOTPView()
}

func (reg *registrator) PassOTP(ctx context.Context, otpPass string) {
	timedCtx, fn := context.WithTimeout(ctx, reg.conf.InterationTimeout)
	defer fn()

	if err := reg.srv.PassRegOTP(timedCtx, otpPass); err != nil {
		reg.view.ShowError(err)
		return
	}

	reg.view.ShowInitMasterKeyView()
}

func (reg *registrator) InitMasterKey(ctx context.Context, mKey *domain.UnencryptedMasterKeyData) {
	timedCtx, fn := context.WithTimeout(ctx, reg.conf.InterationTimeout)
	defer fn()

	if !reg.helper.CheckMasterKeyPasswordComplexityLevel(mKey.MasterKeyPassword) {
		reg.view.ShowError(fmt.Errorf("%w - master key too slow", domain.ErrClientDataIncorrect))
		return
	}

	masterKey := reg.helper.Random32ByteString()
	encryptedMasterKey, err := reg.helper.EncryptMasterKey(mKey.MasterKeyPassword, masterKey)
	if err != nil {
		reg.view.ShowError(err)
		return
	}

	helloStr, err := reg.helper.GenerateHello()
	if err != nil {
		reg.view.ShowError(err)
		return
	}

	helloEncr, err := reg.helper.EncryptShortData([]byte(helloStr), masterKey)
	if err != nil {
		reg.view.ShowError(err)
		return
	}

	mData := &domain.MasterKeyData{
		EncryptedMasterKey: encryptedMasterKey,
		MasterKeyHint:      mKey.MasterKeyHint,
		HelloEncrypted:     helloEncr,
	}

	if err := reg.srv.InitMasterKey(timedCtx, mData); err != nil {
		reg.view.ShowError(err)
		return
	}

	reg.view.ShowLoginView()
}
