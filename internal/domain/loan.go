package domain

import "math"

// AmortizationSystem represents the loan amortization system.
type AmortizationSystem string

const (
	SystemPRICE AmortizationSystem = "PRICE"
	SystemSAC   AmortizationSystem = "SAC"
)

// LoanParams holds the input parameters for a loan simulation.
type LoanParams struct {
	Amount float64
	Rate   float64
	Term   int
	System AmortizationSystem
}

// Installment represents a single loan installment.
type Installment struct {
	Number    int
	Payment   float64
	Principal float64
	Interest  float64
}

// LoanSimulation holds the full result of a loan simulation.
type LoanSimulation struct {
	Amount       float64
	Rate         float64
	Term         int
	System       AmortizationSystem
	TotalAmount  float64
	Installments []Installment
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

	switch p.System {
	case SystemPRICE:
		return simulatePRICE(p), nil
	default:
		return simulateSAC(p), nil
	}
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

		installments[i] = Installment{
			Number:    i + 1,
			Payment:   payment,
			Principal: principal,
			Interest:  interest,
		}

		totalPrincipal += principal
		balance -= principal
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

		installments[i] = Installment{
			Number:    i + 1,
			Payment:   payment,
			Principal: principal,
			Interest:  interest,
		}

		totalPrincipal += principal
		balance -= principal
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
