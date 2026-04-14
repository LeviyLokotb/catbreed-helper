package ml

import "strconv"

type BreedPrediction struct {
	Breed      string
	Confidence float32
}

func (p BreedPrediction) Deconstruct() (string, string) {
	confidence := strconv.FormatFloat(float64(p.Confidence), 'f', 2, 32)
	return p.Breed, confidence
}
