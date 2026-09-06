package domain

import (
	"io"

	sharedDomain "github.com/klik/fcos-kernel/internal/domain"
)

type DocumentParser interface {
	ParseCFDI40(xmlReader io.Reader) (*sharedDomain.ParsedCFDIDTO, error)
	ParseSIPARELine(rawText string) (*sharedDomain.ParsedSIPAREDTO, error)
}
