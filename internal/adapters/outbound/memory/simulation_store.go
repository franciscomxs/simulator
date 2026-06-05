package memory

import (
	"sync"

	"github.com/franciscomxs/simulator/internal/domain"
)

// MaxSimulations is the maximum number of simulations held in memory.
const MaxSimulations = 10_000

// SimulationStore is an in-memory implementation of ports.SimulationStore.
type SimulationStore struct {
	mu          sync.Mutex
	simulations []domain.LoanSimulation
}

// NewSimulationStore creates a new in-memory SimulationStore.
func NewSimulationStore() *SimulationStore {
	return &SimulationStore{}
}

// Save stores a loan simulation in memory, evicting the oldest half when the cap is reached.
func (s *SimulationStore) Save(simulation domain.LoanSimulation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.simulations) >= MaxSimulations {
		s.simulations = s.simulations[MaxSimulations/2:]
	}
	s.simulations = append(s.simulations, simulation)
	return nil
}

// Count returns the current number of stored simulations.
func (s *SimulationStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.simulations)
}
