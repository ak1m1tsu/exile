package log

import (
	"io"
	"log/slog"
)

const (
	development = "development"
	production  = "production"
	local       = "local"
)

// New returns a slog.Logger instance.
//
// If env is development, it will json logger with debug level.
// If env is production, it will json logger with info level.
// If env is local, it will text logger with info level.
// Otherwise, it will text logger with info level.
func New(env string, out io.Writer) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case local:
		logger = slog.New(slog.NewTextHandler(
			out,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case development:
		logger = slog.New(slog.NewJSONHandler(
			out,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case production:
		logger = slog.New(slog.NewJSONHandler(
			out,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	default:
		logger = slog.New(slog.NewTextHandler(
			out,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	return logger
}
