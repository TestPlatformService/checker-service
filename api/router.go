package api

import (
	"checker/api/handler"
	"checker/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Router(logger *slog.Logger, service *service.Service)*gin.Engine{
	router := gin.Default()
	h := handler.NewHandler(logger, service)
	router.POST("/check", h.Check)

	return router
}