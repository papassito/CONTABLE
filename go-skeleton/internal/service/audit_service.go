package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
	"github.com/klik/fcos-kernel/pkg/database"
)

// dbAuditService es una implementación persistente y criptográfica de AuditService.
type dbAuditService struct {
	db *sql.DB
	mu sync.Mutex
}

// NewDBAuditService crea una instancia de auditoría persistente respaldada por la base de datos.
func NewDBAuditService(db *sql.DB) AuditService {
	return &dbAuditService{db: db}
}

func (s *dbAuditService) AppendEvent(ctx context.Context, tenantID, eventType, actorID string, payload any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	canonicalPayload, err := jcs.Transform(payloadJSON)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(canonicalPayload)
	currentPayloadHash := hex.EncodeToString(hash[:])

	var queryRow func(query string, args ...any) *sql.Row
	var exec func(query string, args ...any) (sql.Result, error)

	if tx, ok := database.GetTx(ctx); ok {
		queryRow = func(query string, args ...any) *sql.Row {
			return tx.QueryRowContext(ctx, query, args...)
		}
		exec = func(query string, args ...any) (sql.Result, error) {
			return tx.ExecContext(ctx, query, args...)
		}
	} else {
		queryRow = func(query string, args ...any) *sql.Row {
			return s.db.QueryRowContext(ctx, query, args...)
		}
		exec = func(query string, args ...any) (sql.Result, error) {
			return s.db.ExecContext(ctx, query, args...)
		}
	}

	var previousHash string
	// Consulta determinista por fecha y rowid
	query := "SELECT current_hash FROM audit_events WHERE tenant_id = ? ORDER BY created_at_utc DESC, rowid DESC LIMIT 1"
	err = queryRow(query, tenantID).Scan(&previousHash)
	if err == sql.ErrNoRows {
		previousHash = "0000000000000000000000000000000000000000000000000000000000000000"
	} else if err != nil {
		return err
	}

	chainedInput := currentPayloadHash + previousHash
	chainedHashBytes := sha256.Sum256([]byte(chainedInput))
	currentHash := hex.EncodeToString(chainedHashBytes[:])

	id := uuid.New().String()
	createdAt := time.Now().UTC()

	_, err = exec(`
		INSERT INTO audit_events (id, tenant_id, event_type, payload_json, previous_hash, current_hash, created_at_utc)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, tenantID, eventType, string(canonicalPayload), previousHash, currentHash, createdAt)
	return err
}
