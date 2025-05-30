package list

import (
	requestid "app/internal/api/middleware/requestID"
	"app/internal/lib/api/response"
	"app/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type ListGetter interface {
	List(ctx context.Context) ([]*storage.StorageQuote, error)
	ListByAuthor(ctx context.Context, author string) ([]*storage.StorageQuote, error)
}

func New(log *slog.Logger, listGetter ListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		log := log.With("requestID", reqCtx.Value(requestid.ContextKeyRequestID))

		var list []*storage.StorageQuote

		author := r.URL.Query().Get("author")

		// Search by quert param author
		if author != "" {

			if len(author) < 3 || len(author) > 100 {
				log.Error("author query parameter is not valid", "author", author, "code", http.StatusBadRequest)

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.Error("Author query parameter is not valid. It must be between 3 and 100 characters long."))
				return
			}

			author = strings.ReplaceAll(author, "_", " ")

			list, err := FindByAuthor(reqCtx, author, listGetter)
			if err != nil {
				if errors.Is(err, storage.ErrQuotesListEmpty) {
					log.Info("quotes list is empty", "error", err)

					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response.OKWithPayload([]*storage.StorageQuote{}))
					return
				}

				log.Error("failed to list quotes by author", "error", err, "code", http.StatusInternalServerError)

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response.Error("Internal server error"))
				return
			}

			log.Info("quotes listed successfully by author", "author", author, "count", len(list))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response.OKWithPayload(list))
			return

		}

		// If no author query param, list all quotes
		list, err := SimpleFind(reqCtx, listGetter)
		if err != nil {
			if errors.Is(err, storage.ErrQuotesListEmpty) {
				log.Info("quotes list is empty", "error", err)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response.OKWithPayload(list))

		}

		log.Info("quotes listed successfully", "count", len(list))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.OKWithPayload(list))

	}
}

func FindByAuthor(reqCtx context.Context, author string, listGetter ListGetter) ([]*storage.StorageQuote, error) {
	list, err := listGetter.ListByAuthor(reqCtx, author)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func SimpleFind(reqCtx context.Context, listGetter ListGetter) ([]*storage.StorageQuote, error) {
	list, err := listGetter.List(reqCtx)
	if err != nil {
		return nil, err
	}

	return list, nil
}
