package workers

import (
	"context"
	"go-currency-tracker/client"
	"log"
	"time"
)

type CurrencyData struct {
	Code string
	Rate float64
}

func Work(ctx context.Context, ch chan<- CurrencyData) {
	for {
		log.Println("Робот пошел за валютами...")
		worker, err := client.FetchRates()
		if err != nil {
			log.Printf("Не удалось дозвониться: %v", err)
		}
		targetCurrencies := []string{"EUR", "AED", "RUB"}
		for _, code := range targetCurrencies {
			rate, exists := worker.Rates[code]
			if exists {
				log.Printf("[Фон] Найдена валюта %s | Курс: %.2f.", code, rate)
				ch <- CurrencyData{Code: code, Rate: rate}
			}
		}
		select {
		case <-ctx.Done():
			log.Println("[Фон] Поступил сигнал отмены. Воркер завершает работу.")
			return
		case <-time.After(2 * time.Second):

		}
	}
}
