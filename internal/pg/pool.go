package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type BasicPool struct {
	pool *pgxpool.Pool
}

var _ Pool = (*BasicPool)(nil)

func NewBasicPool(ctx context.Context, config *Config) (Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(config.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	dbConfig.MinConns = config.MinConns
	dbConfig.MaxConns = config.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &BasicPool{
		pool: pool,
	}, nil
}

func (p *BasicPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *BasicPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *BasicPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}
