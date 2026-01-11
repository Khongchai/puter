package main

import (
	"context"
	"os"
	"os/signal"
	"puter/engine"
	lsproto "puter/lsp"
	"syscall"
)

func main() {
	print("Starting vocab-ls...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lsproto.NewBaseReader(os.Stdin)
	outputWriter := lsproto.NewBaseWriter(os.Stdout)
	engine := engine.NewEngine(
		ctx,
		inputReader,
		outputWriter,
	)

	engine.Run(ctx)
}
