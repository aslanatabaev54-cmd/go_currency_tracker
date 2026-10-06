package storage

import (
	"context"
	"fmt"
	"log"
	"os"
)

type FileSaver struct {
	FilePath string
}

func (f FileSaver) SaveCurrency(code string, rate float64) error {
	file, err := os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	row := fmt.Sprintf("Валюта: %s | Курс: %.2f\n", code, rate)
	_, err = file.WriteString(row)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostgresSaver) SaveCurrency(code string, rate float64) error {
	ctx := context.Background()
	sqlQuery := `INSERT INTO rates (code , rate) VALUES($1 , $2)`
	_, err := p.Pool.Exec(ctx, sqlQuery, code, rate)
	return err
}
