//go:build autotest
// +build autotest

package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	postgresUp(ctx)
	defer postgresDown(ctx)

	smtpServiceUp(ctx)
	defer smtpServiceDown(ctx)

	code := m.Run()
	os.Exit(code)
}

func smtpServiceUp(context.Context) {
	smtpServer = smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := smtpServer.Start(); err != nil {
		log.Fatalf("failed to start smtp mock service: %s", err.Error())
	}
}

func smtpServiceDown(context.Context) {
	if err := smtpServer.Stop(); err != nil {
		log.Fatalf("failed to terminate smtp mock service: %s", err)
	}
}

const dbName = "users"
const dbUser = "user"
const dbPassword = "password"

func postgresUp(ctx context.Context) {
	var err error
	postgresContainer, err = postgres.Run(ctx,
		"docker.io/postgres:15.2-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %s", err)
	}
}

func postgresDown(ctx context.Context) {
	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate postgres container: %s", err)
	}
}
