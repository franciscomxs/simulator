package usecases_test

import (
	"errors"
	"testing"

	"github.com/franciscomxs/simulator/internal/application/ports"
	"github.com/franciscomxs/simulator/internal/application/usecases"
	"github.com/franciscomxs/simulator/internal/domain"
)

func TestSimulateLoan_PRICE_Valid(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   0.02,
		Term:   12,
		System: "PRICE",
	}

	out, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(out.Installments))
	}
	if out.System != "PRICE" {
		t.Errorf("expected system PRICE, got %s", out.System)
	}
	if out.Amount != 10000 {
		t.Errorf("expected amount 10000, got %.2f", out.Amount)
	}
	if out.TotalAmount <= 0 {
		t.Errorf("expected positive total amount, got %.2f", out.TotalAmount)
	}
}

func TestSimulateLoan_SAC_Valid(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   0.02,
		Term:   12,
		System: "SAC",
	}

	out, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(out.Installments))
	}
	if out.System != "SAC" {
		t.Errorf("expected system SAC, got %s", out.System)
	}
}

func TestSimulateLoan_InvalidAmount(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 0,
		Rate:   0.02,
		Term:   12,
		System: "PRICE",
	}

	_, err := uc.Execute(input)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestSimulateLoan_InvalidRate(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   -0.01,
		Term:   12,
		System: "PRICE",
	}

	_, err := uc.Execute(input)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulateLoan_InvalidTerm(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   0.02,
		Term:   0,
		System: "PRICE",
	}

	_, err := uc.Execute(input)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulateLoan_InvalidSystem(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   0.02,
		Term:   12,
		System: "INVALID",
	}

	_, err := uc.Execute(input)
	if !errors.Is(err, domain.ErrInvalidSystem) {
		t.Errorf("expected ErrInvalidSystem, got %v", err)
	}
}

func TestSimulateLoan_OutputInstallmentFields(t *testing.T) {
	uc := usecases.NewSimulateLoan()
	input := ports.SimulateLoanInput{
		Amount: 10000,
		Rate:   0.02,
		Term:   12,
		System: "PRICE",
	}

	out, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inst := out.Installments[0]
	if inst.Number != 1 {
		t.Errorf("expected installment Number 1, got %d", inst.Number)
	}
	if inst.Payment <= 0 {
		t.Errorf("expected positive payment, got %.2f", inst.Payment)
	}
	if inst.Principal <= 0 {
		t.Errorf("expected positive principal, got %.2f", inst.Principal)
	}
	if inst.Interest != 200.00 {
		t.Errorf("expected first interest 200.00, got %.2f", inst.Interest)
	}
}
