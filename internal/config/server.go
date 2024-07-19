package config

import (
	"time"

	"github.com/caarlos0/env"
)

const (
	ServerDefaultPort                = ":3200"
	ServerDefaultTLSKey              = ""
	ServerDefaultTLSCert             = ""
	ServerDefaultTokenExp            = 3 * time.Hour
	ServerDefaultTokenSecret         = "secret"
	ServerDefaultAuthStageTimeout    = 5 * time.Minute
	ServerDefaultServerEncryptionKey = "ServerKey!"
	ServerDefaultDomainName          = "localhost"
	ServerDefaultSMTPHost            = "localhost"
	ServerDefaultSMTPPort            = 25
	ServerDefaultServerEMail         = "localhost@localdomain"
)

type ServerConf struct {
	Port                string        `env:"PORT"`
	TLSKey              string        `env:"TLS_KEY"`
	TLSCert             string        `env:"TLS_CERT"`
	TokenExp            time.Duration `env:"JWT_EXP"`
	TokenSecret         string        `env:"JWT_SECRET"`
	AuthStageTimeout    time.Duration `env:"AUTH_STAGE_TIMEOUT"`
	ServerEncryptionKey string        `env:"SERVER_ENCRYPTION_KEY"`
	DomainName          string        `env:"DOMAIN_NAME"`
	SMTPHost            string        `env:"SMTP_HOST"`
	SMTPPort            int           `env:"SMTP_PORT"`
	ServerEMail         string        `env:"SERVER_EMAIL"`
}

func defaultServConf() *ServerConf {
	return &ServerConf{
		Port:                ServerDefaultPort,
		TLSKey:              ServerDefaultTLSKey,
		TLSCert:             ServerDefaultTLSCert,
		TokenExp:            ServerDefaultTokenExp,
		TokenSecret:         ServerDefaultTokenSecret,
		AuthStageTimeout:    ServerDefaultAuthStageTimeout,
		ServerEncryptionKey: ServerDefaultServerEncryptionKey,
		DomainName:          ServerDefaultDomainName,
		SMTPHost:            ServerDefaultSMTPHost,
		SMTPPort:            ServerDefaultSMTPPort,
		ServerEMail:         ServerDefaultServerEMail,
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
