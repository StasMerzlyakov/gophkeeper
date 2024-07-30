package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env"
)

const (
	ServerDefaultPort             = ":3200"
	ServerDefaultTLSKey           = "../../keys/server-key.pem"
	ServerDefaultTLSCert          = "../../keys/server-cert.pem"
	ServerDefaultTokenExp         = 3 * time.Hour
	ServerDefaultTokenSecret      = "secret"
	ServerDefaultAuthStageTimeout = 5 * time.Minute
	ServerDefaultServerSecret     = "ServerKey!12>Au{mL736}"
	ServerDefaultDomainName       = "localhost"
	ServerDefaultSMTPHost         = "localhost"
	ServerDefaultSMTPPort         = 25
	ServerDefaultServerEMail      = "localhost@localdomain"
	ServerDefaultMaxConns         = 5
	ServerDefaultMaxConnLifetime  = 5 * time.Minute
	ServerDefaultMaxConnIdleTime  = 5 * time.Minute
	ServerDefaultDatabaseDN       = "postgres://user:user@localhost:5432/gophkeeper"
)

type ServerConf struct {
	Port             string        `env:"PORT" json:"port,omitempty"`
	TLSKey           string        `env:"TLS_KEY" json:"tlsKey,omitempty"`
	TLSCert          string        `env:"TLS_CERT" json:"tlsCert,omitempty"`
	TokenExp         time.Duration `env:"JWT_EXP" json:"jwtExp,omitempty"`
	TokenSecret      string        `env:"JWT_SECRET" json:"jwtSecret,omitempty"`
	AuthStageTimeout time.Duration `env:"AUTH_STAGE_TIMEOUT" json:"authTimeout,omitempty"`
	ServerSecret     string        `env:"SERVER_SECRET" json:"serverSecret,omitempty"`
	DomainName       string        `env:"DOMAIN_NAME" json:"domainName,omitempty"`
	SMTPHost         string        `env:"SMTP_HOST" json:"smtpHost,omitempty"`
	SMTPPort         int           `env:"SMTP_PORT" json:"smtpPort,omitempty"`
	ServerEMail      string        `env:"SERVER_EMAIL" json:"serverEmail,omitempty"`
	DatabaseDN       string        `env:"DATABASE_DN" json:"dbDN,omitempty"`
	MaxConns         int           `env:"DATABASE_MAX_CONNS" json:"dbMaxConns,omitempty"`
	MaxConnLifetime  time.Duration `env:"DATABASE_MAX_CONN_LIFE_TIME" json:"dbMaxConnLifeTime,omitempty"`
	MaxConnIdleTime  time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" json:"dbMaxConnIdleTime,omitempty"`
}

func defaultServConf() *ServerConf {
	return &ServerConf{
		Port:             ServerDefaultPort,
		TLSKey:           ServerDefaultTLSKey,
		TLSCert:          ServerDefaultTLSCert,
		TokenExp:         ServerDefaultTokenExp,
		TokenSecret:      ServerDefaultTokenSecret,
		AuthStageTimeout: ServerDefaultAuthStageTimeout,
		ServerSecret:     ServerDefaultServerSecret,
		DomainName:       ServerDefaultDomainName,
		SMTPHost:         ServerDefaultSMTPHost,
		SMTPPort:         ServerDefaultSMTPPort,
		ServerEMail:      ServerDefaultServerEMail,
		MaxConns:         ServerDefaultMaxConns,
		DatabaseDN:       ServerDefaultDatabaseDN,
		MaxConnLifetime:  ServerDefaultMaxConnLifetime,
		MaxConnIdleTime:  ServerDefaultMaxConnIdleTime,
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

func (sCnf *ServerConf) UnmarshalJSON(data []byte) (err error) {
	// default json fail on time.Duration
	type ServerConfAlias ServerConf
	aliasValue := &struct {
		*ServerConfAlias
		// redefine field
		TokenExp         string `json:"jwtExp,omitempty"`
		AuthStageTimeout string `json:"authTimeout,omitempty"`
		MaxConnLifetime  string `json:"dbMaxConnLifeTime,omitempty"`
		MaxConnIdleTime  string `json:"dbMaxConnIdleTime,omitempty"`
	}{

		ServerConfAlias: (*ServerConfAlias)(sCnf),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	if aliasValue.TokenExp != "" {
		tm, err := time.ParseDuration(aliasValue.TokenExp)
		if err != nil {
			return err
		}
		sCnf.TokenExp = tm
	}

	if aliasValue.AuthStageTimeout != "" {
		tm, err := time.ParseDuration(aliasValue.AuthStageTimeout)
		if err != nil {
			return err
		}
		sCnf.AuthStageTimeout = tm
	}

	if aliasValue.MaxConnLifetime != "" {
		tm, err := time.ParseDuration(aliasValue.MaxConnLifetime)
		if err != nil {
			return err
		}
		sCnf.MaxConnLifetime = tm
	}

	if aliasValue.MaxConnIdleTime != "" {
		tm, err := time.ParseDuration(aliasValue.MaxConnIdleTime)
		if err != nil {
			return err
		}
		sCnf.MaxConnIdleTime = tm
	}

	return
}
