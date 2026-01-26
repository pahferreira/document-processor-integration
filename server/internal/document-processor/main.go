package documentprocessor

import (
	"os"
)

type DocumentProcessor interface {
	ProcessDocuments(files []os.DirEntry, imagesDirectory string) (map[string][]any, error)
	SaveResultsAsJSON(filename string, results map[string][]any) error
}

const maxConcurrent = 10 // Limit concurrent API requests
