package middlewares

import "net/http"

// Middleware is a standard HTTP middleware function.
type Middleware func(http.Handler) http.Handler

// Chain wraps h with the supplied middlewares so the first listed is the
// outermost wrapper and runs first on the request.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
