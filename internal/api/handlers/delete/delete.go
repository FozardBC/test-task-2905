package delete

import (
	requestid "app/internal/api/middleware/requestID"
	"app/internal/lib/api/response"
	"app/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Deleter interface {
	Delete(ctx context.Context, id string) error
}

func New(log *slog.Logger, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqCtx := r.Context()

		log = log.With(requestid.ContextKeyRequestID, reqCtx.Value(requestid.ContextKeyRequestID))

		vars := mux.Vars(r)

		id, ok := vars["id"]
		if !ok || id == "" {
			log.Error("quote ID is missing in the request", "code", http.StatusBadRequest)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Error("Quote ID is required"))
			return
		}

		if len(id) < 1 || len(id) > 10 {
			log.Error("quote ID is not valid", "id", id, "code", http.StatusBadRequest)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Error("Quote ID must be between 1 and 10 characters long"))
			return
		}

		err := deleter.Delete(reqCtx, id)
		if err != nil {
			if errors.Is(err, storage.ErrQuoteNotFound) {
				log.Info("quote not found", "id", id, "code", http.StatusNotFound)
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response.Error("Quote not found"))
				return
			}
			log.Error("failed to delete quote", "error", err, "code", http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Error("Internal server error"))
			return
		}

		log.Info("quote deleted successfully", "id", id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.OKWithPayload(map[string]string{"message": "Quote deleted successfully"}))

	}
}
