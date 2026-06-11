package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
)

func TestWriteJSON_SetResponseWritesBody(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewares.SetResponse(r, http.StatusCreated, map[string]string{"hello": "world"})
	})

	h := middlewares.WriteJSON()(next)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
	var got map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if got["hello"] != "world" {
		t.Errorf("unexpected body: %v", got)
	}
}

func TestWriteJSON_NoSet_NoWrite(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handler sets nothing
	})

	h := middlewares.WriteJSON()(next)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	// httptest.ResponseRecorder defaults to 200 even when nothing written; check body is empty
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", rr.Body.String())
	}
}

func TestWriteJSON_ErrorSetSkipsBody(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewares.SetResponse(r, http.StatusOK, map[string]string{"k": "v"})
		middlewares.SetError(r, errors.New("boom"), "oops")
	})

	h := middlewares.WriteJSON()(next)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Body.Len() != 0 {
		t.Errorf("expected no body write when err set, got %q", rr.Body.String())
	}
}
