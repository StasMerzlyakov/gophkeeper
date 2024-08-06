package app

import (
	"context"
)

func NewPinger() *pinger {
	return &pinger{}
}

func (p *pinger) SetPinger(srv Pinger) *pinger {
	p.srv = srv
	return p
}

type pinger struct {
	srv Pinger
}

func (p *pinger) Ping(ctx context.Context) error {
	return p.srv.Ping(ctx)
}
