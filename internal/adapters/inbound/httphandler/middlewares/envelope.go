package middlewares

import "net/http"

type envelope struct {
	Status int
	Body   any
	Err    error
	ErrMsg string
}

type envelopeKey struct{}

func envelopeFrom(r *http.Request) (*envelope, bool) {
	e, ok := r.Context().Value(envelopeKey{}).(*envelope)
	return e, ok
}

// SetResponse stores the handler's success response in the envelope so that
// WriteJSON can render it after the handler returns. No-op if the envelope is
// not installed (handler invoked outside the middleware chain).
func SetResponse(r *http.Request, status int, body any) {
	if env, ok := envelopeFrom(r); ok {
		env.Status = status
		env.Body = body
	}
}

// SetError stores the handler's error and a user-facing message in the envelope
// so that ErrorMapper can render them after the handler returns. No-op if the
// envelope is not installed.
func SetError(r *http.Request, err error, message string) {
	if env, ok := envelopeFrom(r); ok {
		env.Err = err
		env.ErrMsg = message
	}
}
