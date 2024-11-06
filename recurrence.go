package cryptomus

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Endpoint constants for recurring payments
const (
	createRecurrenceEndpoint = "/recurrence/create" // Endpoint to create a new recurring payment
	recurrenceInfoEndpoint   = "/recurrence/info"   // Endpoint to retrieve information about a specific recurring payment
	recurrenceListEndpoint   = "/recurrence/list"   // Endpoint to list all recurring payments
	recurrenceCancelEndpoint = "/recurrence/cancel" // Endpoint to cancel a recurring payment
)

// RecurrenceRequest represents the request structure for creating a recurring payment.
type RecurrenceRequest struct {
	Amount         string `json:"amount"`                    // Required: Amount of the payment
	Currency       string `json:"currency"`                  // Required: Currency code (e.g., "USD")
	Name           string `json:"name"`                      // Required: Name or description of the payment
	Period         string `json:"period"`                    // Required: Recurrence period (e.g., "monthly")
	ToCurrency     string `json:"to_currency,omitempty"`     // Optional: Target currency
	OrderId        string `json:"order_id,omitempty"`        // Optional: Order identifier in your system
	UrlCallback    string `json:"url_callback,omitempty"`    // Optional: Callback URL for payment status updates
	DiscountDays   *int   `json:"discount_days,omitempty"`   // Optional: Number of days for discount eligibility
	DiscountAmount string `json:"discount_amount,omitempty"` // Optional: Amount of discount
	AdditionalData string `json:"additional_data,omitempty"` // Optional: Additional data for the payment
}

// Recurrence represents the response structure for a recurring payment.
type Recurrence struct {
	UUID           string     `json:"uuid"`                      // Unique identifier for the recurring payment
	Name           string     `json:"name"`                      // Name or description of the payment
	OrderId        string     `json:"order_id"`                  // Order identifier in your system
	Amount         string     `json:"amount"`                    // Amount of the payment
	Currency       string     `json:"currency"`                  // Currency code (e.g., "USD")
	PayerCurrency  string     `json:"payer_currency"`            // Currency used by the payer
	PayerAmountUSD string     `json:"payer_amount_usd"`          // Payer amount in USD
	PayerAmount    string     `json:"payer_amount"`              // Amount paid by the payer
	UrlCallback    string     `json:"url_callback"`              // Callback URL for payment status updates
	Period         string     `json:"period"`                    // Recurrence period (e.g., "monthly")
	Status         string     `json:"status"`                    // Current status of the payment
	Url            string     `json:"url"`                       // URL for payment processing
	LastPayOff     *time.Time `json:"last_pay_off,omitempty"`    // Optional: Timestamp of the last payment
	DiscountDays   *int       `json:"discount_days,omitempty"`   // Optional: Number of discount days
	DiscountAmount string     `json:"discount_amount,omitempty"` // Optional: Amount of discount
	EndOfDiscount  *time.Time `json:"end_of_discount,omitempty"` // Optional: Timestamp when the discount ends
	AdditionalData string     `json:"additional_data,omitempty"` // Optional: Additional data for the payment
}

// recurrenceRawResponse represents the raw response structure from the API for recurring payments.
type recurrenceRawResponse struct {
	State  int8        `json:"state"`  // State code indicating success or error
	Result *Recurrence `json:"result"` // Resulting Recurrence object on success
}

// RecurrenceInfoRequest represents the request structure for retrieving information about a recurring payment.
type RecurrenceInfoRequest struct {
	UUID    string `json:"uuid,omitempty"`     // Optional: UUID of the recurring payment
	OrderId string `json:"order_id,omitempty"` // Optional: Order identifier in your system
}

// recurrenceInfoRawResponse represents the raw response structure from the API for retrieving recurring payment information.
type recurrenceInfoRawResponse struct {
	State  int8                `json:"state"`            // State code indicating success or error
	Result *Recurrence         `json:"result,omitempty"` // Resulting Recurrence object on success
	Errors map[string][]string `json:"errors,omitempty"` // Validation errors if any
}

// RecurrenceListResponse represents the response structure for listing recurring payments.
type RecurrenceListResponse struct {
	Items    []*Recurrence       `json:"items"`    // List of recurring payments
	Paginate *RecurrencePaginate `json:"paginate"` // Pagination information
}

// RecurrencePaginate represents the pagination information for listing recurring payments.
type RecurrencePaginate struct {
	Count          int    `json:"count"`                    // Total number of items
	HasPages       bool   `json:"hasPages"`                 // Indicates if there are multiple pages
	NextCursor     string `json:"nextCursor,omitempty"`     // Cursor for the next page
	PreviousCursor string `json:"previousCursor,omitempty"` // Cursor for the previous page
	PerPage        int    `json:"perPage"`                  // Number of items per page
}

