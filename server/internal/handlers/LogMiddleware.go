package handlers

import (
	"log"
	"net/http"
)

// Добавляем промежуточное логирование
func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}
