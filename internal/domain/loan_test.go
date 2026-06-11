package domain_test

import (
	"errors"
	"math"
	"testing"

	"github.com/franciscomxs/simulator/internal/domain"
)

const (
	testAmount = 10000.0
	testRate   = 0.02
	testTerm   = 12
)

func priceParams() domain.LoanParams {
	return domain.LoanParams{
		Amount:       testAmount,
		Rate:         testRate,
		Term:         testTerm,
		System:       domain.SystemPRICE,
		CustomerType: domain.CustomerTypePF,
	}
}

func sacParams() domain.LoanParams {
	return domain.LoanParams{
		Amount:       testAmount,
		Rate:         testRate,
		Term:         testTerm,
		System:       domain.SystemSAC,
		CustomerType: domain.CustomerTypePF,
	}
}

// --- PRICE tests ---

func TestSimulatePRICE_InstallmentCount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sim.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(sim.Installments))
	}
}

func TestSimulatePRICE_FixedPayment(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first := sim.Installments[0].Payment
	for i, inst := range sim.Installments[:len(sim.Installments)-1] {
		if inst.Payment != first {
			t.Errorf("installment %d payment %.2f != first %.2f", i+1, inst.Payment, first)
		}
	}
}

func TestSimulatePRICE_DecreasingInterest(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(sim.Installments); i++ {
		if sim.Installments[i].Interest >= sim.Installments[i-1].Interest {
			t.Errorf("installment %d interest %.2f >= installment %d interest %.2f",
				i+1, sim.Installments[i].Interest, i, sim.Installments[i-1].Interest)
		}
	}
}

func TestSimulatePRICE_SumOfPrincipalsEqualsAmount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Principal
	}
	if diff := sum - testAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != amount %.2f", sum, testAmount)
	}
}

func TestSimulatePRICE_SumOfPaymentsEqualsTotalAmount(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Payment
	}
	if diff := sum - sim.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != TotalAmount %.2f", sum, sim.TotalAmount)
	}
}

func TestSimulatePRICE_InstallmentNumbering(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, inst := range sim.Installments {
		if inst.Number != i+1 {
			t.Errorf("installment index %d has Number %d, expected %d", i, inst.Number, i+1)
		}
	}
}

func TestSimulatePRICE_FirstInterest(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 200.00
	if sim.Installments[0].Interest != expected {
		t.Errorf("first installment interest %.2f != %.2f", sim.Installments[0].Interest, expected)
	}
}

// --- SAC tests ---

func TestSimulateSAC_InstallmentCount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sim.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(sim.Installments))
	}
}

func TestSimulateSAC_ConstantPrincipal(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first := sim.Installments[0].Principal
	for i, inst := range sim.Installments[:len(sim.Installments)-1] {
		if inst.Principal != first {
			t.Errorf("installment %d principal %.2f != first %.2f", i+1, inst.Principal, first)
		}
	}
}

func TestSimulateSAC_DecreasingPayments(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(sim.Installments); i++ {
		if sim.Installments[i].Payment >= sim.Installments[i-1].Payment {
			t.Errorf("installment %d payment %.2f >= installment %d payment %.2f",
				i+1, sim.Installments[i].Payment, i, sim.Installments[i-1].Payment)
		}
	}
}

func TestSimulateSAC_SumOfPrincipalsEqualsAmount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Principal
	}
	if diff := sum - testAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of principals %.2f != amount %.2f", sum, testAmount)
	}
}

func TestSimulateSAC_SumOfPaymentsEqualsTotalAmount(t *testing.T) {
	sim, err := domain.Simulate(sacParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := 0.0
	for _, inst := range sim.Installments {
		sum += inst.Payment
	}
	if diff := sum - sim.TotalAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payments %.2f != TotalAmount %.2f", sum, sim.TotalAmount)
	}
}

// --- Validation tests ---

func TestSimulate_InvalidAmount(t *testing.T) {
	p := priceParams()
	p.Amount = 0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestSimulate_NegativeAmount(t *testing.T) {
	p := priceParams()
	p.Amount = -100
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestSimulate_NegativeRate(t *testing.T) {
	p := priceParams()
	p.Rate = -0.01
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestSimulate_InvalidTerm(t *testing.T) {
	p := priceParams()
	p.Term = 0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm, got %v", err)
	}
}

func TestSimulate_InvalidSystem(t *testing.T) {
	p := priceParams()
	p.System = "INVALID"
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidSystem) {
		t.Errorf("expected ErrInvalidSystem, got %v", err)
	}
}

func TestSimulate_RateAtOrAboveOne(t *testing.T) {
	p := priceParams()
	p.Rate = 1.0
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidRate) {
		t.Errorf("expected ErrInvalidRate for rate=1.0, got %v", err)
	}
}

