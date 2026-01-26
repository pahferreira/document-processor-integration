package main

import (
	"fmt"
	"os"
	"path/filepath"

	documentprocessor "github.com/openinfer/regula-integration/internal/document-processor"
)

const imagesDirectory = "../document-images"
const resultsDirectory = "../results"

func main() {
	fmt.Println("🔍 Testing Regula API Connection...")

	cwd, _ := os.Getwd()
	fmt.Println(cwd)

	// Create results directory if it doesn't exist
	if err := os.MkdirAll(resultsDirectory, 0755); err != nil {
		fmt.Printf("❌ Unable to create results directory: %v \n", err)
		return
	}

	files, err := os.ReadDir(imagesDirectory)
	if err != nil {
		fmt.Printf("❌ Unable to read Images directory: %v \n", err)
		return
	}

	documentProcessor := documentprocessor.NewRegulaDocumentProcessor()
	results, err := documentProcessor.ProcessDocuments(files, imagesDirectory)
	if err != nil {
		fmt.Printf("❌ Unable to process documents: %v \n", err)
		return
	}

	jsonPath := filepath.Join(resultsDirectory, "results.json")
	if err = documentProcessor.SaveResultsAsJSON(jsonPath, results); err != nil {
		fmt.Printf("❌ Unable to save JSON results: %v \n", err)
		return
	}

	fmt.Printf("✅ Results saved to %s\n", jsonPath)

	fmt.Println("✅ Test Successful!")
}