package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	externalUrl = "https://redsky-uat.perf.target.com/redsky_aggregations/v1/redsky/case_study_v1?key=%s&tcin=%s"
	apiKeyEnv   = "EXTERNAL_API_KEY"
)

var ErrProductNotFound = errors.New("product data not found")

type apiResponse struct {
	Data apiData `json:"data"`
}
type apiData struct {
	Product productData `json:"product"`
}

type productDescription struct {
	Title                 string `json:"title"`
	DownstreamDescription string `json:"downstream_description"`
}
type images struct {
	PrimaryImageURL string `json:"primary_image_url"`
}
type enrichment struct {
	Images images `json:"images"`
}
type productClassification struct {
	ProductTypeName     string `json:"product_type_name"`
	MerchandiseTypeName string `json:"merchandise_type_name"`
}
type primaryBrand struct {
	Name string `json:"name"`
}
type item struct {
	ProductDescription    productDescription    `json:"product_description"`
	Enrichment            enrichment            `json:"enrichment"`
	ProductClassification productClassification `json:"product_classification"`
	PrimaryBrand          primaryBrand          `json:"primary_brand"`
}
type productData struct {
	Tcin string `json:"tcin"`
	Item item   `json:"item"`
}

type apiErrorResponse struct {
	// Only returned from 401?
	Message string `json:"message"`

	// Always returned
	Errors []apiErrors `json:"errors"`
}
type apiErrors struct {
	// Only returned from 401?
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Reason   string `json:"reason"`

	// Only returned from 404?
	Message string `json:"message"`
}

func getProductDataFromExternal(id string) (*productData, error) {
	apiKey, err := getApiKey()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf(externalUrl, apiKey, id))
	if err != nil {
		return nil, fmt.Errorf("error calling Product API: %w", err)
	}

	switch {
	case resp.StatusCode == 404:
		return nil, fmt.Errorf("%w for id: '%s'", ErrProductNotFound, id)
	case resp.StatusCode > 299:
		return nil, parseApiError(resp)
	}

	return parseProductJSON(resp)
}

func parseProductJSON(resp *http.Response) (*productData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading Product API response body: %w", err)
	}
	apiResponse := apiResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Product API JSON: %w", err)
	}

	return &apiResponse.Data.Product, nil
}

func parseApiError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Product API error response: %w", err)
	}

	errorResponse := apiErrorResponse{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		return fmt.Errorf("error unmarshaling Product API error response: %w", err)
	}

	baseError := fmt.Errorf("received '%d' error code from Product API", resp.StatusCode)
	for _, apiError := range errorResponse.Errors {
		if apiError.Reason != "" {
			baseError = fmt.Errorf("%w - %s", baseError, apiError.Reason)
		} else if apiError.Message != "" {
			baseError = fmt.Errorf("%w - %s", baseError, apiError.Message)
		}
	}

	return baseError
}

func getApiKey() (string, error) {
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		return "", fmt.Errorf("no API key set in env variable '%s'", apiKeyEnv)
	}

	return apiKey, nil
}
