package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vagues/backbone/command"
)

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
