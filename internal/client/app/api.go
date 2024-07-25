package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . RegServer,RegView,RegHelper,LoginServer,LoginView,LoginHelper

type RegServer interface {
	CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error)
	Registrate(ctx context.Context, data *domain.EMailData) error
	PassOTP(ctx context.Context, otpPass string) error
	InitMasterKey(ctx context.Context, mKey *domain.MasterKeyData) error
}

type RegView interface {
	ShowLoginView()
	ShowError(err error)
	ShowMsg(msg string)
	ShowRegForm()
	ShowRegOTPView()
	ShowInitMasterKeyView()
}

type RegHelper interface {
	ParseEMail(address string) bool
	ValidateAuthPassword(pass string) bool
	ValidateEncryptionPassword(pass string) bool
	EncryptData(secretKey string, plaintext string) (string, error)
	GenerateHello() (string, error)
	EncryptAES256(data []byte, passphrase string) (string, error)
	Random32ByteString() string
}

type LoginServer interface {
	Login(ctx context.Context, data *domain.EMailData) error
	PassOTP(ctx context.Context, otpPass string) error
	GetMasterKey(ctx context.Context)
	GetHelloData(ctx context.Context) (*domain.HelloData, error)
}

type LoginView interface {
	ShowLogOTPView()
	ShowError(err error)
	ShowMsg(msg string)
	ShowMasterKeyView(hint string)
	ShowDataAccessView()
}

type LoginHelper interface {
	DecryptData(secretKey string, ciphertext string) (string, error)
	DecryptAES256(ciphertext string, passphrase string) ([]byte, error)
	CheckHello(chk string) (bool, error)
}
