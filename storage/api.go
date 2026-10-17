package storage

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponce struct {
	Status string             `json:"status"`
	Rates  map[string]float64 `json:"rates"`
}

func (p *PostgresSaver) GetRatesHandler(w http.ResponseWriter, r *http.Request) {
	ApiResp := ApiResponce{
		Status: "succeed",
		Rates:  make(map[string]float64),
	}
	ctx := context.Background()
	query := `
	 SELECT code , rate 
	 FROM rates 
	 ORDER BY created_at ASC;
	`
	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		http.Error(w, "произошла ошибка с ", http.StatusInternalServerError)
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		var rate float64
		err = rows.Scan(&code, &rate)
		if err != nil {
			log.Fatal(err)
		}

		ApiResp.Rates[code] = rate
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ApiResp)
}
