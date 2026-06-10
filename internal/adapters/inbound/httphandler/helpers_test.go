package httphandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franciscomxs/simulator/internal/domain"
)

type helperPayload struct {
	Value string `json:"value"`
}

func TestDecodeJSON_OversizedBody(t *testing.T) {
	w := httptest.NewRecorder()
	body := `{"value":"` + strings.Repeat("x", 5000) + `"}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	_, ok := decodeJSON[helperPayload](w, r)

	if ok {
		t.Fatal("expected ok=false for oversized body")
	}
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", w.Code)
	}
	var resp map[string]string
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "request body too large" {
		t.Errorf("unexpected error body: %v", resp)
	}
}

func TestDecodeJSON_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{bad json`))

	_, ok := decodeJSON[helperPayload](w, r)

	if ok {
		t.Fatal("expected ok=false for invalid JSON")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp map[string]string
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "invalid request body" {
		t.Errorf("unexpected error body: %v", resp)
	}
}

func TestDecodeJSON_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"value":"hello"}`))

	got, ok := decodeJSON[helperPayload](w, r)

	if !ok {
		t.Fatal("expected ok=true for valid JSON")
	}
	if got.Value != "hello" {
		t.Errorf("expected Value=hello, got %q", got.Value)
	}
}

func TestWriteJSON_SetsHeadersAndBody(t *testing.T) {
	w := httptest.NewRecorder()

	writeJSON(w, http.StatusCreated, map[string]string{"k": "v"})

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
	var got map[string]string
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected k=v, got %q", got["k"])
	}
}

func TestHTTPStatusForError_ValidationErrors(t *testing.T) {
	cases := []error{
		domain.ErrInvalidAmount,
		domain.ErrInvalidRate,
		domain.ErrInvalidTerm,
		domain.ErrInvalidSystem,
		domain.ErrInvalidInitialAmount,
		domain.ErrInvalidMonthlyContribution,
	}
	for _, err := range cases {
		if got := httpStatusForError(err); got != http.StatusUnprocessableEntity {
			t.Errorf("httpStatusForError(%v) = %d, want 422", err, got)
		}
	}
}

func TestHTTPStatusForError_InvalidSimulationResult(t *testing.T) {
	if got := httpStatusForError(domain.ErrInvalidSimulationResult); got != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", got)
	}
}

func TestHTTPStatusForError_Unknown(t *testing.T) {
	if got := httpStatusForError(errors.New("boom")); got != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", got)
	}
}
