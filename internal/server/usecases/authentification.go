package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAuth(conf *config.ServerConf,
	stflStorage StateFullStorage,
	tempStorage TemporaryStorage,
	emailSender EMailSender,
	regHelper RegistrationHelper) *auth {
	return &auth{
		conf:        conf,
		stflStorage: stflStorage,
		tempStorage: tempStorage,
		emailSender: emailSender,
		regHelper:   regHelper,
	}
}

type auth struct {
	stflStorage StateFullStorage
	tempStorage TemporaryStorage
	emailSender EMailSender
	regHelper   RegistrationHelper
	conf        *config.ServerConf
}

func (ath *auth) Login(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {

	loginData, err := ath.stflStorage.GetLoginData(ctx, data.EMail)
	if err != nil {
		return "", fmt.Errorf("login - GetLoginData err %w", err)
	}

	ok, err := ath.regHelper.CheckPassword(data.Password, loginData.PasswordHash, loginData.PasswordSalt)
	if err != nil {
		return "", fmt.Errorf("login - CheckPassword err %w", err)
	}

	if !ok {
		return "", fmt.Errorf("login - wrong login or pass %w", domain.ErrWrongAuthData)
	}

	sessionID := ath.regHelper.NewSessionID()
	if err = ath.tempStorage.Create(ctx, sessionID, *loginData); err != nil {
		return "", fmt.Errorf("login - can't create data %w", err)
	}
	return sessionID, nil
}
