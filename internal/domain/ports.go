package domain

import "context"

// AccountRepository describe los contratos del almacén persistente de cuentas.
type AccountRepository interface {
	Create(ctx context.Context, acc *Account) error
	GetByID(ctx context.Context, tenantID, id string) (*Account, error)
	GetByCode(ctx context.Context, tenantID, code string) (*Account, error)
	UpdateBalance(ctx context.Context, tenantID, id string, delta Cents) error
	List(ctx context.Context, tenantID string) ([]*Account, error)
}

// JournalRepository describe los contratos del almacén de asientos de diario.
type JournalRepository interface {
	Create(ctx context.Context, entry *JournalEntry) error
	GetByID(ctx context.Context, tenantID, id string) (*JournalEntry, error)
	Update(ctx context.Context, entry *JournalEntry) error
}
