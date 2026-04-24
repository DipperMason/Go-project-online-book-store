package logger

import (
	"io"
	"log/slog"
)

type EnvType string

const (
	EnvLocal EnvType = "local"
	EnvDev   EnvType = "dev"
	EnvProd  EnvType = "prod"
)

func Setup(env EnvType, out io.Writer) {
	var handler slog.Handler

	switch env {
	case EnvLocal:
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelWarn})
	}

	slog.SetDefault(slog.New(handler))
}

func With(args ...any) *slog.Logger {
	return slog.With(args)
}
