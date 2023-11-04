package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

// TODO: cCtxに依存しないようにする。
func SetLogger(cCtx *cli.Context) {
	logLevel := slog.LevelInfo
	if cCtx.Bool(debugFlag) {
		logLevel = slog.LevelDebug
	}

	var programLevel = new(slog.LevelVar)
	programLevel.Set(logLevel)
	var sh slog.Handler
	if cCtx.Bool(jsonFlag) {
		sh = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	} else {
		sh = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	}
	slog.SetDefault(slog.New(sh))
}

func Logger() *slog.Logger {
	return slog.Default()
}
