package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env"
)

const (
	ClientDefaultServerAddres      = "localhost:3200"
	ClientDefaultCACert            = "../../keys/ca-cert.pem"
	ClientDefaultInterationTimeout = 20 * time.Second
)

type ClientConf struct {
	ServerAddress     string        `env:"SERVER_ADDRESS" json:"serverAddress"`
	InterationTimeout time.Duration `env:"INTERACTION_TIMEOUT" json:"interactionTimeout"`
	CACert            string        `env:"CA_CERT" json:"caCert"`
}

func defaultClientConf() *ClientConf {
	return &ClientConf{
		ServerAddress:     ClientDefaultServerAddres,
		CACert:            ClientDefaultCACert,
		InterationTimeout: ClientDefaultInterationTimeout,
	}
}

func LoadClientConf(flagSet *flag.FlagSet) (*ClientConf, error) {
	clntConf := defaultClientConf()

	var configFileName string
	flagSet.StringVar(&configFileName, "c", "", "config file") // config file short format
	flagSet.StringVar(&configFileName, "config", "", "config file")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if configFileName != "" {
		LoadConfigFromFile(configFileName, clntConf)
	}

	err := env.Parse(clntConf)
	if err != nil {
		return nil, err
	}

	return clntConf, nil
}
