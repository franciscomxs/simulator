package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
)

func TestChain_Empty(t *testing.T) {
	called := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	got := middlewares.Chain(h)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	got.ServeHTTP(rr, req)

	if !called {
		t.Fatal("expected base handler to be invoked")
	}
}

func TestChain_OrderFirstListedRunsFirst(t *testing.T) {
	var order []string

	record := func(name string) middlewares.Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}

	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "h")
	})

	chained := middlewares.Chain(base, record("A"), record("B"))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	chained.ServeHTTP(rr, req)

	want := []string{"A", "B", "h"}
	if len(order) != len(want) {
		t.Fatalf("expected order %v, got %v", want, order)
	}
	for i, name := range want {
		if order[i] != name {
			t.Errorf("position %d: expected %q, got %q", i, name, order[i])
		}
	}
}
