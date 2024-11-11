package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresAdapter[T any] struct {
	db *pgx.Conn
}

func NewPostgresAdapter[T any](pgxConn *pgx.Conn) *PostgresAdapter[T] {
	return &PostgresAdapter[T]{db: pgxConn}
}

func (pa PostgresAdapter[T]) Query(context context.Context, sql string, args ...any) ([]T, error) {
	rows, err := pa.db.Query(context, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, fmt.Errorf("collectRows failed: %w", err)
	}

	return results, nil
}

func (pa PostgresAdapter[T]) QueryRow(context context.Context, sql string, args []any, dest ...any) error {
	row := pa.db.QueryRow(context, sql, args...)

	return row.Scan(dest...)
}
