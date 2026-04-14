package main

import (
	"log"
	"net/http"

	"github.com/LeviyLokotb/catbreed-helper-server/internal/config"
	"github.com/LeviyLokotb/catbreed-helper-server/internal/handlers"
	"github.com/LeviyLokotb/catbreed-helper-server/internal/ml"
)

func main() {
	// Загружаем модель
	conf := config.LoadFromEnv()
	if _, err := ml.GetCatBreedPredictor(conf); err != nil {
		log.Fatal(err)
	}

	// Настраиваем обработчики для эндпоинтов
	mux := http.NewServeMux()
	mux.HandleFunc("POST /predict", handlers.PredictHandler)
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// Логирование
	handler := handlers.LogMiddleware(mux)

	// Запускаем сервер
	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
