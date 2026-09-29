package core_pgx_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration // operation timeout
}

func NewPool(
	ctx context.Context,
	cfg Config,
) (*Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pgxCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: cfg.TimeOut,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{cmdTag}, nil
}

func (p *Pool) Begin(ctx context.Context) (core_postgres_pool.Tx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	return pgxTx{tx: tx}, nil
}
