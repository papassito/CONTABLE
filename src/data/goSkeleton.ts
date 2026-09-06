export interface GoFile {
  path: string;
  name: string;
  layer: 'cmd' | 'domain' | 'repository' | 'service' | 'handler' | 'config' | 'pkg' | 'root';
  layerLabel: string;
  description: string;
  code: string;
}

export const GO_SKELETON_FILES: GoFile[] = [
  {
    path: 'cmd/api/main.go',
    name: 'main.go',
    layer: 'cmd',
    layerLabel: 'Punto de Entrada',
    description: 'Arranque del servidor e inyección de dependencias modular.',
    code: `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/klik/contable-fix/config"
	deliveryHttp "github.com/klik/contable-fix/internal/handler/http"
	"github.com/klik/contable-fix/internal/service"
)

func main() {
	fmt.Println("Iniciando Contable Fix by KLIK...")

	// 1. Cargar configuración base
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// 2. Inicializar repositorios (implementaciones de persistencia)
	// TODO: Inyectar repositorios reales (Postgres, MySQL o In-Memory)
	// accountRepo := postgres.NewAccountRepository(db)
	// journalRepo := postgres.NewJournalRepository(db)
	// ledgerRepo := postgres.NewLedgerRepository(db)

	// 3. Inicializar servicios de casos de uso
	accountSvc := service.NewAccountService(nil)
	journalSvc := service.NewJournalService(nil, nil, nil, nil)

	// 4. Inicializar handlers HTTP
	accountHandler := deliveryHttp.NewAccountHandler(accountSvc)
	journalHandler := deliveryHttp.NewJournalHandler(journalSvc)

	// 5. Configurar router
	router := deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		AccountHandler: accountHandler,
		JournalHandler: journalHandler,
	})

	// 6. Arrancar servidor HTTP
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Servidor Contable Fix escuchando en %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Fallo en el servidor HTTP: %v", err)
	}
}
`,
  },
  {
    path: 'config/config.go',
    name: 'config.go',
    layer: 'config',
    layerLabel: 'Configuración',
    description: 'Variables de entorno, base de datos y parámetros del servidor.',
    code: `package config

// Config contiene la configuración general de Contable Fix
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type DatabaseConfig struct {
	Driver string
	Host   string
	Port   int
	User   string
	Pass   string
	Name   string
	SSL    bool
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

// LoadConfig carga la configuración desde variables de entorno o archivo .env
// Estructura vacía lista para ser implementada según el entorno.
func LoadConfig() (*Config, error) {
	// TODO: Cargar variables de entorno / archivo config
	return &Config{
		App: AppConfig{
			Name:        "Contable Fix by KLIK",
			Environment: "development",
			Version:     "1.0.0",
		},
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
	}, nil
}
`,
  },
  {
    path: 'internal/domain/account.go',
    name: 'account.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Plan de cuentas contable, clasificaciones y estados.',
    code: `package domain

import (
	"time"
)

// AccountType clasifica el tipo de cuenta según el Plan Contable
type AccountType string

const (
	AccountTypeActivo     AccountType = "ACTIVO"
	AccountTypePasivo     AccountType = "PASIVO"
	AccountTypePatrimonio AccountType = "PATRIMONIO"
	AccountTypeIngreso    AccountType = "INGRESO"
	AccountTypeGasto      AccountType = "GASTO"
	AccountTypeCosto      AccountType = "COSTO"
)

// AccountStatus define el estado operativo de la cuenta
type AccountStatus string

const (
	AccountStatusActiva    AccountStatus = "ACTIVA"
	AccountStatusInactiva  AccountStatus = "INACTIVA"
	AccountStatusBloqueada AccountStatus = "BLOQUEADA"
)

// Account representa una cuenta del plan contable
type Account struct {
	ID          string        \`json:"id"\`
	Code        string        \`json:"code"\`         // Ej: "110505" (Caja general)
	Name        string        \`json:"name"\`         // Ej: "Caja Principal"
	Type        AccountType   \`json:"type"\`         // Activo, Pasivo, etc.
	ParentID    *string       \`json:"parent_id"\`    // Cuenta padre para estructura en árbol
	Level       int           \`json:"level"\`        // Nivel jerárquico (1 a 5)
	AcceptsMove bool          \`json:"accepts_move"\` // ¿Acepta movimientos directos?
	Status      AccountStatus \`json:"status"\`
	CurrentBal  int64         \`json:"current_balance"\`
	CreatedAt   time.Time     \`json:"created_at"\`
	UpdatedAt   time.Time     \`json:"updated_at"\`
}
`,
  },
  {
    path: 'internal/domain/journal.go',
    name: 'journal.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Comprobantes y asientos contables (encabezado y partidas dobles).',
    code: `package domain

import (
	"time"
)

// EntryStatus define el estado del asiento contable
type EntryStatus string

const (
	EntryStatusBorrador      EntryStatus = "BORRADOR"
	EntryStatusContabilizado EntryStatus = "CONTABILIZADO"
	EntryStatusAnulado       EntryStatus = "ANULADO"
)

// JournalEntry representa el encabezado del asiento/comprobante contable
type JournalEntry struct {
	ID          string        \`json:"id"\`
	Number      string        \`json:"number"\`        // Consecutivo contable único
	Date        time.Time     \`json:"date"\`
	Concept     string        \`json:"concept"\`       // Glosa o concepto
	Reference   string        \`json:"reference"\`     // N° Factura o documento de origen
	Status      EntryStatus   \`json:"status"\`
	Lines       []JournalLine \`json:"lines"\`
	TotalDebit  int64         \`json:"total_debit"\`
	TotalCredit int64         \`json:"total_credit"\`
	CreatedBy   string        \`json:"created_by"\`
	CreatedAt   time.Time     \`json:"created_at"\`
	UpdatedAt   time.Time     \`json:"updated_at"\`
}

// JournalLine representa una línea o partida del asiento contable
type JournalLine struct {
	ID             string    \`json:"id"\`
	JournalEntryID string    \`json:"journal_entry_id"\`
	AccountID      string    \`json:"account_id"\`
	AccountCode    string    \`json:"account_code"\`
	Description    string    \`json:"description"\`
	Debit          int64     \`json:"debit"\`
	Credit         int64     \`json:"credit"\`
	ThirdPartyID   *string   \`json:"third_party_id"\` // Identificación del tercero (NIT / CC)
	CreatedAt      time.Time \`json:"created_at"\`
}
`,
  },
  {
    path: 'internal/domain/ledger.go',
    name: 'ledger.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Libro Mayor y Balance de Comprobación (Sumas y Saldos).',
    code: `package domain

import "time"

// LedgerEntry representa un movimiento en el Libro Mayor
type LedgerEntry struct {
	ID          string    \`json:"id"\`
	AccountID   string    \`json:"account_id"\`
	AccountCode string    \`json:"account_code"\`
	EntryDate   time.Time \`json:"entry_date"\`
	Debit       int64     \`json:"debit"\`
	Credit      int64     \`json:"credit"\`
	Balance     int64     \`json:"balance"\`
	Reference   string    \`json:"reference"\`
}

// TrialBalanceItem representa una fila del Balance de Comprobación
type TrialBalanceItem struct {
	AccountCode    string  \`json:"account_code"\`
	AccountName    string  \`json:"account_name"\`
	InitialBalance int64   \`json:"initial_balance"\`
	TotalDebit     int64   \`json:"total_debit"\`
	TotalCredit    int64   \`json:"total_credit"\`
	FinalBalance   int64   \`json:"final_balance"\`
}
`,
  },
  {
    path: 'internal/domain/invoice.go',
    name: 'invoice.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Facturas de venta, compra y su enlace con asientos contables.',
    code: `package domain

import "time"

// InvoiceType define el tipo de documento soporte
type InvoiceType string

const (
	InvoiceTypeVenta  InvoiceType = "VENTA"
	InvoiceTypeCompra InvoiceType = "COMPRA"
	InvoiceTypeNota   InvoiceType = "NOTA_CREDITO"
)

// InvoiceStatus define el estado de la factura
type InvoiceStatus string

const (
	InvoiceStatusBorrador InvoiceStatus = "BORRADOR"
	InvoiceStatusEmitida  InvoiceStatus = "EMITIDA"
	InvoiceStatusPagada   InvoiceStatus = "PAGADA"
	InvoiceStatusAnulada  InvoiceStatus = "ANULADA"
)

// Invoice representa una factura o documento soporte
type Invoice struct {
	ID             string        \`json:"id"\`
	Number         string        \`json:"number"\`
	Type           InvoiceType   \`json:"type"\`
	Status         InvoiceStatus \`json:"status"\`
	ThirdPartyID   string        \`json:"third_party_id"\`
	IssueDate      time.Time     \`json:"issue_date"\`
	DueDate        time.Time     \`json:"due_date"\`
	Subtotal       int64         \`json:"subtotal"\`
	TaxAmount      int64         \`json:"tax_amount"\`
	Total          int64         \`json:"total"\`
	Items          []InvoiceItem \`json:"items"\`
	JournalEntryID *string       \`json:"journal_entry_id"\` // Asiento contable generado automáticamente
	CreatedAt      time.Time     \`json:"created_at"\`
	UpdatedAt      time.Time     \`json:"updated_at"\`
}

// InvoiceItem representa el detalle de cada línea de la factura
type InvoiceItem struct {
	ID          string  \`json:"id"\`
	InvoiceID   string  \`json:"invoice_id"\`
	Description string  \`json:"description"\`
	AccountID   string  \`json:"account_id"\` // Cuenta contable vinculada
	Quantity    float64 \`json:"quantity"\`
	UnitPrice   int64   \`json:"unit_price"\`
	Discount    int64   \`json:"discount"\`
	TaxRate     float64 \`json:"tax_rate"\`
	Total       int64   \`json:"total"\`
}
`,
  },
  {
    path: 'internal/domain/audit.go',
    name: 'audit.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Trazabilidad y pista de auditoría para operaciones contables.',
    code: `package domain

import "time"

// AuditAction define la acción realizada
type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionPost   AuditAction = "POST_CONTABILIZAR"
	AuditActionVoid   AuditAction = "VOID_ANULAR"
)

// AuditLog registra la trazabilidad y auditoría de cambios
type AuditLog struct {
	ID         string      \`json:"id"\`
	EntityName string      \`json:"entity_name"\` // Ej: "JournalEntry", "Account"
	EntityID   string      \`json:"entity_id"\`
	Action     AuditAction \`json:"action"\`
	UserID     string      \`json:"user_id"\`
	IPAddress  string      \`json:"ip_address"\`
	Details    string      \`json:"details"\`
	Timestamp  time.Time   \`json:"timestamp"\`
}
`,
  },
  {
    path: 'internal/domain/errors.go',
    name: 'errors.go',
    layer: 'domain',
    layerLabel: 'Dominio / Entidades',
    description: 'Definición de errores de negocio contables estándar.',
    code: `package domain

import "errors"

// Errores comunes de la capa de dominio contable
var (
	ErrAccountNotFound    = errors.New("cuenta contable no encontrada")
	ErrAccountInactive    = errors.New("la cuenta contable está inactiva o bloqueada")
	ErrAccountHasChildren = errors.New("no se puede asignar movimientos a una cuenta mayorizadora con subcuentas")
	ErrUnbalancedJournal  = errors.New("el asiento no cumple con la partida doble (Débito != Crédito)")
	ErrEmptyJournalLines  = errors.New("el asiento debe contener al menos dos líneas contables")
	ErrEntryAlreadyPosted = errors.New("el asiento contable ya ha sido contabilizado y no puede modificarse")
	ErrInvoiceNotFound    = errors.New("factura no encontrada")
	ErrInvalidPeriod      = errors.New("el periodo contable se encuentra cerrado")
)
`,
  },
  {
    path: 'internal/repository/account_repository.go',
    name: 'account_repository.go',
    layer: 'repository',
    layerLabel: 'Puertos de Persistencia',
    description: 'Contrato de persistencia (interface) para el plan de cuentas.',
    code: `package repository

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
)

// AccountRepository define el contrato para la persistencia del plan de cuentas
type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByCode(ctx context.Context, code string) (*domain.Account, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id string) error
	UpdateBalance(ctx context.Context, id string, amount int64) error
}
`,
  },
  {
    path: 'internal/repository/journal_repository.go',
    name: 'journal_repository.go',
    layer: 'repository',
    layerLabel: 'Puertos de Persistencia',
    description: 'Contrato de persistencia (interface) para asientos contables.',
    code: `package repository

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
)

// UnitOfWork define el contrato para coordinar transacciones ACID
type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// JournalRepository define el contrato para la persistencia de asientos contables
type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) error
	GetByID(ctx context.Context, id string) (*domain.JournalEntry, error)
	GetByNumber(ctx context.Context, number string) (*domain.JournalEntry, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
	UpdateStatus(ctx context.Context, id string, status domain.EntryStatus) error
	GetLinesByEntryID(ctx context.Context, entryID string) ([]domain.JournalLine, error)
}
`,
  },
  {
    path: 'internal/repository/ledger_repository.go',
    name: 'ledger_repository.go',
    layer: 'repository',
    layerLabel: 'Puertos de Persistencia',
    description: 'Contrato de persistencia para el libro mayor y balances.',
    code: `package repository

import (
	"context"
	"time"
	"github.com/klik/contable-fix/internal/domain"
)

// LedgerRepository define el contrato para consultas de mayor y balances
type LedgerRepository interface {
	RecordMovements(ctx context.Context, entries []domain.LedgerEntry) error
	GetMovementsByAccount(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error)
	GetTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error)
}
`,
  },
  {
    path: 'internal/repository/invoice_repository.go',
    name: 'invoice_repository.go',
    layer: 'repository',
    layerLabel: 'Puertos de Persistencia',
    description: 'Contrato de persistencia para facturación y documentos de soporte.',
    code: `package repository

import (
	"context"
	"github.com/klik/contable-fix/internal/domain"
)

// InvoiceRepository define el contrato para la persistencia de facturas
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*domain.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status domain.InvoiceStatus) error
	LinkJournalEntry(ctx context.Context, invoiceID string, journalEntryID string) error
}
`,
  },
  {
    path: 'internal/service/account_service.go',
    name: 'account_service.go',
    layer: 'service',
    layerLabel: 'Casos de Uso / Servicios',
    description: 'Lógica y validaciones del plan de cuentas contable.',
    code: `package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/repository"
)

// AccountService define los casos de uso para cuentas contables
type AccountService interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	GetAccountByCode(ctx context.Context, code string) (*domain.Account, error)
	ListAccounts(ctx context.Context) ([]*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	DisableAccount(ctx context.Context, id string) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

// NewAccountService crea una nueva instancia del servicio de cuentas
func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: repo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, account *domain.Account) error {
	existing, err := s.accountRepo.GetByCode(ctx, account.Code)
	if err == nil && existing != nil {
		return fmt.Errorf("la cuenta contable con código %s ya existe", account.Code)
	}

	if len(account.Code) == 0 {
		return errors.New("el código contable no puede estar vacío")
	}

	firstChar := account.Code[0]
	var expectedType domain.AccountType
	switch firstChar {
	case '1':
		expectedType = domain.AccountTypeActivo
	case '2':
		expectedType = domain.AccountTypePasivo
	case '3':
		expectedType = domain.AccountTypePatrimonio
	case '4':
		expectedType = domain.AccountTypeIngreso
	case '5':
		expectedType = domain.AccountTypeGasto
	case '6':
		expectedType = domain.AccountTypeCosto
	default:
		return errors.New("código contable inválido: debe iniciar con un dígito de 1 a 6")
	}

	if account.Type != expectedType {
		return fmt.Errorf("tipo de cuenta %s incongruente con el código contable que inicia con %c", account.Type, firstChar)
	}

	if account.ParentID != nil && *account.ParentID != "" {
		parent, err := s.accountRepo.GetByID(ctx, *account.ParentID)
		if err != nil || parent == nil {
			return errors.New("cuenta padre no encontrada")
		}
		if parent.AcceptsMove {
			if parent.CurrentBal != 0 {
				return fmt.Errorf("no se puede convertir la cuenta padre %s en mayorizadora porque posee un saldo de %d centavos", parent.Code, parent.CurrentBal)
			}
			parent.AcceptsMove = false
			_ = s.accountRepo.Update(ctx, parent)
		}
	}

	return s.accountRepo.Create(ctx, account)
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

func (s *accountService) GetAccountByCode(ctx context.Context, code string) (*domain.Account, error) {
	return s.accountRepo.GetByCode(ctx, code)
}

func (s *accountService) ListAccounts(ctx context.Context) ([]*domain.Account, error) {
	return s.accountRepo.List(ctx, nil)
}

func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return s.accountRepo.Update(ctx, account)
}

func (s *accountService) DisableAccount(ctx context.Context, id string) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if account == nil {
		return domain.ErrAccountNotFound
	}

	if account.CurrentBal != 0 {
		return errors.New("no se puede desactivar una cuenta con saldo diferente de cero")
	}

	accounts, err := s.accountRepo.List(ctx, map[string]interface{}{"parent_id": id})
	if err == nil {
		for _, sub := range accounts {
			if sub.Status == domain.AccountStatusActiva {
				return errors.New("no se puede desactivar la cuenta porque tiene subcuentas activas")
			}
		}
	}

	account.Status = domain.AccountStatusInactiva
	return s.accountRepo.Update(ctx, account)
}
`,
  },
  {
    path: 'internal/service/journal_service.go',
    name: 'journal_service.go',
    layer: 'service',
    layerLabel: 'Casos de Uso / Servicios',
    description: 'Lógica contable de partida doble, contabilización y reversiones.',
    code: `package service

import (
	"context"
	"time"

	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/repository"
	"github.com/klik/contable-fix/pkg/validator"
)

// JournalService define los casos de uso para asientos contables
type JournalService interface {
	CreateDraft(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error)
	PostEntry(ctx context.Context, entryID string) error
	ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error)
	GetEntry(ctx context.Context, id string) (*domain.JournalEntry, error)
	ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error)
}

// UnitOfWork define el puerto para coordinar transacciones ACID.
type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type journalService struct {
	journalRepo repository.JournalRepository
	accountRepo repository.AccountRepository
	ledgerRepo  repository.LedgerRepository
	uow         UnitOfWork
}

// NewJournalService crea una nueva instancia del servicio de asientos contables
func NewJournalService(
	journalRepo repository.JournalRepository,
	accountRepo repository.AccountRepository,
	ledgerRepo repository.LedgerRepository,
	uow UnitOfWork,
) JournalService {
	return &journalService{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
		ledgerRepo:  ledgerRepo,
		uow:         uow,
	}
}

func (s *journalService) CreateDraft(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	if len(entry.Lines) < 2 {
		return nil, domain.ErrEmptyJournalLines
	}

	if !validator.ValidateDoubleEntry(entry.Lines) {
		return nil, domain.ErrUnbalancedJournal
	}

	// Validar que cada cuenta exista, esté activa y acepte movimientos
	for _, line := range entry.Lines {
		account, err := s.accountRepo.GetByID(ctx, line.AccountID)
		if err != nil || account == nil {
			return nil, domain.ErrAccountNotFound
		}
		if account.Status != domain.AccountStatusActiva {
			return nil, domain.ErrAccountInactive
		}
		if !account.AcceptsMove {
			return nil, domain.ErrAccountHasChildren
		}
	}

	entry.Status = domain.EntryStatusBorrador

	var totalDebit, totalCredit int64
	for _, line := range entry.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	entry.TotalDebit = totalDebit
	entry.TotalCredit = totalCredit
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	err := s.journalRepo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *journalService) PostEntry(ctx context.Context, entryID string) error {
	entry, err := s.journalRepo.GetByID(ctx, entryID)
	if err != nil {
		return err
	}
	if entry == nil {
		return domain.ErrAccountNotFound
	}

	if entry.Status != domain.EntryStatusBorrador {
		return domain.ErrEntryAlreadyPosted
	}

	lines, err := s.journalRepo.GetLinesByEntryID(ctx, entryID)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		lines = entry.Lines
	}
	if len(lines) < 2 {
		return domain.ErrEmptyJournalLines
	}

	if !validator.ValidateDoubleEntry(lines) {
		return domain.ErrUnbalancedJournal
	}

	return s.uow.WithTransaction(ctx, func(txCtx context.Context) error {
		ledgerEntries := make([]domain.LedgerEntry, len(lines))
		for i, line := range lines {
			ledgerEntries[i] = domain.LedgerEntry{
				AccountID:   line.AccountID,
				AccountCode: line.AccountCode,
				EntryDate:   entry.Date,
				Debit:       line.Debit,
				Credit:      line.Credit,
				Reference:   entry.Number,
			}
		}

		err = s.ledgerRepo.RecordMovements(txCtx, ledgerEntries)
		if err != nil {
			return err
		}

		for _, line := range lines {
			account, err := s.accountRepo.GetByID(txCtx, line.AccountID)
			if err != nil {
				return err
			}
			var amount int64
			switch account.Type {
			case domain.AccountTypeActivo, domain.AccountTypeGasto, domain.AccountTypeCosto:
				amount = line.Debit - line.Credit
			case domain.AccountTypePasivo, domain.AccountTypePatrimonio, domain.AccountTypeIngreso:
				amount = line.Credit - line.Debit
			default:
				amount = line.Debit - line.Credit
			}
			err = s.accountRepo.UpdateBalance(txCtx, line.AccountID, amount)
			if err != nil {
				return err
			}
		}

		return s.journalRepo.UpdateStatus(txCtx, entryID, domain.EntryStatusContabilizado)
	})
}

func (s *journalService) ReverseEntry(ctx context.Context, entryID string, reason string) (*domain.JournalEntry, error) {
	original, err := s.journalRepo.GetByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, domain.ErrAccountNotFound
	}

	originalLines, err := s.journalRepo.GetLinesByEntryID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	if len(originalLines) == 0 {
		originalLines = original.Lines
	}

	reversedLines := make([]domain.JournalLine, len(originalLines))
	for i, line := range originalLines {
		reversedLines[i] = domain.JournalLine{
			AccountID:    line.AccountID,
			AccountCode:  line.AccountCode,
			Description:  "Reversión: " + line.Description,
			Debit:        line.Credit,
			Credit:       line.Debit,
			ThirdPartyID: line.ThirdPartyID,
		}
	}

	reversedEntry := &domain.JournalEntry{
		Number:      "REV-" + original.Number,
		Date:        time.Now(),
		Concept:     "Reversión de " + original.Number + " - Motivo: " + reason,
		Reference:   original.Number,
		Status:      domain.EntryStatusBorrador,
		Lines:       reversedLines,
		TotalDebit:  original.TotalCredit,
		TotalCredit: original.TotalDebit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.journalRepo.Create(ctx, reversedEntry)
	if err != nil {
		return nil, err
	}

	return reversedEntry, nil
}

func (s *journalService) GetEntry(ctx context.Context, id string) (*domain.JournalEntry, error) {
	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	lines, err := s.journalRepo.GetLinesByEntryID(ctx, id)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines
	return entry, nil
}

func (s *journalService) ListEntries(ctx context.Context, filter map[string]interface{}) ([]*domain.JournalEntry, error) {
	return s.journalRepo.List(ctx, filter)
}
`,
  },
  {
    path: 'internal/service/ledger_service.go',
    name: 'ledger_service.go',
    layer: 'service',
    layerLabel: 'Casos de Uso / Servicios',
    description: 'Generación de libro mayor y balances de prueba.',
    code: `package service

import (
	"context"
	"time"
	"github.com/klik/contable-fix/internal/domain"
	"github.com/klik/contable-fix/internal/repository"
)

// LedgerService define los casos de uso para reportes y libro mayor
type LedgerService interface {
	GetAccountLedger(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error)
	GenerateTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error)
}

type ledgerService struct {
	ledgerRepo  repository.LedgerRepository
	accountRepo repository.AccountRepository
}

func NewLedgerService(ledgerRepo repository.LedgerRepository, accountRepo repository.AccountRepository) LedgerService {
	return &ledgerService{
		ledgerRepo:  ledgerRepo,
		accountRepo: accountRepo,
	}
}

func (s *ledgerService) GetAccountLedger(ctx context.Context, accountID string, from, to time.Time) ([]domain.LedgerEntry, error) {
	// TODO: Calcular movimientos y saldos acumulados
	return nil, nil
}

func (s *ledgerService) GenerateTrialBalance(ctx context.Context, from, to time.Time) ([]domain.TrialBalanceItem, error) {
	// TODO: Generar balance de comprobación sumas y saldos
	return nil, nil
}
`,
  },
  {
    path: 'internal/handler/http/account_handler.go',
    name: 'account_handler.go',
    layer: 'handler',
    layerLabel: 'Controladores HTTP',
    description: 'Endpoints REST para el plan de cuentas.',
    code: `package http

import (
	"net/http"
	"github.com/klik/contable-fix/internal/service"
)

// AccountHandler gestiona las peticiones HTTP del plan contable
type AccountHandler struct {
	service service.AccountService
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Decodificar payload JSON, validar y llamar a h.service.CreateAccount
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Extraer ID de la URL y responder JSON
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: Listar cuentas y responder JSON
}
`,
  },
  {
    path: 'internal/handler/http/journal_handler.go',
    name: 'journal_handler.go',
    layer: 'handler',
    layerLabel: 'Controladores HTTP',
    description: 'Endpoints REST para crear borradores, contabilizar y revertir.',
    code: `package http

import (
	"net/http"
	"github.com/klik/contable-fix/internal/service"
)

// JournalHandler gestiona las peticiones HTTP para asientos contables
type JournalHandler struct {
	service service.JournalService
}

func NewJournalHandler(service service.JournalService) *JournalHandler {
	return &JournalHandler{service: service}
}

func (h *JournalHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	// TODO: Recibir asiento contable y llamar a CreateDraft
}

func (h *JournalHandler) Post(w http.ResponseWriter, r *http.Request) {
	// TODO: Contabilizar asiento y asentar en libro mayor
}

func (h *JournalHandler) Reverse(w http.ResponseWriter, r *http.Request) {
	// TODO: Revertir asiento contable
}

func (h *JournalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Consultar asiento específico
}
`,
  },
  {
    path: 'internal/handler/http/router.go',
    name: 'router.go',
    layer: 'handler',
    layerLabel: 'Controladores HTTP',
    description: 'Enrutador HTTP con Go 1.22+ stdlib mux.',
    code: `package http

import (
	"net/http"
)

// RouterConfig contiene dependencias para construir las rutas HTTP
type RouterConfig struct {
	AccountHandler *AccountHandler
	JournalHandler *JournalHandler
}

// NewRouter crea el enrutador HTTP estándar de Go (o adaptable a Chi / Gin / Fiber)
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Plan de cuentas contable
	mux.HandleFunc("POST /api/v1/accounts", cfg.AccountHandler.Create)
	mux.HandleFunc("GET /api/v1/accounts", cfg.AccountHandler.List)
	mux.HandleFunc("GET /api/v1/accounts/{id}", cfg.AccountHandler.GetByID)

	// Asientos y comprobantes contables
	mux.HandleFunc("POST /api/v1/journal-entries", cfg.JournalHandler.CreateDraft)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/post", cfg.JournalHandler.Post)
	mux.HandleFunc("POST /api/v1/journal-entries/{id}/reverse", cfg.JournalHandler.Reverse)
	mux.HandleFunc("GET /api/v1/journal-entries/{id}", cfg.JournalHandler.GetByID)

	return mux
}
`,
  },
  {
    path: 'pkg/response/response.go',
    name: 'response.go',
    layer: 'pkg',
    layerLabel: 'Paquete de Utilidades',
    description: 'Envoltura de respuestas JSON estructuradas y errores.',
    code: `package response

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

// JSON escribe una respuesta estructurada con código de estado HTTP
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{"data": data})
}

// Error escribe un mensaje de error estandarizado
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{"error": message})
}
`,
  },
  {
    path: 'pkg/validator/validator.go',
    name: 'validator.go',
    layer: 'pkg',
    layerLabel: 'Paquete de Utilidades',
    description: 'Validador contable de regla de partida doble (Débitos == Créditos).',
    code: `package validator

import (
	"github.com/klik/contable-fix/internal/domain"
)

// ValidateDoubleEntry valida matemáticamente la partida doble (Débitos == Créditos) usando enteros exactos
func ValidateDoubleEntry(lines []domain.JournalLine) bool {
	var totalDebit, totalCredit int64
	for _, line := range lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}
	return totalDebit == totalCredit
}
`,
  },
  {
    path: 'go.mod',
    name: 'go.mod',
    layer: 'root',
    layerLabel: 'Módulo Go',
    description: 'Definición del módulo Go para Contable Fix by KLIK.',
    code: `module github.com/klik/contable-fix

go 1.22
`,
  },
  {
    path: 'Makefile',
    name: 'Makefile',
    layer: 'root',
    layerLabel: 'Automatización',
    description: 'Comandos de compilación, ejecución y tests.',
    code: `.PHONY: run test build clean lint

run:
	go run cmd/api/main.go

build:
	go build -o bin/contable-fix cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
`,
  },
  {
    path: 'README.md',
    name: 'README.md',
    layer: 'root',
    layerLabel: 'Documentación',
    description: 'Guía de arquitectura y puesta en marcha.',
    code: `# Contable Fix by KLIK

Estructura base limpia y modular en Go (Golang) para el sistema contable **Contable Fix by KLIK**.

## Estructura Modular (Clean Architecture)

\`\`\`
contable-fix/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada e inyección de dependencias
├── config/
│   └── config.go                  # Configuración y variables de entorno
├── internal/
│   ├── domain/                    # Entidades puras y reglas de negocio
│   │   ├── account.go             # Plan de cuentas contables
│   │   ├── journal.go             # Asientos contables y partidas
│   │   ├── ledger.go              # Libro mayor y balances
│   │   ├── invoice.go             # Facturas y documentos soporte
│   │   ├── audit.go               # Pistas de auditoría y trazabilidad
│   │   └── errors.go              # Errores del dominio contable
│   ├── repository/                # Puertos (interfaces) de persistencia
│   │   ├── account_repository.go
│   │   ├── journal_repository.go
│   │   ├── ledger_repository.go
│   │   └── invoice_repository.go
│   ├── service/                   # Casos de uso / Lógica de aplicación
│   │   ├── account_service.go
│   │   ├── journal_service.go
│   │   └── ledger_service.go
│   └── handler/
│       └── http/                  # Controladores REST y enrutamiento
│           ├── account_handler.go
│           ├── journal_handler.go
│           └── router.go
├── pkg/
│   ├── response/                  # Respuestas JSON estandarizadas
│   │   └── response.go
│   └── validator/                 # Validación de partida doble
│       └── validator.go
├── Makefile                       # Comandos make run / build / test
└── go.mod                         # Definición del módulo Go
\`\`\`

## Cómo utilizar este esqueleto

1. Clona o descarga los archivos.
2. Cada función contiene un comentario \`// TODO:\` para implementar según tu motor de base de datos preferido (PostgreSQL, SQLite, MySQL o memoria).
3. Ejecuta \`go run cmd/api/main.go\` para iniciar el servidor.
`,
  },
  {