func TestSimulate_TermExceedsMaximum(t *testing.T) {
	p := priceParams()
	p.Term = 601
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidTerm) {
		t.Errorf("expected ErrInvalidTerm for term=601, got %v", err)
	}
}

func TestSimulate_ZeroRate(t *testing.T) {
	p := domain.LoanParams{
		Amount:       12000,
		Rate:         0,
		Term:         12,
		System:       domain.SystemPRICE,
		CustomerType: domain.CustomerTypePF,
	}
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 1000.0
	for _, inst := range sim.Installments {
		if inst.Payment != expected {
			t.Errorf("expected payment %.2f, got %.2f", expected, inst.Payment)
		}
	}
}

// --- Grace period tests ---

func TestSimulate_NegativeGracePeriod(t *testing.T) {
	p := priceParams()
	p.GracePeriod = -1
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidGracePeriod) {
		t.Errorf("expected ErrInvalidGracePeriod, got %v", err)
	}
}

func TestSimulate_GracePlusTermExceedsMaximum(t *testing.T) {
	p := priceParams()
	p.Term = 600
	p.GracePeriod = 1
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidGracePeriod) {
		t.Errorf("expected ErrInvalidGracePeriod for grace+term > 600, got %v", err)
	}
}

func TestSimulatePRICE_GracePeriod_AdjustedAmountAndGraceEntries(t *testing.T) {
	p := priceParams()
	p.GracePeriod = 3
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sim.AdjustedAmount != 10612.08 {
		t.Errorf("AdjustedAmount: got %.2f want 10612.08", sim.AdjustedAmount)
	}
	if sim.GracePeriod != 3 {
		t.Errorf("GracePeriod: got %d want 3", sim.GracePeriod)
	}
	if sim.TotalDuration != 15 {
		t.Errorf("TotalDuration: got %d want 15", sim.TotalDuration)
	}

	wantGrace := []struct {
		interest float64
		balance  float64
	}{
		{200.00, 10200.00},
		{204.00, 10404.00},
		{208.08, 10612.08},
	}
	for i, want := range wantGrace {
		got := sim.Installments[i]
		if got.Type != domain.InstallmentGrace {
			t.Errorf("grace[%d] Type: got %q want %q", i, got.Type, domain.InstallmentGrace)
		}
		if got.Payment != 0 {
			t.Errorf("grace[%d] Payment: got %.2f want 0", i, got.Payment)
		}
		if got.Principal != 0 {
			t.Errorf("grace[%d] Principal: got %.2f want 0", i, got.Principal)
		}
		if got.Interest != want.interest {
			t.Errorf("grace[%d] Interest: got %.2f want %.2f", i, got.Interest, want.interest)
		}
		if got.Balance != want.balance {
			t.Errorf("grace[%d] Balance: got %.2f want %.2f", i, got.Balance, want.balance)
		}
	}
}

func TestSimulatePRICE_GracePeriod_Numbering(t *testing.T) {
	p := priceParams()
	p.GracePeriod = 3
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sim.Installments) != 15 {
		t.Fatalf("expected 15 installments, got %d", len(sim.Installments))
	}
	for i, inst := range sim.Installments {
		if inst.Number != i+1 {
			t.Errorf("installments[%d].Number: got %d want %d", i, inst.Number, i+1)
		}
	}
	for i := 0; i < 3; i++ {
		if sim.Installments[i].Type != domain.InstallmentGrace {
			t.Errorf("installments[%d].Type: got %q want %q", i, sim.Installments[i].Type, domain.InstallmentGrace)
		}
	}
	paymentCount := 0
	for i := 3; i < 15; i++ {
		if sim.Installments[i].Type != domain.InstallmentPayment {
			t.Errorf("installments[%d].Type: got %q want %q", i, sim.Installments[i].Type, domain.InstallmentPayment)
		}
		paymentCount++
	}
	if paymentCount != p.Term {
		t.Errorf("payment count %d != term %d", paymentCount, p.Term)
	}
}

