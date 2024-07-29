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

func (h *helper) CheckMasterPasswordComplexityLevel(pass string) bool {
	return domain.CheckMasterPasswordComplexityLevel(pass)
}

func (h *helper) Random32ByteString() string {
	return domain.Random32ByteString()
}

func (h *helper) EncryptHello(masterPassword, hello string) (string, error) {
	return domain.EncryptHello(masterPassword, hello)
}

func (h *helper) DecryptHello(masterPassword, helloEncrypted string) error {
	return domain.DecryptHello(masterPassword, helloEncrypted)
}
