package core_pgx_pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapError(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
		pgxUniqueViolationErrorCode    = "23505"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
		if pgErr.Code == pgxUniqueViolationErrorCode {
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrUniqueViolation,
			)
		}
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}

type pgxTx struct {
	tx pgx.Tx
}

func (t pgxTx) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapError(err)
	}

	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	row := t.tx.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (t pgxTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := t.tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapError(err)
	}

	return pgxCommandTag{cmdTag}, nil
}

func (t pgxTx) Commit(ctx context.Context) error {
	return mapError(t.tx.Commit(ctx))
}

func (t pgxTx) Rollback(ctx context.Context) error {
	return mapError(t.tx.Rollback(ctx))
}
