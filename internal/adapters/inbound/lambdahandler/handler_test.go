package lambdahandler_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/lambdahandler"
	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/usecases"
)

func newTestHandler() *lambdahandler.Handler {
	loanUC := usecases.NewSimulateLoan()
	investUC := usecases.NewSimulateInvestment()
	return lambdahandler.NewHandler(loanUC, investUC)
}

func TestLambdaHandler_PRICE_ValidRequest(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.Unmarshal([]byte(resp.Body), &simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(simResp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(simResp.Installments))
	}
	if simResp.GrossValue != 10000 {
		t.Errorf("GrossValue: got %.2f want 10000", simResp.GrossValue)
	}
	if simResp.IOF != 333.20 {
		t.Errorf("IOF: got %.4f want 333.20", simResp.IOF)
	}
	if simResp.NetValue != 9666.80 {
		t.Errorf("NetValue: got %.4f want 9666.80", simResp.NetValue)
	}
	if simResp.FinancedAmount != 10000 {
		t.Errorf("FinancedAmount: got %.2f want 10000", simResp.FinancedAmount)
	}
	if simResp.CustomerType != "PF" {
		t.Errorf("CustomerType: got %q want PF", simResp.CustomerType)
	}
}

func TestLambdaHandler_MissingCustomerType(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_SAC_ValidRequest(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC", "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.Unmarshal([]byte(resp.Body), &simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(simResp.Installments) != 12 {
		t.Errorf("expected 12 installments, got %d", len(simResp.Installments))
	}
}

func TestLambdaHandler_LoanGracePeriod(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "grace_period": 3, "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.Unmarshal([]byte(resp.Body), &simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if simResp.GracePeriod != 3 {
		t.Errorf("GracePeriod: got %d want 3", simResp.GracePeriod)
	}
	if simResp.AdjustedAmount != 10612.08 {
		t.Errorf("AdjustedAmount: got %.2f want 10612.08", simResp.AdjustedAmount)
	}
	if simResp.TotalDuration != 15 {
		t.Errorf("TotalDuration: got %d want 15", simResp.TotalDuration)
	}
	if len(simResp.Installments) != 15 {
		t.Fatalf("expected 15 installments, got %d", len(simResp.Installments))
	}
	for i := 0; i < 3; i++ {
		if simResp.Installments[i].Type != "GRACE" {
			t.Errorf("installments[%d].Type: got %q want \"GRACE\"", i, simResp.Installments[i].Type)
		}
	}
}

func TestLambdaHandler_LoanGracePeriodZero_SafeDefaults(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateResponse
	if err := json.Unmarshal([]byte(resp.Body), &simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if simResp.GracePeriod != 0 {
		t.Errorf("GracePeriod: got %d want 0", simResp.GracePeriod)
	}
	if simResp.AdjustedAmount != simResp.Amount {
		t.Errorf("AdjustedAmount: got %.2f want %.2f", simResp.AdjustedAmount, simResp.Amount)
	}
	if simResp.TotalDuration != simResp.Term {
		t.Errorf("TotalDuration: got %d want %d", simResp.TotalDuration, simResp.Term)
	}
	for i, inst := range simResp.Installments {
		if inst.Type != "PAYMENT" {
			t.Errorf("installments[%d].Type: got %q want \"PAYMENT\"", i, inst.Type)
		}
	}
}

func TestLambdaHandler_UnknownRoute(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE", "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_InvalidSystem(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID", "customer_type": "PF"}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_InvalidJSON(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{invalid json`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_Investment_Valid(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/investment/simulate",
		Body:       `{"initial_amount": 0, "monthly_contribution": 1000, "rate": 0.01, "term": 10}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var simResp dto.SimulateInvestmentResponse
	if err := json.Unmarshal([]byte(resp.Body), &simResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(simResp.Timeline) != 10 {
		t.Errorf("expected 10 timeline entries, got %d", len(simResp.Timeline))
	}
	last := simResp.Timeline[len(simResp.Timeline)-1]
	if last.Balance != simResp.FinalAmount {
		t.Errorf("last timeline balance %.2f != final_amount %.2f", last.Balance, simResp.FinalAmount)
	}
}

func TestLambdaHandler_Investment_InvalidTerm(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/investment/simulate",
		Body:       `{"initial_amount": 0, "monthly_contribution": 1000, "rate": 0.01, "term": 0}`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_Investment_InvalidJSON(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/investment/simulate",
		Body:       `{invalid json`,
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_OversizedLoanBody(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       strings.Repeat("x", 70_000),
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestLambdaHandler_OversizedInvestmentBody(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/investment/simulate",
		Body:       strings.Repeat("x", 70_000),
	}

	resp, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
