package cryptomus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Endpoint constants
const (
	exchangeRateListEndpoint = "exchange-rate/%s/list"
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
	// Проверка обязательного параметра currency
	currency = strings.TrimSpace(currency)
	if currency == "" {
		return nil, errors.New("currency parameter is required")
	}

	// Формируем эндпоинт с указанной валютой
	endpoint := fmt.Sprintf(exchangeRateListEndpoint, currency)

	// Формируем полный URL, корректно объединяя baseURL и endpoint
	fullURL, err := joinURL(c.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL or endpoint: %w", err)
	}

	// Логируем сформированный URL для диагностики
	fmt.Printf("Requesting URL: %s\n", fullURL)

	// Создаём новый HTTP GET-запрос без тела
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Устанавливаем необходимые заголовки
	req.Header.Set("Accept", "application/json") // Опционально, если API требует

	// Отправляем запрос через существующий HTTP-клиент
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	// Проверяем статус-код ответа
	if res.StatusCode != http.StatusOK {
		// Попытка декодировать сообщение об ошибке из тела ответа
		var errResp struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(res.Body).Decode(&errResp) // Игнорируем ошибку декодирования
		if errResp.Message != "" {
			return nil, fmt.Errorf("unexpected status code: %d, message: %s", res.StatusCode, errResp.Message)
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Декодируем JSON-ответ
	response := &exchangeRateListRawResponse{}
	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Проверяем статус ответа от API
	if response.State != 0 {
		return nil, fmt.Errorf("API error: state %d", response.State)
	}

	// Проверяем, что список обменных курсов не пустой
	if len(response.Result) == 0 {
		return nil, errors.New("exchange rate list is empty")
	}

	return response.Result, nil
}
