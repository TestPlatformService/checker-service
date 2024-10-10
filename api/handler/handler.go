package handler

import "log/slog"

type Handler struct {
	Log *slog.Logger
}

func NewHandler(logger *slog.Logger)*Handler{
	return &Handler{
		Log: logger,
	}
}
