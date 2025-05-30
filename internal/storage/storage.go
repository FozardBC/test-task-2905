package storage

import (
	"app/internal/domain/models"
	"context"
	"errors"
)

var (
	ErrConnectStorage = errors.New("failed to connect storage")
	ErrInvalidQuote   = errors.New("invalid quote")
	ErrPingStorage    = errors.New("failed to ping storage")

	ErrQuoteNotFound = errors.New("quote not found")

	ErrFailedToSaveQuote    = errors.New("failed to save quote")
	ErrFailedToDeleteQuote  = errors.New("failed to delete quote")
	ErrFailedToGetQuote     = errors.New("failed to get quote")
	ErrFailedToListQuotes   = errors.New("failed to list quotes")
	ErrFailedToListByAuthor = errors.New("failed to list quotes by author")
	ErrQuotesListEmpty      = errors.New("quotes list is empty")
)

type StorageQuote struct {
	models.Quote
	Id int `json:"id"`
}

type Storage interface {
	Save(ctx context.Context, quote string, author string) (int, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Quote, error)
	List(ctx context.Context) ([]*StorageQuote, error)
	ListByAuthor(ctx context.Context, author string) ([]*StorageQuote, error)
	Random(ctx context.Context) (*models.Quote, error)
	Close()
}
