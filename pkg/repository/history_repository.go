package repository

import (
	"context"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/pkg/models"
	"gorm.io/gorm"
)

type HistoryRepository interface {
	CreateAccountHistory(ctx context.Context, accHist models.AccountHistory, span opentracing.Span) error
	GetAccountHistoryByID(ctx context.Context, id string, span opentracing.Span) (models.AccountHistory, error)
	GetAccountHistoryByOperationID(ctx context.Context, id string, span opentracing.Span) ([]models.AccountHistory, error)
	GetAccountHistoryByClientID(ctx context.Context, cid string, span opentracing.Span) ([]models.AccountHistory, error)
}

type HistoryRepo struct {
	logger promtail.Client
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB, logger promtail.Client) *HistoryRepo {
	return &HistoryRepo{
		logger: logger,
		db: db,
	}
}

func (hc *HistoryRepo) CreateAccountHistory(ctx context.Context, history models.AccountHistory, span opentracing.Span) error {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repository:CreateAccountHistory", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	err := hc.db.Model(models.AccountHistory{}).Create(history).Error
	if err != nil {
		trace.OnError(hc.logger, repoSpan, err)
		return fmt.Errorf("can't create history: %w", err)
	}

	return nil
}

func (hc *HistoryRepo) GetAccountHistoryByID(ctx context.Context, id string, span opentracing.Span) (models.AccountHistory, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repository:GetAccountHistoryByID", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var history models.AccountHistory

	res := hc.db.Model(models.AccountHistory{}).Where("id = ?", id).Find(&history)
	if res.Error != nil {
		trace.OnError(hc.logger, repoSpan, res.Error)
		return history, fmt.Errorf("can't get history by id: %w", res.Error)
	}

	return history, nil
}

func (hc *HistoryRepo) GetAccountHistoryByOperationID(ctx context.Context, id string, span opentracing.Span) ([]models.AccountHistory, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repository:GetAccountHistoryByOperationID", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var history []models.AccountHistory

	res := hc.db.Model(models.AccountHistory{}).Where("operation_id = ?", id).Find(&history)
	if res.Error != nil {
		trace.OnError(hc.logger, repoSpan, res.Error)
		return history, fmt.Errorf("can't get history by operation id: %w", res.Error)
	}

	return history, nil
}

func (hc *HistoryRepo) GetAccountHistoryByClientID(ctx context.Context, cid string, span opentracing.Span) ([]models.AccountHistory, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repository:GetAccountHistoryByClientID", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var history []models.AccountHistory

	res := hc.db.Model(models.AccountHistory{}).Where("client_id = ?", cid).Find(&history)
	if res.Error != nil {
		trace.OnError(hc.logger, repoSpan, res.Error)
		return history, fmt.Errorf("can't get history by client id: %w", res.Error)
	}

	return history, nil
}
