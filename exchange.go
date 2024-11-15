package cryptomus

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Endpoint constants
const (
	exchangeRateListEndpoint = "/v1/exchange-rate/%s/list"
)

// ExchangeRate представляет структуру обменного курса.
type ExchangeRate struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Course string `json:"course"`
}

// exchangeRateListRawResponse представляет структуру ответа API для списка обменных курсов.
type exchangeRateListRawResponse struct {
	State  int8           `json:"state"`
	Result []ExchangeRate `json:"result"`
}

// ListExchangeRates запрашивает список обменных курсов для указанной валюты.
// Параметр currency является обязательным и должен содержать код валюты (например, "ETH").
func (c *Cryptomus) ListExchangeRates(currency string) ([]ExchangeRate, error) {
	if currency == "" {
		return nil, errors.New("currency parameter is required")
	}

	// Формируем эндпоинт с указанной валютой
	endpoint := fmt.Sprintf(exchangeRateListEndpoint, currency)

	// Отправляем GET-запрос без тела запроса
	res, err := c.fetch("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Декодируем ответ
	response := &exchangeRateListRawResponse{}
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, err
	}

	// Проверяем статус ответа
	if response.State != 0 {
		return nil, fmt.Errorf("api error: state %d", response.State)
	}

	return response.Result, nil
}
