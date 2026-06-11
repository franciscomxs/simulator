package dto

// InstallmentResponse is a single installment in the simulate loan response.
type InstallmentResponse struct {
	Number    int     `json:"number"`
	Type      string  `json:"type"`
	Payment   float64 `json:"payment"`
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Balance   float64 `json:"balance"`
}

// SimulateResponse is the JSON body for the simulate loan response.
type SimulateResponse struct {
	Amount         float64               `json:"amount"`
	Rate           float64               `json:"rate"`
	Term           int                   `json:"term"`
	System         string                `json:"system"`
	GracePeriod    int                   `json:"grace_period"`
	AdjustedAmount float64               `json:"adjusted_amount"`
	TotalDuration  int                   `json:"total_duration"`
	CustomerType   string                `json:"customer_type"`
	GrossValue     float64               `json:"gross_value"`
	IOF            float64               `json:"iof"`
	NetValue       float64               `json:"net_value"`
	FinancedAmount float64               `json:"financed_amount"`
	TotalAmount    float64               `json:"total_amount"`
	Installments   []InstallmentResponse `json:"installments"`
}

// ErrorResponse is the JSON body for an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// YearlySnapshotResponse is a single yearly entry in the investment simulation response.
type YearlySnapshotResponse struct {
	Year          int     `json:"year"`
	Contributions float64 `json:"contributions"`
	Balance       float64 `json:"balance"`
}

// SimulateInvestmentResponse is the JSON body for the simulate investment response.
type SimulateInvestmentResponse struct {
	InitialAmount       float64                  `json:"initial_amount"`
	MonthlyContribution float64                  `json:"monthly_contribution"`
	Rate                float64                  `json:"rate"`
	Term                int                      `json:"term"`
	FinalAmount         float64                  `json:"final_amount"`
	Timeline            []YearlySnapshotResponse `json:"timeline"`
}
