package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	imagereader "github.com/openinfer/regula-integration/internal/image-reader"
	regula "github.com/openinfer/regula-integration/internal/regula"
)

const imagesDirectory = "../document-images"
const resultsDirectory = "../results"

var imageReader = imagereader.NewImageReader()
var regulaClient = regula.NewClient()

func saveResultsAsJSON(filename string, results map[string]*regula.DocumentProcessResponse) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func saveResultsAsCSV(filename string, results map[string]*regula.DocumentProcessResponse) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"SourceFile", "FieldName", "Value", "Validity", "Status"}); err != nil {
		return err
	}

	// Write data rows
	for sourceFile, result := range results {
		for _, container := range result.ContainerList.List {
			if container.Text != nil {
				for _, field := range container.Text.FieldList {
					row := []string{
						sourceFile,
						field.FieldName,
						field.Value,
						strconv.Itoa(field.Validity),
						strconv.Itoa(field.Status),
					}
					if err := writer.Write(row); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

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

	// Collect results for export
	results := make(map[string]*regula.DocumentProcessResponse)

	for _, file := range files {
		fmt.Printf("Reading file: %v \n", file.Name())
		filePath := filepath.Join(imagesDirectory, file.Name())
		imageData, err := imageReader.Read(filePath)
		if err != nil {
			fmt.Printf("❌ Unable to read Image File: %v \n", err)
			return
		}

		imageProcessResult, err := regulaClient.ProcessRegulaDocumentRead(imageData)
		if err != nil {
			fmt.Printf("❌ Unable to process Image File with Regula: %v \n", err)
			return
		}

		fmt.Printf("✅ File %v read successfuly: %v \n", file.Name(), imageProcessResult)
		results[file.Name()] = imageProcessResult
	}

	// Export results to JSON
	jsonPath := filepath.Join(resultsDirectory, "results.json")
	if err := saveResultsAsJSON(jsonPath, results); err != nil {
		fmt.Printf("❌ Unable to save JSON results: %v \n", err)
		return
	}
	fmt.Printf("✅ Results saved to %s\n", jsonPath)

	// Export results to CSV
	csvPath := filepath.Join(resultsDirectory, "results.csv")
	if err := saveResultsAsCSV(csvPath, results); err != nil {
		fmt.Printf("❌ Unable to save CSV results: %v \n", err)
		return
	}
	fmt.Printf("✅ Results saved to %s\n", csvPath)

	fmt.Println("✅ Test Successful!")
}