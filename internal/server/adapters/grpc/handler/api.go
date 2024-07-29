package handler

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . Registrator,DataAccessor,AuthService

type Registrator interface {
	GetEMailStatus(ctx context.Context, email string) (domain.EMailStatus, error)
	Registrate(ctx context.Context, data *domain.EMailData) (domain.SessionID, error)
	PassOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.SessionID, error)
	InitMasterKey(ctx context.Context, currentID domain.SessionID, mKey *domain.MasterKeyData) error
}

type DataAccessor interface {
	GetHelloData(ctx context.Context) (*domain.HelloData, error)

	GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error)
	CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	DeleteBankCard(ctx context.Context, number string) error

	GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error)
	CreateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error
	UpdateUserPasswordData(ctx context.Context, bnkCard *domain.EncryptedUserPasswordData) error
	DeleteUserPasswordData(ctx context.Context, hint string) error
}

type AuthService interface {
	Login(ctx context.Context, data *domain.EMailData) (domain.SessionID, error)
	CheckOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error)
}
