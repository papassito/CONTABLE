package repository

import (
	"context"
	"database/sql"

	"github.com/klik/fcos-kernel/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByCode(ctx context.Context, code string) (*domain.Account, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id string) error
	UpdateBalance(ctx context.Context, id string, amount int64) error
}

type sqlAccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &sqlAccountRepository{db: db}
}

func (r *sqlAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	exec := getExecutor(ctx, r.db)
	query := "INSERT INTO accounts (id, code, name, type, parent_id, level, accepts_move, status, current_bal, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := exec.ExecContext(ctx, query, account.ID, account.Code, account.Name, account.Type, account.ParentID, account.Level, account.AcceptsMove, account.Status, account.CurrentBal, account.CreatedAt, account.UpdatedAt)
	return err
}

func (r *sqlAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, code, name, type, parent_id, level, accepts_move, status, current_bal, created_at, updated_at FROM accounts WHERE id = ?", id)
	var acc domain.Account
	err := row.Scan(&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.ParentID, &acc.Level, &acc.AcceptsMove, &acc.Status, &acc.CurrentBal, &acc.CreatedAt, &acc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *sqlAccountRepository) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, code, name, type, parent_id, level, accepts_move, status, current_bal, created_at, updated_at FROM accounts WHERE code = ?", code)
	var acc domain.Account
	err := row.Scan(&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.ParentID, &acc.Level, &acc.AcceptsMove, &acc.Status, &acc.CurrentBal, &acc.CreatedAt, &acc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *sqlAccountRepository) List(ctx context.Context, filter map[string]interface{}) ([]*domain.Account, error) {
	exec := getExecutor(ctx, r.db)
	// Basic implementation without filtering for now.
	rows, err := exec.QueryContext(ctx, "SELECT id, code, name, type, parent_id, level, accepts_move, status, current_bal, created_at, updated_at FROM accounts ORDER BY code")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var acc domain.Account
		err := rows.Scan(&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.ParentID, &acc.Level, &acc.AcceptsMove, &acc.Status, &acc.CurrentBal, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}
	return accounts, nil
}

func (r *sqlAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	exec := getExecutor(ctx, r.db)
	query := "UPDATE accounts SET name = ?, type = ?, parent_id = ?, level = ?, accepts_move = ?, status = ?, updated_at = ? WHERE id = ?"
	_, err := exec.ExecContext(ctx, query, account.Name, account.Type, account.ParentID, account.Level, account.AcceptsMove, account.Status, account.UpdatedAt, account.ID)
	return err
}

func (r *sqlAccountRepository) Delete(ctx context.Context, id string) error {
	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, "DELETE FROM accounts WHERE id = ?", id)
	return err
}

func (r *sqlAccountRepository) UpdateBalance(ctx context.Context, id string, amount int64) error {
	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, "UPDATE accounts SET current_bal = current_bal + ? WHERE id = ?", amount, id)
	return err
}
