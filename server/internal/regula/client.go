package regula

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTPClient *http.Client
}

type DocumentProcessParams struct {
	Scenario string `json:"scenario"`
}

type ImageData struct {
	Image string `json:"image"`
}

type DocumentListItem struct {
	ImageData ImageData
}

type DocumentProcessPayload struct {
	ProcessParam DocumentProcessParams `json:"processParam"`
	List []DocumentListItem `json:"List"`
}

type DocumentProcessResponse struct {
	ChipPage           int            `json:"ChipPage"`
	ProcessingFinished int            `json:"ProcessingFinished"`
	ContainerList      ContainerList  `json:"ContainerList"`
	MorePagesAvailable int            `json:"morePagesAvailable"`
	TransactionInfo    TransactionInfo `json:"TransactionInfo"`
}

type ContainerList struct {
	Count int         `json:"Count"`
	List  []Container `json:"List"`
}

type Container struct {
	ContainerType int         `json:"ContainerType"`
	Text          *TextResult `json:"Text,omitempty"`
	Graphics      *GraphicsResult `json:"Graphics,omitempty"`
	Status        *StatusResult `json:"Status,omitempty"`
}

type TextResult struct {
	FieldList []TextField `json:"fieldList"`
}

type TextField struct {
	FieldType   int    `json:"fieldType"`
	FieldName   string `json:"fieldName"`
	Value       string `json:"value"`
	Validity    int    `json:"validity"`
	Status      int    `json:"status"`
	Comparison  []Comparison `json:"comparison,omitempty"`
}

type Comparison struct {
	SourceType int    `json:"sourceType"`
	Result     int    `json:"result"`
}

type GraphicsResult struct {
	FieldList []GraphicsField `json:"fieldList"`
}

type GraphicsField struct {
	FieldType   int    `json:"fieldType"`
	FieldName   string `json:"fieldName"`
	Value       string `json:"value"` // base64 encoded image
}

type StatusResult struct {
	OverallStatus       int                   `json:"overallStatus"`
	Optical             int                   `json:"optical"`
	DetailsOptical      StatusDetails         `json:"detailsOptical"`
	DetailsRFID         StatusDetails         `json:"detailsRFID,omitempty"`
	DetailsRFIDVisual   StatusDetails         `json:"detailsRFIDVisual,omitempty"`
	AuthenticityChecks  []AuthenticityCheck   `json:"authenticityChecks,omitempty"`
}

type StatusDetails struct {
	OverallStatus int `json:"overallStatus"`
	Reserved      int `json:"reserved"`
}

type AuthenticityCheck struct {
	Type           int    `json:"type"`
	Result         int    `json:"result"`
	CheckName      string `json:"checkName,omitempty"`
	PageIndex      int    `json:"pageIndex,omitempty"`
}

type TransactionInfo struct {
	ComputerName   string `json:"ComputerName"`
	DateTime       string `json:"DateTime"`
	TransactionID  string `json:"TransactionID"`
}

const baseURL = "http://localhost:8080"

func NewClient() *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) ProcessRegulaDocumentRead(image string) ([]TextField, error) {
	payload := DocumentProcessPayload{
		ProcessParam: DocumentProcessParams{
			Scenario: "FullProcess",
		},
		List: []DocumentListItem{
			{
				ImageData: ImageData{
					Image: image,
				},
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL + "/api/process", "application/json", bytes.NewReader(payloadJSON))
	
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DocumentProcessResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return extractTextFields(&result), nil
}

func extractTextFields(response *DocumentProcessResponse) []TextField {
	var fields []TextField
	for _, container := range response.ContainerList.List {
		if container.Text != nil {
			fields = append(fields, container.Text.FieldList...)
		}
	}
	return fields
}