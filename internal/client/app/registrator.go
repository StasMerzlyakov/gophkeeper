package app

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewRegistrator() *registrator {
	return &registrator{}
}

func (reg *registrator) RegServer(srv AppServer) *registrator {
	reg.srv = srv
	return reg
}

func (reg *registrator) RegHelper(helper DomainHelper) *registrator {
	reg.helper = helper
	return reg
}

type registrator struct {
	srv    AppServer
	helper DomainHelper
}

func (reg *registrator) CheckEmail(ctx context.Context, email string) error {
	log := GetMainLogger()
	if !reg.helper.ParseEMail(email) {
		log.Warn("check email - wrong email format")
		return fmt.Errorf("%w - wrong email format", domain.ErrClientDataIncorrect)
	}

	if status, err := reg.srv.CheckEMail(ctx, email); err != nil {
		err := fmt.Errorf("%w check email err", err)
		log.Warn(err.Error())
		return err
	} else {
		if status == domain.EMailBusy {
			log.Warn("check email - email is busy")
			return fmt.Errorf("%w - email is busy", domain.ErrClientDataIncorrect)
		} else {
			log.Warn("email is available")
			return nil
		}
	}
}

func (reg *registrator) Registrate(ctx context.Context, data *domain.EMailData) error {
	log := GetMainLogger()

	log.Debugf("regstration %v start", data.EMail)
	if !reg.helper.ParseEMail(data.EMail) {
		err := fmt.Errorf("%w registration err - wrong email format", domain.ErrClientDataIncorrect)
		log.Warn(err.Error())
		return err
	}

	if !reg.helper.CheckAuthPasswordComplexityLevel(data.Password) {
		err := fmt.Errorf("%w registration err - password too slow", domain.ErrClientDataIncorrect)
		log.Warn(err.Error())
		return err
	}

	if err := reg.srv.Registrate(ctx, data); err != nil {
		err := fmt.Errorf("%w registration err", err)
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (reg *registrator) PassOTP(ctx context.Context, otpPass *domain.OTPPass) error {
	log := GetMainLogger()
	log.Debugf("passOTP start")

	if err := reg.srv.PassRegOTP(ctx, otpPass.Pass); err != nil {
		err := fmt.Errorf("%w - passOTP err", err)
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (reg *registrator) InitMasterKey(ctx context.Context, mKey *domain.UnencryptedMasterKeyData) error {
	log := GetMainLogger()

	if !reg.helper.CheckMasterPasswordComplexityLevel(mKey.MasterPassword) {
		err := fmt.Errorf("%w - initMasterKey err - key too slow", domain.ErrClientDataIncorrect)
		log.Warn(err.Error())
		return err
	}

	helloStr := reg.helper.Random32ByteString()
	helloEncrypted, err := reg.helper.EncryptHello(mKey.MasterPassword, helloStr)
	if err != nil {
		err := fmt.Errorf("%w - initMasterKey err", err)
		log.Warn(err.Error())
		return err
	}

	mData := &domain.MasterKeyData{
		MasterPasswordHint: mKey.MasterPasswordHint,
		HelloEncrypted:     helloEncrypted,
	}

	if err := reg.srv.InitMasterKey(ctx, mData); err != nil {
		err = fmt.Errorf("%w - initMasterKey err", err)
		log.Warn(err.Error())
		return err
	}
	return nil
}
