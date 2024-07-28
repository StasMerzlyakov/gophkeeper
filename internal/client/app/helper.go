package app

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewHelper(salfFn domain.SaltFn) *helper {
	return &helper{
		salfFn: salfFn,
	}
}

var _ RegHelper = (*helper)(nil)
var _ LoginHelper = (*helper)(nil)

type helper struct {
	salfFn domain.SaltFn
}

func (h *helper) ParseEMail(email string) bool {
	return domain.ParseEMail(email)
}

func (h *helper) CheckAuthPasswordComplexityLevel(email string) bool {
	return domain.CheckAuthPasswordComplexityLevel(email)
}

func (h *helper) CheckMasterKeyPasswordComplexityLevel(pass string) bool {
	return domain.CheckMasterKeyPasswordComplexityLevel(pass)
}

func (h *helper) Random32ByteString() string {
	return domain.Random32ByteString()
}

func (h *helper) GenerateHello() (string, error) {
	return domain.GenerateHello(h.salfFn)
}

func (h *helper) CheckHello(chk string) (bool, error) {
	return domain.CheckHello(chk)
}

func (h *helper) EncryptMasterKey(masterKeyPass string, masterKey string) (string, error) {
	return domain.EncryptMasterKey(masterKeyPass, masterKey)
}

func (h *helper) DecryptMasterKey(masterKeyPass string, encryptedMasterKey string) (string, error) {
	return domain.DecryptMasterKey(masterKeyPass, encryptedMasterKey)
}

func (h *helper) DecryptShortData(ciphertext string, masterKey string) ([]byte, error) {
	return domain.DecryptShortData(ciphertext, masterKey)
}

func (h *helper) EncryptShortData(data []byte, masterKey string) (string, error) {
	return domain.EncryptShortData(data, masterKey)
}
