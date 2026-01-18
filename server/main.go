package main

import (
	"context"
	"os"
	"os/signal"
	"puter/engine"
	"puter/interpreter"
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
	interpreter := interpreter.NewInterpreter()
	engine := engine.NewEngine(
		ctx,
		inputReader,
		outputWriter,
		logger,
		interpreter,
	)

	print("Engine running")
	engine.Run(ctx)
	print("Engine stopped")
}
