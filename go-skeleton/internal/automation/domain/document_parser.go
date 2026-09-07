package domain

import "io"

type DocumentParser interface {
	ParseCFDI40(xmlReader io.Reader) (interface{}, error)
	ParseSIPARELine(rawText string) (interface{}, error)
}
