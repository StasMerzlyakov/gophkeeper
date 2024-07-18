package usecases

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/google/uuid"
)

func NewRegistrationHelper(salfFn domain.SaltFn) *regHelper {
	return &regHelper{
		salfFn: salfFn,
	}
}

var _ RegistrationHelper = (*regHelper)(nil)

type regHelper struct {
	salfFn domain.SaltFn
}

func (rg *regHelper) CheckEMailData(data *domain.EMailData) (bool, error) {
	return domain.CheckEMailData(data)
}
func (rg *regHelper) HashPassword(pass string) (*domain.HashData, error) {
	return domain.HashPassword(pass, rg.salfFn)
}

func (rg *regHelper) GenerateQR(issuer string, accountName string) (string, []byte, error) {
	return domain.GenerateQR(issuer, accountName)
}

func (rg *regHelper) EncryptData(secretKey string, plaintext string) (string, error) {
	return domain.EncryptData(secretKey, plaintext, rg.salfFn)
}

func (rg *regHelper) CheckPassword(pass string, hashB64 string, saltB64 string) (bool, error) {
	return domain.CheckPassword(pass, hashB64, saltB64)
}

func (rg *regHelper) DecryptData(secretKey string, ciphertext string) (string, error) {
	return domain.DecryptData(secretKey, ciphertext)
}

func (rg *regHelper) NewSessionID() domain.SessionID {
	return domain.SessionID(uuid.NewString())
}

func (rg *regHelper) ValidatePassCode(keyURL string, passcode string) (bool, error) {
	return domain.ValidatePassCode(keyURL, passcode)
}

func (rg *regHelper) GenerateHello() (string, error) {
	return domain.GenerateHello(rg.salfFn)
}
func (rg *regHelper) CheckHello(toCheck string) (bool, error) {
	return domain.CheckHello(toCheck)
}
