package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
			log.Printf("[UOW-ROLLBACK-FATAL] Failed to rollback transaction: %v. Original error: %v", rbErr, err)
			return fmt.Errorf("error de persistencia interna: la operación no pudo ser completada")
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
