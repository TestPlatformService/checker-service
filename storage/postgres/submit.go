package postgres

import (
	"checker/logs"
	"checker/storage/repo"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	pb "checker/genproto/checker"
	"github.com/google/uuid"
)

type checkRepo struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewCheckRepo(DB *sql.DB) repo.ICheckStorage {
	return &checkRepo{
		DB:  DB,
		Log: logs.NewLogger(),
	}
}

func (c *checkRepo) Submit(ctx context.Context, req *pb.SubmitReq) (*pb.SubmitResp, error) {
	submissionID, err := uuid.NewUUID()
	if err != nil {
		c.Log.Error("Failed to generate UUID for submission", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	query := `INSERT INTO submited 
	(id, code, user_task_id, submited_at) 
	VALUES ($1, $2, $3, $4)`

	_, err = c.DB.ExecContext(ctx, query, submissionID, req.Code, req.Lang, time.Now())
	if err != nil {
		c.Log.Error("Failed to insert submission", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to insert submission: %w", err)
	}

	c.Log.Info("Successfully inserted submission", slog.String("submission_id", submissionID.String()))

	return &pb.SubmitResp{}, nil
}
