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
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`,
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

func TestLambdaHandler_SAC_ValidRequest(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/loan/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "SAC"}`,
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

func TestLambdaHandler_UnknownRoute(t *testing.T) {
	h := newTestHandler()

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/simulate",
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "PRICE"}`,
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
		Body:       `{"amount": 10000, "rate": 0.02, "term": 12, "system": "INVALID"}`,
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
