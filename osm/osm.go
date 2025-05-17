package osm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var client *Client
var params url.Values

// NewClient creates a new LocationIQ client
func Init(accessToken string) {
	params := url.Values{}
	params.Set("key", accessToken)
	params.Set("format", "json")
	params.Set("limit", "1")
	client = &Client{
		BaseURL:    "https://us1.locationiq.com/v1/search.php",
		HTTPClient: &http.Client{},
	}
}

func SearchText(query string) (*LocationIQResponse, error) {
	params.Set("q", query)
	reqURL := client.BaseURL + "?" + params.Encode()
	resp, err := client.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	var results []LocationIQResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &results[0], nil
}
