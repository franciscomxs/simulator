package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"

	"github.com/franciscomxs/simulator/internal/adapters/inbound/lambdahandler"
	"github.com/franciscomxs/simulator/internal/infrastructure/container"
)

func main() {
	loanUC := container.NewSimulateLoanUseCase()
	investUC := container.NewSimulateInvestmentUseCase()
	handler := lambdahandler.NewHandler(loanUC, investUC)
	awslambda.Start(handler.Handle)
}
