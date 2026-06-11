package dto

import "errors"

// ErrInvalidCustomerType is returned when customer_type is missing or unknown.
var ErrInvalidCustomerType = errors.New("invalid customer_type")

// SimulateRequest is the JSON body for the simulate loan request.
type SimulateRequest struct {
	Amount       float64 `json:"amount"`
	Rate         float64 `json:"rate"`
	Term         int     `json:"term"`
	System       string  `json:"system"`
	GracePeriod  int     `json:"grace_period"`
	CustomerType string  `json:"customer_type"`
}

// Validate enforces SimulateRequest invariants enforced at the transport
// boundary (currently: customer_type must be "PF" or "PJ").
func (r SimulateRequest) Validate() error {
	if r.CustomerType != "PF" && r.CustomerType != "PJ" {
		return ErrInvalidCustomerType
	}
	return nil
}

// SimulateInvestmentRequest is the JSON body for the simulate investment request.
type SimulateInvestmentRequest struct {
	InitialAmount       float64 `json:"initial_amount"`
	MonthlyContribution float64 `json:"monthly_contribution"`
	Rate                float64 `json:"rate"`
	Term                int     `json:"term"`
}
