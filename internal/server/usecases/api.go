package usecases

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"

	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . StateFullStorage,TemporaryStorage

type StateFullStorage interface {
	IsEMailBusy(ctx context.Context, email string) (bool, error)
	Registrate(ctx context.Context, data *domain.FullRegistrationData) (domain.UserID, error)
	GetLoginData(ctx context.Context, email string) (*domain.LoginData, error)
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
}
