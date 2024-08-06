package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env"
)

const (
	ClientDefaultServerAddres      = "localhost:3200"
	ClientDefaultCACert            = "../../keys/ca-cert.pem"
	ClientDefaultLogFile           = "./client.log"
	ClientDefaultInterationTimeout = 3 * time.Second
)

type ClientConf struct {
	ServerAddress     string        `env:"SERVER_ADDRESS" json:"serverAddress,omitempty"`
	InterationTimeout time.Duration `env:"INTERACTION_TIMEOUT" json:"interactionTimeout,omitempty"`
	CACert            string        `env:"CA_CERT" json:"caCert"`
	LogFile           string        `env:"LOG_FILE" json:"logFile"`
}

func defaultClientConf() *ClientConf {
	return &ClientConf{
		ServerAddress:     ClientDefaultServerAddres,
		CACert:            ClientDefaultCACert,
		InterationTimeout: ClientDefaultInterationTimeout,
		LogFile:           ClientDefaultLogFile,
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

func (cCnf *ClientConf) UnmarshalJSON(data []byte) (err error) {
	// default json fail on time.Duration
	type ClientConfAlias ClientConf
	aliasValue := &struct {
		*ClientConfAlias
		// redefine field
		InterationTimeout string `json:"interactionTimeout,omitempty"`
	}{

		ClientConfAlias: (*ClientConfAlias)(cCnf),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	if aliasValue.InterationTimeout != "" {
		tm, err := time.ParseDuration(aliasValue.InterationTimeout)
		if err != nil {
			return err
		}
		cCnf.InterationTimeout = tm
	}
	return
}
