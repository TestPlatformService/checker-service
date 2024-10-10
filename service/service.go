package service

import (
	pb "checker/genproto/checker"
	"log/slog"
)

type Service struct{
	Log *slog.Logger
	pb.UnimplementedCheckerServiceServer
}

func NewService(logger *slog.Logger)*Service{
	return &Service{
		Log: logger,
	}
}



