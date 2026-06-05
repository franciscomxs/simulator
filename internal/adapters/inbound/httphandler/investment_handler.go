package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/ports"
)

// InvestmentHandler holds the HTTP handler dependencies for investment simulation.
type InvestmentHandler struct {
	useCase ports.SimulateInvestmentUseCase
}

// NewInvestmentHandler creates a new InvestmentHandler.
func NewInvestmentHandler(uc ports.SimulateInvestmentUseCase) *InvestmentHandler {
	return &InvestmentHandler{useCase: uc}
}

// SimulateInvestment handles POST /investment/simulate requests.
func (h *InvestmentHandler) SimulateInvestment(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req dto.SimulateInvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "request body too large"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid request body"})
		}
		return
	}

	input := ports.SimulateInvestmentInput{
		InitialAmount:       req.InitialAmount,
		MonthlyContribution: req.MonthlyContribution,
		Rate:                req.Rate,
		Term:                req.Term,
	}

	out, err := h.useCase.Execute(input)
	if err != nil {
		log.Printf("investment simulate: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid simulation parameters"})
		return
	}

	timeline := make([]dto.YearlySnapshotResponse, len(out.Timeline))
	for i, snap := range out.Timeline {
		timeline[i] = dto.YearlySnapshotResponse{
			Year:          snap.Year,
			Contributions: snap.Contributions,
			Balance:       snap.Balance,
		}
	}

	resp := dto.SimulateInvestmentResponse{
		InitialAmount:       out.InitialAmount,
		MonthlyContribution: out.MonthlyContribution,
		Rate:                out.Rate,
		Term:                out.Term,
		FinalAmount:         out.FinalAmount,
		Timeline:            timeline,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
