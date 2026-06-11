package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSONRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"ankur","extra":"x"}`))
	rec := httptest.NewRecorder()

	var payload struct {
		Name string `json:"name"`
	}

	err := readJSON(rec, req, &payload)
	if err == nil {
		t.Fatal("expected readJSON to reject unknown fields")
	}
}

func TestWriteJSONSetsContentTypeAndStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	err := writeJSON(rec, http.StatusCreated, map[string]string{"status": "ok"})
	if err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}
	if got := res.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json, got %q", got)
	}
}
