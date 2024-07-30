package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . AppServer,AppView,Pinger,DomainHelper,AppStorage

type Pinger interface {
	Ping(ctx context.Context) error
}

type AppStorage interface {
	SetMasterPassword(masterPassword string)
	GetMasterPassword() string
	AddBankCard(bankCard *domain.BankCard) error
	AddUserPasswordData(data *domain.UserPasswordData) error
	UpdateBankCard(bankCard *domain.BankCard) error
	UpdatePasswordData(data *domain.UserPasswordData) error
	DeleteBankCard(number string) error
	DeleteUpdatePasswordData(hint string) error
	GetBankCard(number string) (*domain.BankCard, error)
	GetUpdatePasswordData(hint string) (*domain.UserPasswordData, error)
	GetBankCardNumberList() []string
	GetUserPasswordDataList() []string
}

type DomainHelper interface {
	ParseEMail(address string) bool
	CheckAuthPasswordComplexityLevel(pass string) bool
	CheckMasterPasswordComplexityLevel(pass string) bool
	EncryptHello(masterPass string, hello string) (string, error)
	Random32ByteString() string
	DecryptHello(masterPassword string, helloEncrypted string) error
	CheckBankCardData(data *domain.BankCard) error
	CheckUserPasswordData(data *domain.UserPasswordData) error
}

type AppView interface {
	ShowError(err error)
	ShowMsg(msg string)
	ShowLogOTPView()
	ShowMasterKeyView(hint string)
	ShowDataAccessView()
	ShowLoginView()
	ShowRegView()
	ShowRegOTPView()
	ShowRegMasterKeyView()
	ShowBankCardListView(cardsNumber []string)
	ShowBankCardView(bankCard *domain.BankCard)
}

type AppServer interface {
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
