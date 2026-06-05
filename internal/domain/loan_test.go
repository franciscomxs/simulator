package domain_test

import (
	"errors"
	"testing"

	"github.com/franciscomxs/simulator/internal/domain"
)

const (
	testAmount = 10000.0
	testRate   = 0.02
	testTerm   = 12
)

func priceParams() domain.LoanParams {
	return domain.LoanParams{
		Amount: testAmount,
		Rate:   testRate,
		Term:   testTerm,
		System: domain.SystemPRICE,
	}
}

func sacParams() domain.LoanParams {
	return domain.LoanParams{
		Amount: testAmount,
		Rate:   testRate,
		Term:   testTerm,
		System: domain.SystemSAC,
	}
}

// --- PRICE tests ---

func TestSimulatePRICE_InstallmentCount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sim.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(sim.Installments))
	}
}

func TestSimulatePRICE_FixedPayment(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first := sim.Installments[0].Payment
	for i, inst := range sim.Installments[:len(sim.Installments)-1] {
		if inst.Payment != first {
			t.Errorf("installment %d payment %.2f != first %.2f", i+1, inst.Payment, first)
		}
	}
}

func TestSimulatePRICE_DecreasingInterest(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(sim.Installments); i++ {
		if sim.Installments[i].Interest >= sim.Installments[i-1].Interest {
			t.Errorf("installment %d interest %.2f >= installment %d interest %.2f",
				i+1, sim.Installments[i].Interest, i, sim.Installments[i-1].Interest)
		}
	}
}

func TestSimulatePRICE_SumOfPrincipalsEqualsAmount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Principal
	}
	if diff := sum - testAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != amount %.2f", sum, testAmount)
	}
}

func TestSimulatePRICE_SumOfPaymentsEqualsTotalAmount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Payment
	}
	if diff := sum - sim.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != TotalAmount %.2f", sum, sim.TotalAmount)
	}
}

func TestSimulatePRICE_InstallmentNumbering(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, inst := range sim.Installments {
		if inst.Number != i+1 {
			t.Errorf("installment index %d has Number %d, expected %d", i, inst.Number, i+1)
		}
	}
}

func TestSimulatePRICE_FirstInterest(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 200.00
	if sim.Installments[0].Interest != expected {
		t.Errorf("first installment interest %.2f != %.2f", sim.Installments[0].Interest, expected)
	}
}

// --- SAC tests ---

func TestSimulateSAC_InstallmentCount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sim.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(sim.Installments))
	}
}

func TestSimulateSAC_ConstantPrincipal(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first := sim.Installments[0].Principal
	for i, inst := range sim.Installments[:len(sim.Installments)-1] {
		if inst.Principal != first {
			t.Errorf("installment %d principal %.2f != first %.2f", i+1, inst.Principal, first)
		}
	}
}

func TestSimulateSAC_DecreasingPayments(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(sim.Installments); i++ {
		if sim.Installments[i].Payment >= sim.Installments[i-1].Payment {
			t.Errorf("installment %d payment %.2f >= installment %d payment %.2f",
				i+1, sim.Installments[i].Payment, i, sim.Installments[i-1].Payment)
		}
	}
}

func TestSimulateSAC_SumOfPrincipalsEqualsAmount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Principal
	}
	if diff := sum - testAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != amount %.2f", sum, testAmount)
	}
}

func TestSimulateSAC_SumOfPaymentsEqualsTotalAmount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Payment
	}
	if diff := sum - sim.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != TotalAmount %.2f", sum, sim.TotalAmount)
	}
}

// --- Validation tests ---

func TestSimulate_InvalidAmount(t *testing.T) {
	p := priceParams()
	p.Amount = 0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestSimulate_NegativeAmount(t *testing.T) {
	p := priceParams()
	p.Amount = -100
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestSimulate_NegativeRate(t *testing.T) {
	p := priceParams()
	p.Rate = -0.01
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulate_InvalidTerm(t *testing.T) {
	p := priceParams()
	p.Term = 0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulate_InvalidSystem(t *testing.T) {
	p := priceParams()
	p.System = "INVALID"
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidSystem) {
		t.Errorf("expected ErrInvalidSystem, got %v", err)
	}
}

func TestSimulate_RateAtOrAboveOne(t *testing.T) {
	p := priceParams()
	p.Rate = 1.0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate for rate=1.0, got %v", err)
	}
}

func TestSimulate_TermExceedsMaximum(t *testing.T) {
	p := priceParams()
	p.Term = 601
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm for term=601, got %v", err)
	}
}

func TestSimulate_ZeroRate(t *testing.T) {
	p := domain.LoanParams{
		Amount: 12000,
		Rate:   0,
		Term:   12,
		System: domain.SystemPRICE,
	}
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 1000.0
	for _, inst := range sim.Installments {
		if inst.Payment != expected {
			t.Errorf("expected payment %.2f, got %.2f", expected, inst.Payment)
		}
	}
}
