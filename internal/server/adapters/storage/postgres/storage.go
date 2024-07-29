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
	create table if not exists user_info (
		user_id bigserial,
		email text not null,
		pass_hash text not null,
		pass_salt text not null,
		otp_key text not null,
		master_hint text not null,
		hello_encrypted text not null,
		primary key(user_id)
	);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create unique index if not exists idx_user_info_email on user_info(email);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `
	create table if not exists bank_card (
		id bigserial,
		number text not null,
		content text not null,
		user_id bigint not null references user_info (user_id) on delete cascade,
		primary key(id)
	);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create unique index if not exists idx_bank_card_number on bank_card(number,user_id);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create index if not exists idx_bank_card_user_id on bank_card(user_id);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `
	create table if not exists user_password_data (
		id bigserial,
		hint text not null,
		content text not null,
		user_id bigint not null references user_info (user_id) on delete cascade,
		primary key(id)
	);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create unique index if not exists idx_user_password_data_hint on user_password_data(hint,user_id);`); err != nil {
		panic(err)
	}

	if _, err = tx.Exec(ctx, `create index if not exists idx_user_password_data_user_id on user_password_data(user_id);`); err != nil {
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
