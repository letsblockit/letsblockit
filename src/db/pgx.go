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
	format = fmt.Sprintf("db[%s]: %s", m.db, format)
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

func (d *pgxStore) RunTx(e echo.Context, f TxFunc) error {
	c := e.Request().Context()
	return d.pool.BeginFunc(c, func(tx pgx.Tx) error {
		return f(c, New(tx))
	})
}

func Connect(host, database string) (Store, error) {
	dsn := fmt.Sprintf("postgresql:///%s?host=%s", database, host)
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &pgxStore{
		Queries: New(pool),
		pool:    pool,
	}, nil
}

func Migrate(host, database string) error {
	db, err := (&mpgx.Postgres{}).Open(fmt.Sprintf("pgx:///%s?host=%s", database, host))
	if err != nil {
		return err
	}
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	instance, err := migrate.NewWithInstance("embedded", source, "pgx", db)
	if err != nil {
		return err
	}
	instance.Log = &migrateLogger{db: database}
	if err = instance.Up(); err != nil {
		if err == migrate.ErrNoChange {
			instance.Log.Printf("No database migration to run\n")
		} else {
			return err
		}
	}
	if err = source.Close(); err != nil {
		return err
	}
	return db.Close()
}
