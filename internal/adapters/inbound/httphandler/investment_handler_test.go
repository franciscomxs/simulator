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

func newInvestmentHandler() *httphandler.InvestmentHandler {
	uc := usecases.NewSimulateInvestment()
	return httphandler.NewInvestmentHandler(uc)
}

func TestHTTPInvestmentHandler_Valid(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": 0, "monthly_contribution": 1000, "rate": 0.01, "term": 10}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp dto.SimulateInvestmentResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Timeline) != 10 {
		t.Errorf("expected 10 timeline entries, got %d", len(resp.Timeline))
	}
	if resp.FinalAmount <= 0 {
		t.Errorf("expected positive final_amount, got %.2f", resp.FinalAmount)
	}
	last := resp.Timeline[len(resp.Timeline)-1]
	if last.Balance != resp.FinalAmount {
		t.Errorf("last timeline balance %.2f != final_amount %.2f", last.Balance, resp.FinalAmount)
	}
}

func TestHTTPInvestmentHandler_NegativeInitialAmount(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": -100, "monthly_contribution": 1000, "rate": 0.01, "term": 1}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_NegativeMonthlyContribution(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": 0, "monthly_contribution": -1, "rate": 0.01, "term": 1}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_NegativeRate(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": 0, "monthly_contribution": 1000, "rate": -0.01, "term": 1}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_ZeroTerm(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": 0, "monthly_contribution": 1000, "rate": 0.01, "term": 0}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_InvalidJSON(t *testing.T) {
	h := newInvestmentHandler()

	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_OversizedBody(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount":0,"monthly_contribution":1000,"rate":0.01,"term":1,"extra":"` + strings.Repeat("x", 2<<20) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected status 413, got %d", rr.Code)
	}
}

func TestHTTPInvestmentHandler_ErrorIsGeneric(t *testing.T) {
	h := newInvestmentHandler()

	body := `{"initial_amount": -1, "monthly_contribution": 1000, "rate": 0.01, "term": 1}`
	req := httptest.NewRequest(http.MethodPost, "/investment/simulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.SimulateInvestment(rr, req)

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error == "invalid initial_amount" {
		t.Error("error response must not leak internal domain error messages")
	}
	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}
