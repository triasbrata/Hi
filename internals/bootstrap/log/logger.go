package log

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

func LoadLogger() fx.Option {
	return fx.Module("bootstrap/logger",
		fx.Provide(func() *slog.Logger {
			handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
			logger := slog.New(handler)
			slog.SetDefault(logger)
			return logger
		}),
	)
}
