package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewRegistrator(conf *config.ServerConf,
	stflStorage StateFullStorage,
	tempStorage TemporaryStorage,
	emailSender EMailSender,
	regHelper RegistrationHelper,
) *registrator {

	return &registrator{
		conf:        conf,
		stflStorage: stflStorage,
		tempStorage: tempStorage,
		emailSender: emailSender,
		regHelper:   regHelper,
	}
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

func (reg *registrator) Register(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {

	sessionID := reg.regHelper.NewSessionID()
	log := domain.GetApplicationLogger()

	action := domain.GetAction(1)

	if _, err := reg.regHelper.CheckEMailData(data); err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	hashPasswordData, err := reg.regHelper.HashPassword(data.Password)
	if err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	key, image, err := reg.regHelper.GenerateQR(reg.conf.DomainName, data.EMail)
	if err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	encryptedKey, err := reg.regHelper.EncryptData(reg.conf.ServerKey, key)
	if err != nil {
		log.Warnf(action, "err", fmt.Sprintf("encrypt data err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}

	regData := domain.RegistrationData{
		EMail:           data.EMail,
		PasswordHash:    hashPasswordData.Hash,
		PasswordSalt:    hashPasswordData.Salt,
		EncryptedOTPKey: encryptedKey,
		State:           domain.RegistrationStateInit,
	}

	if err := reg.tempStorage.Create(ctx, sessionID, regData); err != nil {
		log.Warnf(action, "err", fmt.Sprintf("create data err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}

	if err := reg.emailSender.Send(ctx, data.EMail, image); err != nil {
		log.Warnf(action, "err", fmt.Sprintf("send email err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}
	return sessionID, nil
}

func (reg *registrator) PassOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error) {

	log := domain.GetApplicationLogger()
	action := domain.GetAction(1)

	data, err := reg.tempStorage.Load(ctx, currentID)
	if err != nil {
		log.Warnf(action, "err", fmt.Sprintf("encrypt data err %s", err.Error()))
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	regData, ok := data.(domain.RegistrationData)
	if !ok {
		err := fmt.Errorf("%w unexpected data by id %s", domain.ErrClientDataIncorrect, currentID)
		log.Warnf(action, "err", err.Error())
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	if regData.State != domain.RegistrationStateInit {
		err := fmt.Errorf("%w - wrong registartionState by id %s", domain.ErrClientDataIncorrect, currentID)
		log.Warnf(action, "err", err.Error())
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	otpKeyUrl, err := reg.regHelper.DecryptData(reg.conf.ServerKey, regData.EncryptedOTPKey)
	if err != nil {
		log.Warnf(action, "err", fmt.Sprintf("decrypt data err %s", err.Error()))
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	ok, err = reg.regHelper.ValidatePassCode(otpKeyUrl, otpPass)
	if err != nil {
		log.Warnf(action, "err", fmt.Sprintf("validate pass err %s", err.Error()))
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	if !ok {
		err := fmt.Errorf("%w - wrong otp pass", domain.ErrClientDataIncorrect)
		log.Warnf(action, "err", err.Error())
		return "", fmt.Errorf("passOTP err - %w", err)
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
		log.Warnf(action, "err", fmt.Sprintf("can't refresh reg data - %s", err.Error()))
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	return newSessionID, nil
}
