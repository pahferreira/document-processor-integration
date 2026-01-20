package main

import (
	"fmt"
	"net/http"
	"time"
)

const regulaBaseURL = "http://localhost:8080"

func main() {
	fmt.Println("🔍 Testing Regula API Connection...")

	err := testGetHomepage()

	if err != nil {
		fmt.Printf("❌ Test failed: %v", err)
		return
	}

	fmt.Println("✅ Test Successful!")

}

func testGetHomepage() error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(regulaBaseURL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}