package dto_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
)

func TestSimulateRequest_JSONRoundTrip(t *testing.T) {
	raw := `{"amount":10000.5,"rate":0.02,"term":12,"system":"PRICE","customer_type":"PF"}`
	var req dto.SimulateRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Amount != 10000.5 {
		t.Errorf("Amount: got %v want 10000.5", req.Amount)
	}
	if req.Rate != 0.02 {
		t.Errorf("Rate: got %v want 0.02", req.Rate)
	}
	if req.Term != 12 {
		t.Errorf("Term: got %v want 12", req.Term)
	}
	if req.System != "PRICE" {
		t.Errorf("System: got %q want PRICE", req.System)
	}
	if req.CustomerType != "PF" {
		t.Errorf("CustomerType: got %q want PF", req.CustomerType)
	}
}

func TestSimulateRequest_Validate_MissingCustomerType(t *testing.T) {
	raw := `{"amount":10000,"rate":0.02,"term":12,"system":"PRICE"}`
	var req dto.SimulateRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := req.Validate(); !errors.Is(err, dto.ErrInvalidCustomerType) {
		t.Errorf("expected ErrInvalidCustomerType, got %v", err)
	}
}

func TestSimulateRequest_Validate_UnknownCustomerType(t *testing.T) {
	req := dto.SimulateRequest{Amount: 10000, Rate: 0.02, Term: 12, System: "PRICE", CustomerType: "XX"}
	if err := req.Validate(); !errors.Is(err, dto.ErrInvalidCustomerType) {
		t.Errorf("expected ErrInvalidCustomerType, got %v", err)
	}
}

func TestSimulateRequest_Validate_AcceptsPFAndPJ(t *testing.T) {
	for _, ct := range []string{"PF", "PJ"} {
		req := dto.SimulateRequest{Amount: 10000, Rate: 0.02, Term: 12, System: "PRICE", CustomerType: ct}
		if err := req.Validate(); err != nil {
			t.Errorf("CustomerType=%q: unexpected error %v", ct, err)
		}
	}
}

func TestSimulateRequest_GracePeriodDefaultsToZero(t *testing.T) {
	raw := `{"amount":10000,"rate":0.02,"term":12,"system":"PRICE"}`
	var req dto.SimulateRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.GracePeriod != 0 {
		t.Errorf("GracePeriod: got %d want 0", req.GracePeriod)
	}
}

func TestSimulateRequest_GracePeriodExplicit(t *testing.T) {
	raw := `{"amount":10000,"rate":0.02,"term":12,"system":"PRICE","grace_period":3}`
	var req dto.SimulateRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.GracePeriod != 3 {
		t.Errorf("GracePeriod: got %d want 3", req.GracePeriod)
	}
}

func TestSimulateRequest_GracePeriodNegativePreserved(t *testing.T) {
	raw := `{"amount":10000,"rate":0.02,"term":12,"system":"PRICE","grace_period":-1}`
	var req dto.SimulateRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.GracePeriod != -1 {
		t.Errorf("GracePeriod: got %d want -1 (DTO must not validate)", req.GracePeriod)
	}
}

func TestSimulateInvestmentRequest_JSONRoundTrip(t *testing.T) {
	raw := `{"initial_amount":5000,"monthly_contribution":200,"rate":0.01,"term":10}`
	var req dto.SimulateInvestmentRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.InitialAmount != 5000 {
		t.Errorf("InitialAmount: got %v want 5000", req.InitialAmount)
	}
	if req.MonthlyContribution != 200 {
		t.Errorf("MonthlyContribution: got %v want 200", req.MonthlyContribution)
	}
	if req.Rate != 0.01 {
		t.Errorf("Rate: got %v want 0.01", req.Rate)
	}
	if req.Term != 10 {
		t.Errorf("Term: got %v want 10", req.Term)
	}
}
