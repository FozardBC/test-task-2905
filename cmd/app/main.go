package main

import (
	"app/internal/config"
	"app/internal/storage/postgres"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"app/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	ctx := context.Background()

	cfg := config.MustRead()

	logger := logger.New(cfg.Log)

	storage, err := postgres.New(ctx, logger, cfg.DbConnString)
	if err != nil {
		logger.Error("failed to create storage", "error", err)
		return
	}

	err = startMigrations(logger, cfg.DbConnString)
	if err != nil {
		logger.Error("failet to start migrations", "error", err)
		return
	}

	_ = storage

	//server

	//graceful shutdown

}

func startMigrations(log *slog.Logger, connString string) error {
	m, err := migrate.New("file://migrations", connString) // DEBUG: ../../migrations"
	if err != nil {
		return fmt.Errorf("can't start migration driver:%w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Warn("Migarate didn't run. Nothing to change")
			return nil
		}
		return fmt.Errorf("failed to do migrations:%w", err)

	}

	return nil
}
