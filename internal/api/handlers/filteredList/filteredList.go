package filteredlist

import (
	requestid "app/internal/api/middleware/requestID"
	"app/internal/domain/models"
	"app/internal/lib/api/response"
	"app/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type FileteredListGetter interface {
	ListByAuthor(ctx context.Context, author string) ([]*storage.StorageQuote, error)
}

func New(log *slog.Logger, listGetter FileteredListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		log = log.With(requestid.ContextKeyRequestID, reqCtx.Value(requestid.ContextKeyRequestID))

		author := r.URL.Query().Get("author")

		if len(author) < 3 || len(author) > 100 {
			log.Error("author query parameter is not valid", "author", author, "code", http.StatusBadRequest)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Error("Author query parameter is not valid. It must be between 3 and 100 characters long."))
			return
		}

		author = strings.ReplaceAll(author, "_", " ")

		list, err := listGetter.ListByAuthor(reqCtx, author)
		if err != nil {
			if errors.Is(err, storage.ErrQuotesListEmpty) {
				log.Info("quotes list is empty", "error", err)

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response.OKWithPayload([]*models.Quote{}))
				return
			}

			log.Error("failed to list quotes", "error", err, "code", http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Error("Internal server error"))
			return
		}

		log.Info("quotes listed successfully", "count", len(list))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.OKWithPayload(list))

	}
}
