package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
)

// WriteJSON returns a middleware that installs a response envelope into the
// request context. After the next handler returns, the middleware writes the
// envelope's body as JSON with the envelope's status, unless an error is set or
// no body was provided.
func WriteJSON() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			env := &envelope{}
			ctx := context.WithValue(r.Context(), envelopeKey{}, env)
			next.ServeHTTP(w, r.WithContext(ctx))
			if env.Err != nil || env.Body == nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(env.Status)
			_ = json.NewEncoder(w).Encode(env.Body)
		})
	}
}
