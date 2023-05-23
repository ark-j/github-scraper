package internal

import (
	"log"
	"os"
)

// NewLogger returns the preconfigured instance of logger
func NewLogger() *log.Logger {
	return log.New(os.Stderr, "github-scrape ", log.Ltime|log.LstdFlags|log.Lshortfile)
}
