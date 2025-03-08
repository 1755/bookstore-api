package api

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	configPath ConfigPath
	server     *Server
}

func NewApplication(
	configPath ConfigPath,
	server *Server,
) *Application {
	return &Application{configPath, server}
}

func (a *Application) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Run server
	errChan := make(chan error, 1)
	go func() {
		errChan <- a.server.Run()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		return a.server.Shutdown(shutdownCtx)
	}
}
