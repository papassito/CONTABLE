package repository

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error)
}