// recurrenceListRawResponse represents the raw response structure from the API for listing recurring payments.
type recurrenceListRawResponse struct {
	State  int                     `json:"state"`  // State code indicating success or error
	Result *RecurrenceListResponse `json:"result"` // Resulting RecurrenceListResponse object on success
}

// RecurrenceCancelRequest represents the request structure for canceling a recurring payment.
type RecurrenceCancelRequest struct {
	UUID    string `json:"uuid,omitempty"`     // Optional: UUID of the recurring payment to cancel
	OrderId string `json:"order_id,omitempty"` // Optional: Order identifier in your system
}

// recurrenceCancelRawResponse represents the raw response structure from the API for canceling a recurring payment.
type recurrenceCancelRawResponse struct {
	State  int8                `json:"state"`            // State code indicating success or error
	Result *Recurrence         `json:"result,omitempty"` // Resulting Recurrence object on success
	Errors map[string][]string `json:"errors,omitempty"` // Validation errors if any
}

// CreateRecurrence creates a new recurring payment.
func (c *Cryptomus) CreateRecurrence(recReq *RecurrenceRequest) (*Recurrence, error) {
	if recReq == nil {
		return nil, errors.New("recurrence request cannot be nil")
	}

	// Send a POST request to create a recurring payment
	res, err := c.fetch("POST", createRecurrenceEndpoint, recReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Check for unexpected HTTP status codes
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	// Decode the JSON response
	response := &recurrenceRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check the state of the response
	if response.State != 0 {
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	// Ensure the result is not nil
	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

// GetRecurrenceInfo retrieves information about a specific recurring payment using UUID or OrderId.
func (c *Cryptomus) GetRecurrenceInfo(infoReq *RecurrenceInfoRequest) (*Recurrence, error) {
	if infoReq == nil {
		return nil, errors.New("recurrence info request cannot be nil")
	}

	if infoReq.UUID == "" && infoReq.OrderId == "" {
		return nil, errors.New("either uuid or order_id must be provided")
	}

	// Send a POST request to retrieve recurring payment information
	res, err := c.fetch("POST", recurrenceInfoEndpoint, infoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Handle non-200 HTTP status codes by attempting to decode validation errors
	if res.StatusCode != 200 {
		var errorResponse recurrenceInfoRawResponse
		if decodeErr := json.NewDecoder(res.Body).Decode(&errorResponse); decodeErr == nil && errorResponse.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", errorResponse.Errors)
		}
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	// Decode the JSON response
	response := &recurrenceInfoRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check the state of the response and handle validation errors
	if response.State != 0 {
		if response.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", response.Errors)
		}
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	// Ensure the result is not nil
	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

// ListRecurrences retrieves a list of all recurring payments with optional pagination using a cursor.
func (c *Cryptomus) ListRecurrences(cursor string) (*RecurrenceListResponse, error) {
	payload := make(map[string]interface{})
	if cursor != "" {
		payload["cursor"] = cursor
	}

	// Send a POST request to list recurring payments
	res, err := c.fetch("POST", recurrenceListEndpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Check for unexpected HTTP status codes
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	// Decode the JSON response
	response := &recurrenceListRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check the state of the response
	if response.State != 0 {
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	// Ensure the result is not nil
	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

// CancelRecurrence cancels a recurring payment using UUID or OrderId.
func (c *Cryptomus) CancelRecurrence(cancelReq *RecurrenceCancelRequest) (*Recurrence, error) {
	if cancelReq == nil {
		return nil, errors.New("recurrence cancel request cannot be nil")
	}

	if cancelReq.UUID == "" && cancelReq.OrderId == "" {
		return nil, errors.New("either uuid or order_id must be provided")
	}

	// Send a POST request to cancel the recurring payment
	res, err := c.fetch("POST", recurrenceCancelEndpoint, cancelReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Handle non-200 HTTP status codes by attempting to decode validation errors
	if res.StatusCode != 200 {
		var errorResponse recurrenceCancelRawResponse
		if decodeErr := json.NewDecoder(res.Body).Decode(&errorResponse); decodeErr == nil && errorResponse.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", errorResponse.Errors)
		}
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	// Decode the JSON response
	response := &recurrenceCancelRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check the state of the response and handle validation errors
	if response.State != 0 {
		if response.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", response.Errors)
		}
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	// Ensure the result is not nil
	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}
