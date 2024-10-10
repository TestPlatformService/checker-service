package pkg

import (
	"checker/config"
	"checker/genproto/question"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct{
	Question question.QuestionServiceClient
	Input question.InputServiceClient
	Output question.OutputServiceClient
	Case question.TestCaseServiceClient
}

func QuestionServiceClient(cfg config.Config)*Clients{
	client, err := grpc.NewClient(cfg.QUESTION_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err)
	}

	return &Clients{
		Question: question.NewQuestionServiceClient(client),
		Input: question.NewInputServiceClient(client),
		Output: question.NewOutputServiceClient(client),
		Case: question.NewTestCaseServiceClient(client),
	} 
}
