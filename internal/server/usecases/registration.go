package usecases

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/domain"
)

type StateFullStorage interface {
	IsEMailBusy(email string) (bool, error)
}

func NewRegistratrator(conf *config.ServerConf) *registrator {
	return &registrator{}
}

type registrator struct {
}

func (reg *registrator) Register(ctx context.Context, data *domain.EMailData) (*domain.Claims, error) {
	return nil, nil
}
