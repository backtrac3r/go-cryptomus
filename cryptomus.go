// cryptomus.go
package cryptomus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BaseURL is the default API endpoint for Cryptomus.
// It can be overridden by providing a different baseURL when initializing the client.
const BaseURL = "https://api.cryptomus.com/v1"

// Cryptomus represents the Cryptomus API client.
type Cryptomus struct {
	baseURL       string       // Base URL for the API endpoints
	merchantID    string       // Merchant identifier
	paymentApiKey string       // API key for payment operations
	payoutApiKey  string       // API key for payout operations
	client        *http.Client // HTTP client used to make requests
}

// NewCryptomus creates a new Cryptomus API client.
// Parameters:
// - client: An instance of http.Client. If nil, http.DefaultClient is used.
// - merchantID: Your merchant identifier.
// - paymentApiKey: Your API key for payment-related operations.
// - payoutApiKey: Your API key for payout-related operations.
func New(client *http.Client, merchantID, paymentApiKey, payoutApiKey string) *Cryptomus {
	if client == nil {
		client = http.DefaultClient
	}

	return &Cryptomus{
		baseURL:       BaseURL,
		merchantID:    merchantID,
		paymentApiKey: paymentApiKey,
		payoutApiKey:  payoutApiKey,
		client:        client,
	}
}

// SetBaseURL allows overriding the default BaseURL.
// This can be useful for testing or if the API endpoint changes.
func (c *Cryptomus) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// fetch performs an HTTP request to the specified endpoint with the given method and payload.
// It sets the necessary headers, including merchant ID and signature.
// Parameters:
// - method: HTTP method (e.g., "POST").
// - endpoint: API endpoint (e.g., "/recurrence/create").
// - payload: Request payload to be sent as JSON.
// Returns:
// - *http.Response: The HTTP response from the API.
// - error: Error if the request failed.
func (c *Cryptomus) fetch(method, endpoint string, payload interface{}) (*http.Response, error) {
	// Marshal the payload into JSON.
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Generate the signature using the payment API key.
	// Assumes that the signRequest method is implemented in sign.go.
	sign, err := c.signRequest(c.paymentApiKey, body)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %w", err)
	}

	// Create a new HTTP request.
	req, err := http.NewRequest(method, c.baseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the required headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("merchant", c.merchantID)
	req.Header.Set("sign", sign)

	// Execute the HTTP request.
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return res, nil
}
