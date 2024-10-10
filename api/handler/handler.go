package handler

import (
	"checker/service"
	"log/slog"
)

type Handler struct {
	Log *slog.Logger
	Service *service.Service
}

func NewHandler(logger *slog.Logger, service *service.Service)*Handler{
	return &Handler{
		Log: logger,
		Service: service,
	}
}
