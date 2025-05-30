package postgres

import (
	"app/internal/domain/models"
	"app/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	QuoteTable      = "quotes"
	IdColumn        = "id"
	quoteColumn     = "quote"
	authorColumn    = "author"
	isDeletedColumn = "is_deleted"
)

type PostgreStorage struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

var (
	ErrTxBegin  = errors.New("can't start transaction")
	ErrTxCommit = errors.New("can't commit transaction")
	ErrQuery    = errors.New("can't do query")
)

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

func (p *PostgreStorage) Save(ctx context.Context, quote string, author string) (int, error) {

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s) VALUES ($1,$2) RETURNING %s",
		QuoteTable,
		quoteColumn,
		authorColumn,
		IdColumn,
	)

	var id int

	err := p.conn.QueryRow(ctx, query, quote, author).Scan(&id)
	if err != nil {
		p.log.Error(storage.ErrFailedToSaveQuote.Error(), "error", err)

		return 0, fmt.Errorf("%w: %w", storage.ErrFailedToSaveQuote, err)
	}

	p.log.Debug("Quote saved successfully", "id", id, "quote", quote, "author", author)
	return id, nil

}

func (p *PostgreStorage) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM quotes WHERE id = $1;`

	result, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		p.log.Error("Failed to delete quote", "error", err, "id", id)
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		p.log.Warn("Quote not found", "id", id)
		return storage.ErrQuoteNotFound // Ваша кастомная ошибка (например, "quote not found")
	}

	p.log.Debug("Quote deleted successfully", "id", id)
	return nil
}

func (p *PostgreStorage) Get(ctx context.Context, id int) (*models.Quote, error) {
	return nil, nil
}

func (p *PostgreStorage) List(ctx context.Context) ([]*storage.StorageQuote, error) {

	tx, err := p.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		p.log.Error(ErrTxBegin.Error(), "err", err.Error())

		return nil, fmt.Errorf("%w:%w", ErrTxBegin, err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s",
		IdColumn,
		quoteColumn,
		authorColumn,
		QuoteTable,
	)

	var quotes []*storage.StorageQuote

	rows, err := tx.Query(ctx, query)
	if err != nil {
		p.log.Error("Failed to query quotes", "error", err)
		return nil, fmt.Errorf("failed to query quotes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var q storage.StorageQuote
		if err := rows.Scan(&q.Id, &q.Text, &q.Author); err != nil {
			p.log.Error("Failed to scan quote", "error", err)
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, &q)
	}

	tx.Commit(ctx)

	if len(quotes) == 0 {
		p.log.Info(storage.ErrQuotesListEmpty.Error())
		return nil, storage.ErrQuotesListEmpty
	}

	return quotes, nil

}

func (p *PostgreStorage) ListByAuthor(ctx context.Context, author string) ([]*storage.StorageQuote, error) {
	tx, err := p.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		p.log.Error(ErrTxBegin.Error(), "err", err.Error())

		return nil, fmt.Errorf("%w:%w", ErrTxBegin, err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s WHERE %s = $1",
		IdColumn,
		quoteColumn,
		authorColumn,
		QuoteTable,
		authorColumn,
	)

	var quotes []*storage.StorageQuote

	rows, err := tx.Query(ctx, query, author)
	if err != nil {
		p.log.Error("Failed to query quotes", "error", err)
		return nil, fmt.Errorf("failed to query quotes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var q storage.StorageQuote
		if err := rows.Scan(&q.Id, &q.Text, &q.Author); err != nil {
			p.log.Error("Failed to scan quote", "error", err)
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, &q)
	}

	tx.Commit(ctx)

	if len(quotes) == 0 {
		p.log.Info("No quotes found for author", "author", author)
		return nil, storage.ErrQuotesListEmpty
	}

	return quotes, nil
}

func (p *PostgreStorage) Random(ctx context.Context) (*models.Quote, error) {

	query := fmt.Sprintf(
		"SELECT %s, %s FROM %s ORDER BY RANDOM() LIMIT 1",
		quoteColumn,
		authorColumn,
		QuoteTable,
	)

	var quote models.Quote
	err := p.conn.QueryRow(ctx, query).Scan(&quote.Text, &quote.Author)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrQuotesListEmpty
		}
		p.log.Error("Failed to get random quote", "error", err)
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &quote, nil
}

func (p *PostgreStorage) Ping(ctx context.Context) error {
	if err := p.conn.Ping(ctx); err != nil {
		p.log.Error("Failed to ping database", "error", err)
		return fmt.Errorf("%w:%w", storage.ErrPingStorage, err)
	}
	p.log.Debug("Database ping successful")
	return nil
}

func (p *PostgreStorage) Close() {
	if p.conn != nil {
		p.log.Debug("Closing database connection")
		p.conn.Close()
	} else {
		p.log.Warn("No database connection to close")
	}
	p.log.Debug("Database connection closed")
}
