package middlewares

import (
	"net/http"
)

// JSONMiddleware appends the "Content-Type" headers on the response paylaod
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
