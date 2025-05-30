package quteos

import (
	"app/internal/domain/models"
	"app/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrQuoteIsNil    = fmt.Errorf("quote is nil")
	ErrValidateQuote = fmt.Errorf("validation failed for quote")

	ErrSaveQuoteFailed   = fmt.Errorf("failed to save quote")
	ErrDeleteQuoteFailed = fmt.Errorf("failed to delete quote")
	ErrGetQuoteFailed    = fmt.Errorf("failed to get quote")
	ErrInvalidQuoteID    = fmt.Errorf("invalid quote ID, must be a positive integer")
	ErrInvalidAuthorName = fmt.Errorf("invalid author name, must be at least 2 characters long")
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Service struct {
	storage storage.Storage
	log     *slog.Logger
}

func New(storage *storage.Storage, log *slog.Logger) *Service {
	return &Service{
		storage: *storage,
		log:     log,
	}
}

func (s *Service) Save(ctx context.Context, q *models.Quote) (int, error) {
	s.log.Debug("Saving quote", "quote", q)

	if q == nil {
		s.log.Error(ErrQuoteIsNil.Error())

		return 0, ErrQuoteIsNil
	}

	// validate quote
	if err := validate.Struct(q); err != nil {
		var errMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			switch err.ActualTag() {
			case "required":
				errMsg.WriteString(err.Field() + " cannot be empty; ")
			case "max":
				errMsg.WriteString(err.Field() + " length exceeds " + err.Param() + " characters; ")
			case "printascii":
				errMsg.WriteString(err.Field() + " contains invalid characters; ")
			case "min":
				errMsg.WriteString(err.Field() + " must be at least " + err.Param() + " characters long; ")
			}
		}

		s.log.Error(ErrValidateQuote.Error(), "error", errMsg.String())

		return 0, fmt.Errorf("%w:%s", ErrValidateQuote, errMsg.String())
	}

	// save quote to storage

	id, err := s.storage.Save(ctx, q.Text, q.Author)
	if err != nil {
		s.log.Error(ErrSaveQuoteFailed.Error(), "error", err)

		return 0, fmt.Errorf("%w: %w", ErrSaveQuoteFailed, err)
	}

	s.log.Debug("Quote saved successfully", "quote", q)

	return id, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	s.log.Debug("Deleting quote", "id", id)

	intID, err := validateQuoteID(id)
	if err != nil {
		s.log.Error(ErrInvalidQuoteID.Error(), "error", err)
		return fmt.Errorf("%w: %s", ErrInvalidQuoteID, id)
	}

	if err := s.storage.Delete(ctx, intID); err != nil {
		if errors.Is(err, storage.ErrQuoteNotFound) {
			s.log.Error(ErrDeleteQuoteFailed.Error(), "err", storage.ErrQuoteNotFound, "id", id)

			return fmt.Errorf("%w: %w", ErrDeleteQuoteFailed, storage.ErrQuoteNotFound)
		}
		s.log.Error(ErrDeleteQuoteFailed.Error(), "error", err)

		return fmt.Errorf("%w: %w", ErrDeleteQuoteFailed, err)
	}

	s.log.Debug("Quote deleted successfully", "id", id)

	return nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Quote, error) {
	s.log.Debug("Getting quote", "id", id)

	intID, err := validateQuoteID(id)
	if err != nil {
		s.log.Error(ErrInvalidQuoteID.Error(), "error", err)
		return nil, fmt.Errorf("%w: %s", ErrInvalidQuoteID, id)
	}

	quote, err := s.storage.Get(ctx, intID)
	if err != nil {
		if errors.Is(err, storage.ErrQuoteNotFound) {
			s.log.Error(storage.ErrQuoteNotFound.Error(), "id", id)

			return nil, fmt.Errorf("%w: %s", storage.ErrQuoteNotFound, id)
		}
		s.log.Error(ErrGetQuoteFailed.Error(), "error", err)

		return nil, fmt.Errorf("%w: %w", ErrGetQuoteFailed, err)
	}

	s.log.Debug("Quote retrieved successfully", "quote", quote)

	return quote, nil
}

func (s *Service) List(ctx context.Context) ([]*storage.StorageQuote, error) {
	s.log.Debug("Listing all quotes")

	quotes, err := s.storage.List(ctx)
	if err != nil {
		s.log.Error(ErrGetQuoteFailed.Error(), "error", err)

		return nil, fmt.Errorf("%w: %w", ErrGetQuoteFailed, err)
	}

	s.log.Debug("Quotes retrieved successfully", "count", len(quotes))

	return quotes, nil
}

func (s *Service) ListByAuthor(ctx context.Context, author string) ([]*storage.StorageQuote, error) {
	s.log.Debug("Listing all quotes")

	if len(author) < 2 {
		s.log.Error(ErrInvalidAuthorName.Error(), "author", author)

		return nil, fmt.Errorf("%w: %s", ErrInvalidAuthorName, author)
	}

	quotes, err := s.storage.ListByAuthor(ctx, author)
	if err != nil {
		s.log.Error(ErrGetQuoteFailed.Error(), "error", err)

		return nil, fmt.Errorf("%w: %w", ErrGetQuoteFailed, err)
	}

	s.log.Debug("Quotes retrieved successfully", "count", len(quotes))

	return quotes, nil
}

func validateQuoteID(id string) (int, error) {
	if id == "" {
		return 0, fmt.Errorf("%w: %s", ErrInvalidQuoteID, id)
	}

	intID, err := strconv.Atoi(id)
	if err != nil {

		return 0, fmt.Errorf("%w: %s", ErrInvalidQuoteID, id)
	}

	return intID, nil
}
