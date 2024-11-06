package cryptomus

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	createRecurrenceEndpoint = "/recurrence/create"
	recurrenceInfoEndpoint   = "/recurrence/info"
	recurrenceListEndpoint   = "/recurrence/list"
	recurrenceCancelEndpoint = "/recurrence/cancel"
)

type RecurrenceRequest struct {
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Name           string `json:"name"`
	Period         string `json:"period"`
	ToCurrency     string `json:"to_currency,omitempty"`
	OrderId        string `json:"order_id,omitempty"`
	UrlCallback    string `json:"url_callback,omitempty"`
	DiscountDays   *int   `json:"discount_days,omitempty"`
	DiscountAmount string `json:"discount_amount,omitempty"`
	AdditionalData string `json:"additional_data,omitempty"`
}

type Recurrence struct {
	UUID           string     `json:"uuid"`
	Name           string     `json:"name"`
	OrderId        string     `json:"order_id"`
	Amount         string     `json:"amount"`
	Currency       string     `json:"currency"`
	PayerCurrency  string     `json:"payer_currency"`
	PayerAmountUSD string     `json:"payer_amount_usd"`
	PayerAmount    string     `json:"payer_amount"`
	UrlCallback    string     `json:"url_callback"`
	Period         string     `json:"period"`
	Status         string     `json:"status"`
	Url            string     `json:"url"`
	LastPayOff     *time.Time `json:"last_pay_off,omitempty"`
	DiscountDays   *int       `json:"discount_days,omitempty"`
	DiscountAmount string     `json:"discount_amount,omitempty"`
	EndOfDiscount  *time.Time `json:"end_of_discount,omitempty"`
	AdditionalData string     `json:"additional_data,omitempty"`
}

type recurrenceRawResponse struct {
	State  int8        `json:"state"`
	Result *Recurrence `json:"result"`
}

type RecurrenceInfoRequest struct {
	UUID    string `json:"uuid,omitempty"`
	OrderId string `json:"order_id,omitempty"`
}

type recurrenceInfoRawResponse struct {
	State  int8                `json:"state"`
	Result *Recurrence         `json:"result,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

type RecurrenceListResponse struct {
	Items    []*Recurrence       `json:"items"`
	Paginate *RecurrencePaginate `json:"paginate"`
}

type RecurrencePaginate struct {
	Count          int    `json:"count"`
	HasPages       bool   `json:"hasPages"`
	NextCursor     string `json:"nextCursor,omitempty"`
	PreviousCursor string `json:"previousCursor,omitempty"`
	PerPage        int    `json:"perPage"`
}

type recurrenceListRawResponse struct {
	State  int                     `json:"state"`
	Result *RecurrenceListResponse `json:"result"`
}

type RecurrenceCancelRequest struct {
	UUID    string `json:"uuid,omitempty"`
	OrderId string `json:"order_id,omitempty"`
}

type recurrenceCancelRawResponse struct {
	State  int8                `json:"state"`
	Result *Recurrence         `json:"result,omitempty"`
	Errors map[string][]string `json:"errors,omitempty"`
}

func (c *Cryptomus) CreateRecurrence(recReq *RecurrenceRequest) (*Recurrence, error) {
	if recReq == nil {
		return nil, errors.New("recurrence request cannot be nil")
	}

	res, err := c.fetch("POST", createRecurrenceEndpoint, recReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	response := &recurrenceRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.State != 0 {
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

func (c *Cryptomus) GetRecurrenceInfo(infoReq *RecurrenceInfoRequest) (*Recurrence, error) {
	if infoReq == nil {
		return nil, errors.New("recurrence info request cannot be nil")
	}

	if infoReq.UUID == "" && infoReq.OrderId == "" {
		return nil, errors.New("either uuid or order_id must be provided")
	}

	res, err := c.fetch("POST", recurrenceInfoEndpoint, infoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errorResponse recurrenceInfoRawResponse
		if decodeErr := json.NewDecoder(res.Body).Decode(&errorResponse); decodeErr == nil && errorResponse.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", errorResponse.Errors)
		}
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	response := &recurrenceInfoRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.State != 0 {
		if response.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", response.Errors)
		}
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

func (c *Cryptomus) ListRecurrences(cursor string) (*RecurrenceListResponse, error) {
	payload := map[string]interface{}{}
	if cursor != "" {
		payload["cursor"] = cursor
	}

	res, err := c.fetch("POST", recurrenceListEndpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	response := &recurrenceListRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.State != 0 {
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}

func (c *Cryptomus) CancelRecurrence(cancelReq *RecurrenceCancelRequest) (*Recurrence, error) {
	if cancelReq == nil {
		return nil, errors.New("recurrence cancel request cannot be nil")
	}

	if cancelReq.UUID == "" && cancelReq.OrderId == "" {
		return nil, errors.New("either uuid or order_id must be provided")
	}

	res, err := c.fetch("POST", recurrenceCancelEndpoint, cancelReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errorResponse recurrenceCancelRawResponse
		if decodeErr := json.NewDecoder(res.Body).Decode(&errorResponse); decodeErr == nil && errorResponse.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", errorResponse.Errors)
		}
		return nil, fmt.Errorf("unexpected HTTP status: %s", res.Status)
	}

	response := &recurrenceCancelRawResponse{}
	if err = json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.State != 0 {
		if response.Errors != nil {
			return nil, fmt.Errorf("validation errors: %v", response.Errors)
		}
		return nil, fmt.Errorf("API returned non-zero state: %d", response.State)
	}

	if response.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return response.Result, nil
}
