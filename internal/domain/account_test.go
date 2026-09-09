package domain

import "testing"

func TestAccount_ValidacionEstado(t *testing.T) {
	acc := &Account{
		ID:      "1",
		Status:  AccountActive,
		Balance: 1500, // Saldo activo
	}

	err := acc.ValidateTransition(AccountInactive)
	if err == nil {
		t.Error("Se esperaba error al desactivar una cuenta con saldo activo")
	}

	acc.Balance = 0
	err = acc.ValidateTransition(AccountInactive)
	if err != nil {
		t.Errorf("No se esperaba error desactivando cuenta sin saldo: %v", err)
	}
}

func TestAccount_DeteccionCiclosJerarquicos(t *testing.T) {
	// Generamos un catálogo de pruebas aislado.
	catalog := map[string]*Account{
		"A": {ID: "A", ParentID: "B"},
		"B": {ID: "B", ParentID: "C"},
		"C": {ID: "C", ParentID: ""}, // Estructura sana
	}

	if DetectCycle(catalog, "A") {
		t.Error("La jerarquía es sana, no se debía reportar ciclo")
	}

	// Introducción intencional de un bucle cíclico: C -> A
	catalog["C"].ParentID = "A"

	if !DetectCycle(catalog, "A") {
		t.Error("Se esperaba la detección de ciclo infinito (A -> B -> C -> A)")
	}

	if !DetectCycle(catalog, "B") {
		t.Error("Se esperaba la detección de ciclo comenzando desde el nodo medio B")
	}

	// Test con raíz vacía.
	emptyCycle := &Account{ID: "X", ParentID: ""}
	if DetectCycle(map[string]*Account{"X": emptyCycle}, "X") {
		t.Error("Cuenta sin padre no puede tener un ciclo")
	}
}
