package usecases

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . StateFullStorage,TemporaryStorage,EMailSender,RegistrationHelper

type StateFullStorage interface {
	IsEMailAvailable(ctx context.Context, email string) (bool, error)
	Registrate(ctx context.Context, data *domain.FullRegistrationData) error
	GetLoginData(ctx context.Context, email string) (*domain.LoginData, error)
	GetHelloData(ctx context.Context) (*domain.HelloData, error)
}

type TemporaryStorage interface {
	Create(ctx context.Context, sessionID domain.SessionID, data any) error
	DeleteAndCreate(ctx context.Context,
		oldSessionID domain.SessionID,
		sessionID domain.SessionID,
		data any,
	) error
	LoadAndDelete(ctx context.Context,
		sessionID domain.SessionID,
	) (any, error)
	Load(ctx context.Context, sessionID domain.SessionID) (any, error)
}

type EMailSender interface {
	Send(ctx context.Context, email string, png []byte) error
}

type RegistrationHelper interface {
	CheckEMailData(data *domain.EMailData) (bool, error)
	HashPassword(pass string) (*domain.HashData, error)
	GenerateQR(accountName string) (string, []byte, error)
	EncryptData(plaintext string) (string, error)
	DecryptData(ciphertext string) (string, error)
	ValidateAccountPass(pass string, hashB64 string, saltB64 string) (bool, error)
	NewSessionID() domain.SessionID
	ValidatePassCode(keyURL string, passcode string) (bool, error)
	GenerateHello() (string, error)
	CreateJWTToken(userID domain.UserID) (domain.JWTToken, error)
	ParseJWTToken(jwtToken domain.JWTToken) (domain.UserID, error)
}
