package usecases

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/google/uuid"
)

func NewRegistrationHelper(conf *config.ServerConf, salfFn domain.SaltFn) *regHelper {
	return &regHelper{
		salfFn: salfFn,
		conf:   conf,
	}
}

var _ RegistrationHelper = (*regHelper)(nil)

type regHelper struct {
	salfFn domain.SaltFn
	conf   *config.ServerConf
}

func (rg *regHelper) CheckEMailData(data *domain.EMailData) (bool, error) {
	return domain.CheckEMailData(data)
}
func (rg *regHelper) HashPassword(pass string) (*domain.HashData, error) {
	return domain.HashPassword(pass, rg.salfFn)
}

func (rg *regHelper) GenerateQR(userEmail string) (string, []byte, error) {
	return domain.GenerateQR(rg.conf.DomainName, userEmail)
}

func (rg *regHelper) EncryptOTPKey(plaintext string) (string, error) {
	return domain.EncryptOTPKey(rg.conf.ServerSecret, plaintext)
}

func (rg *regHelper) ValidateAccountPass(pass string, hashB64 string, saltB64 string) (bool, error) {
	return domain.ValidateAccountPass(pass, hashB64, saltB64)
}

func (rg *regHelper) DecryptOTPKey(ciphertext string) (string, error) {
	return domain.DecryptOTPKey(rg.conf.ServerSecret, ciphertext)
}

func (rg *regHelper) NewSessionID() domain.SessionID {
	return domain.SessionID(uuid.NewString())
}

func (rg *regHelper) ValidateOTPCode(keyURL string, passcode string) (bool, error) {
	return domain.ValidateOTPCode(keyURL, passcode)
}

func (rg *regHelper) CreateJWTToken(userID domain.UserID) (domain.JWTToken, error) {
	return domain.CreateJWTToken([]byte(rg.conf.TokenSecret), rg.conf.TokenExp, userID)
}

func (rg *regHelper) ParseJWTToken(jwtToken domain.JWTToken) (domain.UserID, error) {
	return domain.ParseJWTToken([]byte(rg.conf.TokenSecret), jwtToken)
}
