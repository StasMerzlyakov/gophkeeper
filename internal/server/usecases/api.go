package usecases

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . StateFullStorage,TemporaryStorage,EMailSender,RegistrationHelper,FileStorage

type StateFullStorage interface {
	IsEMailAvailable(ctx context.Context, email string) (bool, error)
	Registrate(ctx context.Context, data *domain.FullRegistrationData) error
	GetLoginData(ctx context.Context, email string) (*domain.LoginData, error)
	GetHelloData(ctx context.Context) (*domain.HelloData, error)

	GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error)
	CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	DeleteBankCard(ctx context.Context, number string) error

	GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error)
	CreateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error
	UpdateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error
	DeleteUserPasswordData(ctx context.Context, hint string) error

	GetUserFilesBucket(ctx context.Context) (string, error)
}

type FileStorage interface {
	GetFileInfoList(ctx context.Context, bucket string) ([]domain.FileInfo, error)
	DeleteFileInfo(ctx context.Context, bucket string, name string) error
	StoreFile(ctx context.Context, bucket string, name string) (domain.StreamSender, error)
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
	Delete(ctx context.Context, sessionID domain.SessionID)
	Load(ctx context.Context, sessionID domain.SessionID) (any, error)
}

type EMailSender interface {
	Send(ctx context.Context, email string, png []byte) error
}

type RegistrationHelper interface {
	CheckEMailData(data *domain.EMailData) (bool, error)
	HashPassword(pass string) (*domain.HashData, error)
	GenerateQR(accountName string) (string, []byte, error)
	EncryptOTPKey(plaintext string) (string, error)
	DecryptOTPKey(ciphertext string) (string, error)
	ValidateAccountPass(pass string, hashB64 string, saltB64 string) (bool, error)
	NewSessionID() domain.SessionID
	ValidateOTPCode(keyURL string, passcode string) (bool, error)
	CreateJWTToken(userID domain.UserID) (domain.JWTToken, error)
	ParseJWTToken(jwtToken domain.JWTToken) (domain.UserID, error)
}
