package repository

import (
	"context"
	"database/sql"

	"github.com/klik/fcos-kernel/pkg/database"
)

// sqlExecutor defines an interface that is satisfied by both *sql.DB and *sql.Tx.
// This allows repository methods to work with or without a transaction.
type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// getExecutor retrieves a transactional executor from the context if one exists,
// otherwise it returns the base database connection.
func getExecutor(ctx context.Context, db *sql.DB) sqlExecutor {
	if tx, ok := database.GetTx(ctx); ok {
		return tx
	}
	return db
}
