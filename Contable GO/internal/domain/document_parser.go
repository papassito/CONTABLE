package domain

import "io"

// DocumentParser define la interfaz para el análisis de documentos contables.
type DocumentParser interface {
	ParseCFDI40(xmlReader io.Reader) (*ParsedCFDIDTO, error)
	ParseSIPARELine(rawText string) (*ParsedSIPAREDTO, error)
}
