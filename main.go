package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}
