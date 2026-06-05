package ports

import "github.com/franciscomxs/simulator/internal/domain"

// SimulationStore is the output port for persisting loan simulations.
type SimulationStore interface {
	Save(simulation domain.LoanSimulation) error
}
