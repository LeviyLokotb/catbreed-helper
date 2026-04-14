package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/LeviyLokotb/catbreed-helper-server/internal/config"
	"github.com/LeviyLokotb/catbreed-helper-server/internal/ml"
	"github.com/LeviyLokotb/catbreed-helper-server/pkg/responseform"
)

// Главная функциональность.
// Принимает POST запрос с файлом,
// запрашивает анализ у ИНС,
// и возвращает json с результатом классификации
func PredictHandler(w http.ResponseWriter, r *http.Request) {
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

	// Декодируем в image.Image
	img, format, err := image.Decode(file)
	if err != nil {
		err = fmt.Errorf("image decoding error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Received %s image: %s, format: %s", header.Filename, img.Bounds(), format)

	// Получение доступа к модели
	conf := config.LoadFromEnv()
	model, err := ml.GetCatBreedPredictor(conf)
	if err != nil {
		err = fmt.Errorf("loading model error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Предсказание
	pred, err := model.Predict(img)
	if err != nil {
		err = fmt.Errorf("prediction error: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	breed, confidence := pred.Deconstruct()

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