func TestSimulatePRICE_GracePeriod_PaymentExceedsNoGrace(t *testing.T) {
	noGrace, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error (no grace): %v", err)
	}

	withGrace := priceParams()
	withGrace.GracePeriod = 3
	graced, err := domain.Simulate(withGrace)
	if err != nil {
		t.Fatalf("unexpected error (grace): %v", err)
	}

	noGracePay := noGrace.Installments[0].Payment
	gracedFirstPayment := graced.Installments[3].Payment
	if !(gracedFirstPayment > noGracePay) {
		t.Errorf("expected graced payment %.2f > non-graced %.2f", gracedFirstPayment, noGracePay)
	}
	if !(graced.TotalAmount > noGrace.TotalAmount) {
		t.Errorf("expected graced total %.2f > non-graced %.2f", graced.TotalAmount, noGrace.TotalAmount)
	}
}

func TestSimulateSAC_GracePeriod_Behavior(t *testing.T) {
	p := sacParams()
	p.GracePeriod = 3
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sim.AdjustedAmount != 10612.08 {
		t.Errorf("AdjustedAmount: got %.2f want 10612.08", sim.AdjustedAmount)
	}
	if len(sim.Installments) != 15 {
		t.Fatalf("expected 15 installments, got %d", len(sim.Installments))
	}

	paymentInstallments := sim.Installments[3:]
	if len(paymentInstallments) != p.Term {
		t.Errorf("payment installment count %d != term %d", len(paymentInstallments), p.Term)
	}

	sumPrincipal := 0.0
	for _, inst := range paymentInstallments {
		sumPrincipal += inst.Principal
	}
	if diff := sumPrincipal - sim.AdjustedAmount; diff > 0.01 || diff < -0.01 {
		t.Errorf("sum of payment principals %.2f != AdjustedAmount %.2f", sumPrincipal, sim.AdjustedAmount)
	}

	for i := 1; i < len(paymentInstallments); i++ {
		if paymentInstallments[i].Payment >= paymentInstallments[i-1].Payment {
			t.Errorf("payment %d %.2f >= payment %d %.2f",
				i+1, paymentInstallments[i].Payment, i, paymentInstallments[i-1].Payment)
		}
	}
}

func TestSimulate_GracePeriodZero_RegressionMatchesPreFeature(t *testing.T) {
	preFeature, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error (baseline): %v", err)
	}

	p := priceParams()
	p.GracePeriod = 0
	graced, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error (grace=0): %v", err)
	}

	if graced.GracePeriod != 0 {
		t.Errorf("GracePeriod: got %d want 0", graced.GracePeriod)
	}
	if graced.AdjustedAmount != p.Amount {
		t.Errorf("AdjustedAmount: got %.2f want %.2f", graced.AdjustedAmount, p.Amount)
	}
	if graced.TotalDuration != p.Term {
		t.Errorf("TotalDuration: got %d want %d", graced.TotalDuration, p.Term)
	}
	if graced.TotalAmount != preFeature.TotalAmount {
		t.Errorf("TotalAmount: got %.2f want %.2f", graced.TotalAmount, preFeature.TotalAmount)
	}
	if len(graced.Installments) != len(preFeature.Installments) {
		t.Fatalf("installment count: got %d want %d", len(graced.Installments), len(preFeature.Installments))
	}
	for i, inst := range graced.Installments {
		base := preFeature.Installments[i]
		if inst.Number != base.Number {
			t.Errorf("inst[%d].Number: got %d want %d", i, inst.Number, base.Number)
		}
		if inst.Payment != base.Payment {
			t.Errorf("inst[%d].Payment: got %.2f want %.2f", i, inst.Payment, base.Payment)
		}
		if inst.Principal != base.Principal {
			t.Errorf("inst[%d].Principal: got %.2f want %.2f", i, inst.Principal, base.Principal)
		}
		if inst.Interest != base.Interest {
			t.Errorf("inst[%d].Interest: got %.2f want %.2f", i, inst.Interest, base.Interest)
		}
		if inst.Type != domain.InstallmentPayment {
			t.Errorf("inst[%d].Type: got %q want %q", i, inst.Type, domain.InstallmentPayment)
		}
	}
}

// --- customer_type tests ---

func TestSimulate_RejectsMissingCustomerType(t *testing.T) {
	p := priceParams()
	p.CustomerType = ""
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidCustomerType) {
		t.Errorf("expected ErrInvalidCustomerType, got %v", err)
	}
}

func TestSimulate_RejectsUnknownCustomerType(t *testing.T) {
	p := priceParams()
	p.CustomerType = "XX"
	_, err := domain.Simulate(p)
	if !errors.Is(err, domain.ErrInvalidCustomerType) {
		t.Errorf("expected ErrInvalidCustomerType, got %v", err)
	}
}

