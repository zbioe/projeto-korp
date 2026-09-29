package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRoutes(t *testing.T) {
	t.Parallel()
	const (
		textType   = "text/plain; charset=utf-8"
		jsonType   = "application/json"
		metricType = "text/plain; version=0.0.4; charset=utf-8; escaping=underscores"
	)
	tests := map[string]struct {
		method, target string
		wantCode       int
		wantType       string
	}{
		"get projeto-korp":       {http.MethodGet, "/projeto-korp", http.StatusOK, jsonType},
		"head projeto-korp":      {http.MethodHead, "/projeto-korp", http.StatusOK, jsonType},
		"post projeto-korp":      {http.MethodPost, "/projeto-korp", http.StatusMethodNotAllowed, textType},
		"put projeto-korp":       {http.MethodPut, "/projeto-korp", http.StatusMethodNotAllowed, textType},
		"get projeto-korp slash": {http.MethodGet, "/projeto-korp/", http.StatusNotFound, textType},
		"get metrics":            {http.MethodGet, "/metrics", http.StatusOK, metricType},
		"head metrics":           {http.MethodHead, "/metrics", http.StatusOK, metricType},
		"post metrics":           {http.MethodPost, "/metrics", http.StatusMethodNotAllowed, textType},
		"get metrics slash":      {http.MethodGet, "/metrics/", http.StatusNotFound, textType},
		"get healthz":            {http.MethodGet, "/healthz", http.StatusOK, ""},
		"head healthz":           {http.MethodHead, "/healthz", http.StatusOK, ""},
		"post healthz":           {http.MethodPost, "/healthz", http.StatusMethodNotAllowed, textType},
		"get healthz slash":      {http.MethodGet, "/healthz/", http.StatusNotFound, textType},
		"get root path":          {http.MethodGet, "/", http.StatusNotFound, textType},
		"get unknown path":       {http.MethodGet, "/unknownpath", http.StatusNotFound, textType},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			newMux().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.target, nil))

			if rec.Code != tc.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantCode)
			}

			if got := rec.Header().Get("Content-Type"); got != tc.wantType {
				t.Errorf("Content-Type = %q, want %q", got, tc.wantType)
			}
		})
	}
}

func TestBody(t *testing.T) {
	t.Parallel()
	base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		elapsed time.Duration
		want    string
	}{
		"start":                    {0, "2000-01-01T00:00:00Z"},
		"drop fractions of second": {1500 * time.Millisecond, "2000-01-01T00:00:01Z"},
		"next day":                 {24 * time.Hour, "2000-01-02T00:00:00Z"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			want := `{"nome":"Projeto Korp","horario":"` + tc.want + `"}`
			got, err := json.Marshal(newResponse(base.Add(tc.elapsed)))
			if err != nil {
				t.Fatalf("error marshal response: %v", err)
			}
			if string(got) != want {
				t.Errorf("body = %s, want %s", got, want)
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	t.Parallel()
	mux := newMux()
	for _, method := range []string{http.MethodGet, http.MethodGet, http.MethodPost} {
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(method, "/projeto-korp", nil))
	}

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	for _, want := range []string{
		`http_requests_total{code="200",method="get"} 2`,
		`http_requests_total{code="405",method="post"} 1`,
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("/metrics does not contain %q", want)
		}
	}
}
