package database

import (
	"context"
	"database/sql"
	"fmt"
)

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}

type sqlUnitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) UnitOfWork {
	return &sqlUnitOfWork{db: db}
}

type contextKey string

const txKey contextKey = "tx_context_key"

func (uow *sqlUnitOfWork) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := uow.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	txCtx := context.WithValue(ctx, txKey, tx)
	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error en ejecucion: %v (rollback fallido: %v)", err, rbErr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción (commit): %w", err)
	}
	return nil
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}
