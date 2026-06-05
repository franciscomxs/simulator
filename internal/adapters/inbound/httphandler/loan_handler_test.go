package httphandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/usecases"
)

func newTestLoanHandler() *httphandler.LoanHandler {
	uc := usecases.NewSimulateLoan()
	return httphandler.NewLoanHandler(uc)
}

func TestHTTPHandler_SimulatePRICE(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp dto.SimulateResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(resp.Installments))
	}
}

func TestHTTPHandler_SimulateSAC(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp dto.SimulateResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(resp.Installments))
	}
}

func TestHTTPHandler_InvalidSystem(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestHTTPHandler_NegativeAmount(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": -100, "rate": 0.02, "term": 12, "system": "PRICE"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPHandler_InvalidJSON(t *testing.T) {
	h := newTestLoanHandler()

	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPHandler_OversizedBody(t *testing.T) {
	h := newTestLoanHandler()

	// Valid JSON prefix forces the decoder to read past the 4 KB limit
	// before hitting MaxBytesReader, so we get *http.MaxBytesError (not json.SyntaxError).
	body := `{"amount":1,"rate":0.01,"term":1,"system":"` + strings.Repeat("x", 2<<20) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status 413, got %d", rr.Code)
	}
}

func TestHTTPHandler_ErrorMessageIsGeneric(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": -100, "rate": 0.02, "term": 12, "system": "PRICE"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error == "invalid amount" {
		t.Error("error response must not leak internal domain error messages")
	}
	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}
