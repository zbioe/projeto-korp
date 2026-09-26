package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	nome        = "Projeto Korp"
	defaultAddr = ":8080"
)

type response struct {
	Nome    string  `json:"nome"`
	Horario horario `json:"horario"`
}

type horario time.Time

func (h horario) MarshalText() ([]byte, error) {
	return time.Time(h).UTC().AppendFormat(nil, time.RFC3339), nil
}

func newResponse(now time.Time) response {
	return response{
		Nome:    nome,
		Horario: horario(now),
	}
}

func handleProjetoKorp(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newResponse(time.Now())); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func newMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("GET /projeto-korp", handleProjetoKorp)
	return m
}

func main() {
	s := &http.Server{
		Handler:           newMux(),
		Addr:              defaultAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
