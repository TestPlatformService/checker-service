package service

import (
	pb "checker/genproto/checker"
	"checker/model"
	"checker/storage"
	"log/slog"
)

type Service struct{
	Log *slog.Logger
	pb.UnimplementedCheckerServiceServer
	Storage storage.Istorage
}

func NewService(logger *slog.Logger, storage storage.Istorage)*Service{
	return &Service{
		Log: logger,
		Storage: storage,
	}
}

func(S *Service) QuestionInfo(id string)(model.QuestionInfo, error) {
	
}



