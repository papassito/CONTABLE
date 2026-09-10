package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/klik/fcos-kernel/internal/domain"
)

type LedgerRepository interface {
	RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error
	GetMovementsByAccount(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error)
	GetTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error)
}

type sqlLedgerRepository struct {
	db *sql.DB
}

func NewLedgerRepository(db *sql.DB) LedgerRepository {
	return &sqlLedgerRepository{db: db}
}

func (r *sqlLedgerRepository) RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error {
	exec := getExecutor(ctx, r.db)
	for _, entry := range entries {
		_, err := exec.ExecContext(ctx,
			"INSERT INTO ledger_entries (account_id, account_code, entry_date, debit, credit, reference) VALUES (?, ?, ?, ?, ?, ?)",
			entry.AccountID, entry.AccountCode, entry.EntryDate, entry.Debit, entry.Credit, entry.Reference)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *sqlLedgerRepository) GetMovementsByAccount(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error) {
	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx,
		"SELECT id, account_id, account_code, entry_date, debit, credit, reference FROM ledger_entries WHERE account_id = ? AND entry_date BETWEEN ? AND ? ORDER BY entry_date ASC, id ASC",
		accountID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []domain.LedgerEntry
	for rows.Next() {
		var mov domain.LedgerEntry
		if err := rows.Scan(&mov.ID, &mov.AccountID, &mov.AccountCode, &mov.EntryDate, &mov.Debit, &mov.Credit, &mov.Reference); err != nil {
			return nil, err
		}
		movements = append(movements, mov)
	}
	return movements, nil
}

func (r *sqlLedgerRepository) GetTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error) {
	exec := getExecutor(ctx, r.db)
	query := `
		WITH initial_balances AS (
			SELECT
				account_id,
				SUM(debit) as total_debit,
				SUM(credit) as total_credit
			FROM ledger_entries
			WHERE entry_date < ?
			GROUP BY account_id
		),
		period_movements AS (
			SELECT
				account_id,
				SUM(debit) as total_debit,
				SUM(credit) as total_credit
			FROM ledger_entries
			WHERE entry_date BETWEEN ? AND ?
			GROUP BY account_id
		)
		SELECT
			a.code, a.name,
			COALESCE(ib.total_debit, 0), COALESCE(ib.total_credit, 0),
			COALESCE(pm.total_debit, 0), COALESCE(pm.total_credit, 0)
		FROM accounts a
		LEFT JOIN initial_balances ib ON a.id = ib.account_id
		LEFT JOIN period_movements pm ON a.id = pm.account_id
		WHERE a.accepts_move = TRUE AND (ib.total_debit IS NOT NULL OR pm.total_debit IS NOT NULL)
		ORDER BY a.code;
`
	rows, err := exec.QueryContext(ctx, query, from, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.TrialBalanceItem
	for rows.Next() {
		var item domain.TrialBalanceItem
		var initialDebit, initialCredit, periodDebit, periodCredit int64

		if err := rows.Scan(&item.AccountCode, &item.AccountName, &initialDebit, &initialCredit, &periodDebit, &periodCredit); err != nil {
			return nil, err
		}

		item.TotalDebit = periodDebit
		item.TotalCredit = periodCredit

		item.InitialBalance = initialDebit - initialCredit
		item.FinalBalance = item.InitialBalance + item.TotalDebit - item.TotalCredit
		items = append(items, item)
	}
	return items, nil
}
