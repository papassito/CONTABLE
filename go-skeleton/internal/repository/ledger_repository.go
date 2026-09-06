package repository

import (
    "context"
    "github.com/klik/fcos-kernel/internal/domain"
)

// LedgerRepository define las operaciones de persistencia del Libro Mayor.
type LedgerRepository interface {
    RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error
    GetBalance(ctx context.Context, accountID string) (int64, error)
    SaveEntry(ctx context.Context, entry *domain.LedgerEntry) error
}
