package config

import (
	"time"

	"github.com/caarlos0/env"
)

const (
	ServerDefaultPort        = ":3200"
	ServerDefaultTLSKey      = ""
	ServerDefaultTLSCert     = ""
	ServerDefaultTokenExp    = 3 * time.Hour
	ServerDefaultTokenSecret = "secret"
	ServerDefaultAuthTimeout = 10 * time.Minute
	ServerDefaultMasterKey   = "MasterKey!"
	ServerDefaultDomainName  = "localhost"
)

type ServerConf struct {
	Port        string        `env:"PORT"`
	TLSKey      string        `env:"TLS_KEY"`
	TLSCert     string        `env:"TLS_CERT"`
	TokenExp    time.Duration `env:"JWT_EXP"`
	TokenSecret string        `env:"JWT_SECRET"`
	AuthTimeout time.Duration `env:"AUTH_TIMEOUT"`
	MasterKey   string        `env:"MASTER_KEY"`
	DomainName  string        `env:"DOMAIN_NAME"`
}

func defaultServConf() *ServerConf {
	return &ServerConf{
		Port:        ServerDefaultPort,
		TLSKey:      ServerDefaultTLSKey,
		TLSCert:     ServerDefaultTLSCert,
		TokenExp:    ServerDefaultTokenExp,
		TokenSecret: ServerDefaultTokenSecret,
		AuthTimeout: ServerDefaultAuthTimeout,
		MasterKey:   ServerDefaultMasterKey,
		DomainName:  ServerDefaultDomainName,
	}
}

func LoadServConf() (*ServerConf, error) {
	srvConf := defaultServConf()
	err := env.Parse(srvConf)
	if err != nil {
		return nil, err
	}

	return srvConf, nil
}
