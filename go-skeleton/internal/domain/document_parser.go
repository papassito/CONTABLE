package infrastructure

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/klik/fcos-kernel/internal/domain"
)

// Estructura XML Canónica CFDI 4.0 (Anexo 20 SAT)
type CFDIComprobante struct {
	XMLName           xml.Name      `xml:"Comprobante"`
	Version           string        `xml:"Version,attr"`
	Serie             string        `xml:"Serie,attr"`
	Folio             string        `xml:"Folio,attr"`
	Fecha             string        `xml:"Fecha,attr"`
	SubTotal          string        `xml:"SubTotal,attr"`
	Total             string        `xml:"Total,attr"`
	Moneda            string        `xml:"Moneda,attr"`
	TipoDeComprobante string        `xml:"TipoDeComprobante,attr"`
	Exportacion       string        `xml:"Exportacion,attr"`
	LugarExpedicion   string        `xml:"LugarExpedicion,attr"`
	Emisor            CFDIEmisor    `xml:"Emisor"`
	Receptor          CFDIReceptor  `xml:"Receptor"`
	Conceptos         CFDIConceptos `xml:"Conceptos"`
}

type CFDIEmisor struct {
	Rfc           string `xml:"Rfc,attr"`
	Nombre        string `xml:"Nombre,attr"`
	RegimenFiscal string `xml:"RegimenFiscal,attr"`
}

type CFDIReceptor struct {
	Rfc                     string `xml:"Rfc,attr"`
	Nombre                  string `xml:"Nombre,attr"`
	DomicilioFiscalReceptor string `xml:"DomicilioFiscalReceptor,attr"`
	RegimenFiscalReceptor   string `xml:"RegimenFiscalReceptor,attr"`
	UsoCFDI                 string `xml:"UsoCFDI,attr"`
}

type CFDIConceptos struct {
	Conceptos []CFDIConcepto `xml:"Concepto"`
}

type CFDIConcepto struct {
	ClaveProdServ string `xml:"ClaveProdServ,attr"`
	Cantidad      string `xml:"Cantidad,attr"`
	ClaveUnidad   string `xml:"ClaveUnidad,attr"`
	Descripcion   string `xml:"Descripcion,attr"`
	ValorUnitario string `xml:"ValorUnitario,attr"`
	Importe       string `xml:"Importe,attr"`
}

type DocumentParser interface {
	ParseCFDI40(xmlReader io.Reader) (*domain.ParsedCFDIDTO, error)
	ParseSIPARELine(rawText string) (*domain.ParsedSIPAREDTO, error)
}

type documentParser struct{}

func NewDocumentParser() DocumentParser {
	return &documentParser{}
}

func (p *documentParser) ParseCFDI40(xmlReader io.Reader) (*domain.ParsedCFDIDTO, error) {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(xmlReader, buf)

	// 1. Calcular Hash SHA-256 del archivo físico
	hash := sha256.New()
	if _, err := io.Copy(hash, tee); err != nil {
		return nil, fmt.Errorf("error leyendo stream para hash SHA-256: %w", err)
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	// 2. Deserializar XML
	var cfdi CFDIComprobante
	if err := xml.NewDecoder(buf).Decode(&cfdi); err != nil {
		return nil, fmt.Errorf("error decodificando XML CFDI 4.0: %w", err)
	}

	if cfdi.Version != "4.0" {
		return nil, fmt.Errorf("versión CFDI no soportada: %s (se requiere 4.0)", cfdi.Version)
	}

	// 3. Convertir importes a céntimos enteros (int64) para evitar float64
	subtotalCents, err := parseToCents(cfdi.SubTotal)
	if err != nil {
		return nil, fmt.Errorf("subtotal CFDI inválido: %w", err)
	}

	totalCents, err := parseToCents(cfdi.Total)
	if err != nil {
		return nil, fmt.Errorf("total CFDI inválido: %w", err)
	}

	parsedDate, err := time.Parse("2006-01-02T15:04:05", cfdi.Fecha)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha CFDI inválido: %w", err)
	}

	return &domain.ParsedCFDIDTO{
		Version:           cfdi.Version,
		Serie:             cfdi.Serie,
		Folio:             cfdi.Folio,
		EmisorRFC:         cfdi.Emisor.Rfc,
		EmisorNombre:      cfdi.Emisor.Nombre,
		ReceptorRFC:       cfdi.Receptor.Rfc,
		ReceptorNombre:    cfdi.Receptor.Nombre,
		SubTotalCents:     subtotalCents,
		TotalCents:        totalCents,
		FechaEmisionUTC:   parsedDate.UTC(),
		TipoDeComprobante: cfdi.TipoDeComprobante,
		FileSHA256:        fileHash,
	}, nil
}

func (p *documentParser) ParseSIPARELine(rawText string) (*domain.ParsedSIPAREDTO, error) {
	// Parseo de texto de archivo o PDF SIPARE
	lines := strings.Split(rawText, "\n")
	var dto domain.ParsedSIPAREDTO

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "LINEA_CAPTURA:") {
			dto.LineaCaptura = strings.TrimSpace(strings.TrimPrefix(line, "LINEA_CAPTURA:"))
		} else if strings.HasPrefix(line, "IMPORTE_TOTAL:") {
			valStr := strings.TrimSpace(strings.TrimPrefix(line, "IMPORTE_TOTAL:"))
			cents, err := parseToCents(valStr)
			if err != nil {
				return nil, fmt.Errorf("importe SIPARE inválido: %w", err)
			}
			dto.TotalAmountCents = cents
		} else if strings.HasPrefix(line, "VENCIMIENTO:") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "VENCIMIENTO:"))
			d, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return nil, fmt.Errorf("fecha vencimiento SIPARE inválida: %w", err)
			}
			dto.DueDate = d
		}
	}

	if dto.LineaCaptura == "" || dto.TotalAmountCents == 0 {
		return nil, fmt.Errorf("falla al parsear información esencial de la ficha SIPARE")
	}

	return &dto, nil
}

func parseToCents(valStr string) (int64, error) {
	valStr = strings.ReplaceAll(valStr, ",", "")
	parts := strings.Split(valStr, ".")

	if len(parts) == 1 {
		v, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, err
		}
		return v * 100, nil
	}

	if len(parts) == 2 {
		intPart, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, err
		}
		decStr := parts[1]
		if len(decStr) == 1 {
			decStr += "0"
		} else if len(decStr) > 2 {
			decStr = decStr[:2]
		}
		decPart, err := strconv.ParseInt(decStr, 10, 64)
		if err != nil {
			return 0, err
		}
		if intPart < 0 {
			return intPart*100 - decPart, nil
		}
		return intPart*100 + decPart, nil
	}

	return 0, fmt.Errorf("formato numérico no reconocido: %s", valStr)
}
