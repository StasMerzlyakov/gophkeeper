package config

import "github.com/caarlos0/env"

const (
	ClientDefaultServerAddres = "localhost:3200"
	ClientDefaultCACert       = ""
)

type ClientConf struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	CACert        string `env:"CA_CERT"`
}

func defaultClientConf() *ClientConf {
	return &ClientConf{
		ServerAddress: ClientDefaultServerAddres,
		CACert:        ClientDefaultCACert,
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
