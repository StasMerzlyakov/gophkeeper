// Package integrations contains integration server tests
package main

import (
	"flag"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	flagGophKeeperTlsKey           string
	flagGophKeeperTlsCert          string
	flagGophKeeperTlsCaCert        string
	flagGophKeeperServerSecret     string
	flagGophKeeperServerBinaryPath string
)

var postgresContainer *postgres.PostgresContainer
var smtpServer *smtpmock.Server

func init() {
	flag.StringVar(&flagGophKeeperServerBinaryPath, "gophkeeper-binary-path", "", "GophKeeper binary")
	flag.StringVar(&flagGophKeeperTlsKey, "gophkeeper-tls-key", "", "GophKeeper tls key file")
	flag.StringVar(&flagGophKeeperTlsCert, "gophkeeper-tls-cert", "", "GophKeeper tls cert file")
	flag.StringVar(&flagGophKeeperTlsCaCert, "gophkeeper-tls-ca-cert", "", "GophKeeper tls ca file")
	flag.StringVar(&flagGophKeeperServerSecret, "gophkeeper-server-secret", "", "GophKeeper server secret")
}
