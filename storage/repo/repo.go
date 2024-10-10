package repo

import (
	pb "checker/genproto/checker"
	"checker/model"
	"context"
)

type IStorage interface {
	Check() ICheckStorage
	Close()
}

type ICheckStorage interface {
	// Submit(context.Context, *pb.SubmitReq) (*pb.SubmitResp, error)
	GetSubmits(context.Context, *pb.GetSubmitsRequest) (*pb.GetSubmitsResponse, error)
	Submit(context.Context, *model.Request) (string, error)
}