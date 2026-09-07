package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gowebpki/jcs"
)

// loggingAuditService es una implementación de AuditService que registra eventos en el log estándar.
type loggingAuditService struct{}

// NewLoggingAuditService crea una nueva instancia del servicio de auditoría basado en logs.
func NewLoggingAuditService() AuditService {
	return &loggingAuditService{}
}

// AppendEvent registra un evento de auditoría en la salida estándar, garantizando formato canónico.
func (s *loggingAuditService) AppendEvent(ctx context.Context, tenantID, eventType, actorID string, payload any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	canonicalPayload, _ := jcs.Transform(payloadJSON)
	log.Printf("[AUDIT] Tenant: %s, Actor: %s, Event: %s, Payload: %s", tenantID, actorID, eventType, string(canonicalPayload))
	return nil
}
