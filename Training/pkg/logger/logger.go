package logger

import (
	"log/slog"
	"os"
)

var Logger = slog.New(
	slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		}),
)
