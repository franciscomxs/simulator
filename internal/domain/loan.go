package domain

import "math"

// AmortizationSystem represents the loan amortization system.
type AmortizationSystem string

const (
	SystemPRICE AmortizationSystem = "PRICE"
	SystemSAC   AmortizationSystem = "SAC"
)

// InstallmentType marks an installment as a grace month or a payment month.
type InstallmentType string

const (
	InstallmentGrace   InstallmentType = "GRACE"
	InstallmentPayment InstallmentType = "PAYMENT"
)

// LoanParams holds the input parameters for a loan simulation.
type LoanParams struct {
	Amount      float64
	Rate        float64
	Term        int
	System      AmortizationSystem
	GracePeriod int
}

// Installment represents a single loan installment.
type Installment struct {
	Number    int
	Type      InstallmentType
	Payment   float64
	Principal float64
	Interest  float64
	Balance   float64
}

// LoanSimulation holds the full result of a loan simulation.
type LoanSimulation struct {
	Amount         float64
	Rate           float64
	Term           int
	System         AmortizationSystem
	GracePeriod    int
	AdjustedAmount float64
	TotalDuration  int
	TotalAmount    float64
	Installments   []Installment
}

// round2 rounds a value to 2 decimal places.
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

const maxTerm = 600

// Simulate validates params and runs the appropriate amortization simulation.
func Simulate(p LoanParams) (LoanSimulation, error) {
	if p.Amount <= 0 {
		return LoanSimulation{}, ErrInvalidAmount
	}
	if p.Rate < 0 || p.Rate >= 1.0 {
		return LoanSimulation{}, ErrInvalidRate
	}
	if p.Term <= 0 || p.Term > maxTerm {
		return LoanSimulation{}, ErrInvalidTerm
	}
	if p.System != SystemPRICE && p.System != SystemSAC {
		return LoanSimulation{}, ErrInvalidSystem
	}
	if p.GracePeriod < 0 {
		return LoanSimulation{}, ErrInvalidGracePeriod
	}
	if p.GracePeriod+p.Term > maxTerm {
		return LoanSimulation{}, ErrInvalidGracePeriod
	}

	graceEntries, adjustedAmount := capitalizeGrace(p)

	amortParams := p
	amortParams.Amount = adjustedAmount

	var sim LoanSimulation
	switch p.System {
	case SystemPRICE:
		sim = simulatePRICE(amortParams)
	default:
		sim = simulateSAC(amortParams)
	}

	for i := range sim.Installments {
		sim.Installments[i].Number = i + 1 + p.GracePeriod
	}

	sim.Amount = p.Amount
	sim.GracePeriod = p.GracePeriod
	sim.AdjustedAmount = adjustedAmount
	sim.TotalDuration = p.GracePeriod + p.Term
	sim.Installments = append(graceEntries, sim.Installments...)

	return sim, nil
}

// capitalizeGrace builds the grace installments and returns the capitalized balance.
// When GracePeriod is zero, returns an empty slice and the original amount.
func capitalizeGrace(p LoanParams) ([]Installment, float64) {
	if p.GracePeriod == 0 {
		return nil, p.Amount
	}

	entries := make([]Installment, p.GracePeriod)
	balance := p.Amount
	for i := 0; i < p.GracePeriod; i++ {
		interest := round2(balance * p.Rate)
		balance = round2(balance + interest)
		entries[i] = Installment{
			Number:   i + 1,
			Type:     InstallmentGrace,
			Interest: interest,
			Balance:  balance,
		}
	}
	return entries, balance
}

// simulatePRICE calculates the PRICE (French) amortization schedule.
func simulatePRICE(p LoanParams) LoanSimulation {
	var pmt float64
	if p.Rate == 0 {
		pmt = p.Amount / float64(p.Term)
	} else {
		pmt = p.Amount * p.Rate / (1 - math.Pow(1+p.Rate, -float64(p.Term)))
	}

	installments := make([]Installment, p.Term)
	balance := p.Amount
	totalPrincipal := 0.0

	for i := 0; i < p.Term; i++ {
		interest := round2(balance * p.Rate)

		var principal, payment float64
		if i == p.Term-1 {
			principal = round2(p.Amount - totalPrincipal)
			payment = round2(principal + interest)
		} else {
			principal = round2(pmt - interest)
			payment = round2(pmt)
		}

		totalPrincipal += principal
		balance -= principal

		installments[i] = Installment{
			Number:    i + 1,
			Type:      InstallmentPayment,
			Payment:   payment,
			Principal: principal,
			Interest:  interest,
			Balance:   round2(balance),
		}
	}

	totalAmount := 0.0
	for _, inst := range installments {
		totalAmount += inst.Payment
	}

	return LoanSimulation{
		Amount:       p.Amount,
		Rate:         p.Rate,
		Term:         p.Term,
		System:       p.System,
		TotalAmount:  round2(totalAmount),
		Installments: installments,
	}
}

// simulateSAC calculates the SAC (Constant Amortization) schedule.
func simulateSAC(p LoanParams) LoanSimulation {
	principalAmort := round2(p.Amount / float64(p.Term))

	installments := make([]Installment, p.Term)
	balance := p.Amount
	totalPrincipal := 0.0

	for i := 0; i < p.Term; i++ {
		interest := round2(balance * p.Rate)

		var principal float64
		if i == p.Term-1 {
			principal = round2(p.Amount - totalPrincipal)
		} else {
			principal = principalAmort
		}

		payment := round2(principal + interest)

		totalPrincipal += principal
		balance -= principal

		installments[i] = Installment{
			Number:    i + 1,
			Type:      InstallmentPayment,
			Payment:   payment,
			Principal: principal,
			Interest:  interest,
			Balance:   round2(balance),
		}
	}

	totalAmount := 0.0
	for _, inst := range installments {
		totalAmount += inst.Payment
	}

	return LoanSimulation{
		Amount:       p.Amount,
		Rate:         p.Rate,
		Term:         p.Term,
		System:       p.System,
		TotalAmount:  round2(totalAmount),
		Installments: installments,
	}
}
