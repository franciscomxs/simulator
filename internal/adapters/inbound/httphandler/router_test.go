package httphandler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler"
	"github.com/franciscomxs/simulator/internal/application/usecases"
)

func newTestInvestmentHandler() *httphandler.InvestmentHandler {
	uc := usecases.NewSimulateInvestment()
	return httphandler.NewInvestmentHandler(uc)
}

func newTestRouter() http.Handler {
	return httphandler.NewRouter(newTestLoanHandler(), newTestInvestmentHandler())
}

func TestHTTPRouter_HealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	newTestRouter().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", resp["status"])
	}
}

func TestHTTPRouter_OpenAPISpec(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rr := httptest.NewRecorder()

	newTestRouter().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/yaml" {
		t.Errorf("expected Content-Type application/yaml, got %q", ct)
	}
	if rr.Body.Len() == 0 {
		t.Error("expected non-empty OpenAPI spec body")
	}
}

func TestHTTPRouter_Docs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rr := httptest.NewRecorder()

	newTestRouter().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html; charset=utf-8, got %q", ct)
	}
	if rr.Body.Len() == 0 {
		t.Error("expected non-empty docs body")
	}
}
