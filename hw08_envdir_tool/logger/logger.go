package logger

import (
	"log/slog"
	"os"

	cslog "github.com/phsym/console-slog"
)

var Logger *slog.Logger

func init() {
	logHandler := cslog.NewHandler(os.Stderr, &cslog.HandlerOptions{
		Theme: cslog.NewBrightTheme(),
		Level: slog.LevelDebug,
	})

	Logger = slog.New(logHandler)
}
