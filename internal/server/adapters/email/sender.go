package email

import (
	"context"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	mail "github.com/xhit/go-simple-mail"
)

func NewSender(conf *config.ServerConf) *sender {

	return &sender{
		conf: conf,
	}
}

type sender struct {
	client *mail.SMTPClient
	conf   *config.ServerConf
	once   sync.Once
}

func (snd *sender) Connect(ctx context.Context) error {
	snd.once.Do(func() {

	})
}

func (snd *sender) Close() {

}
