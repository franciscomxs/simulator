package container_test

import (
	"testing"

	"github.com/franciscomxs/simulator/internal/infrastructure/container"
)

func TestNewSimulateLoanUseCase_ReturnsNonNil(t *testing.T) {
	uc := container.NewSimulateLoanUseCase()
	if uc == nil {
		t.Error("expected non-nil SimulateLoanUseCase")
	}
}

func TestNewSimulateInvestmentUseCase_ReturnsNonNil(t *testing.T) {
	uc := container.NewSimulateInvestmentUseCase()
	if uc == nil {
		t.Error("expected non-nil SimulateInvestmentUseCase")
	}
}
