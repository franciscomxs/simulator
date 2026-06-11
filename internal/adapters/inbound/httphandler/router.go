package httphandler

import (
	"net/http"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/handlers"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler/middlewares"
)

func NewRouter(lh *handlers.LoanHandler, ih *handlers.InvestmentHandler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /loan/simulate", middlewares.Chain(
		http.HandlerFunc(lh.Simulate),
		middlewares.WriteJSON(),
		middlewares.ErrorMapper(),
		middlewares.DecodeJSON[dto.SimulateRequest](),
	))

	mux.Handle("POST /investment/simulate", middlewares.Chain(
		http.HandlerFunc(ih.SimulateInvestment),
		middlewares.WriteJSON(),
		middlewares.ErrorMapper(),
		middlewares.DecodeJSON[dto.SimulateInvestmentRequest](),
	))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(openAPISpec)
	})

	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(swaggerUIHTML)
	})

	return middlewares.Chain(mux,
		middlewares.MaxBodyBytes(middlewares.DefaultMaxBodyBytes),
	)
}
