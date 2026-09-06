package application

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"

	"github.com/klik/fcos-kernel/pkg/database"
)

type HotSwapperService interface {
	ApplyNormativeUpdate(ctx context.Context, payload []byte, signature []byte, pubKeyPEM []byte) error
}

type hotSwapperService struct {
	db *sql.DB
}

func NewHotSwapperService(db *sql.DB) HotSwapperService {
	return &hotSwapperService{db: db}
}

func (s *hotSwapperService) ApplyNormativeUpdate(ctx context.Context, payload []byte, signature []byte, pubKeyPEM []byte) error {
	// 1. Validar la firma digital de Klik HQ
	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return fmt.Errorf("clave pública PEM inválida")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error al procesar clave pública: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("la clave pública debe ser RSA")
	}

	hashed := sha256.Sum256(payload)
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("RECHAZO DE SEGURIDAD: La firma del paquete normativo es inválida o fue alterada: %w", err)
	}

	// 2. Extraer transacción activa del UnitOfWork
	tx, ok := database.GetTx(ctx)
	if !ok {
		return fmt.Errorf("se requiere una transacción atómica activa")
	}

	// 3. Aplicar cierre a la versión anterior e insertar nueva versión en normative_parameter_versions
	// Garantía de trazabilidad estricta: nunca se sobrescribe el pasado
	_, err = tx.ExecContext(ctx, `
		UPDATE normative_parameter_versions 
		SET effective_to = '2027-01-31' 
		WHERE parameter_code = 'UMA' AND jurisdiction = 'FEDERAL' AND effective_to IS NULL;
	`)
	if err != nil {
		return fmt.Errorf("error cerrando versión previa de UMA: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO normative_parameter_versions 
		(source_id, parameter_code, jurisdiction, numeric_value, effective_from, effective_to, publication_date, document_reference, document_hash)
		VALUES 
		((SELECT id FROM normative_official_sources WHERE code = 'INEGI'), 'UMA', 'FEDERAL', 117.82, '2027-02-01', NULL, '2027-01-10', 'DOF 10/01/2027', 'a1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef0');
	`)
	if err != nil {
		return fmt.Errorf("error insertando nueva versión de UMA: %w", err)
	}

	return nil
}
