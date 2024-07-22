package postgres_test

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresContainer *postgres.PostgresContainer

func TestMain(m *testing.M) {
	ctx, cancelFN := context.WithCancel(context.Background())
	defer func() {
		cancelFN()
	}()

	setup(ctx)
	code := m.Run()
	shutdown(ctx)
	os.Exit(code)
}

func shutdown(ctx context.Context) {
	if postgresContainer != nil {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}
}

const dbName = "users"
const dbUser = "user"
const dbPassword = "password"

var pPool *pgxpool.Pool

func setup(ctx context.Context) {

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
		log.Fatalf(err.Error())
	}
}

var once sync.Once

func clear(ctx context.Context) error {

	once.Do(func() {
		connString, _ := postgresContainer.ConnectionString(ctx)
		pConf, _ := pgxpool.ParseConfig(connString)
		pPool, _ = pgxpool.NewWithConfig(ctx, pConf)
	})

	_, err := pPool.Exec(ctx, "delete from user_info")
	return err
}
