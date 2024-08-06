package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

const portSmtpServSMTP = "1025"
const portSmtpServHTTP = "1080"

func smtpServiceUp(ctx context.Context) {
	req := testcontainers.ContainerRequest{
		Image:        "haravich/fake-smtp-server",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", portSmtpServSMTP), fmt.Sprintf("%s/tcp", portSmtpServHTTP)},
		WaitingFor:   wait.ForLog("http://0.0.0.0:1080"),
	}
	var err error
	smtpServer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("Could not start smtpServer: %s", err)
	}
}

func smtpServiceDown(ctx context.Context) {

	if err := smtpServer.Terminate(ctx); err != nil {
		log.Fatalf("Could not stop smtpServer: %s", err)
	}
}
