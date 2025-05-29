package postgres

import (
	"app/internal/storage"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	QuoteTable      = "quotes"
	IdColumn        = "id"
	quoteColumn     = "quote"
	isDeletedColumn = "is_deleted"
)

type PostgreStorage struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, connString string) (*PostgreStorage, error) {
	log.Debug("Connecting to database", "Connect String", connString)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Error(storage.ErrConnectStorage.Error(), "err", err.Error())

		return nil, fmt.Errorf("%w:%w", storage.ErrConnectStorage, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Error(storage.ErrConnectStorage.Error(), "err", err.Error())

		return nil, fmt.Errorf("%w:%w", storage.ErrConnectStorage, err)
	}

	log.Debug("Database is connected")

	return &PostgreStorage{
		conn: conn,
		log:  log,
	}, nil
}
