package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSaver struct {
	Pool *pgxpool.Pool
}

func (p *PostgresSaver) Close() {
	p.Pool.Close()
}
func NewPostgresSaver() (*PostgresSaver, error) {
	ctx := context.Background()
	connString := "postgres://postgres:secretpassword@localhost:5432/metrics_db?sslmode=disable"
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(ctx, connString)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("[DB] Успешно подключено к PostgreSQL!")
				query := `
				CREATE TABLE IF NOT EXISTS rates (
				id SERIAL PRIMARY KEY ,
				code TEXT NOT NULL ,
				rate NUMERIC NOT NULL , 
				created_at TIMESTAMP DEFAULT NOW()
				);
				`

				_, err = pool.Exec(ctx, query)
				if err == nil {
					return &PostgresSaver{Pool: pool}, nil
				}
				return &PostgresSaver{}, err
			}
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("не удалось подключиться к базе данных после 5 попыток: %v", err)
}
