package adapters

import (
	"context"
)

type RelationalDatabaseTable[T any] interface {
	Query(context.Context, string, ...any) ([]T, error)
	QueryRow(context.Context, string, []any, ...any) error
}
