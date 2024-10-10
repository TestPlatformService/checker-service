package postgres

import (
	"checker/logs"
	"checker/model"
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

func (c *checkRepo) Submit(ctx context.Context, req *model.Request) (string, error) {

	query := `
				INSERT INTO submitted (
					question_id, user_id, question_name, status, lang, compiled_time, compiled_memory, code, user_task_id)
				VALUES 
					($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.DB.ExecContext(ctx, query, 
		req.QuestionId,
		req.UserId,
		req.QuestionName,
		req.Status,
		req.Language,
		req.CompiledTime,
		req.CompiledMemory,
		req.Code,
		req.UserTaskId)
	if err != nil {
		c.Log.Error(err.Error())
		return "", err
	}

	return "Accepted", nil
}

func (c *checkRepo) GetSubmits(ctx context.Context, req *pb.GetSubmitsRequest) (*pb.GetSubmitsResponse, error) {
	if req.QuestionId == "" || req.UserId == "" {
		return nil, errors.New("question_id and user_id cannot be empty")
	}

	var submits []*pb.GetSubmit


	rows, err := c.DB.QueryContext(ctx, `
        SELECT id, question_name, status, lang, compiled_time, compiled_memory, submitted_at
        FROM submitted
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
