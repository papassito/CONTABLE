package domain

import "errors"

var (
    ErrAccountNotFound    = errors.New("cuenta contable no encontrada")
    ErrJournalNotFound    = errors.New("asiento contable no encontrado")
    ErrUnbalancedEntry    = errors.New("el asiento contable no está balanceado")
    ErrUnbalancedJournal  = errors.New("el asiento contable está desbalanceado")
    ErrEmptyJournalLines  = errors.New("el asiento contable debe contener al menos una línea")
    ErrAccountInactive    = errors.New("la cuenta contable se encuentra inactiva")
    ErrAccountHasChildren = errors.New("no se pueden imputar movimientos a una cuenta con subcuentas")
    ErrEntryAlreadyPosted = errors.New("el asiento contable ya ha sido contabilizado")
    ErrClosedPeriod       = errors.New("el periodo contable correspondiente está cerrado")
)
