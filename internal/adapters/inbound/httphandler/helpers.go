package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/domain"
)

const maxRequestBodyBytes = 4 << 10 // 4 KB

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		w.Header().Set("Content-Type", "application/json")
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "request body too large"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "invalid request body"})
		}
		return v, false
	}
	return v, true
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func httpStatusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrInvalidRate),
		errors.Is(err, domain.ErrInvalidTerm),
		errors.Is(err, domain.ErrInvalidSystem),
		errors.Is(err, domain.ErrInvalidInitialAmount),
		errors.Is(err, domain.ErrInvalidMonthlyContribution):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
