package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init khởi tạo logger cho toàn bộ service
func Init(service string) {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = log.Output(os.Stdout).
		With().
		Str("service", service).
		Logger()
}
