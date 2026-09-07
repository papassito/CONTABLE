package repository

import (
	"context"
	"database/sql"

	"github.com/klik/fcos-kernel/internal/domain"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error)
	CheckPeriodStatus(ctx context.Context, year int, month int) (string, error)
}

type sqlJournalRepository struct {
	db *sql.DB
}

func NewJournalRepository(db *sql.DB) JournalRepository {
	return &sqlJournalRepository{db: db}
}

func (r *sqlJournalRepository) Create(ctx context.Context, entry *domain.JournalEntry) error {
	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, "INSERT INTO journal_entries (id, number, date, concept, reference, status, total_debit, total_credit) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		entry.ID, entry.Number, entry.Date, entry.Concept, entry.Reference, entry.Status, entry.TotalDebit, entry.TotalCredit)
	if err != nil {
		return err
	}
	for _, line := range entry.Lines {
		_, err = exec.ExecContext(ctx, "INSERT INTO journal_lines (id, journal_entry_id, account_id, account_code, description, debit, credit, third_party_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			line.ID, entry.ID, line.AccountID, line.AccountCode, line.Description, line.Debit, line.Credit, line.ThirdPartyID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *sqlJournalRepository) GetByID(ctx context.Context, id string) (*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries WHERE id = ?", id)
	var entry domain.JournalEntry
	err := row.Scan(&entry.ID, &entry.Number, &entry.Date, &entry.Concept, &entry.Reference, &entry.Status, &entry.TotalDebit, &entry.TotalCredit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *sqlJournalRepository) GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries WHERE number = ?", number)
	var entry domain.JournalEntry
	err := row.Scan(&entry.ID, &entry.Number, &entry.Date, &entry.Concept, &entry.Reference, &entry.Status, &entry.TotalDebit, &entry.TotalCredit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *sqlJournalRepository) List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, "SELECT id, number, date, concept, reference, status, total_debit, total_credit FROM journal_entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.JournalEntry
	for rows.Next() {
		var entry domain.JournalEntry
		if err := rows.Scan(&entry.ID, &entry.Number, &entry.Date, &entry.Concept, &entry.Reference, &entry.Status, &entry.TotalDebit, &entry.TotalCredit); err != nil {
			return nil, err
		}
		list = append(list, &entry)
	}
	return list, nil
}

func (r *sqlJournalRepository) UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error {
	exec := getExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, "UPDATE journal_entries SET status = ? WHERE id = ?", status, id)
	return err
}

func (r *sqlJournalRepository) GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error) {
	exec := getExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, "SELECT id, journal_entry_id, account_id, account_code, description, debit, credit, third_party_id FROM journal_lines WHERE journal_entry_id = ?", entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []domain.JournalLine
	for rows.Next() {
		var line domain.JournalLine
		if err := rows.Scan(&line.ID, &line.JournalEntryID, &line.AccountID, &line.AccountCode, &line.Description, &line.Debit, &line.Credit, &line.ThirdPartyID); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func (r *sqlJournalRepository) CheckPeriodStatus(ctx context.Context, year int, month int) (string, error) {
	exec := getExecutor(ctx, r.db)
	var status string
	err := exec.QueryRowContext(ctx, "SELECT status FROM compliance_periods WHERE period_year = ? AND period_month = ?", year, month).Scan(&status)
	if err == sql.ErrNoRows {
		return "OPEN", nil
	}
	return status, err
}
