package config

import (
	"time"

	"github.com/caarlos0/env"
)

const (
	ClientDefaultServerAddres      = "localhost:3200"
	ClientDefaultCACert            = ""
	ClientDefaultInterationTimeout = 2 * time.Second
)

type ClientConf struct {
	ServerAddress     string        `env:"SERVER_ADDRESS"`
	InterationTimeout time.Duration `env:"INTERACTION_TIMEOUT"`
	CACert            string        `env:"CA_CERT"`
}

func defaultClientConf() *ClientConf {
	return &ClientConf{
		ServerAddress:     ClientDefaultServerAddres,
		CACert:            ClientDefaultCACert,
		InterationTimeout: ClientDefaultInterationTimeout,
	}
}

func LoadClientConf() (*ClientConf, error) {
	srvConf := defaultClientConf()
	err := env.Parse(srvConf)
	if err != nil {
		return nil, err
	}

	return srvConf, nil
}
