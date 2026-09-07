package domain

import "errors"

var (
	ErrAccountNotFound    = errors.New("cuenta contable no encontrada")
	ErrJournalNotFound    = errors.New("asiento contable no encontrado")
	ErrUnbalancedJournal  = errors.New("el asiento no cumple con la partida doble (Débito != Crédito)")
	ErrEmptyJournalLines  = errors.New("el asiento debe contener al menos dos líneas contables")
	ErrAccountInactive    = errors.New("la cuenta contable está inactiva o bloqueada")
	ErrAccountHasChildren = errors.New("no se pueden imputar movimientos a una cuenta mayorizadora")
	ErrEntryAlreadyPosted = errors.New("el asiento contable ya ha sido contabilizado y no puede ser modificado")
	ErrClosedPeriod       = errors.New("el periodo contable se encuentra cerrado o declarado")
)
