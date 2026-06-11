package lambdahandler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/dto"
	"github.com/franciscomxs/simulator/internal/application/ports"
)

const maxLambdaBodyBytes = 1 << 16 // 64 KB

// Handler holds the Lambda handler dependencies.
type Handler struct {
	loanUseCase       ports.SimulateLoanUseCase
	investmentUseCase ports.SimulateInvestmentUseCase
}

// NewHandler creates a new Lambda Handler.
func NewHandler(loanUC ports.SimulateLoanUseCase, investUC ports.SimulateInvestmentUseCase) *Handler {
	return &Handler{
		loanUseCase:       loanUC,
		investmentUseCase: investUC,
	}
}

func jsonResponse(statusCode int, v any) events.APIGatewayProxyResponse {
	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("lambda: json.Marshal failed: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"internal server error"}`,
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

// Handle routes API Gateway proxy requests to the appropriate handler.
func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch {
	case req.HTTPMethod == "POST" && req.Path == "/loan/simulate":
		return h.handleLoanSimulate(req)
	case req.HTTPMethod == "POST" && req.Path == "/investment/simulate":
		return h.handleInvestmentSimulate(req)
	default:
		return jsonResponse(404, dto.ErrorResponse{Error: "not found"}), nil
	}
}

func (h *Handler) handleLoanSimulate(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(req.Body) > maxLambdaBodyBytes {
		return jsonResponse(400, dto.ErrorResponse{Error: "request body too large"}), nil
	}

	var simReq dto.SimulateRequest
	if err := json.Unmarshal([]byte(req.Body), &simReq); err != nil {
		return jsonResponse(400, dto.ErrorResponse{Error: "invalid request body"}), nil
	}

	out, err := h.loanUseCase.Execute(ports.SimulateLoanInput{
		Amount:      simReq.Amount,
		Rate:        simReq.Rate,
		Term:        simReq.Term,
		System:      simReq.System,
		GracePeriod: simReq.GracePeriod,
	})
	if err != nil {
		return jsonResponse(400, dto.ErrorResponse{Error: "invalid simulation parameters"}), nil
	}

	installments := make([]dto.InstallmentResponse, len(out.Installments))
	for i, inst := range out.Installments {
		installments[i] = dto.InstallmentResponse{
			Number:    inst.Number,
			Type:      inst.Type,
			Payment:   inst.Payment,
			Principal: inst.Principal,
			Interest:  inst.Interest,
			Balance:   inst.Balance,
		}
	}

	return jsonResponse(200, dto.SimulateResponse{
		Amount:         out.Amount,
		Rate:           out.Rate,
		Term:           out.Term,
		System:         out.System,
		GracePeriod:    out.GracePeriod,
		AdjustedAmount: out.AdjustedAmount,
		TotalDuration:  out.TotalDuration,
		TotalAmount:    out.TotalAmount,
		Installments:   installments,
	}), nil
}

func (h *Handler) handleInvestmentSimulate(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(req.Body) > maxLambdaBodyBytes {
		return jsonResponse(400, dto.ErrorResponse{Error: "request body too large"}), nil
	}

	var simReq dto.SimulateInvestmentRequest
	if err := json.Unmarshal([]byte(req.Body), &simReq); err != nil {
		return jsonResponse(400, dto.ErrorResponse{Error: "invalid request body"}), nil
	}

	out, err := h.investmentUseCase.Execute(ports.SimulateInvestmentInput{
		InitialAmount:       simReq.InitialAmount,
		MonthlyContribution: simReq.MonthlyContribution,
		Rate:                simReq.Rate,
		Term:                simReq.Term,
	})
	if err != nil {
		return jsonResponse(400, dto.ErrorResponse{Error: "invalid simulation parameters"}), nil
	}

	timeline := make([]dto.YearlySnapshotResponse, len(out.Timeline))
	for i, snap := range out.Timeline {
		timeline[i] = dto.YearlySnapshotResponse{
			Year:          snap.Year,
			Contributions: snap.Contributions,
			Balance:       snap.Balance,
		}
	}

	return jsonResponse(200, dto.SimulateInvestmentResponse{
		InitialAmount:       out.InitialAmount,
		MonthlyContribution: out.MonthlyContribution,
		Rate:                out.Rate,
		Term:                out.Term,
		FinalAmount:         out.FinalAmount,
		Timeline:            timeline,
	}), nil
}
