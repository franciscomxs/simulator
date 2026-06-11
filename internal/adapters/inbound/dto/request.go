package dto

// SimulateRequest is the JSON body for the simulate loan request.
type SimulateRequest struct {
	Amount      float64 `json:"amount"`
	Rate        float64 `json:"rate"`
	Term        int     `json:"term"`
	System      string  `json:"system"`
	GracePeriod int     `json:"grace_period"`
}

// SimulateInvestmentRequest is the JSON body for the simulate investment request.
type SimulateInvestmentRequest struct {
	InitialAmount       float64 `json:"initial_amount"`
	MonthlyContribution float64 `json:"monthly_contribution"`
	Rate                float64 `json:"rate"`
	Term                int     `json:"term"`
}
