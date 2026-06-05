package dto_test

import (
	"encoding/json"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
)

func TestSimulateResponse_JSONMarshal(t *testing.T) {
	resp := dto.SimulateResponse{
		Amount:      10000,
		Rate:        0.02,
		Term:        12,
		System:      "PRICE",
		TotalAmount: 11234.56,
		Installments: []dto.InstallmentResponse{
			{Number: 1, Payment: 936.21, Principal: 736.21, Interest: 200.00},
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}

	for _, key := range []string{"amount", "rate", "term", "system", "total_amount", "installments"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON key %q", key)
		}
	}
}

func TestSimulateInvestmentResponse_JSONMarshal(t *testing.T) {
	resp := dto.SimulateInvestmentResponse{
		InitialAmount:       1000,
		MonthlyContribution: 100,
		Rate:                0.01,
		Term:                5,
		FinalAmount:         8234.56,
		Timeline: []dto.YearlySnapshotResponse{
			{Year: 1, Contributions: 2200, Balance: 2310.50},
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}

	for _, key := range []string{"initial_amount", "monthly_contribution", "rate", "term", "final_amount", "timeline"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON key %q", key)
		}
	}
}

func TestErrorResponse_JSONMarshal(t *testing.T) {
	resp := dto.ErrorResponse{Error: "invalid amount"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}
	if m["error"] != "invalid amount" {
		t.Errorf("error field: got %q want %q", m["error"], "invalid amount")
	}
}
