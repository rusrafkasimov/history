package usecases

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/history/internal/convert"
	"github.com/rusrafkasimov/history/internal/errs"
	"github.com/rusrafkasimov/history/internal/queue"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/internal/types"
	"github.com/rusrafkasimov/history/pkg/dto"
	"github.com/rusrafkasimov/history/pkg/repository"
	"github.com/satori/go.uuid"
	"time"
)

type EventQueue interface {
	Subscribe() (<-chan queue.Event, error)
	Publish(ah dto.AccountHistory) error
}

type HistoryUseCases interface {
	AddHistoryToQueue(ctx context.Context, history dto.AccountHistory, span opentracing.Span) types.ErrorWithCode
	RecordAccountHistory(ctx context.Context, history dto.AccountHistory, span opentracing.Span) (dto.AccountHistory, types.ErrorWithCode)
	GetAccountHistoryByID(ctx context.Context, id string, span opentracing.Span) (dto.AccountHistory, types.ErrorWithCode)
	GetAccountHistoryByOperationID(ctx context.Context, opId string, span opentracing.Span) ([]dto.AccountHistory, types.ErrorWithCode)
	GetAccountHistoryByClientID(ctx context.Context, cId string, span opentracing.Span) ([]dto.AccountHistory, types.ErrorWithCode)
}

type HistoryUC struct {
	logger  promtail.Client
	queue   EventQueue
	accHist repository.HistoryRepository
}

func NewAccountOperationUseCases(queue EventQueue, history repository.HistoryRepository, logger promtail.Client) *HistoryUC {
	return &HistoryUC{
		logger:  logger,
		queue:   queue,
		accHist: history,
	}
}

// AddHistoryToQueue is example queue function, not for release version. Event publication must be initialized from another service
func (h *HistoryUC) AddHistoryToQueue(ctx context.Context, history dto.AccountHistory, span opentracing.Span) types.ErrorWithCode {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:AddHistoryToQueue", opentracing.ChildOf(span.Context()))
	defer func() {
		ctx.Done()
		useCaseSpan.Finish()
	}()

	err := h.queue.Publish(history)
	if err != nil {
		trace.OnError(h.logger, useCaseSpan, err)
		return errs.NewInternalServerError(err.Error())
	}

	return nil
}

func (h *HistoryUC) RecordAccountHistory(ctx context.Context, history dto.AccountHistory, span opentracing.Span) (dto.AccountHistory, types.ErrorWithCode) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:RecordAccountHistory", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()

	historyModel, errConv := convert.AccountHistoryDto(history)
	if errConv != nil {
		trace.OnError(h.logger, useCaseSpan, errConv)
		return history, errs.NewBadRequestError(errConv.Error())
	}

	historyModel.ID = uuid.NewV4()
	historyModel.CreatedAt = time.Now()
	historyModel.UpdatedAt = time.Now()

	err := h.accHist.CreateAccountHistory(ctx, historyModel, span)
	if err != nil {
		trace.OnError(h.logger, useCaseSpan, err)
		return history, errs.NewInternalServerError(err.Error())
	}

	historyDto := convert.AccountHistoryModel(historyModel)

	return historyDto, nil
}

func (h *HistoryUC) GetAccountHistoryByID(ctx context.Context, id string, span opentracing.Span) (dto.AccountHistory, types.ErrorWithCode) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetAccountHistoryByID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()

	historyModel, err := h.accHist.GetAccountHistoryByID(ctx, id, span)
	if err != nil {
		trace.OnError(h.logger, useCaseSpan, err)
		return dto.AccountHistory{}, errs.NewInternalServerError(err.Error())
	}

	historyDto := convert.AccountHistoryModel(historyModel)

	return historyDto, nil
}

func (h *HistoryUC) GetAccountHistoryByOperationID(ctx context.Context, opId string, span opentracing.Span) ([]dto.AccountHistory, types.ErrorWithCode) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetAccountHistoryByOperationID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result []dto.AccountHistory

	historyModels, err := h.accHist.GetAccountHistoryByOperationID(ctx, opId, span)
	if err != nil {
		trace.OnError(h.logger, useCaseSpan, err)
		return result, errs.NewInternalServerError(err.Error())
	}

	for _, model := range historyModels {
		historyDto := convert.AccountHistoryModel(model)
		result = append(result, historyDto)
	}

	return result, nil
}

func (h *HistoryUC) GetAccountHistoryByClientID(ctx context.Context, cId string, span opentracing.Span) ([]dto.AccountHistory, types.ErrorWithCode) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetAccountHistoryByClientID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()

	var result []dto.AccountHistory

	historyModels, err := h.accHist.GetAccountHistoryByClientID(ctx, cId, span)
	if err != nil {
		trace.OnError(h.logger, useCaseSpan, err)
		return []dto.AccountHistory{}, errs.NewInternalServerError(err.Error())
	}

	for _, model := range historyModels {
		historyDto := convert.AccountHistoryModel(model)
		result = append(result, historyDto)
	}

	return result, nil
}

func (h *HistoryUC) SubscribeOnEvents(ctx context.Context) (<-chan queue.Event, types.ErrorWithCode) {
	span := trace.MakeSpan(ctx, opentracing.GlobalTracer(), "SubscribeOnEvents")
	defer span.Finish()
	eventCh, err := h.queue.Subscribe()
	if err != nil {
		trace.OnError(h.logger, span, err)
		return nil, errs.NewInternalServerError(err.Error())
	}
	return eventCh, nil
}

func (h *HistoryUC) RunReceiveEventsLoop(ctx context.Context, eventCh <-chan queue.Event) {
	span := trace.MakeSpan(ctx, opentracing.GlobalTracer(), "RunReceiveEventsLoop")
	defer span.Finish()
	go func() {
		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					return
				}
				h.logger.Infof("%v", event.Operation().OperationCode)

				_, err := h.RecordAccountHistory(ctx, event.Operation(), span)
				if err != nil {
					trace.OnError(h.logger, span, err)
				}

				errAck := event.Ack()
				if errAck != nil {
					trace.OnError(h.logger, span, errAck)
				}

			case <-ctx.Done():
				return

			default:
				if ctx.Err() != nil {
					trace.OnError(h.logger, span, ctx.Err())
					return
				}
			}
		}
	}()
}
