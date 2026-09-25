package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const DefaultNome = "Projeto Korp"
const DefaultAddr = ":8080"

type Response struct {
	Nome    string    `json:"nome"`
	Horario time.Time `json:"horario"`
}

func NewResponse() Response {
	return Response{
		Nome:    DefaultNome,
		Horario: time.Now().UTC(),
	}
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("GET /projeto-korp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(NewResponse()); err != nil {
			log.Printf("encode response: %s", err)
		}
	})
	s := http.Server{
		Handler: m,
		Addr:    DefaultAddr,
	}
	log.Fatal(s.ListenAndServe())
}
