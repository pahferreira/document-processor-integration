package recognito

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	AuthorizationToken string
	BaseURL            string
	HTTPClient         *http.Client
}

const baseURL = "https://privateid-document-api-255329405686.us-central1.run.app"
const authToken = "ZsBbWuhXTrWOZ4hpb7NG"

func NewClient() *Client {
	return &Client{
		AuthorizationToken: authToken,
		BaseURL:            baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type MRZData struct {
	DocumentType   string `json:"document_type"`
	CountryCode    string `json:"country_code"`
	Surname        string `json:"surname"`
	GivenNames     string `json:"given_names"`
	DocumentNumber string `json:"document_number"`
	Nationality    string `json:"nationality"`
	DateOfBirth    string `json:"date_of_birth"`
	Sex            string `json:"sex"`
	ExpiryDate     string `json:"expiry_date"`
}

type OCRData struct {
	FullName       string `json:"full_name"`
	Address        string `json:"address"`
	DateOfIssue    string `json:"date_of_issue"`
	AdditionalData any    `json:"additional_data"`
}

type ConfidenceScores struct {
	MRZ float64 `json:"mrz"`
	OCR float64 `json:"ocr"`
}

type IDCardOCRResponse struct {
	Status       string            `json:"status"`
	DocumentType string            `json:"document_type"`
	MRZData      *MRZData          `json:"mrz_data,omitempty"`
	OCRData      *OCRData          `json:"ocr_data,omitempty"`
	Confidence   *ConfidenceScores `json:"confidence,omitempty"`
	RawResult    any               `json:"raw_result,omitempty"`
	Error        string            `json:"error,omitempty"`
}

type OCRDocumentResponse struct {
	Status       string            `json:"status"`
	DocumentType string            `json:"document_type"`
	MRZData      *MRZData          `json:"mrz_data,omitempty"`
	OCRData      *OCRData          `json:"ocr_data,omitempty"`
	Confidence   *ConfidenceScores `json:"confidence,omitempty"`
	Error        string            `json:"error,omitempty"`
}

type BarcodeData struct {
	Format     string `json:"format"`
	RawData    string `json:"raw_data"`
	ParsedData any    `json:"parsed_data"`
}

type BarcodeOCRResponse struct {
	Status       string       `json:"status"`
	DocumentType string       `json:"document_type"`
	BarcodeData  *BarcodeData `json:"barcode_data,omitempty"`
	Confidence   float64      `json:"confidence,omitempty"`
	RawResult    any          `json:"raw_result,omitempty"`
	Error        string       `json:"error,omitempty"`
}

type LivenessScores struct {
	ScreenReplay        float64 `json:"screen_replay"`
	PrintedCopy         float64 `json:"printed_copy"`
	PortraitReplacement float64 `json:"portrait_replacement"`
}

type LivenessResponse struct {
	Status         string          `json:"status"`
	IsGenuine      bool            `json:"is_genuine,omitempty"`
	LivenessScores *LivenessScores `json:"liveness_scores,omitempty"`
	OverallScore   float64         `json:"overall_score,omitempty"`
	Threshold      float64         `json:"threshold"`
	Error          string          `json:"error,omitempty"`
}

func (c *Client) OCRIDCard(frontImage []byte, backImage []byte) (*OCRDocumentResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add front image (required)
	part, err := writer.CreateFormFile("front_image", "front.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(frontImage); err != nil {
		return nil, fmt.Errorf("failed to write front image: %w", err)
	}

	// Add back image (optional)
	if backImage != nil {
		part, err := writer.CreateFormFile("back_image", "back.jpg")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := part.Write(backImage); err != nil {
			return nil, fmt.Errorf("failed to write back image: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/ocr/id-card", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.AuthorizationToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var jsonResult IDCardOCRResponse
	if err := json.Unmarshal(body, &jsonResult); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := OCRDocumentResponse{
		Status:       jsonResult.Status,
		DocumentType: jsonResult.DocumentType,
		MRZData:      jsonResult.MRZData,
		OCRData:      jsonResult.OCRData,
		Confidence:   jsonResult.Confidence,
		Error:        jsonResult.Error,
	}

	return &result, nil
}

func (c *Client) OCRBarcode(image []byte) (*BarcodeOCRResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("image", "barcode.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(image); err != nil {
		return nil, fmt.Errorf("failed to write image: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/ocr/barcode", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.AuthorizationToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result BarcodeOCRResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (c *Client) LivenessDetect(image []byte) (*LivenessResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("image", "document.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(image); err != nil {
		return nil, fmt.Errorf("failed to write image: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/liveness/detect", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.AuthorizationToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result LivenessResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
