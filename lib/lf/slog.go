package lf

import (
	"log/slog"
	"os"
)

const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelTrace            = slog.Level(-8)
)

func newSlogLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	// formatter := slogformatter.FormatByKind(slog.KindAny, func(v slog.Value) slog.Value {
	// 	return slog.StringValue(fmt.Sprintf("%T%+v", v.Any(), v.Any()))
	// })
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	// return slog.Default()
}
