package usecases

import (
	"github.com/franciscomxs/simulator/internal/application/ports"
	"github.com/franciscomxs/simulator/internal/domain"
)

type simulateInvestmentUseCase struct{}

// NewSimulateInvestment creates a new SimulateInvestmentUseCase.
func NewSimulateInvestment() ports.SimulateInvestmentUseCase {
	return &simulateInvestmentUseCase{}
}

// Execute runs the simulate investment use case.
func (uc *simulateInvestmentUseCase) Execute(input ports.SimulateInvestmentInput) (ports.SimulateInvestmentOutput, error) {
	params := domain.InvestmentParams{
		InitialAmount:       input.InitialAmount,
		MonthlyContribution: input.MonthlyContribution,
		Rate:                input.Rate,
		Term:                input.Term,
	}

	sim, err := domain.SimulateInvestment(params)
	if err != nil {
		return ports.SimulateInvestmentOutput{}, err
	}

	timeline := make([]ports.YearlySnapshotOutput, len(sim.Timeline))
	for i, snap := range sim.Timeline {
		timeline[i] = ports.YearlySnapshotOutput{
			Year:          snap.Year,
			Contributions: snap.Contributions,
			Balance:       snap.Balance,
		}
	}

	return ports.SimulateInvestmentOutput{
		InitialAmount:       sim.InitialAmount,
		MonthlyContribution: sim.MonthlyContribution,
		Rate:                sim.Rate,
		Term:                sim.Term,
		FinalAmount:         sim.FinalAmount,
		Timeline:            timeline,
	}, nil
}
