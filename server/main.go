package main

import (
	"context"
	"os"
	"os/signal"
	"puter/engine"
	"puter/interpreter"
	"puter/logging"
	lsproto "puter/lsp"
	"puter/unit"
	"syscall"
)

func main() {
	logger := logging.NewLogger(os.Stderr)
	logger.Info("Starting puter...\n")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inputReader := lsproto.NewBaseReader(os.Stdin)
	outputWriter := lsproto.NewBaseWriter(os.Stdout)

	currencyConverter := unit.GetCurrencyConverter()
	measurementConverter := unit.GetMeasurementConverter()
	converters := &unit.Converters{
		ConvertCurrency:    currencyConverter,
		ConvertMeasurement: measurementConverter,
	}

	interpreter := interpreter.NewInterpreter(ctx, converters)

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
