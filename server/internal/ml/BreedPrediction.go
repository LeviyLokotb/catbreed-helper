package ml

import "strconv"

type BreedPrediction struct {
	Breed      string
	Confidence float32
}

func (p BreedPrediction) Deconstruct() (string, string) {
	conf := p.Confidence
	if conf < 0 {
		conf = 0
	}
	if conf > 1 {
		conf = 1
	}
	confidence := strconv.FormatFloat(float64(conf), 'f', 4, 32)
	return p.Breed, confidence
}
