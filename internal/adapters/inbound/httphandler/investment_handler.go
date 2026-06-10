package httphandler

import (
	"log/slog"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/ports"
)

// InvestmentHandler holds the HTTP handler dependencies for investment simulation.
type InvestmentHandler struct {
	useCase ports.SimulateInvestmentUseCase
	log     *slog.Logger
}

// NewInvestmentHandler creates a new InvestmentHandler.
func NewInvestmentHandler(uc ports.SimulateInvestmentUseCase, log *slog.Logger) *InvestmentHandler {
	return &InvestmentHandler{useCase: uc, log: log}
}

// SimulateInvestment handles POST /investment/simulate requests.
func (h *InvestmentHandler) SimulateInvestment(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeJSON[dto.SimulateInvestmentRequest](w, r)
	if !ok {
		return
	}

	out, err := h.useCase.Execute(ports.SimulateInvestmentInput{
		InitialAmount:       req.InitialAmount,
		MonthlyContribution: req.MonthlyContribution,
		Rate:                req.Rate,
		Term:                req.Term,
	})
	if err != nil {
		h.log.Error("investment simulate", "error", err)
		writeJSON(w, httpStatusForError(err), dto.ErrorResponse{Error: "invalid simulation parameters"})
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

	writeJSON(w, http.StatusOK, dto.SimulateInvestmentResponse{
		InitialAmount:       out.InitialAmount,
		MonthlyContribution: out.MonthlyContribution,
		Rate:                out.Rate,
		Term:                out.Term,
		FinalAmount:         out.FinalAmount,
		Timeline:            timeline,
	})
}
