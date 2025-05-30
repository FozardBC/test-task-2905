package main

import (
	"app/internal/api"
	"app/internal/config"
	"app/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	InfoDbClosed = "Storage is closed. App is shuting down"
)

func main() {

	ctx := context.Background()

	cfg := config.MustRead()

	logger := logger.New(cfg.Log)

	// Initialize storage
	storage, err := postgres.New(ctx, logger, cfg.DbConnString)
	if err != nil {
		logger.Error("failed to create storage", "error", err)
		return
	}

	// start migrations
	err = startMigrations(logger, cfg.DbConnString)
	if err != nil {
		logger.Error("failed to start migrations", "error", err)
		return
	}

	// initialize HTTP API
	API := api.New(storage, logger)

	srv := http.Server{
		Addr:    cfg.ServerHost + ":" + cfg.ServerPort,
		Handler: &API.Router,
	}

	// Graceful shutdown
	chanErrors := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("HTTP server started", "Addres", srv.Addr)
		chanErrors <- srv.ListenAndServe()
	}()

	go func() {
		logger.Info("Started to ping databse")
		for {
			time.Sleep(5 * time.Second)
			err := storage.Ping(ctx)
			if err != nil {
				chanErrors <- err
				break
			}
		}

	}()

	logger.Info("HTTP server is runned", "addres", srv.Addr)

	logger.Info("App is started")

	select {
	case err := <-chanErrors:
		logger.Error("Shutting down. Critical error:", "err", err)

		shutdown <- syscall.SIGTERM
	case sig := <-shutdown:
		logger.Error("received signal, starting graceful shutdown", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("server graceful shutdown failed", "err", err)
			err = srv.Close()
			if err != nil {
				logger.Error("forced shutdown failed", "err", err)
			}
		}

		storage.Close()

		logger.Info(InfoDbClosed)

		logger.Info("shutdown completed")

	}

}

func startMigrations(logger *slog.Logger, connString string) error {
	m, err := migrate.New("file://migrations", connString) // DEBUG:"
	if err != nil {
		return fmt.Errorf("can't start migration driver:%w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Warn("Migarate didn't run. Nothing to change")
			return nil
		}
		return fmt.Errorf("failed to do migrations:%w", err)

	}

	return nil
}
