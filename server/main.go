package main

import (
	"context"
	"os"
	"os/signal"
	"puter/engine"
	"puter/logging"
	lsproto "puter/lsp"
	"syscall"
)

func main() {
	print("Starting vocab-ls...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lsproto.NewBaseReader(os.Stdin)
	outputWriter := lsproto.NewBaseWriter(os.Stdout)
	logger := logging.NewLogger(os.Stdout)
	engine := engine.NewEngine(
		ctx,
		inputReader,
		outputWriter,
		logger,
	)

	engine.Run(ctx)
}
