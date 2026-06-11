package domain

import "errors"

var (
	ErrInvalidAmount               = errors.New("invalid amount")
	ErrInvalidRate                 = errors.New("invalid rate")
	ErrInvalidTerm                 = errors.New("invalid term")
	ErrInvalidGracePeriod          = errors.New("invalid grace_period")
	ErrInvalidSystem               = errors.New("invalid amortization system")
	ErrInvalidCustomerType         = errors.New("invalid customer_type")
	ErrInvalidInitialAmount        = errors.New("invalid initial_amount")
	ErrInvalidMonthlyContribution  = errors.New("invalid monthly_contribution")
	ErrInvalidSimulationResult     = errors.New("simulation produced invalid result")
)
