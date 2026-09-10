package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/klik/fcos-kernel/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, acc *domain.Account) error
	GetByID(ctx context.Context, tenantID, id string) (*domain.Account, error)
	GetByCode(ctx context.Context, tenantID, code string) (*domain.Account, error)
	UpdateBalance(ctx context.Context, tenantID, id string, delta domain.Cents) error
	List(ctx context.Context, tenantID string) ([]*domain.Account, error)
}

type sqlAccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &sqlAccountRepository{db: db}
}

func (r *sqlAccountRepository) Create(ctx context.Context, acc *domain.Account) error {
	exec := getExecutor(ctx, r.db)
	const query = `
		INSERT INTO accounts (
			tenant_id,
			id,
			code,
			name,
			parent_id,
			accepts_move,
			status,
			balance_cents
		)
		VALUES (?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)
	`
	_, err := exec.ExecContext(
		ctx,
		query,
		acc.TenantID,
		acc.ID,
		acc.Code,
		acc.Name,
		acc.ParentID,
		acc.AcceptsMove,
		acc.Status,
		acc.Balance,
	)
	return err
}

func (r *sqlAccountRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	const query = `
		SELECT
			id,
			tenant_id,
			code,
			name,
			COALESCE(parent_id, ''),
			accepts_move,
			status,
			balance_cents
		FROM accounts
		WHERE tenant_id = ?
		  AND id = ?
	`
	var acc domain.Account
	err := exec.QueryRowContext(ctx, query, tenantID, id).Scan(
		&acc.ID,
		&acc.TenantID,
		&acc.Code,
		&acc.Name,
		&acc.ParentID,
		&acc.AcceptsMove,
		&acc.Status,
		&acc.Balance,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *sqlAccountRepository) GetByCode(ctx context.Context, tenantID, code string) (*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	const query = `
		SELECT
			id,
			tenant_id,
			code,
			name,
			COALESCE(parent_id, ''),
			accepts_move,
			status,
			balance_cents
		FROM accounts
		WHERE tenant_id = ?
		  AND code = ?
	`
	var acc domain.Account
	err := exec.QueryRowContext(ctx, query, tenantID, code).Scan(
		&acc.ID,
		&acc.TenantID,
		&acc.Code,
		&acc.Name,
		&acc.ParentID,
		&acc.AcceptsMove,
		&acc.Status,
		&acc.Balance,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *sqlAccountRepository) List(ctx context.Context, tenantID string) ([]*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	const query = `
		SELECT
			id,
			tenant_id,
			code,
			name,
			COALESCE(parent_id, ''),
			accepts_move,
			status,
			balance_cents
		FROM accounts
		WHERE tenant_id = ?
		ORDER BY code
	`
	rows, err := exec.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var acc domain.Account
		err := rows.Scan(
			&acc.ID,
			&acc.TenantID,
			&acc.Code,
			&acc.Name,
			&acc.ParentID,
			&acc.AcceptsMove,
			&acc.Status,
			&acc.Balance,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}
	return accounts, nil
}

func (r *sqlAccountRepository) UpdateBalance(ctx context.Context, tenantID, id string, delta domain.Cents) error {
	exec := getExecutor(ctx, r.db)
	const query = `
		UPDATE accounts
		SET balance_cents = balance_cents + ?
		WHERE tenant_id = ?
		  AND id = ?
	`
	_, err := exec.ExecContext(ctx, query, delta, tenantID, id)
	return err
}
