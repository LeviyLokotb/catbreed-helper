package ml

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"sync"

	"github.com/nfnt/resize"
	ort "github.com/yalue/onnxruntime_go"
)

type CatBreedPredictor struct {
	modelPath  string
	labels     []string
	inputSize  int
	inputName  string
	outputName string
	mu         sync.Mutex
	session    *ort.AdvancedSession // кешированная сессия
}

func Initialize() error {
	// Инициализируем ONNX Runtime (нужно вызвать до любых операций с тензорами)
	return ort.InitializeEnvironment()
}

func newCatBreedPredictor(modelPath, labelsPath string) (*CatBreedPredictor, error) {
	// Загружаем лейблы из txt файла
	labelsData, err := os.ReadFile(labelsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read labels: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(labelsData)), "\n")
	var labels []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			labels = append(labels, line)
		}
	}

	if len(labels) == 0 {
		return nil, errors.New("no labels found in file")
	}

	return &CatBreedPredictor{
		modelPath:  modelPath,
		labels:     labels,
		inputSize:  224,
		inputName:  "input",
		outputName: "output",
	}, nil
}

func (p *CatBreedPredictor) Predict(img image.Image) (*BreedPrediction, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	inputData := p.preprocessImage(img)

	inputShape := ort.NewShape(1, 3, 224, 224)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}
	defer inputTensor.Destroy()

	outputShape := ort.NewShape(1, int64(len(p.labels)))
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}
	defer outputTensor.Destroy()

	// Закрываем старую сессию, если есть
	if p.session != nil {
		p.session.Destroy()
	}

	// Создаём новую сессию
	p.session, err = ort.NewAdvancedSession(
		p.modelPath,
		[]string{p.inputName},
		[]string{p.outputName},
		[]ort.Value{inputTensor},
		[]ort.Value{outputTensor},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	err = p.session.Run()
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	outputData := outputTensor.GetData()

	maxIdx := 0
	maxVal := outputData[0]
	for i, v := range outputData {
		if v > maxVal {
			maxVal = v
			maxIdx = i
		}
	}

	return &BreedPrediction{
		Breed:      p.labels[maxIdx],
		Confidence: maxVal,
	}, nil
}

func (p *CatBreedPredictor) preprocessImage(img image.Image) []float32 {
	// Ресайз до 224x224
	resized := resize.Resize(224, 224, img, resize.Lanczos3)

	// NCHW формат: [1][3][224][224]
	input := make([]float32, 1*3*224*224)

	bounds := resized.Bounds()

	// Разделяем по каналам
	planeSize := 224 * 224
	rPlane := make([]float32, 0, planeSize)
	gPlane := make([]float32, 0, planeSize)
	bPlane := make([]float32, 0, planeSize)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Конвертация из 0-65535 в 0-1
			rPlane = append(rPlane, float32(r>>8)/255.0)
			gPlane = append(gPlane, float32(g>>8)/255.0)
			bPlane = append(bPlane, float32(b>>8)/255.0)
		}
	}

	// Копируем в порядке каналов: R, G, B
	copy(input[0*planeSize:1*planeSize], rPlane)
	copy(input[1*planeSize:2*planeSize], gPlane)
	copy(input[2*planeSize:3*planeSize], bPlane)

	return input
}

func (p *CatBreedPredictor) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.session != nil {
		p.session.Destroy()
		p.session = nil
	}
}
