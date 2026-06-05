package container

import (
	"github.com/franciscomxs/simulator/internal/application/ports"
	"github.com/franciscomxs/simulator/internal/application/usecases"
)

// NewSimulateLoanUseCase creates and returns a SimulateLoanUseCase.
func NewSimulateLoanUseCase() ports.SimulateLoanUseCase {
	return usecases.NewSimulateLoan()
}

// NewSimulateInvestmentUseCase creates and returns a SimulateInvestmentUseCase.
func NewSimulateInvestmentUseCase() ports.SimulateInvestmentUseCase {
	return usecases.NewSimulateInvestment()
}
