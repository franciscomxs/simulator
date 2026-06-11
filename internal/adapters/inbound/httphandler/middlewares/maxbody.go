package middlewares

import "net/http"

// DefaultMaxBodyBytes is the canonical request body size limit for HTTP routes.
const DefaultMaxBodyBytes int64 = 4 << 10 // 4 KB

// MaxBodyBytes returns a middleware that caps the request body size at limit.
// Reads past limit return *http.MaxBytesError.
func MaxBodyBytes(limit int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
