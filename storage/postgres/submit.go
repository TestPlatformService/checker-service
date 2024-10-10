package postgres

import (
	"checker/logs"
	"checker/storage/repo"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	pb "checker/genproto/checker"
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

// func (c *checkRepo) Submit(ctx context.Context, req *pb.SubmitReq) (*pb.SubmitResp, error) {
// 	submissionID, err := uuid.NewUUID()
// 	if err != nil {
// 		c.Log.Error("Failed to generate UUID for submission", slog.String("error", err.Error()))
// 		return nil, fmt.Errorf("failed to generate UUID: %w", err)
// 	}

// 	query := `INSERT INTO submited
// 	(id, code, user_task_id, submited_at)
// 	VALUES ($1, $2, $3, $4)`

// 	_, err = c.DB.ExecContext(ctx, query, submissionID, req.Code, req.Lang, time.Now())
// 	if err != nil {
// 		c.Log.Error("Failed to insert submission", slog.String("error", err.Error()))
// 		return nil, fmt.Errorf("failed to insert submission: %w", err)
// 	}

// 	c.Log.Info("Successfully inserted submission", slog.String("submission_id", submissionID.String()))

// 	return &pb.SubmitResp{}, nil
// }

func (c *checkRepo) GetSubmits(ctx context.Context, req *pb.GetSubmitsRequest) (*pb.GetSubmitsResponse, error) {
	// Validate the request parameters (optional).
	if req.QuestionId == "" || req.UserId == "" {
		return nil, errors.New("question_id and user_id cannot be empty")
	}

	var submits []*pb.GetSubmit


	rows, err := c.DB.QueryContext(ctx, `
        SELECT id, question_name, status, language, compiled_time, compiled_memory, submitted_at
        FROM submits
        WHERE question_id = $1 AND user_id = $2`, req.QuestionId, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	for rows.Next() {
		var submit pb.GetSubmit
		if err := rows.Scan(
			&submit.Id,
			&submit.QuestionName,
			&submit.Status,
			&submit.Language,
			&submit.CompiledTime,
			&submit.CompiledMemory,
			&submit.SubmittedAt,
		); err != nil {
			return nil, err
		}
		submits = append(submits, &submit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetSubmitsResponse{
		Submits: submits,
	}, nil
}
