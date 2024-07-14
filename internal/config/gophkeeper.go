package config

import "github.com/caarlos0/env"

const (
	ServerDefaultPort    = ":3200"
	ServerDefaultTLSKey  = ""
	ServerDefaultTLSCert = ""
)

type GophKeeperConf struct {
	Port    string `env:"PORT"`
	TLSKey  string `env:"TLS_KEY"`
	TLSCert string `env:"TLS_CERT"`
}

func defaultServConf() *GophKeeperConf {
	return &GophKeeperConf{
		Port:    ServerDefaultPort,
		TLSKey:  ServerDefaultTLSKey,
		TLSCert: ServerDefaultTLSCert,
	}
}

func LoadServConf() (*GophKeeperConf, error) {
	srvConf := defaultServConf()
	err := env.Parse(srvConf)
	if err != nil {
		return nil, err
	}

	return srvConf, nil
}
