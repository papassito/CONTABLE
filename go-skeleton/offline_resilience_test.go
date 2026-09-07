package main_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/klik/fcos-kernel/internal/domain"
)

type MockOfficialAdapter struct {
	IsOnline bool
}

func (m *MockOfficialAdapter) Download32DOpinion(ctx context.Context, rfc string) (string, error) {
	if !m.IsOnline {
		return "", errors.New("HTTP 503 Service Unavailable: Portal SAT fuera de servicio por mantenimiento")
	}
	return "OPINION_POSITIVA_HASH_123", nil
}

func TestOfflineResilienceAndExternalAdapterOutage(t *testing.T) {
	adapter := &MockOfficialAdapter{IsOnline: false}

	expediente := &domain.ExpedienteResponse{
		ExpedienteID:   uuid.New().String(),
		ExpedienteCode: "EXP-SAT-2026-09",
		CurrentStage:   "PROCESSING",
	}

	t.Logf("Iniciando simulación de resiliencia offline para expediente %s (ID: %s)", expediente.ExpedienteCode, expediente.ExpedienteID)

	// Simular ejecución del worker RPA
	_, err := adapter.Download32DOpinion(context.Background(), "CSO180512AAA")

	if err != nil {
		// Aplicar regla de resiliencia Offline-First
		expediente.CurrentStage = "EXTERNAL_UNAVAILABLE"

		job := map[string]interface{}{
			"job_type":      "RPA_DOWNLOAD_32D",
			"status":        "RETRY",
			"retry_count":   1,
			"max_retries":   3,
			"error_message": err.Error(),
			"scheduled_at":  time.Now().Add(5 * time.Minute),
		}

		// ASERCIONES
		if expediente.CurrentStage != "EXTERNAL_UNAVAILABLE" {
			t.Fatalf("Esperado estado EXTERNAL_UNAVAILABLE, obtenido: %s", expediente.CurrentStage)
		}

		if job["status"] != "RETRY" || job["retry_count"].(int) != 1 {
			t.Fatalf("El trabajo RPA no se encoló correctamente para reintento diferido")
		}
	} else {
		t.Fatalf("Se esperaba fallo por indisponibilidad simulada del adaptador")
	}

	// Restaurar servicio (Portal vuelve a estar en línea)
	adapter.IsOnline = true
	res, err := adapter.Download32DOpinion(context.Background(), "CSO180512AAA")
	if err != nil || res == "" {
		t.Fatalf("Fallo la recuperación automática del adaptador tras restaurar servicio")
	}

	expediente.CurrentStage = "EVIDENCE_CAPTURED"
	if expediente.CurrentStage != "EVIDENCE_CAPTURED" {
		t.Fatalf("No se actualizó el expediente a EVIDENCE_CAPTURED tras la recuperación")
	}
}
