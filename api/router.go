package api

import (
	"checker/api/handler"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Router(logger *slog.Logger)*gin.Engine{
	router := gin.Default()
	h := handler.NewHandler(logger)
	router.POST("/check", h.Check)

	return router
}