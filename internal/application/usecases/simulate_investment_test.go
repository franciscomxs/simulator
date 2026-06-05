package usecases_test

import (
	"errors"
	"testing"

	"github.com/franciscomxs/simulator/internal/application/ports"
	"github.com/franciscomxs/simulator/internal/application/usecases"
	"github.com/franciscomxs/simulator/internal/domain"
)

func newInvestmentUseCase() ports.SimulateInvestmentUseCase {
	return usecases.NewSimulateInvestment()
}

func TestSimulateInvestmentUseCase_Valid(t *testing.T) {
	uc := newInvestmentUseCase()
	input := ports.SimulateInvestmentInput{
		InitialAmount:       0,
		MonthlyContribution: 1000,
		Rate:                0.01,
		Term:                1,
	}

	out, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Timeline) != 1 {
		t.Errorf("expected 1 timeline entry, got %d", len(out.Timeline))
	}
	if out.Timeline[0].Balance != 12682.50 {
		t.Errorf("year 1 balance = %.2f, want 12682.50", out.Timeline[0].Balance)
	}
	if out.FinalAmount != out.Timeline[0].Balance {
		t.Errorf("FinalAmount %.2f != last timeline balance %.2f", out.FinalAmount, out.Timeline[0].Balance)
	}
}

func TestSimulateInvestmentUseCase_InvalidInitialAmount(t *testing.T) {
	uc := newInvestmentUseCase()
	_, err := uc.Execute(ports.SimulateInvestmentInput{InitialAmount: -1, Term: 1})
	if !errors.Is(err, domain.ErrInvalidInitialAmount) {
		t.Errorf("expected ErrInvalidInitialAmount, got %v", err)
	}
}

func TestSimulateInvestmentUseCase_InvalidMonthlyContribution(t *testing.T) {
	uc := newInvestmentUseCase()
	_, err := uc.Execute(ports.SimulateInvestmentInput{MonthlyContribution: -1, Term: 1})
	if !errors.Is(err, domain.ErrInvalidMonthlyContribution) {
		t.Errorf("expected ErrInvalidMonthlyContribution, got %v", err)
	}
}

func TestSimulateInvestmentUseCase_InvalidRate(t *testing.T) {
	uc := newInvestmentUseCase()
	_, err := uc.Execute(ports.SimulateInvestmentInput{Rate: -0.01, Term: 1})
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulateInvestmentUseCase_InvalidTerm(t *testing.T) {
	uc := newInvestmentUseCase()
	_, err := uc.Execute(ports.SimulateInvestmentInput{Term: 0})
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}
