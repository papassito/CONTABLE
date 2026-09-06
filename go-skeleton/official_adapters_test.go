package integration_test

import (
	"context"
	"strings"
	"testing"

	automationApp "github.com/klik/fcos-kernel/internal/automation/infrastructure"
	"github.com/klik/fcos-kernel/internal/domain"
)

func TestCFDI40XMLParsingAndAmountConversion(t *testing.T) {
	parser := automationApp.NewDocumentParser()

	rawXML := `<?xml version="1.0" encoding="UTF-8"?>
	<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" Version="4.0" Serie="F" Folio="1024" Fecha="2026-09-06T10:00:00" SubTotal="1000.50" Total="1160.58" Moneda="MXN" TipoDeComprobante="I" Exportacion="01" LugarExpedicion="83000">
		<cfdi:Emisor Rfc="CSO180512AAA" Nombre="COMERCIALIZADORA SONORA SA DE CV" RegimenFiscal="601"/>
		<cfdi:Receptor Rfc="XAXX010101000" Nombre="PUBLICO EN GENERAL" DomicilioFiscalReceptor="83000" RegimenFiscalReceptor="616" UsoCFDI="S01"/>
		<cfdi:Conceptos>
			<cfdi:Concepto ClaveProdServ="84111506" Cantidad="1" ClaveUnidad="E48" Descripcion="Servicios Contables FCOS" ValorUnitario="1000.50" Importe="1000.50"/>
		</cfdi:Conceptos>
	</cfdi:Comprobante>`

	dto, err := parser.ParseCFDI40(strings.NewReader(rawXML))
	if err != nil {
		t.Fatalf("Fallo en ParseCFDI40: %v", err)
	}

	// Aserción 1: Subtotal en céntimos ($1,000.50 = 100050 céntimos)
	if dto.SubTotalCents != 100050 {
		t.Fatalf("Subtotal esperado 100050 céntimos, obtenido: %d", dto.SubTotalCents)
	}

	// Aserción 2: Total en céntimos ($1,160.58 = 116058 céntimos)
	if dto.TotalCents != 116058 {
		t.Fatalf("Total esperado 116058 céntimos, obtenido: %d", dto.TotalCents)
	}

	// Aserción 3: Validación de RFCs
	if dto.EmisorRFC != "CSO180512AAA" || dto.ReceptorRFC != "XAXX010101000" {
		t.Fatalf("RFC Emisor o Receptor no coinciden con el XML")
	}

	// Aserción 4: Presencia de Hash SHA-256
	if dto.FileSHA256 == "" {
		t.Fatalf("El parser debe calcular obligatoriamente el Hash SHA-256 del archivo")
	}
}

func TestSIPARELineParsing(t *testing.T) {
	parser := automationApp.NewDocumentParser()

	rawText := `FICHA DE PAGO SIPARE IMSS
	REGISTRO_PATRONAL: A123456710
	LINEA_CAPTURA: 0123456789012345678901234
	IMPORTE_TOTAL: 4,285.18
	VENCIMIENTO: 2026-09-17`

	dto, err := parser.ParseSIPARELine(rawText)
	if err != nil {
		t.Fatalf("Fallo en ParseSIPARELine: %v", err)
	}

	if dto.LineaCaptura != "0123456789012345678901234" {
		t.Fatalf("Línea de captura incorrecta: %s", dto.LineaCaptura)
	}

	// $4,285.18 = 428518 céntimos
	if dto.TotalAmountCents != 428518 {
		t.Fatalf("Importe total esperado 428518 céntimos, obtenido: %d", dto.TotalAmountCents)
	}
}

func TestSATAdapterBarrierHandling(t *testing.T) {
	parser := automationApp.NewDocumentParser()
	adapter := automationApp.NewSATAdapter(parser)

	// Simular detección de barrera CAPTCHA usando el contexto de prueba
	ctx := context.WithValue(context.Background(), "simulate_barrier", true)

	_, err := adapter.Download32DOpinion(ctx, "CSO180512AAA", []byte("secret_mock_data"))

	if err == nil {
		t.Fatalf("Se esperaba fallo por barrera CAPTCHA activa")
	}

	barrierErr, ok := err.(*domain.RPABarrierError)
	if !ok {
		t.Fatalf("El error devuelto debe ser de tipo RPABarrierError")
	}

	if barrierErr.BarrierType != "CAPTCHA" || barrierErr.WorkerID == "" {
		t.Fatalf("Propiedades del error de barrera RPA incorrectas: %v", barrierErr)
	}
}
