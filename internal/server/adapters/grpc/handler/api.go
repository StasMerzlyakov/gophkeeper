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
}

type AuthService interface {
	Login(ctx context.Context, data *domain.EMailData) (domain.SessionID, error)
	CheckOTP(ctx context.Context, currentID domain.SessionID, otpPass string) (domain.JWTToken, error)
}
