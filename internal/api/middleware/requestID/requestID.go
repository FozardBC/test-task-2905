package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const ContextKeyRequestID contextKey = "requestID"

func RequestIdMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := uuid.New()

		ctx = context.WithValue(ctx, ContextKeyRequestID, id.String())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
