package domain

import "errors"

var (
	ErrInvalidAmount               = errors.New("invalid amount")
	ErrInvalidRate                 = errors.New("invalid rate")
	ErrInvalidTerm                 = errors.New("invalid term")
	ErrInvalidSystem               = errors.New("invalid amortization system")
	ErrInvalidInitialAmount        = errors.New("invalid initial_amount")
	ErrInvalidMonthlyContribution  = errors.New("invalid monthly_contribution")
	ErrInvalidSimulationResult     = errors.New("simulation produced invalid result")
)
