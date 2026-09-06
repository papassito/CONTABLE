package domain

import "time"

// LedgerEntry representa un movimiento individual registrado en el Libro Mayor.
type LedgerEntry struct {
    ID          string    `json:"id"`
    JournalID   string    `json:"journal_id"`
    AccountID   string    `json:"account_id"`
    AccountCode string    `json:"account_code"`
    EntryDate   time.Time `json:"entry_date"`
    Debit       int64     `json:"debit"`
    Credit      int64     `json:"credit"`
    AmountCents int64     `json:"amount_cents"`
    IsDebit     bool      `json:"is_debit"`
    Reference   string    `json:"reference"`
    PostedAtUTC time.Time `json:"posted_at_utc"`
}

// TrialBalanceItem representa la estructura de balance de comprobación.
type TrialBalanceItem struct {
    AccountID    string `json:"account_id"`
    AccountName  string `json:"account_name"`
    InitialCents int64  `json:"initial_cents"`
    DebitCents   int64  `json:"debit_cents"`
    CreditCents  int64  `json:"credit_cents"`
    FinalCents   int64  `json:"final_cents"`
}
