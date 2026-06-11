package handlers

import (
	"log/slog"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
	"github.com/franciscomxs/simulator/internal/application/ports"
)

// LoanHandler holds the HTTP handler dependencies for loan simulation.
type LoanHandler struct {
	useCase ports.SimulateLoanUseCase
	log     *slog.Logger
}

// NewLoanHandler creates a new LoanHandler.
func NewLoanHandler(uc ports.SimulateLoanUseCase, log *slog.Logger) *LoanHandler {
	return &LoanHandler{useCase: uc, log: log}
}

// Simulate handles POST /loan/simulate requests. Expects the request to flow
// through the per-route middleware chain (DecodeJSON, WriteJSON, ErrorMapper).
func (h *LoanHandler) Simulate(w http.ResponseWriter, r *http.Request) {
	req, _ := middlewares.Decoded[dto.SimulateRequest](r)

	if err := req.Validate(); err != nil {
		middlewares.SetResponse(r, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	out, err := h.useCase.Execute(ports.SimulateLoanInput{
		Amount:       req.Amount,
		Rate:         req.Rate,
		Term:         req.Term,
		System:       req.System,
		GracePeriod:  req.GracePeriod,
		CustomerType: req.CustomerType,
	})
	if err != nil {
		h.log.Error("loan simulate", "error", err)
		middlewares.SetError(r, err, "invalid simulation parameters")
		return
	}

	installments := make([]dto.InstallmentResponse, len(out.Installments))
	for i, inst := range out.Installments {
		installments[i] = dto.InstallmentResponse{
			Number:    inst.Number,
			Type:      inst.Type,
			Payment:   inst.Payment,
			Principal: inst.Principal,
			Interest:  inst.Interest,
			Balance:   inst.Balance,
		}
	}

	middlewares.SetResponse(r, http.StatusOK, dto.SimulateResponse{
		Amount:         out.Amount,
		Rate:           out.Rate,
		Term:           out.Term,
		System:         out.System,
		GracePeriod:    out.GracePeriod,
		AdjustedAmount: out.AdjustedAmount,
		TotalDuration:  out.TotalDuration,
		CustomerType:   out.CustomerType,
		GrossValue:     out.GrossValue,
		IOF:            out.IOF,
		NetValue:       out.NetValue,
		FinancedAmount: out.FinancedAmount,
		TotalAmount:    out.TotalAmount,
		Installments:   installments,
	})
}
