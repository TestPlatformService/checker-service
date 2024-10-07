package repo

import (
	pb "checker/genproto/checker"
	"context"
)

type IStorage interface {
	Check() ICheckStorage
	Close()
}

type ICheckStorage interface {
	Submit(context.Context, *pb.SubmitReq) (*pb.SubmitResp, error)
}