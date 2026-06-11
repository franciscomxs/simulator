package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
)

type decodedKey[T any] struct{}

// DecodeJSON returns a middleware that decodes the request body into a value of
// type T and stashes it in the request context. On a *http.MaxBytesError the
// middleware writes 413 and stops; on any other decode error it writes 400 and
// stops.
func DecodeJSON[T any]() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				return
			}
			ctx := context.WithValue(r.Context(), decodedKey[T]{}, v)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Decoded retrieves the value previously installed by DecodeJSON[T]() from the
// request context. The boolean indicates whether a value was present.
func Decoded[T any](r *http.Request) (T, bool) {
	v, ok := r.Context().Value(decodedKey[T]{}).(T)
	return v, ok
}
