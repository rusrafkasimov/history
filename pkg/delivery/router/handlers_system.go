package router

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

// PingHandler Return pong
// @Summary Ping
// @Tags testing
// @Description Just ping-pong endpoint, can be used as health indicator
// @ID ping
// @Produce json
// @Success 200 {string} string "pong"
// @Failure 500 {string} Error
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	span := opentracing.GlobalTracer().StartSpan(
		"Handler:PingHandler",
		jaeger.SelfRef(jaeger.SpanContext{}),
	)
	defer span.Finish()

	c.JSON(http.StatusOK, "pong")
}

// HealthHandler Return Status
// @Summary Health check
// @Tags system
// @Description Return answer from server for checking what server is stay alive
// @ID health-check
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 500 {string} Error
// @Router /health [get]
func HealthHandler (c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}