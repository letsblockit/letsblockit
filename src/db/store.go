package db

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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
	return pgx.BeginFunc(c, s.pool, func(tx pgx.Tx) error {
		return f(c, New(tx))
	})
}

func Connect(databaseUrl, poolOptions string, dsd statsd.ClientInterface) (Store, error) {
	pool, err := pgxpool.New(context.Background(), databaseUrl+poolOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot connect: %w", err)
	}
	return &pgxStore{
		Queries: New(&instrumentedDB{pool, dsd}),
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

type instrumentedDB struct {
	db  DBTX
	dsd statsd.ClientInterface
}

func (i instrumentedDB) Exec(ctx context.Context, q string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()
	tag, err := i.db.Exec(ctx, q, args...)
	_ = i.dsd.Distribution("letsblockit.pg_request_duration", float64(time.Since(start).Nanoseconds()),
		[]string{"type:exec"}, 1)
	return tag, err
}

func (i instrumentedDB) Query(ctx context.Context, q string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()
	rows, err := i.db.Query(ctx, q, args...)
	_ = i.dsd.Distribution("letsblockit.pg_request_duration", float64(time.Since(start).Nanoseconds()),
		[]string{"type:query"}, 1)
	return rows, err
}

func (i instrumentedDB) QueryRow(ctx context.Context, q string, args ...interface{}) pgx.Row {
	start := time.Now()
	row := i.db.QueryRow(ctx, q, args...)
	_ = i.dsd.Distribution("letsblockit.pg_request_duration", float64(time.Since(start).Nanoseconds()),
		[]string{"type:query_row"}, 1)
	return row
}
