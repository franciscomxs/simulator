package middlewares_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
)

type decodePayload struct {
	Value string `json:"value"`
}

func TestDecodeJSON_ValidBodyReachesNext(t *testing.T) {
	var got decodePayload
	var ok bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok = middlewares.Decoded[decodePayload](r)
	})

	h := middlewares.DecodeJSON[decodePayload]()(next)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"value":"hello"}`))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if !ok {
		t.Fatal("expected Decoded ok=true")
	}
	if got.Value != "hello" {
		t.Errorf("expected Value=hello, got %q", got.Value)
	}
}

func TestDecodeJSON_OversizedBody(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	h := middlewares.DecodeJSON[decodePayload]()(next)

	body := `{"value":"` + strings.Repeat("x", 5000) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rr := httptest.NewRecorder()
	req.Body = http.MaxBytesReader(rr, req.Body, 1024)

	h.ServeHTTP(rr, req)

	if called {
		t.Fatal("expected next NOT to be called")
	}
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", rr.Code)
	}
	var resp map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if resp["error"] != "request body too large" {
		t.Errorf("unexpected error body: %v", resp)
	}
}

func TestDecodeJSON_MalformedBody(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	h := middlewares.DecodeJSON[decodePayload]()(next)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{not json`))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if called {
		t.Fatal("expected next NOT to be called")
	}
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
	var resp map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	if resp["error"] != "invalid request body" {
		t.Errorf("unexpected error body: %v", resp)
	}
}

func TestDecoded_NotInstalled(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	_, ok := middlewares.Decoded[decodePayload](req)
	if ok {
		t.Error("expected ok=false when middleware not in chain")
	}
}
