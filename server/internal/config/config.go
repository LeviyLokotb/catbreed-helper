package config

import "os"

type Config struct {
	ModelPath  string
	LabelsPath string
}

func LoadFromEnv() Config {
	return Config{
		ModelPath:  os.Getenv("MODEL_PATH"),
		LabelsPath: os.Getenv("LABELS_PATH"),
	}
}
