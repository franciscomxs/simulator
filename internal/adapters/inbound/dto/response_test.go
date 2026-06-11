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

func TestSimulateResponse_GracePeriodFieldsMarshal(t *testing.T) {
	resp := dto.SimulateResponse{
		Amount:         10000,
		Rate:           0.02,
		Term:           12,
		System:         "PRICE",
		GracePeriod:    3,
		AdjustedAmount: 10612.08,
		TotalDuration:  15,
		TotalAmount:    12041.99,
		Installments: []dto.InstallmentResponse{
			{Number: 1, Type: "GRACE", Payment: 0, Principal: 0, Interest: 200.00, Balance: 10200.00},
			{Number: 4, Type: "PAYMENT", Payment: 1003.47, Principal: 791.23, Interest: 212.24, Balance: 9820.85},
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
	for _, key := range []string{"grace_period", "adjusted_amount", "total_duration"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON key %q", key)
		}
	}

	insts, ok := m["installments"].([]any)
	if !ok {
		t.Fatalf("installments not an array: %T", m["installments"])
	}
	for i, raw := range insts {
		inst, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("installment %d not an object", i)
		}
		for _, key := range []string{"type", "balance"} {
			if _, ok := inst[key]; !ok {
				t.Errorf("installment %d missing JSON key %q", i, key)
			}
		}
	}
}

func TestSimulateResponse_GracePeriodZero_SafeDefaults(t *testing.T) {
	resp := dto.SimulateResponse{
		Amount:         10000,
		Rate:           0.02,
		Term:           12,
		System:         "PRICE",
		GracePeriod:    0,
		AdjustedAmount: 10000,
		TotalDuration:  12,
		TotalAmount:    11347.20,
		Installments: []dto.InstallmentResponse{
			{Number: 1, Type: "PAYMENT", Payment: 945.60, Principal: 745.60, Interest: 200.00, Balance: 9254.40},
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
	if got := m["grace_period"]; got != float64(0) {
		t.Errorf("grace_period: got %v want 0", got)
	}
	if got := m["adjusted_amount"]; got != float64(10000) {
		t.Errorf("adjusted_amount: got %v want 10000", got)
	}
	if got := m["total_duration"]; got != float64(12) {
		t.Errorf("total_duration: got %v want 12", got)
	}
	insts := m["installments"].([]any)
	first := insts[0].(map[string]any)
	if first["type"] != "PAYMENT" {
		t.Errorf("installments[0].type: got %v want \"PAYMENT\"", first["type"])
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
