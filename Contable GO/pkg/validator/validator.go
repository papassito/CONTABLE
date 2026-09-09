package validator

import (
	"github.com/klik/fcos-kernel/internal/domain"
)

func ValidateDoubleEntry(lines []domain.JournalLine) bool {
	var totalDebit, totalCredit int64
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	return totalDebit == totalCredit
}
