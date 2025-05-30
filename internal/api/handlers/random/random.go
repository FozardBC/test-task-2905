package random

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

type RandomGetter interface {
	RandomQuote(ctx context.Context) (*models.Quote, error)
}

func New(log *slog.Logger, getter RandomGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()

		log = log.With(requestid.ContextKeyRequestID, reqCtx.Value(requestid.ContextKeyRequestID))

		quote, err := getter.RandomQuote(reqCtx)
		if err != nil {
			if errors.Is(err, storage.ErrQuotesListEmpty) {
				log.Info("quotes list is empty", "error", err)

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response.OKWithPayload(nil))
				return
			}
			log.Error("failed to get random quote", "error", err, "code", http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Error("Internal server error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.OKWithPayload(quote))
		log.Info("random quote retrieved successfully", "quote", quote)

	}
}
