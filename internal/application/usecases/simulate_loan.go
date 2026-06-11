package usecases

import (
	"github.com/franciscomxs/simulator/internal/application/ports"
	"github.com/franciscomxs/simulator/internal/domain"
)

type simulateLoanUseCase struct{}

// NewSimulateLoan creates a new SimulateLoanUseCase.
func NewSimulateLoan() ports.SimulateLoanUseCase {
	return &simulateLoanUseCase{}
}

// Execute runs the simulate loan use case.
func (uc *simulateLoanUseCase) Execute(input ports.SimulateLoanInput) (ports.SimulateLoanOutput, error) {
	params := domain.LoanParams{
		Amount:      input.Amount,
		Rate:        input.Rate,
		Term:        input.Term,
		System:      domain.AmortizationSystem(input.System),
		GracePeriod: input.GracePeriod,
	}

	sim, err := domain.Simulate(params)
	if err != nil {
		return ports.SimulateLoanOutput{}, err
	}

	installments := make([]ports.InstallmentOutput, len(sim.Installments))
	for i, inst := range sim.Installments {
		installments[i] = ports.InstallmentOutput{
			Number:    inst.Number,
			Type:      string(inst.Type),
			Payment:   inst.Payment,
			Principal: inst.Principal,
			Interest:  inst.Interest,
			Balance:   inst.Balance,
		}
	}

	return ports.SimulateLoanOutput{
		Amount:         sim.Amount,
		Rate:           sim.Rate,
		Term:           sim.Term,
		System:         string(sim.System),
		GracePeriod:    sim.GracePeriod,
		AdjustedAmount: sim.AdjustedAmount,
		TotalDuration:  sim.TotalDuration,
		TotalAmount:    sim.TotalAmount,
		Installments:   installments,
	}, nil
}
