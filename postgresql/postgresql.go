package postgresql

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func New() *PostgreSQL {
	return &PostgreSQL{}
}

func (pg *PostgreSQL) Connect(ctx context.Context, connString string) error {
	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return err
	}
	pg.pool = pool

	return nil
}

func (pg PostgreSQL) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	ct, err := pg.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (pg PostgreSQL) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := pg.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (pg PostgreSQL) ScanAll(dst any, rows pgx.Rows) error {
	if err := pgxscan.ScanAll(dst, rows); err != nil {
		return err
	}

	return nil
}

var ErrNotFound = errors.New("item not found")

func (pg PostgreSQL) ScanOne(dst any, rows pgx.Rows) error {
	if err := pgxscan.ScanOne(dst, rows); err != nil {
		switch {
		case pgxscan.NotFound(err):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
