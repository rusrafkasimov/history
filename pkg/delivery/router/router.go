package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapUrl(router *gin.Engine, appCtx *ApplicationContext) {
	router.Use(SetMiddleware())

	// Set Middleware
	authorized := router.Group("/")
	authorized.Use(SetMiddlewareAuthentication())

	// Base test routes
	authorized.POST("/acc_history", appCtx.HistoryController.CreateHistory)
	authorized.GET("/acc_history/:id", appCtx.HistoryController.GetHistoryByID)
	authorized.GET("/acc_history_opn/:id", appCtx.HistoryController.GetHistoryByOpnID)
	authorized.GET("/acc_history_client/:id", appCtx.HistoryController.GetHistoryByClientID)

	// System Routes
	router.GET("/ping", PingHandler)
	router.GET("/health", HealthHandler)

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
