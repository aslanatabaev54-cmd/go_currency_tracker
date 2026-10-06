package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
)

type CurrencyResponce struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func FetchRates() (*CurrencyResponce, error) {
	file, err := os.ReadFile("rates_mock.json")
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать локальный json-файл: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(file)
	}))

	url := ts.URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка сети при запросе: %v", err)

	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сайт вернул плохой статус: %d", resp.StatusCode)
	}
	var result CurrencyResponce
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("не удалось расшифровать JSON: %v", err)
	}
	return &result, nil
}
