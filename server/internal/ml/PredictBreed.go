package ml

import "mime/multipart"

func PredictBreed(file multipart.File) (string, string) {
	return "some cat (idk)", "0.5"
}
