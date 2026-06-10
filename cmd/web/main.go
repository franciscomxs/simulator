package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/httphandler"
	"github.com/franciscomxs/simulator/internal/infrastructure/config"
	"github.com/franciscomxs/simulator/internal/infrastructure/container"
)

func main() {
	cfg := config.Load()
	logger := slog.Default()

	loanUC := container.NewSimulateLoanUseCase()
	investUC := container.NewSimulateInvestmentUseCase()
	loanHandler := httphandler.NewLoanHandler(loanUC, logger)
	investHandler := httphandler.NewInvestmentHandler(investUC, logger)
	router := httphandler.NewRouter(loanHandler, investHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      http.TimeoutHandler(router, 9*time.Second, `{"error":"request timeout"}`),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting web server on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
