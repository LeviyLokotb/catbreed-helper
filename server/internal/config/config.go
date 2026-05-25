package config

import "os"

type Config struct {
	ModelPath  string
	LabelsPath string
}

func LoadFromEnv() Config {
	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		modelPath = "./model.onnx"
	}

	labelsPath := os.Getenv("LABELS_PATH")
	if labelsPath == "" {
		labelsPath = "./labels.txt"
	}

	return Config{
		ModelPath:  modelPath,
		LabelsPath: labelsPath,
	}
}
