package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LeviyLokotb/catbreed-helper-server/internal/ml"
	"github.com/LeviyLokotb/catbreed-helper-server/pkg/responseform"
)

// Главная функциональность.
// Принимает POST запрос с файлом,
// запрашивает анализ у ИНС,
// и возвращает json с результатом классификации
func PredictHandlerMock(w http.ResponseWriter, r *http.Request) {
	// Ограничение файла в 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// Парсим форму, получаем файл
	file, header, err := r.FormFile("file")
	if err != nil {
		err = fmt.Errorf("file reading error: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("File: %s, size: %d", header.Filename, header.Size)

	// Запрос к модели
	breed, confidence := ml.PredictBreed(file)

	// Отдаём JSON ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := responseform.CatBreedResponse{
		Breed:      breed,
		Confidence: confidence,
		FileName:   header.Filename,
	}

	json.NewEncoder(w).Encode(resp)
}
