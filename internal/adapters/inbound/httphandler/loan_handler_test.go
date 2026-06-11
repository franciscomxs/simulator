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

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
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
	if resp.GrossValue != 10000 {
		t.Errorf("GrossValue: got %.2f want 10000", resp.GrossValue)
	}
	if resp.IOF != 333.20 {
		t.Errorf("IOF: got %.4f want 333.20", resp.IOF)
	}
	if resp.NetValue != 9666.80 {
		t.Errorf("NetValue: got %.4f want 9666.80", resp.NetValue)
	}
	if resp.FinancedAmount != 10000 {
		t.Errorf("FinancedAmount: got %.2f want 10000", resp.FinancedAmount)
	}
	if resp.CustomerType != "PF" {
		t.Errorf("CustomerType: got %q want PF", resp.CustomerType)
	}
}

func TestHTTPHandler_MissingCustomerType(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestHTTPHandler_UnknownCustomerType(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "XX"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPHandler_SimulateSAC(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC", "customer_type": "PF"}`
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

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID", "customer_type": "PF"}`
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

	body := `{"amount": -100, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
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

func TestHTTPHandler_SimulateWithGracePeriod(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "grace_period": 3, "customer_type": "PF"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp dto.SimulateResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.GracePeriod != 3 {
		t.Errorf("GracePeriod: got %d want 3", resp.GracePeriod)
	}
	if resp.AdjustedAmount != 10612.08 {
		t.Errorf("AdjustedAmount: got %.2f want 10612.08", resp.AdjustedAmount)
	}
	if resp.TotalDuration != 15 {
		t.Errorf("TotalDuration: got %d want 15", resp.TotalDuration)
	}
	if len(resp.Installments) != 15 {
		t.Fatalf("expected 15 installments, got %d", len(resp.Installments))
	}
	for i := 0; i < 3; i++ {
		if resp.Installments[i].Type != "GRACE" {
			t.Errorf("installments[%d].Type: got %q want \"GRACE\"", i, resp.Installments[i].Type)
		}
		if resp.Installments[i].Payment != 0 {
			t.Errorf("installments[%d].Payment: got %.2f want 0", i, resp.Installments[i].Payment)
		}
	}
}

func TestHTTPHandler_SimulateWithoutGracePeriod_BackwardCompat(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
	req := httptest.NewRequest(http.MethodPost, "/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Simulate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp dto.SimulateResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.GracePeriod != 0 {
		t.Errorf("GracePeriod: got %d want 0", resp.GracePeriod)
	}
	if resp.AdjustedAmount != resp.Amount {
		t.Errorf("AdjustedAmount: got %.2f want %.2f", resp.AdjustedAmount, resp.Amount)
	}
	if resp.TotalDuration != resp.Term {
		t.Errorf("TotalDuration: got %d want %d", resp.TotalDuration, resp.Term)
	}
	for i, inst := range resp.Installments {
		if inst.Type != "PAYMENT" {
			t.Errorf("installments[%d].Type: got %q want \"PAYMENT\"", i, inst.Type)
		}
	}
}

func TestHTTPHandler_NegativeGracePeriod(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "grace_period": -1, "customer_type": "PF"}`
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

func TestHTTPHandler_ErrorMessageIsGeneric(t *testing.T) {
	h := newTestLoanHandler()

	body := `{"amount": -100, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
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
