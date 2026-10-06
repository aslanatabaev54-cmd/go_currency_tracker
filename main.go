package main

import (
	"context"
	"go-currency-tracker/storage"
	"go-currency-tracker/workers"
	"log"
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
	for {
		select {
		case <-ctx.Done():
			log.Println("[Main] Время вышло! Директор закрывает офис. Программа завершена.")
			return
		case data := <-chanel:
			log.Printf("[Main] Поймали из канала: %s = %.2f. Сохраняем в файл...", data.Code, data.Rate)
			err := saver.SaveCurrency(data.Code, data.Rate)
			if err != nil {
				log.Printf("[Main Ошибка] Не удалось записать: %v", err)
			}
		}
	}

}
