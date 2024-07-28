package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAuth(conf *config.ServerConf) *auth {
	auth := &auth{
		conf: conf,
	}

	return auth
}

func (ath *auth) StateFullStorage(stflStorage StateFullStorage) *auth {
	ath.stflStorage = stflStorage
	return ath
}

func (ath *auth) TemporaryStorage(tempStorage TemporaryStorage) *auth {
	ath.tempStorage = tempStorage
	return ath
}

func (ath *auth) RegistrationHelper(regHelper RegistrationHelper) *auth {
	ath.regHelper = regHelper
	return ath
}

type auth struct {
	stflStorage StateFullStorage
	tempStorage TemporaryStorage
	regHelper   RegistrationHelper
	conf        *config.ServerConf
}

// First part of the authentification process. Check login and password.
func (auth *auth) Login(ctx context.Context, data *domain.EMailData) (domain.SessionID, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "email", data.EMail, "msg", "authentification started")

	loginData, err := auth.stflStorage.GetLoginData(ctx, data.EMail)
	if err != nil {
		err = fmt.Errorf("login err - GetLoginData err %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	ok, err := auth.regHelper.ValidateAccountPass(data.Password, loginData.PasswordHash, loginData.PasswordSalt)
	if err != nil {
		err = fmt.Errorf("login err -ValidateAccountPass err %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if !ok {
		err := fmt.Errorf("login err - wrong login or pass %w", domain.ErrAuthDataIncorrect)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	sessionID := auth.regHelper.NewSessionID()
	if err = auth.tempStorage.Create(ctx, sessionID, *loginData); err != nil {
		err = fmt.Errorf("login err - can't create data %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	log.Debugw(action, "email", data.EMail, "msg", "email checked")
	return sessionID, nil
}

// Second part of the authentification process. Check OTP code.
func (auth *auth) CheckOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)
	log.Debugw(action, "msg", "chek otp started")

	data, err := auth.tempStorage.Load(ctx, currentID)
	if err != nil {
		err = fmt.Errorf("checkOTP err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	authData, ok := data.(domain.LoginData)
	if !ok {
		err := fmt.Errorf("%w checkOTP err - unexpected data type by id %s", domain.ErrAuthDataIncorrect, currentID)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	otpKey, err := auth.regHelper.DecryptOTPKey(authData.EncryptedOTPKey)
	if err != nil {
		err = fmt.Errorf("checkOTP decryptOTPKey err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	ok, err = auth.regHelper.ValidateOTPCode(otpKey, otpPass)
	if err != nil {
		err = fmt.Errorf("checkOTP validateOTPCode err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	if !ok {
		err := fmt.Errorf("%w - checkOTP err - wrong otp pass", domain.ErrAuthDataIncorrect)
		log.Infow(action, "err", err.Error())
		return "", err
	}

	jwtTok, err := auth.regHelper.CreateJWTToken(authData.UserID)
	if err != nil {
		err = fmt.Errorf("checkOTP createJWT err - %w", err)
		log.Infow(action, "err", err.Error())
		return "", err
	}
	auth.tempStorage.Delete(ctx, currentID)
	log.Debugw(action, "msg", "autehntification complete")
	return jwtTok, nil
}
