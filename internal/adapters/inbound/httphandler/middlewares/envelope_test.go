package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetResponse_WithEnvelope(t *testing.T) {
	env := &envelope{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), envelopeKey{}, env))

	SetResponse(req, http.StatusCreated, map[string]string{"k": "v"})

	if env.Status != http.StatusCreated {
		t.Errorf("expected status 201, got %d", env.Status)
	}
	body, ok := env.Body.(map[string]string)
	if !ok || body["k"] != "v" {
		t.Errorf("expected body to contain k=v, got %v", env.Body)
	}
}

func TestSetError_WithEnvelope(t *testing.T) {
	env := &envelope{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), envelopeKey{}, env))

	boom := errors.New("boom")
	SetError(req, boom, "oh no")

	if env.Err != boom {
		t.Errorf("expected err=%v, got %v", boom, env.Err)
	}
	if env.ErrMsg != "oh no" {
		t.Errorf("expected msg=oh no, got %q", env.ErrMsg)
	}
}

func TestSetResponse_NoEnvelope_NoOp(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// Must not panic
	SetResponse(req, http.StatusOK, "ignored")
}

func TestSetError_NoEnvelope_NoOp(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// Must not panic
	SetError(req, errors.New("boom"), "ignored")
}
