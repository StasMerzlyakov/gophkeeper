package config

import "time"

type GophKeeperConf struct {
	TokenExp    time.Duration `env:"JWT_EXP"`
	TokenSecret string        `env:"JWT_SECRET"`
	RunAddress  string        `env:"RUN_ADDRESS"`
	TLSKey      string        `env:"TLS_KEY"`
}
