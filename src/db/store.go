package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

var (
	//go:embed migrations/*.up.sql
	migrations embed.FS
	NotFound   = pgx.ErrNoRows
)

type migrateLogger struct {
	db string
}

func (m migrateLogger) Printf(format string, v ...interface{}) {
	format = fmt.Sprintf("migrate[%s]: %s", m.db, format)
	fmt.Printf(format, v...)
}

func (m migrateLogger) Verbose() bool {
	return true
}

type Store interface {
	Querier
	RunTx(e echo.Context, f TxFunc) error
}

type TxFunc func(context.Context, Querier) error

type pgxStore struct {
	*Queries
	pool *pgxpool.Pool
}

func (s *pgxStore) RunTx(e echo.Context, f TxFunc) error {
	c := e.Request().Context()
	return s.pool.BeginFunc(c, func(tx pgx.Tx) error {
		return f(c, New(tx))
	})
}

func Connect(databaseUrl, poolOptions string) (Store, error) {
	pool, err := pgxpool.Connect(context.Background(), databaseUrl+poolOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot connect: %w", err)
	}
	return &pgxStore{
		Queries: New(pool),
		pool:    pool,
	}, nil
}

func Migrate(databaseUrl string) error {
	db, err := (&mpgx.Postgres{}).Open(databaseUrl)
	if err != nil {
		return fmt.Errorf("cannot open db for migration: %w", err)
	}
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("cannot open migration source: %w", err)
	}

	instance, err := migrate.NewWithInstance("embedded", source, "pgx", db)
	if err != nil {
		return fmt.Errorf("cannot init migrator: %w", err)
	}
	if c, err := pgx.ParseConfig(databaseUrl); err == nil {
		instance.Log = &migrateLogger{db: c.Database}
	}

	if err = instance.Up(); err != nil {
		if err == migrate.ErrNoChange {
			instance.Log.Printf("No database migration to run\n")
		} else {
			return fmt.Errorf("migration error: %w", err)
		}
	}
	if err = source.Close(); err != nil {
		return fmt.Errorf("cannot close migration source: %w", err)
	}
	return db.Close()
}
