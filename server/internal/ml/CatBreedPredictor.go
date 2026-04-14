package ml

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"sort"
	"sync"

	tflite "github.com/mattn/go-tflite"
	"github.com/mattn/go-tflite/delegates/xnnpack"
	"github.com/nfnt/resize"
)

type CatBreedPredictor struct {
	interpreter *tflite.Interpreter
	labels      []string
	inputSize   int
	mu          sync.Mutex
}

func newCatBreedPredictor(modelPath, labelsPath string) (*CatBreedPredictor, error) {
	// Загружаем модель
	model := tflite.NewModelFromFile(modelPath)
	if model == nil {
		return nil, errors.New("failed to load model")
	}

	// Делегат для ускорения на CPU
	// (плагин xnnpack, перехватывает операции и оптимизирует их)
	options := tflite.NewInterpreterOptions()
	defer options.Delete()

	delegate := xnnpack.New(xnnpack.DelegateOptions{})
	options.AddDelegate(delegate)

	// Интерпретатор
	interpreter := tflite.NewInterpreter(model, options)
	if interpreter == nil {
		return nil, errors.New("failed to create interpreter")
	}

	// Аллоцируем тензоры (библиотека использует cgo)
	status := interpreter.AllocateTensors()
	if status != tflite.OK {
		return nil, errors.New("allocate tensors error")
	}

	// Проверяем входной тензор
	input := interpreter.GetInputTensor(0)
	if input.NumDims() != 4 {
		return nil, fmt.Errorf("expected 4D tensor as input")
	}
	inputSize := input.Dim(1)

	// Загружаем заголовки классов
	labelsData, err := os.ReadFile(labelsPath)
	if err != nil {
		return nil, err
	}
	var labels []string
	json.Unmarshal(labelsData, &labels)

	return &CatBreedPredictor{
		interpreter: interpreter,
		labels:      labels,
		inputSize:   inputSize,
	}, nil
}

func (p *CatBreedPredictor) Predict(img image.Image) (*BreedPrediction, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Получаем входной тензор
	inputTensor := p.interpreter.GetInputTensor(0)

	// Заполняем его данными
	inputData := p.preprocessImage(img)
	inputTensor.SetFloat32s(inputData)

	// Инференс
	p.interpreter.Invoke()

	// Получаем выходной тензор
	outputTensor := p.interpreter.GetOutputTensor(0)
	outputData := outputTensor.Float32s()

	// Находим самое вероятное предсказание
	if len(outputData) != len(p.labels) {
		return nil, fmt.Errorf("unknown label")
	}

	// Сортируем индексы по убыванию вероятности
	indices := make([]int, len(outputData))
	for i := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		return outputData[indices[i]] > outputData[indices[j]]
	})

	result := &BreedPrediction{
		Breed:      p.labels[indices[0]],
		Confidence: outputData[indices[0]],
	}

	return result, nil
}

func (p *CatBreedPredictor) preprocessImage(img image.Image) []float32 {
	// Ресайз до нужного размера
	resized := resize.Resize(uint(p.inputSize), uint(p.inputSize), img, resize.Lanczos3)

	// Создаем массив float32 размером [1][inputSize][inputSize][3]
	input := make([]float32, 1*p.inputSize*p.inputSize*3)

	bounds := resized.Bounds()
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Конвертация из 0-65535 в 0-1 и нормализация
			// Уточните у Python-разработчика точные параметры нормализации!
			input[idx] = float32(r>>8) / 255.0
			input[idx+1] = float32(g>>8) / 255.0
			input[idx+2] = float32(b>>8) / 255.0
			idx += 3
		}
	}

	return input
}

func (p *CatBreedPredictor) Close() {
	p.interpreter.Delete()
}
