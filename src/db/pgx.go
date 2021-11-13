package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

var NotFound = pgx.ErrNoRows

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
