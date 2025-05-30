package save

import (
	requestid "app/internal/api/middleware/requestID"
	"app/internal/domain/models"
	"app/internal/lib/api/response"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Saver interface {
	Save(ctx context.Context, q *models.Quote) (int, error)
}

func New(log *slog.Logger, saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqCtx := r.Context()

		log = log.With(requestid.ContextKeyRequestID, reqCtx.Value(requestid.ContextKeyRequestID))

		var req models.Quote

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Error("Invalid request body"))
			return
		}
		defer r.Body.Close()

		err := validator.New().Struct(&req)
		if err != nil {
			validatorErr := err.(validator.ValidationErrors)

			log.Error("Invalid request body", "error", validatorErr)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.ValidationError(validatorErr))
			return
		}

		id, err := saver.Save(reqCtx, &req)
		if err != nil {
			log.Error("failed to save quote", "error", err, "code", http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Error("Internal server error"))

			return

		}

		log.Info("quote saved", "id", id)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.OKWithPayload(map[string]int{"id": id}))

	}
}
