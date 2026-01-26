package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	documentprocessor "github.com/openinfer/regula-integration/internal/document-processor"
)

const imagesDirectory = "../document-images"
const resultsDirectory = "../results"
type Integrator string
const (
	Regula    Integrator = "regula"
	Recognito Integrator = "recognito"
)

func (i *Integrator) String() string {
	return string(*i)
}

func (i *Integrator) Set(value string) error {
	switch Integrator(value) {
		case Regula, Recognito:
			*i = Integrator(value)
			return nil
		default:
			return fmt.Errorf("must be one of: %s, %s", Regula, Recognito)
	}
}

func main() {
	var integrator Integrator
	flag.Var(&integrator, "integrator", "Document Processor Tool to integrate")
	flag.Parse()

	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)

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


	switch integrator {
	case "regula":
		runRegulaTest(files, imagesDirectory)
	case "recognito":
		runRecognitoTest(files, imagesDirectory)
	}

	fmt.Println("✅ Test Successful!")
}

func runRegulaTest(files []os.DirEntry, imagesDirectory string) {
	fmt.Println("🔍 Testing Regula API Connection...")
	documentProcessor := documentprocessor.NewRegulaDocumentProcessor()
	
	results, err := documentProcessor.ProcessDocuments(files, imagesDirectory)
	if err != nil {
		fmt.Printf("❌ Unable to process documents: %v \n", err)
		return
	}

	jsonPath := filepath.Join(resultsDirectory, "results-regula.json")
	if err = documentProcessor.SaveResultsAsJSON(jsonPath, results); err != nil {
		fmt.Printf("❌ Unable to save JSON results: %v \n", err)
		return
	}

	fmt.Printf("✅ Results saved to %s\n", jsonPath)
}

func runRecognitoTest(files []os.DirEntry, imagesDirectory string) {
	fmt.Println("🔍 Testing Recognito API Connection...")
	documentprocessor := documentprocessor.NewRecognitoDocumentProcessor()

	results, err := documentprocessor.ProcessDocuments(files, imagesDirectory)
	if err != nil {
		fmt.Printf("❌ Unable to process documents: %v \n", err)
		return
	}

	jsonPath := filepath.Join(resultsDirectory, "results-recognito.json")
	if err = documentprocessor.SaveResultsAsJSON(jsonPath, results); err != nil {
		fmt.Printf("❌ Unable to save JSON results: %v \n", err)
		return
	}

	fmt.Printf("✅ Results saved to %s\n", jsonPath)
}