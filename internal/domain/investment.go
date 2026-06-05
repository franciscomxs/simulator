package domain

import "math"

const maxInvestmentTerm = 50 // years

// InvestmentParams holds the input parameters for an investment simulation.
type InvestmentParams struct {
	InitialAmount       float64
	MonthlyContribution float64
	Rate                float64
	Term                int // years
}

// YearlySnapshot represents the portfolio state at the end of a year.
type YearlySnapshot struct {
	Year          int
	Contributions float64
	Balance       float64
}

// InvestmentSimulation holds the full result of an investment simulation.
type InvestmentSimulation struct {
	InitialAmount       float64
	MonthlyContribution float64
	Rate                float64
	Term                int
	FinalAmount         float64
	Timeline            []YearlySnapshot
}

// SimulateInvestment validates params and calculates compound growth month by month.
func SimulateInvestment(p InvestmentParams) (InvestmentSimulation, error) {
	if p.InitialAmount < 0 {
		return InvestmentSimulation{}, ErrInvalidInitialAmount
	}
	if p.MonthlyContribution < 0 {
		return InvestmentSimulation{}, ErrInvalidMonthlyContribution
	}
	if p.Rate < 0 || p.Rate >= 1.0 {
		return InvestmentSimulation{}, ErrInvalidRate
	}
	if p.Term <= 0 || p.Term > maxInvestmentTerm {
		return InvestmentSimulation{}, ErrInvalidTerm
	}

	totalMonths := p.Term * 12
	balance := p.InitialAmount
	timeline := make([]YearlySnapshot, 0, p.Term)

	for month := 1; month <= totalMonths; month++ {
		balance = balance*(1+p.Rate) + p.MonthlyContribution
		if math.IsInf(balance, 0) || math.IsNaN(balance) {
			return InvestmentSimulation{}, ErrInvalidSimulationResult
		}
		if month%12 == 0 {
			year := month / 12
			contributions := round2(p.InitialAmount + p.MonthlyContribution*float64(month))
			timeline = append(timeline, YearlySnapshot{
				Year:          year,
				Contributions: contributions,
				Balance:       round2(balance),
			})
		}
	}

	return InvestmentSimulation{
		InitialAmount:       p.InitialAmount,
		MonthlyContribution: p.MonthlyContribution,
		Rate:                p.Rate,
		Term:                p.Term,
		FinalAmount:         round2(balance),
		Timeline:            timeline,
	}, nil
}
