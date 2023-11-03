package main

import (
	"log/slog"
	"os"
)

const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
)

func SetLogger(lvl slog.Level) {
	var programLevel = new(slog.LevelVar)
	programLevel.Set(lvl)
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
}

func Logger() *slog.Logger {
	return slog.Default()
}
