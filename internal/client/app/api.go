package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . Server,InfoView,Pinger,RegServer,RegView,RegHelper,LoginServer,LoginView,LoginHelper,LoginStorage

type Pinger interface {
	Ping(ctx context.Context) error
}

type LoginStorage interface {
	SetMasterKey(masterKey string)
	GetMasterKey() string
}

type RegServer interface {
	CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error)
	Registrate(ctx context.Context, data *domain.EMailData) error
	PassRegOTP(ctx context.Context, otpPass string) error
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
	CheckAuthPasswordComplexityLevel(pass string) bool
	CheckMasterKeyPasswordComplexityLevel(pass string) bool
	EncryptMasterKey(masterKeyPass string, masterKey string) (string, error)
	GenerateHello() (string, error)
	EncryptShortData(data []byte, masterKey string) (string, error)
	Random32ByteString() string
}

type LoginServer interface {
	Login(ctx context.Context, data *domain.EMailData) error
	PassLoginOTP(ctx context.Context, otpPass string) error
	GetHelloData(ctx context.Context) (*domain.HelloData, error)
}

type InfoView interface {
	ShowError(err error)
	ShowMsg(msg string)
	ShowLogOTPView()
	ShowMasterKeyView(hint string)
	ShowDataAccessView()
	ShowLoginView()
	ShowRegForm()
	ShowRegOTPView()
	ShowInitMasterKeyView()
}

type LoginView interface {
	ShowLogOTPView()
	ShowError(err error)
	ShowMsg(msg string)
	ShowMasterKeyView(hint string)
	ShowDataAccessView()
}

type LoginHelper interface {
	DecryptMasterKey(masterKeyPass string, encryptedMasterKey string) (string, error)
	DecryptShortData(ciphertext string, masterKey string) ([]byte, error)
	CheckHello(chk string) (bool, error)
}

type Server interface {
	Stop()
	Ping(ctx context.Context) error
	CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error)
	Registrate(ctx context.Context, data *domain.EMailData) error
	PassRegOTP(ctx context.Context, otpPass string) error
	InitMasterKey(ctx context.Context, mKey *domain.MasterKeyData) error
	Login(ctx context.Context, data *domain.EMailData) error
	PassLoginOTP(ctx context.Context, otpPass string) error
	GetHelloData(ctx context.Context) (*domain.HelloData, error)
}
