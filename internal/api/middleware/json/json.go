package json

import "net/http"

func JSONContentTypeMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)
	})
}
