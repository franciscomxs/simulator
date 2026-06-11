package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/domain"
)

// ErrorMapper returns a middleware that, after the next handler returns,
// inspects the response envelope and writes a JSON error response if the
// handler set an error. The HTTP status is derived from httpStatusForError.
func ErrorMapper() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			env, ok := envelopeFrom(r)
			if !ok || env.Err == nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpStatusForError(env.Err))
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: env.ErrMsg})
		})
	}
}

func httpStatusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrInvalidRate),
		errors.Is(err, domain.ErrInvalidTerm),
		errors.Is(err, domain.ErrInvalidSystem),
		errors.Is(err, domain.ErrInvalidGracePeriod),
		errors.Is(err, domain.ErrInvalidCustomerType),
		errors.Is(err, domain.ErrInvalidInitialAmount),
		errors.Is(err, domain.ErrInvalidMonthlyContribution):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
