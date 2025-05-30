package list

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
)

type ListGetter interface {
	List(ctx context.Context) ([]*storage.StorageQuote, error)
}

func New(log *slog.Logger, listGetter ListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		log = log.With(requestid.ContextKeyRequestID, reqCtx.Value(requestid.ContextKeyRequestID))

		list, err := listGetter.List(reqCtx)
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
