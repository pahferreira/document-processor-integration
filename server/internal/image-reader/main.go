package imagereader

import (
	"encoding/base64"
	"os"
)

type FileReader interface {
	Read(path string) (string, error)
	EncodeFile(fileBytes []byte) (string, error)
}

type ImageReader struct {}

func NewImageReader() FileReader {
	return &ImageReader{}
}

func (ir *ImageReader) Read(path string) (string, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(fileData), nil
}

func (ir *ImageReader) EncodeFile(fileBytes []byte) (string, error) {
	encodedFile := base64.StdEncoding.EncodeToString(fileBytes)
	return encodedFile, nil
}

