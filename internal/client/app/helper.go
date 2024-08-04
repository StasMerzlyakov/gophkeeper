package app

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

// NewHelper create domain method warpper object. Useful for covering controller logic tests.
func NewHelper(salfFn domain.SaltFn) *helper {
	return &helper{
		salfFn: salfFn,
	}
}

var _ DomainHelper = (*helper)(nil)

type helper struct {
	salfFn domain.SaltFn
}

func (h *helper) ParseEMail(email string) bool {
	return domain.ParseEMail(email)
}

func (h *helper) CheckAuthPasswordComplexityLevel(email string) bool {
	return domain.CheckAuthPasswordComplexityLevel(email)
}

func (h *helper) CheckBankCardData(data *domain.BankCard) error {
	return domain.CheckBankCardData(data)
}
func (h *helper) CheckUserPasswordData(data *domain.UserPasswordData) error {
	return domain.CheckUserPasswordData(data)
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

func (h *helper) EncryptShortData(masterKey string, data string) (string, error) {
	return domain.EncryptShortData(masterKey, data)
}
func (h *helper) DecryptShortData(masterKey string, ciphertext string) (string, error) {
	return domain.DecryptShortData(masterKey, ciphertext)
}

func (h *helper) CheckFileForRead(info *domain.FileInfo) error {
	return domain.CheckFileForRead(info)
}

func (h *helper) CheckFileForWrite(info *domain.FileInfo) error {
	return domain.CheckFileForWrite(info)
}

func (h *helper) CreateFileStreamer(info *domain.FileInfo) (domain.StreamFileReader, error) {
	return domain.CreateFileStreamer(info.Path)
}
