package memory_test

import (
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/outbound/memory"
	"github.com/franciscomxs/simulator/internal/domain"
)

func TestSimulationStore_Save(t *testing.T) {
	store := memory.NewSimulationStore()
	err := store.Save(domain.LoanSimulation{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Count() != 1 {
		t.Errorf("expected 1 entry, got %d", store.Count())
	}
}

func TestSimulationStore_MaxEntriesEnforced(t *testing.T) {
	store := memory.NewSimulationStore()
	for i := 0; i < memory.MaxSimulations+10; i++ {
		if err := store.Save(domain.LoanSimulation{}); err != nil {
			t.Fatalf("unexpected error at i=%d: %v", i, err)
		}
	}
	if store.Count() > memory.MaxSimulations {
		t.Errorf("store has %d entries, max is %d", store.Count(), memory.MaxSimulations)
	}
}
