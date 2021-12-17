package controllers

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/pkg/dto"
	"github.com/rusrafkasimov/history/pkg/usecases"
	"github.com/rusrafkasimov/history/proto"
)

type RPCHistoryController struct {
	logger    promtail.Client
	historyUC usecases.HistoryUseCases
	operation_history.UnimplementedOperationHistoryServer
}

func NewRPCHistoryController(service usecases.HistoryUseCases, logger promtail.Client) *RPCHistoryController {
	return &RPCHistoryController{
		logger:    logger,
		historyUC: service,
	}
}


func (h *RPCHistoryController) CreateHistory(ctx context.Context, request *operation_history.OperationHistoryRequest) (*operation_history.OperationHistoryResponse, error) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:CreateHistory")
	defer controllerSpan.Finish()

	history := dto.AccountHistory{
		OperationID: request.GetOperationId(),
		ClientID: request.GetClientId(),
		AccountType: request.GetAccountType(),
		OperationCode: request.GetOperationCode(),
		OperationDesc: request.GetOperationDesc(),
		TransactionCode: request.GetOperationCode(),
		TransactionDesc: request.GetTransactionDesc(),
		Comment: request.GetComment(),
		State: request.GetState(),
		StateDesc: request.GetStateDesc(),
		ScoreSum: request.GetScoreSum(),
		AvailableBalance: request.GetAvailableBalance(),
		TransactionSum: request.GetTransactionSum(),
		TransactionTimestamp: request.GetTransactionTimestamp(),
	}

	err := h.historyUC.AddHistoryToQueue(ctx, history, controllerSpan)
	if err != nil {
		trace.OnError(h.logger, controllerSpan, err)
		return &operation_history.OperationHistoryResponse{
			ServerCode:    400,
			ServerMessage: err.Error(),
		}, err
	}

	return &operation_history.OperationHistoryResponse{
		ServerCode:    200,
		ServerMessage: "Ok",
	}, nil
}