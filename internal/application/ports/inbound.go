package ports

// SimulateLoanInput holds the input for the simulate loan use case.
type SimulateLoanInput struct {
	Amount      float64
	Rate        float64
	Term        int
	System      string
	GracePeriod int
}

// InstallmentOutput holds a single installment in the use case output.
type InstallmentOutput struct {
	Number    int
	Type      string
	Payment   float64
	Principal float64
	Interest  float64
	Balance   float64
}

// SimulateLoanOutput holds the result of the simulate loan use case.
type SimulateLoanOutput struct {
	Amount         float64
	Rate           float64
	Term           int
	System         string
	GracePeriod    int
	AdjustedAmount float64
	TotalDuration  int
	TotalAmount    float64
	Installments   []InstallmentOutput
}

// SimulateLoanUseCase is the input port for the simulate loan use case.
type SimulateLoanUseCase interface {
	Execute(input SimulateLoanInput) (SimulateLoanOutput, error)
}

// SimulateInvestmentInput holds the input for the simulate investment use case.
type SimulateInvestmentInput struct {
	InitialAmount       float64
	MonthlyContribution float64
	Rate                float64
	Term                int // years
}

// YearlySnapshotOutput holds yearly portfolio state.
type YearlySnapshotOutput struct {
	Year          int
	Contributions float64
	Balance       float64
}

// SimulateInvestmentOutput holds the result of the simulate investment use case.
type SimulateInvestmentOutput struct {
	InitialAmount       float64
	MonthlyContribution float64
	Rate                float64
	Term                int
	FinalAmount         float64
	Timeline            []YearlySnapshotOutput
}

// SimulateInvestmentUseCase is the input port for the simulate investment use case.
type SimulateInvestmentUseCase interface {
	Execute(input SimulateInvestmentInput) (SimulateInvestmentOutput, error)
}
