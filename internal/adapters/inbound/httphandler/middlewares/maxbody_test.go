package middlewares_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
)

func TestMaxBodyBytes_WithinLimit(t *testing.T) {
	var got []byte

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected read error: %v", err)
		}
		got = b
	})

	h := middlewares.MaxBodyBytes(1024)(next)

	body := strings.Repeat("a", 512)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if string(got) != body {
		t.Errorf("expected body %d bytes, got %d bytes", len(body), len(got))
	}
}

func TestMaxBodyBytes_ExceedsLimit(t *testing.T) {
	var readErr error

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	})

	h := middlewares.MaxBodyBytes(1024)(next)

	body := strings.Repeat("a", 2048)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	var maxBytesErr *http.MaxBytesError
	if !errors.As(readErr, &maxBytesErr) {
		t.Fatalf("expected *http.MaxBytesError, got %T (%v)", readErr, readErr)
	}
}

func TestDefaultMaxBodyBytes(t *testing.T) {
	if middlewares.DefaultMaxBodyBytes != 4096 {
		t.Errorf("expected DefaultMaxBodyBytes=4096, got %d", middlewares.DefaultMaxBodyBytes)
	}
}
