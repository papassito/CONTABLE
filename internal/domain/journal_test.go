package domain

import "testing"

func TestJournalEntry_PartidaDobleExigida(t *testing.T) {
	// Caso 1: Asiento perfectamente equilibrado ($100.00 de debe y haber).
	entry := &JournalEntry{
		ID: "E1",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: 10000, Credit: 0},
			{AccountID: "ACC-02", Debit: 0, Credit: 10000},
		},
	}

	if err := entry.Validate(); err != nil {
		t.Errorf("Se esperaba validación exitosa de asiento balanceado, error: %v", err)
	}

	// Caso 2: Desbalance contable por 1 céntimo.
	unbalanced := &JournalEntry{
		ID: "E2",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: 10001, Credit: 0},
			{AccountID: "ACC-02", Debit: 0, Credit: 10000},
		},
	}

	if err := unbalanced.Validate(); err != ErrUnbalancedJournal {
		t.Errorf("Se esperaba error ErrUnbalancedJournal, obtenido: %v", err)
	}
}

func TestJournalEntry_ValidacionInvariantesLineas(t *testing.T) {
	// Caso 1: Menos de dos líneas
	badLinesCount := &JournalEntry{
		ID: "E3",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: 100, Credit: 0},
		},
	}
	if err := badLinesCount.Validate(); err != ErrInvalidLinesCount {
		t.Errorf("Se esperaba ErrInvalidLinesCount, obtenido: %v", err)
	}

	// Caso 2: Monto negativo inyectado
	negativeAmount := &JournalEntry{
		ID: "E4",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: -100, Credit: 0},
			{AccountID: "ACC-02", Debit: 0, Credit: -100},
		},
	}
	if err := negativeAmount.Validate(); err != ErrNegativeAmount {
		t.Errorf("Se esperaba ErrNegativeAmount, obtenido: %v", err)
	}

	// Caso 3: Doble afectación en una única línea (Debe y Haber activos)
	doubleSided := &JournalEntry{
		ID: "E5",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: 100, Credit: 100},
			{AccountID: "ACC-02", Debit: 0, Credit: 0},
		},
	}
	if err := doubleSided.Validate(); err != ErrDoubleSidedLine {
		t.Errorf("Se esperaba ErrDoubleSidedLine, obtenido: %v", err)
	}

	// Caso 4: Línea nula
	zeroMove := &JournalEntry{
		ID: "E6",
		Lines: []JournalLine{
			{AccountID: "ACC-01", Debit: 0, Credit: 0},
			{AccountID: "ACC-02", Debit: 0, Credit: 0},
		},
	}
	if err := zeroMove.Validate(); err != ErrZeroMove {
		t.Errorf("Se esperaba ErrZeroMove, obtenido: %v", err)
	}
}
