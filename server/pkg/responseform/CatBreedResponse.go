package responseform

type CatBreedResponse struct {
	Breed      string `json:"breed"`
	Confidence string `json:"confidence"`
	FileName   string `json:"filename"`
}
