package postgres

import (
	"context"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
)

var once sync.Once
var st *storage

func NewStorage(appCnf context.Context, conf *config.ServerConf) *storage {
	once.Do(func() {
		st = initializePGXConf(appCnf, conf)
	})

	return st
}

func initializePGXConf(ctx context.Context, conf *config.ServerConf) *storage {

	logger := domain.GetApplicationLogger()

	logger.Infow("initializePGXConf", "status", "start")

	pConf, err := pgxpool.ParseConfig(conf.DatabaseURI)
	if err != nil {
		panic(err)
	}

	// Конфигурация по мотивам
	// https://habr.com/ru/companies/oleg-bunin/articles/461935/
	pConf.MaxConns = int32(conf.MaxConns)
	pConf.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	pConf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pConf.MaxConnLifetime = conf.MaxConnIdleTime
	pConf.MaxConnIdleTime = conf.MaxConnIdleTime

	pConf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewLogAdapter(logger),
		LogLevel: tracelog.LogLevelError,
	}

	pPool, err := pgxpool.NewWithConfig(ctx, pConf)

	if err != nil {
		panic(err)
	}

	st = &storage{
		pPool: pPool,
	}

	if err := st.init(ctx, logger); err != nil {
		panic(err)
	}

	logger.Infow("initializePGXConf", "status", "complete")
	return st
}

type storage struct {
	pPool *pgxpool.Pool
}

func (st *storage) init(ctx context.Context, logger domain.Logger) error {
	logger.Infow("init", "status", "start")

	tx, err := st.pPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err = tx.Exec(ctx, `
	create table user_info (
		userId bigserial,
		email text not null,
		pass_hash text not null,
		pass_salt text not null,
		otp_key text not null,
		master_key text not null,
		master_hint text not null,
		hello_encrypted text not null,
		primary key(userId)
	);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create unique index idx_user_info_email on user_info(email);`); err != nil {
		panic(err)
	}

	return tx.Commit(ctx)
}

func (st *storage) Ping(ctx context.Context) error {
	return st.pPool.Ping(ctx)
}

func (st *storage) Close() {
	st.pPool.Close()
}