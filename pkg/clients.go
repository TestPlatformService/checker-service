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

func NewClients(cfg config.Config)*Clients{
	questionClient, err := grpc.NewClient(cfg.QUESTION_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		panic(err)
	}

	return &Clients{
		Question: question.NewQuestionServiceClient(questionClient),
		Input: question.NewInputServiceClient(questionClient),
		Output: question.NewOutputServiceClient(questionClient),
		Case: question.NewTestCaseServiceClient(questionClient),
	} 
}
