package console

import (
	"log"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteResult(url string, statusCode int, message string) error {
	log.Printf(`
URL: %s
Status: %d
Message: %s

`, url, statusCode, message)
	return nil
}
