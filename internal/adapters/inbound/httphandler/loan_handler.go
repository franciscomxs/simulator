package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/ports"
)

const maxRequestBodyBytes = 4 << 10 // 4 kb

// LoanHandler holds the HTTP handler dependencies for loan simulation.
type LoanHandler struct {
	useCase ports.SimulateLoanUseCase
}

// NewLoanHandler creates a new LoanHandler.
func NewLoanHandler(uc ports.SimulateLoanUseCase) *LoanHandler {
	return &LoanHandler{useCase: uc}
}

// Simulate handles POST /loan/simulate requests.
func (h *LoanHandler) Simulate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req dto.SimulateRequest
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

	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	input := ports.SimulateLoanInput{
		Amount:       req.Amount,
		Rate:         req.Rate,
		Term:         req.Term,
		System:       req.System,
		GracePeriod:  req.GracePeriod,
		CustomerType: req.CustomerType,
	}

	out, err := h.useCase.Execute(input)
	if err != nil {
		log.Printf("loan simulate: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid simulation parameters"})
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

	resp := dto.SimulateResponse{
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
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
