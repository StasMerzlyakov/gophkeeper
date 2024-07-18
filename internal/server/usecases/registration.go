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
		return "", fmt.Errorf("register err - %w", err)
	}

	if err := reg.emailSender.Send(ctx, data.EMail, image); err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}
	return sessionID, nil
}

func (reg *registrator) PassOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error) {

	data, err := reg.tempStorage.Load(ctx, currentID)
	if err != nil {
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	regData, ok := data.(domain.RegistrationData)
	if !ok {
		err := fmt.Errorf("%w unexpected data by id %s", domain.ErrClientDataIncorrect, currentID)
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	if regData.State != domain.RegistrationStateInit {
		err := fmt.Errorf("%w - wrong registartionState by id %s", domain.ErrClientDataIncorrect, currentID)
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	otpKeyUrl, err := reg.regHelper.DecryptData(reg.conf.ServerKey, regData.EncryptedOTPKey)
	if err != nil {
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	ok, err = reg.regHelper.ValidatePassCode(otpKeyUrl, otpPass)
	if err != nil {
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	if !ok {
		err := fmt.Errorf("%w - wrong otp pass", domain.ErrClientDataIncorrect)
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
		return "", fmt.Errorf("passOTP err - %w", err)
	}

	return newSessionID, nil
}

func (reg *registrator) InitMasterKey(ctx context.Context, currentID domain.SessionID, mKey *domain.MasterKeyData) error {

	data, err := reg.tempStorage.LoadAndDelete(ctx, currentID)
	if err != nil {
		return fmt.Errorf("initMasterKey err - %w", err)
	}

	regData, ok := data.(domain.RegistrationData)
	if !ok {
		err := fmt.Errorf("%w unexpected data by id %s", domain.ErrClientDataIncorrect, currentID)
		return fmt.Errorf("initMasterKey err - %w", err)
	}

	if regData.State != domain.RegistrationStateAuth {
		err := fmt.Errorf("%w - wrong registartionState by id %s", domain.ErrClientDataIncorrect, currentID)
		return fmt.Errorf("initMasterKey err - %w", err)
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
		return fmt.Errorf("initMasterKey err - %w", err)
	} else {
		return nil
	}
}
