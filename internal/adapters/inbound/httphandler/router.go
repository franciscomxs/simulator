package httphandler

import "net/http"

func NewRouter(lh *LoanHandler, ih *InvestmentHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /loan/simulate", lh.Simulate)
	mux.HandleFunc("POST /investment/simulate", ih.SimulateInvestment)

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

	return mux
}
