package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env"
)

const (
	ServerDefaultPort                = ":3200"
	ServerDefaultTLSKey              = "../../keys/server-key.pem"
	ServerDefaultTLSCert             = "../../keys/server-cert.pem"
	ServerDefaultTokenExp            = 3 * time.Hour
	ServerDefaultTokenSecret         = "secret"
	ServerDefaultAuthStageTimeout    = 5 * time.Minute
	ServerDefaultServerEncryptionKey = "ServerKey!12>Au{mL736}"
	ServerDefaultDomainName          = "localhost"
	ServerDefaultSMTPHost            = "localhost"
	ServerDefaultSMTPPort            = 25
	ServerDefaultServerEMail         = "localhost@localdomain"
	ServerDefaultMaxConns            = 5
	ServerDefaultMaxConnLifetime     = 5 * time.Minute
	ServerDefaultMaxConnIdleTime     = 5 * time.Minute
	ServerDefaultDatabaseURI         = "postgres://user:user@localhost:5432/gophkeeper"
)

type ServerConf struct {
	Port                string        `env:"PORT" json:"port"`
	TLSKey              string        `env:"TLS_KEY" json:"tlsKey"`
	TLSCert             string        `env:"TLS_CERT" json:"tlsCert"`
	TokenExp            time.Duration `env:"JWT_EXP" json:"jwtExp"`
	TokenSecret         string        `env:"JWT_SECRET" json:"jwtSecret"`
	AuthStageTimeout    time.Duration `env:"AUTH_STAGE_TIMEOUT" json:"authTimeout"`
	ServerEncryptionKey string        `env:"SERVER_ENCRYPTION_KEY" json:"serverSecret"`
	DomainName          string        `env:"DOMAIN_NAME" json:"domainName"`
	SMTPHost            string        `env:"SMTP_HOST" json:"smtpHost"`
	SMTPPort            int           `env:"SMTP_PORT" json:"smtpPort"`
	ServerEMail         string        `env:"SERVER_EMAIL" json:"serverEmail"`
	DatabaseURI         string        `env:"DATABASE_URI" json:"dbDsn"`
	MaxConns            int           `env:"DATABASE_MAX_CONNS" json:"dbMaxConns"`
	MaxConnLifetime     time.Duration `env:"DATABASE_MAX_CONN_LIFE_TIME" json:"dbMaxConnLifeTime"`
	MaxConnIdleTime     time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" json:"dbMaxConnIdleTime"`
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
		MaxConns:            ServerDefaultMaxConns,
		DatabaseURI:         ServerDefaultDatabaseURI,
		MaxConnLifetime:     ServerDefaultMaxConnLifetime,
		MaxConnIdleTime:     ServerDefaultMaxConnIdleTime,
	}
}

func LoadConfigFromFile(fileName string, srcConf any) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err := json.NewDecoder(f).Decode(srcConf); err != nil {
		panic(err)
	}
}

func LoadServConf(flagSet *flag.FlagSet) (*ServerConf, error) {
	srvConf := defaultServConf()

	var configFileName string
	flagSet.StringVar(&configFileName, "c", "", "config file") // config file short format
	flagSet.StringVar(&configFileName, "config", "", "config file")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if configFileName != "" {
		LoadConfigFromFile(configFileName, srvConf)
	}

	err := env.Parse(srvConf)
	if err != nil {
		return nil, err
	}

	return srvConf, nil
}
