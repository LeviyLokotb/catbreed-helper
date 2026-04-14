package main

import (
	"log"
	"net/http"

	"github.com/LeviyLokotb/catbreed-helper-server/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /predict", handlers.PredictHandlerMock)
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	handler := handlers.LogMiddleware(mux)

	log.Println("Server started at :8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
