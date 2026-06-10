package e2e_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler"
	"github.com/franciscomxs/simulator/internal/infrastructure/container"
)

func newTestServer() *httptest.Server {
	loanUC := container.NewSimulateLoanUseCase()
	investUC := container.NewSimulateInvestmentUseCase()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	loanHandler := httphandler.NewLoanHandler(loanUC, logger)
	investHandler := httphandler.NewInvestmentHandler(investUC, logger)
	router := httphandler.NewRouter(loanHandler, investHandler)
	return httptest.NewServer(router)
}

func TestE2E_SimulatePRICE(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
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

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC", "customer_type": "PF"}`
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

func TestE2E_SimulatePRICE_WithGracePeriod(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "grace_period": 3, "customer_type": "PF"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if simResp.GracePeriod != 3 {
		t.Errorf("GracePeriod: got %d want 3", simResp.GracePeriod)
	}
	if simResp.AdjustedAmount != 10612.08 {
		t.Errorf("AdjustedAmount: got %.2f want 10612.08", simResp.AdjustedAmount)
	}
	if simResp.TotalDuration != 15 {
		t.Errorf("TotalDuration: got %d want 15", simResp.TotalDuration)
	}
	if len(simResp.Installments) != 15 {
		t.Fatalf("expected 15 installments, got %d", len(simResp.Installments))
	}
	for i := 0; i < 3; i++ {
		if simResp.Installments[i].Type != "GRACE" {
			t.Errorf("installments[%d].Type: got %q want \"GRACE\"", i, simResp.Installments[i].Type)
		}
		if simResp.Installments[i].Payment != 0 {
			t.Errorf("installments[%d].Payment: got %.2f want 0", i, simResp.Installments[i].Payment)
		}
	}
	for i := 3; i < 15; i++ {
		if simResp.Installments[i].Type != "PAYMENT" {
			t.Errorf("installments[%d].Type: got %q want \"PAYMENT\"", i, simResp.Installments[i].Type)
		}
	}
}

func TestE2E_SimulatePRICE_NoGracePeriod_BackwardCompat(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if simResp.GracePeriod != 0 {
		t.Errorf("GracePeriod: got %d want 0", simResp.GracePeriod)
	}
	if simResp.AdjustedAmount != simResp.Amount {
		t.Errorf("AdjustedAmount: got %.2f want %.2f", simResp.AdjustedAmount, simResp.Amount)
	}
	if simResp.TotalDuration != simResp.Term {
		t.Errorf("TotalDuration: got %d want %d", simResp.TotalDuration, simResp.Term)
	}
	if len(simResp.Installments) != 12 {
		t.Fatalf("expected 12 installments, got %d", len(simResp.Installments))
	}
	for i, inst := range simResp.Installments {
		if inst.Type != "PAYMENT" {
			t.Errorf("installments[%d].Type: got %q want \"PAYMENT\"", i, inst.Type)
		}
	}
}

func TestE2E_ValidationError(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID", "customer_type": "PF"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", resp.StatusCode)
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

func TestE2E_SimulateLoan_IOF_PF(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.NewDecoder(resp.Body).Decode(&simResp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if simResp.GrossValue != 10000 {
		t.Errorf("GrossValue: got %.2f want 10000", simResp.GrossValue)
	}
	if simResp.IOF != 333.20 {
		t.Errorf("IOF: got %.4f want 333.20", simResp.IOF)
	}
	if simResp.NetValue != 9666.80 {
		t.Errorf("NetValue: got %.4f want 9666.80", simResp.NetValue)
	}
	if simResp.FinancedAmount != 10000 {
		t.Errorf("FinancedAmount: got %.2f want 10000", simResp.FinancedAmount)
	}
	if simResp.CustomerType != "PF" {
		t.Errorf("CustomerType: got %q want PF", simResp.CustomerType)
	}
}

func TestE2E_SimulateLoan_IOF_PJ_MatchesPF(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	post := func(ct string) dto.SimulateResponse {
		body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "` + ct + `"}`
		resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("request failed for %s: %v", ct, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status for %s: got %d, want 200", ct, resp.StatusCode)
		}
		var sr dto.SimulateResponse
		if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
			t.Fatalf("decode for %s: %v", ct, err)
		}
		return sr
	}

	pf := post("PF")
	pj := post("PJ")

	if pf.IOF != pj.IOF {
		t.Errorf("IOF PF=%.2f PJ=%.2f, expected equal", pf.IOF, pj.IOF)
	}
	if pf.NetValue != pj.NetValue {
		t.Errorf("NetValue PF=%.2f PJ=%.2f, expected equal", pf.NetValue, pj.NetValue)
	}
	if pf.FinancedAmount != pj.FinancedAmount {
		t.Errorf("FinancedAmount PF=%.2f PJ=%.2f, expected equal", pf.FinancedAmount, pj.FinancedAmount)
	}
}

func TestE2E_SimulateLoan_MissingCustomerType(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestE2E_SimulateLoan_UnknownCustomerType(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "XX"}`
	resp, err := http.Post(srv.URL+"/loan/simulate", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
