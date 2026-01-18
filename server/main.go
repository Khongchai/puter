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
	print("Starting puter...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lsproto.NewBaseReader(os.Stdin)
	outputWriter := lsproto.NewBaseWriter(os.Stdout)
	logger := logging.NewLogger(os.Stderr)
	engine := engine.NewEngine(
		ctx,
		inputReader,
		outputWriter,
		logger,
	)

	print("Engine running")
	engine.Run(ctx)
	print("Engine stopped")
}
