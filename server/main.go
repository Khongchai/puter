package main

import (
	"context"
	"os"
	"os/signal"
	"puter/engine"
	"puter/engine/extensions"
	"puter/interpreter"
	"puter/logging"
	lsproto "puter/lsp"
	"syscall"
)

func main() {
	logger := logging.NewLogger(os.Stderr)
	logger.Info("Starting puter...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lsproto.NewBaseReader(os.Stdin)
	outputWriter := lsproto.NewBaseWriter(os.Stdout)

	currencyConverter := extensions.GetCurrencyConverter()
	interpreter := interpreter.NewInterpreter(ctx, currencyConverter)

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
