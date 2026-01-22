package imagereader

import (
	"encoding/base64"
	"os"
)

type FileReader interface {
	Read(path string) (string, error)
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

	encodedFile := base64.StdEncoding.EncodeToString(fileData)
	return encodedFile, nil
}


