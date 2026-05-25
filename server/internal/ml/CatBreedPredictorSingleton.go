package ml

import (
	"sync"

	"github.com/LeviyLokotb/catbreed-helper-server/internal/config"
)

var (
	instance *CatBreedPredictor
	once     sync.Once
	initErr  error
)

func GetCatBreedPredictor(conf config.Config) (*CatBreedPredictor, error) {
	err := Initialize()
	if err != nil {
		// return nil, fmt.Errorf("failed to initialize tensors: %w", err)
	}

	once.Do(func() {
		instance, initErr = newCatBreedPredictor(conf.ModelPath, conf.LabelsPath)
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}
