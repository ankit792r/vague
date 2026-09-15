package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vague/backbone/command"
)

//go:embed "frontend/dist/*"
var uiDist embed.FS

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := command.Execute(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
