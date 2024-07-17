package usecases

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"image/png"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

func NewRegistratrator(conf *config.ServerConf,
	stflStorage StateFullStorage,
	tempStorage TemporaryStorage,
) *registrator {

	return &registrator{
		stflStorage: stflStorage,
		serverKey:   conf.ServerKey,
		domainName:  conf.DomainName,
		tempStorage: tempStorage,
	}
}

type registrator struct {
	stflStorage StateFullStorage
	tempStorage TemporaryStorage
	serverKey   string
	domainName  string
}

func (reg *registrator) CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error) {
	if isBusy, err := reg.stflStorage.IsEMailBusy(ctx, email); err != nil {
		return domain.EMailBusy, fmt.Errorf("checkEMail err - %w", err)
	} else {
		if isBusy {
			return domain.EMailBusy, nil
		} else {
			return domain.EMailAvailable, nil
		}
	}
}

func (reg *registrator) Register(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {

	if _, err := domain.CheckEMailData(data); err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	sessionID := domain.SessionID(uuid.NewString())

	action := domain.GetAction(1)

	hashPasswordData, err := domain.HashPassword(data.Password, rand.Read)
	if err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      reg.domainName,
		AccountName: data.EMail,
	})

	if err != nil {
		log := domain.GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("TOTP key generation err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}

	var buf bytes.Buffer
	img, err := key.Image(450, 450)
	if err != nil {
		log := domain.GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("TOTP image generation err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}
	if err = png.Encode(&buf, img); err != nil {
		log := domain.GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("png ecode err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}

	keyURL := key.URL()

	encryptedKey, err := domain.EncryptData(reg.serverKey, keyURL, rand.Read)
	if err != nil {
		log := domain.GetApplicationLogger()
		log.Warnf(action, "err", fmt.Sprintf("encrypt data err %s", err.Error()))
		return "", fmt.Errorf("register err - %w", err)
	}

	regData := domain.RegistrationData{
		EMail:            data.EMail,
		PasswordHash:     hashPasswordData.Hash,
		PasswordSalt:     hashPasswordData.Salt,
		EncryptedOTPPass: encryptedKey,
		State:            domain.RegistrationStateInit,
	}

	if err := reg.tempStorage.Create(ctx, sessionID, regData); err != nil {
		return "", fmt.Errorf("register err - %w", err)
	}

	// generate OTP pass
	return sessionID, nil
}
