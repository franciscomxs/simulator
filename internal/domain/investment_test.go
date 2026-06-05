package domain_test

import (
	"errors"
	"testing"

	"github.com/franciscomxs/simulator/internal/domain"
)

func investmentParams() domain.InvestmentParams {
	return domain.InvestmentParams{
		InitialAmount:       0,
		MonthlyContribution: 1000,
		Rate:                0.01,
		Term:                10,
	}
}

func TestSimulateInvestment_TimelineLength(t *testing.T) {
	sim, err := domain.SimulateInvestment(investmentParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sim.Timeline) != 10 {
		t.Errorf("expected 10 timeline entries, got %d", len(sim.Timeline))
	}
}

func TestSimulateInvestment_YearNumbering(t *testing.T) {
	sim, err := domain.SimulateInvestment(investmentParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, snap := range sim.Timeline {
		if snap.Year != i+1 {
			t.Errorf("timeline[%d].Year = %d, want %d", i, snap.Year, i+1)
		}
	}
}

func TestSimulateInvestment_ContributionsIncreaseEachYear(t *testing.T) {
	sim, err := domain.SimulateInvestment(investmentParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(sim.Timeline); i++ {
		if sim.Timeline[i].Contributions <= sim.Timeline[i-1].Contributions {
			t.Errorf("year %d contributions %.2f <= year %d contributions %.2f",
				sim.Timeline[i].Year, sim.Timeline[i].Contributions,
				sim.Timeline[i-1].Year, sim.Timeline[i-1].Contributions)
		}
	}
}

func TestSimulateInvestment_FinalTimelineBalanceEqualsFinalAmount(t *testing.T) {
	sim, err := domain.SimulateInvestment(investmentParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	last := sim.Timeline[len(sim.Timeline)-1]
	if last.Balance != sim.FinalAmount {
		t.Errorf("last timeline balance %.2f != FinalAmount %.2f", last.Balance, sim.FinalAmount)
	}
}

func TestSimulateInvestment_BalanceGeContributionsWhenRatePositive(t *testing.T) {
	sim, err := domain.SimulateInvestment(investmentParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, snap := range sim.Timeline {
		if snap.Balance < snap.Contributions {
			t.Errorf("year %d: balance %.2f < contributions %.2f", snap.Year, snap.Balance, snap.Contributions)
		}
	}
}

func TestSimulateInvestment_Year1Values(t *testing.T) {
	p := domain.InvestmentParams{
		InitialAmount:       0,
		MonthlyContribution: 1000,
		Rate:                0.01,
		Term:                1,
	}
	sim, err := domain.SimulateInvestment(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim.Timeline[0].Balance != 12682.50 {
		t.Errorf("year 1 balance = %.2f, want 12682.50", sim.Timeline[0].Balance)
	}
	if sim.Timeline[0].Contributions != 12000.00 {
		t.Errorf("year 1 contributions = %.2f, want 12000.00", sim.Timeline[0].Contributions)
	}
}

func TestSimulateInvestment_InitialAmountOnly(t *testing.T) {
	p := domain.InvestmentParams{
		InitialAmount:       10000,
		MonthlyContribution: 0,
		Rate:                0.01,
		Term:                1,
	}
	sim, err := domain.SimulateInvestment(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim.Timeline[0].Balance != 11268.25 {
		t.Errorf("year 1 balance = %.2f, want 11268.25", sim.Timeline[0].Balance)
	}
	if sim.Timeline[0].Contributions != 10000.00 {
		t.Errorf("year 1 contributions = %.2f, want 10000.00", sim.Timeline[0].Contributions)
	}
}

func TestSimulateInvestment_ZeroRate(t *testing.T) {
	p := domain.InvestmentParams{
		InitialAmount:       1000,
		MonthlyContribution: 100,
		Rate:                0,
		Term:                1,
	}
	sim, err := domain.SimulateInvestment(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim.FinalAmount != 2200.00 {
		t.Errorf("final amount = %.2f, want 2200.00", sim.FinalAmount)
	}
	if sim.Timeline[0].Balance != 2200.00 {
		t.Errorf("year 1 balance = %.2f, want 2200.00", sim.Timeline[0].Balance)
	}
}

func TestSimulateInvestment_NegativeInitialAmount(t *testing.T) {
	p := investmentParams()
	p.InitialAmount = -1
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidInitialAmount) {
		t.Errorf("expected ErrInvalidInitialAmount, got %v", err)
	}
}

func TestSimulateInvestment_NegativeMonthlyContribution(t *testing.T) {
	p := investmentParams()
	p.MonthlyContribution = -1
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidMonthlyContribution) {
		t.Errorf("expected ErrInvalidMonthlyContribution, got %v", err)
	}
}

func TestSimulateInvestment_NegativeRate(t *testing.T) {
	p := investmentParams()
	p.Rate = -0.01
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulateInvestment_ZeroTerm(t *testing.T) {
	p := investmentParams()
	p.Term = 0
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulateInvestment_NegativeTerm(t *testing.T) {
	p := investmentParams()
	p.Term = -1
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulateInvestment_TermExceedsMax(t *testing.T) {
	p := investmentParams()
	p.Term = 51
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulateInvestment_RateGeOne(t *testing.T) {
	p := investmentParams()
	p.Rate = 1.0
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulateInvestment_ExtremeValues_InfBalance(t *testing.T) {
	p := domain.InvestmentParams{
		InitialAmount:       1e300,
		MonthlyContribution: 0,
		Rate:                0.1,
		Term:                50,
	}
	_, err := domain.SimulateInvestment(p)
	if !errors.Is(err, domain.ErrInvalidSimulationResult) {
		t.Errorf("expected ErrInvalidSimulationResult, got %v", err)
	}
}
