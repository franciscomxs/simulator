package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
	"github.com/franciscomxs/simulator/internal/domain"
)

func serve(t *testing.T, handlerErr error, msg string) *httptest.ResponseRecorder {
	t.Helper()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handlerErr != nil {
			middlewares.SetError(r, handlerErr, msg)
		}
	})

	// WriteJSON outer installs envelope; ErrorMapper inner reads it on return.
	h := middlewares.Chain(next, middlewares.WriteJSON(), middlewares.ErrorMapper())

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestErrorMapper_ValidationErrors_422(t *testing.T) {
	cases := []error{
		domain.ErrInvalidAmount,
		domain.ErrInvalidRate,
		domain.ErrInvalidTerm,
		domain.ErrInvalidSystem,
		domain.ErrInvalidInitialAmount,
		domain.ErrInvalidMonthlyContribution,
	}
	for _, err := range cases {
		rr := serve(t, err, "invalid simulation parameters")
		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("err=%v: expected 422, got %d", err, rr.Code)
		}
		var body map[string]string
		_ = json.NewDecoder(rr.Body).Decode(&body)
		if body["error"] != "invalid simulation parameters" {
			t.Errorf("err=%v: unexpected body %v", err, body)
		}
	}
}

func TestErrorMapper_InvalidSimulationResult_500(t *testing.T) {
	rr := serve(t, domain.ErrInvalidSimulationResult, "invalid simulation parameters")
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestErrorMapper_Unknown_500(t *testing.T) {
	rr := serve(t, errors.New("boom"), "internal")
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["error"] != "internal" {
		t.Errorf("expected error=internal, got %v", body)
	}
}

func TestErrorMapper_NoError_NoWrite(t *testing.T) {
	rr := serve(t, nil, "")
	if rr.Body.Len() != 0 {
		t.Errorf("expected no write, got %q", rr.Body.String())
	}
}
