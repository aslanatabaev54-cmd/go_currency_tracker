package main

import (
	"context"
	"go-currency-tracker/storage"
	"go-currency-tracker/workers"
	"log"
	"net/http"
	"time"
)

func main() {
	chanel := make(chan workers.CurrencyData)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Программа запущена ...")

	postgresObj, err := storage.NewPostgresSaver()
	if err != nil {
		log.Fatalf("Не удалось запустить Postgres: %v", err)
	}
	defer postgresObj.Close()
	// Надеваем на него универсальную маску интерфейса DataSaver
	var saver storage.DataSaver = postgresObj
	go workers.Work(ctx, chanel)
	go func() {
		for Data := range chanel {
			log.Printf("[Main] Поймали из канала: %s = %.2f. Сохраняем в Postgres...", Data.Code, Data.Rate)
			err := saver.SaveCurrency(Data.Code, Data.Rate)
			if err != nil {
				log.Printf("[Main Ошибка] Не удалось записать: %v", err)
			}
		}
	}()

	http.HandleFunc("/rates", postgresObj.GetRatesHandler)
	log.Println("REST API Сервер успешно запущен на http://localhost:8080/rates")
	http.ListenAndServe(":8080", nil)
}
