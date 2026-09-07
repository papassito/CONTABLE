package domain

import "time"

// LedgerEntry represents a movement in the General Ledger.
type LedgerEntry struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	AccountCode string    `json:"account_code"`
	EntryDate   time.Time `json:"entry_date"`
	Debit       int64     `json:"debit"`
	Credit      int64     `json:"credit"`
	Balance     int64     `json:"balance"`
	Reference   string    `json:"reference"`
}

// TrialBalanceItem represents a row of the Trial Balance report.
type TrialBalanceItem struct {
	AccountCode    string `json:"account_code"`
	AccountName    string `json:"account_name"`
	InitialBalance int64  `json:"initial_balance"`
	TotalDebit     int64  `json:"total_debit"`
	TotalCredit    int64  `json:"total_credit"`
	FinalBalance   int64  `json:"final_balance"`
}
