package domain

import "time"

// DTO para CFDI 4.0 procesado
type ParsedCFDIDTO struct {
	Version           string    `json:"version"`
	Serie             string    `json:"serie"`
	Folio             string    `json:"folio"`
	EmisorRFC         string    `json:"emisor_rfc"`
	EmisorNombre      string    `json:"emisor_nombre"`
	ReceptorRFC       string    `json:"receptor_rfc"`
	ReceptorNombre    string    `json:"receptor_nombre"`
	SubTotalCents     int64     `json:"subtotal_cents"`
	TotalCents        int64     `json:"total_cents"`
	FechaEmisionUTC   time.Time `json:"fecha_emision_utc"`
	TipoDeComprobante string    `json:"tipo_de_comprobante"`
	FileSHA256        string    `json:"file_sha256"`
}

// DTO para Ficha SIPARE procesada
type ParsedSIPAREDTO struct {
	LineaCaptura     string    `json:"linea_captura"`
	TotalAmountCents int64     `json:"total_amount_cents"`
	DueDate          time.Time `json:"due_date"`
}

// DTO para respuesta SAT 32-D
type SAT32DResponseDTO struct {
	DocumentID      string    `json:"document_id"`
	RFC             string    `json:"rfc"`
	OpinionResult   string    `json:"opinion_result"` // 'POSITIVA', 'NEGATIVA', 'SIN_OBLIGACIONES'
	PdfRawBytes     []byte    `json:"-"`
	PdfSHA256       string    `json:"pdf_sha256"`
	DownloadedAtUTC time.Time `json:"downloaded_at_utc"`
}

// Error estructurado para barrera RPA
type RPABarrierError struct {
	WorkerID    string `json:"worker_id"`
	PortalName  string `json:"portal_name"`
	BarrierType string `json:"barrier_type"` // 'CAPTCHA', 'MFA', 'LAYOUT_CHANGE'
	CaptchaB64  string `json:"captcha_b64,omitempty"`
	Message     string `json:"message"`
}

func (e *RPABarrierError) Error() string {
	return e.Message
}
