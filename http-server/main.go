package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const nome = "Projeto Korp"
const defaultAddr = ":8080"

type response struct {
	Nome    string  `json:"nome"`
	Horario horario `json:"horario"`
}

type horario time.Time

func (h horario) MarshalText() ([]byte, error) {
	return time.Time(h).UTC().AppendFormat(nil, time.RFC3339), nil
}

func newResponse() response {
	return response{
		Nome:    nome,
		Horario: horario(time.Now()),
	}
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("GET /projeto-korp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(newResponse()); err != nil {
			log.Printf("encode response: %s", err)
		}
	})
	s := http.Server{
		Handler:           m,
		Addr:              defaultAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
