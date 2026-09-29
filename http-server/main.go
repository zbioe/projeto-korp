package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type metrics struct {
	requests *prometheus.CounterVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "número total de requests HTTP",
			},
			[]string{"code", "method"}),
	}
	reg.MustRegister(m.requests)
	return m
}

func handleProjetoKorp(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newResponse(time.Now())); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func newHandler() http.Handler {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)

	app := http.NewServeMux()
	app.HandleFunc("GET /projeto-korp", handleProjetoKorp)

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.Handle("/", promhttp.InstrumentHandlerCounter(m.requests, app))
	return mux
}

func main() {
	s := &http.Server{
		Handler:           newHandler(),
		Addr:              defaultAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
