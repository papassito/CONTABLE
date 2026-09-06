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
    Emisor            CFDIEmisor    `xml:"Emisor"`
    Receptor          CFDIReceptor  `xml:"Receptor"`
}

type CFDIEmisor struct {
    Rfc    string `xml:"Rfc,attr"`
    Nombre string `xml:"Nombre,attr"`
}

type CFDIReceptor struct {
    Rfc    string `xml:"Rfc,attr"`
    Nombre string `xml:"Nombre,attr"`
}

type documentParser struct{}

// NewDocumentParser retorna una nueva instancia del analizador de documentos.
func NewDocumentParser() domain.DocumentParser {
    return &documentParser{}
}

func (p *documentParser) ParseCFDI40(xmlReader io.Reader) (*domain.ParsedCFDIDTO, error) {
    buf := new(bytes.Buffer)
    tee := io.TeeReader(xmlReader, buf)
    hash := sha256.New()
    if _, err := io.Copy(hash, tee); err != nil {
        return nil, fmt.Errorf("error leyendo flujo de datos: %w", err)
    }
    fileHash := hex.EncodeToString(hash.Sum(nil))

    var cfdi CFDIComprobante
    if err := xml.NewDecoder(buf).Decode(&cfdi); err != nil {
        return nil, fmt.Errorf("error decodificando XML CFDI: %w", err)
    }

    subtotalCents, err := parseToCents(cfdi.SubTotal)
    if err != nil {
        return nil, err
    }
    totalCents, err := parseToCents(cfdi.Total)
    if err != nil {
        return nil, err
    }
    parsedDate, err := time.Parse("2006-01-02T15:04:05", cfdi.Fecha)
    if err != nil {
        parsedDate = time.Now()
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
                return nil, err
            }
            dto.TotalAmountCents = cents
        } else if strings.HasPrefix(line, "VENCIMIENTO:") {
            dateStr := strings.TrimSpace(strings.TrimPrefix(line, "VENCIMIENTO:"))
            d, err := time.Parse("2006-01-02", dateStr)
            if err == nil {
                dto.DueDate = d
            }
        }
    }
    if dto.LineaCaptura == "" || dto.TotalAmountCents == 0 {
        return nil, fmt.Errorf("fallo al procesar linea SIPARE")
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
        return intPart*100 + decPart, nil
    }
    return 0, fmt.Errorf("formato no reconocido: %s", valStr)
}
