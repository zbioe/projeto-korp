package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRoutes(t *testing.T) {
	t.Parallel()
	const (
		textType = "text/plain; charset=utf-8"
		jsonType = "application/json"
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
