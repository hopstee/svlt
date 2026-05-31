package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hopstee/svlt/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go systemSygnals(cancel)

	app := app.InitApp(ctx)
	defer func() {
		if err := app.Stop(); err != nil {
			slog.Info("Failed to stop app", slog.Any("error", err))
		}
	}()

	if err := app.Run(); err != nil {
		slog.Error("App did not terminate correctly", slog.Any("error", err))
		os.Exit(1)
	}
}

func systemSygnals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	cancel()
}
