package logger

import (
	requestid "app/internal/api/middleware/requestID"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			reqId := r.Context().Value(requestid.ContextKeyRequestID)
			if reqId == nil {
				log.Error("requestID not found in context")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			requestID, ok := reqId.(string)
			if !ok {
				log.Error("requestID is not a string")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestID),
			)

			entry.Info("request received")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)

	}
}