// --- IOF tests ---

func TestSimulate_IOF_PF_WorkedExample(t *testing.T) {
	sim, err := domain.Simulate(priceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sim.GrossValue != 10000 {
		t.Errorf("GrossValue: got %.2f, want 10000", sim.GrossValue)
	}
	if sim.IOF != 333.20 {
		t.Errorf("IOF: got %.4f, want 333.20", sim.IOF)
	}
	if sim.NetValue != 9666.80 {
		t.Errorf("NetValue: got %.4f, want 9666.80", sim.NetValue)
	}
	if sim.FinancedAmount != 10000 {
		t.Errorf("FinancedAmount: got %.2f, want 10000", sim.FinancedAmount)
	}
	if sim.CustomerType != domain.CustomerTypePF {
		t.Errorf("CustomerType: got %q, want %q", sim.CustomerType, domain.CustomerTypePF)
	}
}

func TestSimulate_IOF_PJ_SameAsPF(t *testing.T) {
	pfParams := priceParams()
	pjParams := priceParams()
	pjParams.CustomerType = domain.CustomerTypePJ

	pfSim, err := domain.Simulate(pfParams)
	if err != nil {
		t.Fatalf("PF unexpected error: %v", err)
	}
	pjSim, err := domain.Simulate(pjParams)
	if err != nil {
		t.Fatalf("PJ unexpected error: %v", err)
	}

	if pfSim.IOF != pjSim.IOF {
		t.Errorf("IOF PF=%.2f PJ=%.2f, expected equal", pfSim.IOF, pjSim.IOF)
	}
	if pfSim.NetValue != pjSim.NetValue {
		t.Errorf("NetValue PF=%.2f PJ=%.2f, expected equal", pfSim.NetValue, pjSim.NetValue)
	}
	if pfSim.FinancedAmount != pjSim.FinancedAmount {
		t.Errorf("FinancedAmount PF=%.2f PJ=%.2f, expected equal", pfSim.FinancedAmount, pjSim.FinancedAmount)
	}
}

func TestSimulate_IOF_DailyDaysCappedAt365(t *testing.T) {
	p := priceParams()
	p.Term = 24 // 720 days -> capped at 365
	sim, err := domain.Simulate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFixed := 10000 * 0.0038
	expectedDaily := 10000 * 0.000082 * 365
	expectedIOF := math.Round((expectedFixed+expectedDaily)*100) / 100

	if sim.IOF != expectedIOF {
		t.Errorf("IOF: got %.4f, want %.4f (capped at 365 days)", sim.IOF, expectedIOF)
	}
}

func TestSimulate_IOF_AmortizationUnchanged(t *testing.T) {
	pfParams := priceParams()
	pjParams := priceParams()
	pjParams.CustomerType = domain.CustomerTypePJ

	pfSim, err := domain.Simulate(pfParams)
	if err != nil {
		t.Fatalf("PF unexpected error: %v", err)
	}
	pjSim, err := domain.Simulate(pjParams)
	if err != nil {
		t.Fatalf("PJ unexpected error: %v", err)
	}

	if pfSim.TotalAmount != pjSim.TotalAmount {
		t.Errorf("TotalAmount PF=%.2f PJ=%.2f, expected equal", pfSim.TotalAmount, pjSim.TotalAmount)
	}
	if len(pfSim.Installments) != len(pjSim.Installments) {
		t.Fatalf("installments len mismatch: PF=%d PJ=%d", len(pfSim.Installments), len(pjSim.Installments))
	}
	for i := range pfSim.Installments {
		if pfSim.Installments[i] != pjSim.Installments[i] {
			t.Errorf("installment %d differs: PF=%+v PJ=%+v", i, pfSim.Installments[i], pjSim.Installments[i])
		}
	}

	// Frozen reference: pre-feature PRICE result for the canonical inputs.
	expectedFirstPayment := 945.60
	expectedFirstInterest := 200.00
	if pfSim.Installments[0].Payment != expectedFirstPayment {
		t.Errorf("first payment: got %.2f, want %.2f", pfSim.Installments[0].Payment, expectedFirstPayment)
	}
	if pfSim.Installments[0].Interest != expectedFirstInterest {
		t.Errorf("first interest: got %.2f, want %.2f", pfSim.Installments[0].Interest, expectedFirstInterest)
	}
}
