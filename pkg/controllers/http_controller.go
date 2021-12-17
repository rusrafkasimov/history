package controllers

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/history/internal/errs"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/pkg/dto"
	"github.com/rusrafkasimov/history/pkg/usecases"
	"net/http"
	"time"
)

const (
	errInvID      = "empty or invalid id parameter"
	errInvJSON    = "invalid json body"
	serverTimeout = 10 * time.Second
)

type HistoryController struct {
	ctx                   context.Context
	logger                promtail.Client
	accOperationHistoryUC usecases.HistoryUseCases
}

func NewAccOperationController(service usecases.HistoryUseCases, logger promtail.Client) *HistoryController {
	return &HistoryController{
		logger:                logger,
		accOperationHistoryUC: service,
	}
}

// CreateHistory godoc
// @Summary Create history
// @Description Get JSON AccountHistory, return created JSON AccountHistory
// @Tags History
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param data body dto.AccountHistory true "AccountHistory"
// @Success 200 {object} dto.AccountHistory
// @Failure 400 {object} dto.Error Invalid JSON or AccountHistory not Valid
// @Failure 500 {object} dto.Error Can't create AccountHistory in DB
// @Router /acc_history [post]
func (hst *HistoryController) CreateHistory(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:CreateHistory")
	ctx, cancel := context.WithTimeout(c, serverTimeout)
	defer func() {
		cancel()
		controllerSpan.Finish()
	}()

	historyDto := dto.AccountHistory{}
	if err := c.ShouldBindJSON(&historyDto); err != nil {
		trace.OnError(hst.logger, controllerSpan, err)
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvJSON+err.Error()))
		return
	}

	errQueue := hst.accOperationHistoryUC.AddHistoryToQueue(ctx, historyDto, controllerSpan)
	if errQueue != nil {
		trace.OnError(hst.logger, controllerSpan, errQueue)
		errs.ErrorHandler(c, errQueue)
		return
	}

	c.JSON(http.StatusCreated, "created")
}

// GetHistoryByID godoc
// @Summary Get history by ID
// @Description Get id from param, return JSON AccountHistory
// @Tags History
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param id path int true "AccountHistory ID"
// @Success 200 {object} dto.AccountHistory
// @Failure 400 {object} dto.Error Invalid JSON or AccountHistory not Valid
// @Failure 500 {object} dto.Error Can't create AccountHistory in DB
// @Router /acc_history/{id} [get]
func (hst *HistoryController) GetHistoryByID(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetHistoryByID")
	defer controllerSpan.Finish()
	ctx := context.Background()

	id, ok := c.Params.Get("id")
	if !ok || id == "" {
		trace.OnError(hst.logger, controllerSpan, errs.NewBadRequestError(errInvID))
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvID))
		return
	}

	historyDto, err := hst.accOperationHistoryUC.GetAccountHistoryByID(ctx, id, controllerSpan)
	if err != nil {
		trace.OnError(hst.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, historyDto)
}

// GetHistoryByOpnID godoc
// @Summary Get history by operation ID
// @Description Get id from param, return JSON AccountHistory
// @Tags History
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param id path int true "AccountHistory OperationID"
// @Success 200 {array} dto.AccountHistory
// @Failure 400 {object} dto.Error Invalid JSON or AccountHistory not Valid
// @Failure 500 {object} dto.Error Can't create AccountHistory in DB
// @Router /acc_history_opn/{id} [get]
func (hst *HistoryController) GetHistoryByOpnID(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetHistoryByOpnID")
	defer controllerSpan.Finish()
	ctx := context.Background()

	id, ok := c.Params.Get("id")
	if !ok || id == "" {
		trace.OnError(hst.logger, controllerSpan, errs.NewBadRequestError(errInvID))
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvID))
		return
	}

	historyDto, err := hst.accOperationHistoryUC.GetAccountHistoryByOperationID(ctx, id, controllerSpan)
	if err != nil {
		trace.OnError(hst.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, historyDto)

}

// GetHistoryByClientID godoc
// @Summary Get history by client ID
// @Description Get id from param, return JSON AccountHistory
// @Tags History
// @Produce  json
// @Content application/json
// @Security TokenJWT
// @Param id path int true "AccountHistory ClientID"
// @Success 200 {array} dto.AccountHistory
// @Failure 400 {object} dto.Error Invalid JSON or AccountHistory not Valid
// @Failure 500 {object} dto.Error Can't create AccountHistory in DB
// @Router /acc_history_client/{id} [get]
func (hst *HistoryController) GetHistoryByClientID(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	controllerSpan := tracer.StartSpan("Controller:GetHistoryByClientID")
	defer controllerSpan.Finish()
	ctx := context.Background()

	id, ok := c.Params.Get("id")
	if !ok || id == "" {
		trace.OnError(hst.logger, controllerSpan, errs.NewBadRequestError(errInvID))
		errs.ErrorHandler(c, errs.NewBadRequestError(errInvID))
		return
	}

	historyDto, err := hst.accOperationHistoryUC.GetAccountHistoryByClientID(ctx, id, controllerSpan)
	if err != nil {
		trace.OnError(hst.logger, controllerSpan, err)
		errs.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, historyDto)
}
