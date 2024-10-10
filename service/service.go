package service

import (
	pb "checker/genproto/checker"
	pb2 "checker/genproto/question"
	"checker/model"
	"checker/pkg"
	"checker/storage"
	"context"
	"fmt"
	"log/slog"
)

type Service struct{
	Log *slog.Logger
	pb.UnimplementedCheckerServiceServer
	Storage storage.Istorage
	Client *pkg.Clients
}

func NewService(logger *slog.Logger, storage storage.Istorage, client *pkg.Clients)*Service{
	return &Service{
		Log: logger,
		Storage: storage,
		Client: client,
	}
}

func(S *Service) QuestionInfo(id string)(*model.QuestionInfo, error){
	var inputsOutputs = []model.InputOutput{}
	question, err := S.Client.Question.GetQuestion(context.Background(), &pb2.QuestionId{Id: id})
	if err != nil{
		S.Log.Error(fmt.Sprintf("Question ma'lumotlarini olishda xatolik: %v", err))
		return nil, err
	}
	inputs, err := S.Client.Input.GetAllQuestionInputsByQuestionId(context.Background(), &pb2.GetAllQuestionInputsByQuestionIdRequest{QuestionId: id})
	if err != nil{
		S.Log.Error(fmt.Sprintf("Questionning test-case ya'ni inputlarini olishda xatolik: %v", err))
		return nil, err
	}
	for _, i := range inputs.QuestionInputs{
		outputs, err := S.Client.Output.GetQUestionOutPutByInputId(context.Background(), &pb2.GetQUestionOutPutByInputIdRequest{InputId: i.Id})
		if err != nil{
			S.Log.Error(fmt.Sprintf("Inputga mos outputlarni olishda xatolik: %v", err))
			return nil, err
		}
		for _, o := range outputs.QuestionOutputs{
			inputsOutputs = append(inputsOutputs, model.InputOutput{
				In: i.Input,
				Out: o.Answer,
			})
		}
	}


	return &model.QuestionInfo{
		QuestionId: id,
		MemoryLimit: question.MemoryLimit,
		TimeLimit: int32(question.TimeLimit),
		IO: inputsOutputs,
	}, nil
}


func (S *Service) GetSubmits(ctx context.Context, req *pb.GetSubmitsRequest) (*pb.GetSubmitsResponse, error) {

	resp, err := S.Storage.Check().GetSubmits(ctx, req)
	if err != nil {
		S.Log.Error(fmt.Sprintf("Error while getting submits: %v", err))
		return nil, err
	}

	return resp, err
}

func (S *Service) Submit(ctx context.Context, req *model.Request) (string, error) {

	resp, err := S.Storage.Check().Submit(ctx, req)
	if err != nil {
		S.Log.Error(fmt.Sprintf("Error adding submit: %v", err))
		return "", err
	}

	return resp, err
}