package documentprocessor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	imagereader "github.com/openinfer/regula-integration/internal/image-reader"
	"github.com/openinfer/regula-integration/internal/regula"
)

type RegulaDocumentProcessor struct {
	RegulaClient regula.Client
	ImageReader	imagereader.FileReader
}

func NewRegulaDocumentProcessor() DocumentProcessor {
	return &RegulaDocumentProcessor{
		RegulaClient: *regula.NewClient(),
		ImageReader:  imagereader.NewImageReader(),
	}
}

const maxConcurrent = 10 // Limit concurrent API requests

func (r *RegulaDocumentProcessor) ProcessDocuments(files []os.DirEntry, imagesDirectory string) (map[string][]any, error) {
	results := make(map[string][]regula.TextField)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent) // Semaphore to limit concurrency

	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			fmt.Printf("Reading file: %v \n", file.Name())
			filePath := filepath.Join(imagesDirectory, file.Name())
			imageData, err := r.ImageReader.Read(filePath)
			if err != nil {
				fmt.Printf("❌ Unable to read Image File: %v \n", err)
				return
			}

			encodedImageData, err := r.ImageReader.EncodeFile([]byte(imageData))
			if err != nil {
				fmt.Printf("❌ Unable to encode Image File: %v \n", err)
				return
			}

			imageProcessResult, err := r.RegulaClient.ProcessRegulaDocumentRead(encodedImageData)
			if err != nil {
				fmt.Printf("❌ Unable to process Image File with Regula: %v \n", err)
				return
			}

			mu.Lock()
			results[file.Name()] = imageProcessResult
			mu.Unlock()
			fmt.Printf("✅ File %v processed\n", file.Name())
		}(file)
	}

	wg.Wait() // Wait for all goroutines to complete

	convertedResults := make(map[string][]any)

	for filename, textFields := range results {
		anySlice := make([]any, len(textFields))
		for i, textField := range textFields {
			anySlice[i] = textField
		}
		convertedResults[filename] = anySlice
	}

	return convertedResults, nil
}

func (r *RegulaDocumentProcessor) SaveResultsAsJSON(filename string, results map[string][]any) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}