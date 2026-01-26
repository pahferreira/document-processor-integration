package documentprocessor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	imagereader "github.com/openinfer/regula-integration/internal/image-reader"
	"github.com/openinfer/regula-integration/internal/recognito"
)

type RecognitoDocumentProcessor struct {
	RecognitoClient recognito.Client
	ImageReader imagereader.FileReader
}

type RecognitoResult struct {
	Ocr recognito.OCRDocumentResponse `json:"ocr"`
	Liveness recognito.LivenessResponse `json:"liveness"`
}

func NewRecognitoDocumentProcessor() DocumentProcessor {
	return &RecognitoDocumentProcessor{
		RecognitoClient: *recognito.NewClient(),
		ImageReader:  imagereader.NewImageReader(),
	}
}

func (r *RecognitoDocumentProcessor) ProcessDocuments(files []os.DirEntry, imagesDirectory string) (map[string][]any, error) {
	results := make(map[string]RecognitoResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent) // Semaphore to limit concurrency

	for _, file := range files {
		wg.Add(1)
		go func (file os.DirEntry) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			fmt.Printf("Reading file: %v \n", file.Name())
			filePath := filepath.Join(imagesDirectory, file.Name())
			imageData, err := r.ImageReader.Read(filePath)

			ocrResult, err := r.RecognitoClient.OCRIDCard([]byte(imageData), nil)
			if err != nil {
				fmt.Printf("❌ Unable to process OCR with Recognito: %v \n", err)
				return
			}

			// livenessResult, err := r.RecognitoClient.LivenessDetect([]byte(imageData))
			// if err != nil {
			// 	fmt.Printf("❌ Unable to process Liveness with Recognito: %v \n", err)
			// 	return
			// }

			mu.Lock()
			results[file.Name()] = RecognitoResult{
				Ocr: *ocrResult,
				// Liveness: *livenessResult,
			}
			mu.Unlock()
			fmt.Printf("✅ File %v processed\n", file.Name())

		}(file)		
	}

	wg.Wait()

	convertedResults := make(map[string][]any)

	for filename, recognitoResult := range results {
		convertedResults[filename] = []any{recognitoResult}
	}

	return convertedResults, nil
}

func (r *RecognitoDocumentProcessor) SaveResultsAsJSON(filename string, results map[string][]any) error {
	formattedResults := make(map[string]any)
	for key, value := range results {
		formattedResults[key] = value[0]
	}
	data, err := json.MarshalIndent(formattedResults, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}