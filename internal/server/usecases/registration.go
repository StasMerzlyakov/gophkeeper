package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewRegistrator(conf *config.ServerConf) *registrator {
	reg := &registrator{
		conf: conf,
	}
	return reg
}

func (reg *registrator) StateFullStorage(stflStorage StateFullStorage) *registrator {
	reg.stflStorage = stflStorage
	return reg
}

func (reg *registrator) TemporaryStorage(tempStorage TemporaryStorage) *registrator {
	reg.tempStorage = tempStorage
	return reg
}

func (reg *registrator) EMailSender(emailSender EMailSender) *registrator {
	reg.emailSender = emailSender
	return reg
}

func (reg *registrator) RegistrationHelper(regHelper RegistrationHelper) *registrator {
	reg.regHelper = regHelper
	return reg
}

type registrator struct {
	stflStorage StateFullStorage
	tempStorage TemporaryStorage
	emailSender EMailSender
	regHelper   RegistrationHelper
	conf        *config.ServerConf
}

func (reg *registrator) GetEMailStatus(ctx context.Context, email string) (domain.EMailStatus, error) {
	if isAvailable, err := reg.stflStorage.IsEMailAvailable(ctx, email); err != nil {
		return domain.EMailBusy, fmt.Errorf("checkEMail err - %w", err)
	} else {
		if isAvailable {
			return domain.EMailAvailable, nil
		} else {
			return domain.EMailBusy, nil
		}
	}
}

// First part of the registration process. Store email and password data. Generate and send email with OTP QR.
func (reg *registrator) Registrate(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)
	log.Debugw(action, "email", data.EMail, "msg", "registration started")
	if _, err := reg.regHelper.CheckEMailData(data); err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if isAvailable, err := reg.stflStorage.IsEMailAvailable(ctx, data.EMail); err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	} else {
		if !isAvailable {
			return "", fmt.Errorf("%w register err - email %v is busy", domain.ErrClientDataIncorrect, data.EMail)
		}
	}

	sessionID := reg.regHelper.NewSessionID()

	hashPasswordData, err := reg.regHelper.HashPassword(data.Password)
	if err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	key, image, err := reg.regHelper.GenerateQR(data.EMail)
	if err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	encryptedKey, err := reg.regHelper.EncryptOTPKey(key)
	if err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	regData := domain.RegistrationData{
		EMail:           data.EMail,
		PasswordHash:    hashPasswordData.Hash,
		PasswordSalt:    hashPasswordData.Salt,
		EncryptedOTPKey: encryptedKey,
		State:           domain.RegistrationStateInit,
	}

	if err := reg.tempStorage.Create(ctx, sessionID, regData); err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if err := reg.emailSender.Send(ctx, data.EMail, image); err != nil {
		err = fmt.Errorf("register err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	log.Debugw(action, "email", data.EMail, "msg", "email registerd")
	return sessionID, nil
}

// Second part of the registration process. Check OTP password.
func (reg *registrator) PassOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "passOTP start")

	data, err := reg.tempStorage.Load(ctx, currentID)
	if err != nil {
		err = fmt.Errorf("passOTP err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	regData, ok := data.(domain.RegistrationData)
	if !ok {
		err := fmt.Errorf("%w passOTP err - unexpected data by id %s", domain.ErrAuthDataIncorrect, currentID)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if regData.State != domain.RegistrationStateInit {
		err := fmt.Errorf("%w passOTP err - wrong registartionState by id %s", domain.ErrAuthDataIncorrect, currentID)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	otpKey, err := reg.regHelper.DecryptOTPKey(regData.EncryptedOTPKey)
	if err != nil {
		err = fmt.Errorf("passOTP err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	ok, err = reg.regHelper.ValidateOTPCode(otpKey, otpPass)
	if err != nil {
		err = fmt.Errorf("passOTP err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if !ok {
		err := fmt.Errorf("%w - passOTP err - wrong otp pass", domain.ErrAuthDataIncorrect)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	regDataNew := domain.RegistrationData{
		EMail:           regData.EMail,
		PasswordHash:    regData.PasswordHash,
		PasswordSalt:    regData.PasswordSalt,
		EncryptedOTPKey: regData.EncryptedOTPKey,
		State:           domain.RegistrationStateAuth,
	}

	newSessionID := reg.regHelper.NewSessionID()

	if err := reg.tempStorage.DeleteAndCreate(ctx, currentID, newSessionID, regDataNew); err != nil {
		err = fmt.Errorf("passOTP err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	log.Debugw(action, "msg", "passOTP complete")
	return newSessionID, nil
}

// The last part of the registration process - store MasterKeyData.
func (reg *registrator) InitMasterKey(ctx context.Context, currentID domain.SessionID, mKey *domain.MasterKeyData) error {

	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "initMasterKey start")

	data, err := reg.tempStorage.LoadAndDelete(ctx, currentID)
	if err != nil {
		err = fmt.Errorf("initMasterKey err - %w", err)
		log.Infow(action, "err", err.Error())
		return err
	}

	regData, ok := data.(domain.RegistrationData)
	if !ok {
		err := fmt.Errorf("initMasterKey err - %w unexpected data by id %s", domain.ErrClientDataIncorrect, currentID)
		log.Infow(action, "err", err.Error())
		return err
	}

	if regData.State != domain.RegistrationStateAuth {
		err := fmt.Errorf("initMasterKey err - %w - wrong registartionState by id %s", domain.ErrClientDataIncorrect, currentID)
		log.Infow(action, "err", err.Error())
		return err
	}

	fullData := &domain.FullRegistrationData{
		EMail:              regData.EMail,
		PasswordHash:       regData.PasswordHash,
		PasswordSalt:       regData.PasswordSalt,
		EncryptedOTPKey:    regData.EncryptedOTPKey,
		EncryptedMasterKey: mKey.EncryptedMasterKey,
		MasterKeyHint:      mKey.MasterKeyHint,
		HelloEncrypted:     mKey.HelloEncrypted,
	}

	if err := reg.stflStorage.Registrate(ctx, fullData); err != nil {
		err = fmt.Errorf("initMasterKey err - %w", err)
		log.Infow(action, "err", err.Error())
		return err
	}

	log.Debugw(action, "msg", "initMasterKey success")

	return nil

}
