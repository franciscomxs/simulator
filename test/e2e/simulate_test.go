package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/infrastructure/container"
)


func newTestServer() *httptest.Server {
	loanUC := container.NewSimulateLoanUseCase()
	investUC := container.NewSimulateInvestmentUseCase()
	loanHandler := httphandler.NewLoanHandler(loanUC)
	investHandler := httphandler.NewInvestmentHandler(investUC)
	router := httphandler.NewRouter(loanHandler, investHandler)
	return httptest.NewServer(router)
}

func TestE2E_SimulatePRICE(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(simResp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(simResp.Installments))
	}

	sumPrincipal := 0.0
	sumPayment := 0.0
	for _, inst := range simResp.Installments {
		sumPrincipal += inst.Principal
		sumPayment += inst.Payment
	}

	if diff := sumPrincipal - 10000; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != 10000", sumPrincipal)
	}

	if diff := sumPayment - simResp.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != total_amount %.2f", sumPayment, simResp.TotalAmount)
	}
}

func TestE2E_SimulateSAC(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(simResp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(simResp.Installments))
	}

	sumPrincipal := 0.0
	sumPayment := 0.0
	for _, inst := range simResp.Installments {
		sumPrincipal += inst.Principal
		sumPayment += inst.Payment
	}

	if diff := sumPrincipal - 10000; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != 10000", sumPrincipal)
	}

	if diff := sumPayment - simResp.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != total_amount %.2f", sumPayment, simResp.TotalAmount)
	}
}

func TestE2E_ValidationError(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var errResp dto.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestE2E_SimulateInvestment(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"initial_amount": 0, "monthly_contribution": 1000, "rate": 0.01, "term": 10}`
	resp, err := http.Post(srv.URL+"/investment/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateInvestmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(simResp.Timeline) != 10 {
		t.Errorf("expected 10 timeline entries, got %d", len(simResp.Timeline))
	}

	last := simResp.Timeline[len(simResp.Timeline)-1]
	if last.Balance != simResp.FinalAmount {
		t.Errorf("last timeline balance %.2f != final_amount %.2f", last.Balance, simResp.FinalAmount)
	}

	for i := 1; i < len(simResp.Timeline); i++ {
		if simResp.Timeline[i].Contributions <= simResp.Timeline[i-1].Contributions {
			t.Errorf("year %d contributions %.2f not greater than year %d contributions %.2f",
				simResp.Timeline[i].Year, simResp.Timeline[i].Contributions,
				simResp.Timeline[i-1].Year, simResp.Timeline[i-1].Contributions)
		}
	}

	for _, snap := range simResp.Timeline {
		if snap.Balance < snap.Contributions {
			t.Errorf("year %d: balance %.2f < contributions %.2f", snap.Year, snap.Balance, snap.Contributions)
		}
	}
}
