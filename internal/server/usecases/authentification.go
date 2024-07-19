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
	regHelper RegistrationHelper) *auth {
	return &auth{
		conf:        conf,
		stflStorage: stflStorage,
		tempStorage: tempStorage,
		regHelper:   regHelper,
	}
}

type auth struct {
	stflStorage StateFullStorage
	tempStorage TemporaryStorage
	regHelper   RegistrationHelper
	conf        *config.ServerConf
}

func (auth *auth) Login(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {

	loginData, err := auth.stflStorage.GetLoginData(ctx, data.EMail)
	if err != nil {
		return "", fmt.Errorf("login - GetLoginData err %w", err)
	}

	ok, err := auth.regHelper.CheckPassword(data.Password, loginData.PasswordHash, loginData.PasswordSalt)
	if err != nil {
		return "", fmt.Errorf("login - CheckPassword err %w", err)
	}

	if !ok {
		return "", fmt.Errorf("login - wrong login or pass %w", domain.ErrAuthDataIncorrect)
	}

	sessionID := auth.regHelper.NewSessionID()
	if err = auth.tempStorage.Create(ctx, sessionID, *loginData); err != nil {
		return "", fmt.Errorf("login - can't create data %w", err)
	}
	return sessionID, nil
}

func (auth *auth) CheckOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error) {

	data, err := auth.tempStorage.LoadAndDelete(ctx, currentID)
	if err != nil {
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	authData, ok := data.(domain.LoginData)
	if !ok {
		err := fmt.Errorf("%w unexpected data type by id %s", domain.ErrAuthDataIncorrect, currentID)
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	otpKeyUrl, err := auth.regHelper.DecryptData(authData.EncryptedOTPKey)
	if err != nil {
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	ok, err = auth.regHelper.ValidatePassCode(otpKeyUrl, otpPass)
	if err != nil {
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	if !ok {
		err := fmt.Errorf("%w - wrong otp pass", domain.ErrAuthDataIncorrect)
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	jwtTok, err := auth.regHelper.CreateJWTToken(authData.UserID)
	if err != nil {
		return "", fmt.Errorf("checkOTP err - %w", err)
	}

	return jwtTok, nil
}
