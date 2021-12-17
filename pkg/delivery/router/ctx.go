package router

import (

	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/history/pkg/controllers"
	"github.com/rusrafkasimov/history/pkg/repository"
	"github.com/rusrafkasimov/history/pkg/usecases"
	"gorm.io/gorm"
)

type RepositoryContext struct {
	HistoryRepo *repository.HistoryRepo
}

type UseCaseContext struct {
	HistoryUse *usecases.HistoryUC
}

type ApplicationContext struct {
	HistoryController    *controllers.HistoryController
	RPCHistoryController *controllers.RPCHistoryController
}

func BuildRepositoryContext(db *gorm.DB) *RepositoryContext {
	return &RepositoryContext{
		HistoryRepo: repository.NewHistoryRepository(db),
	}
}

func BuildUseCaseContext(queue usecases.EventQueue, repoCtx *RepositoryContext, logger promtail.Client) *UseCaseContext {
	return &UseCaseContext{
		HistoryUse: usecases.NewAccountOperationUseCases(queue, repoCtx.HistoryRepo, logger),
	}
}

func BuildApplicationContext(ucCtx *UseCaseContext, logger promtail.Client) *ApplicationContext {
	return &ApplicationContext{
		HistoryController:    controllers.NewAccOperationController(ucCtx.HistoryUse, logger),
		RPCHistoryController: controllers.NewRPCHistoryController(ucCtx.HistoryUse, logger),
	}
}
