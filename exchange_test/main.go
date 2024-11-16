package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/backtrac3r/go-cryptomus"
)

func main() {
	// Создаём HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Инициализируем клиента Cryptomus без необходимости авторизации
	// Передаём пустые строки для merchantID и ключей, так как они не требуются
	apiClient := cryptomus.New(client, "", "", "")

	// (Опционально) Переопределяем baseURL, если необходимо
	// Например, для тестирования или использования другого окружения
	// apiClient.SetBaseURL("https://api.cryptomus.com/v1")

	// Указываем валюту, для которой хотим получить обменные курсы
	currency := "USDT"

	// Вызываем метод ListExchangeRates
	rates, err := apiClient.ListExchangeRates(currency)
	if err != nil {
		log.Fatalf("Error fetching exchange rates: %v", err)
	}

	// Выводим полученные обменные курсы
	fmt.Printf("Exchange Rates for %s:\n", currency)
	for _, rate := range rates {
		fmt.Printf("1 %s = %s %s\n", rate.From, rate.Course, rate.To)
	}
}
