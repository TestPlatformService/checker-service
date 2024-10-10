package storage

import (
	"checker/storage/postgres"
	"checker/storage/repo"
	"database/sql"
	"log/slog"
)

type Istorage interface {
	Check() repo.ICheckStorage
}

type StoragePro struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewStoragePro(DB *sql.DB, logger *slog.Logger) Istorage {
	return &StoragePro{
		DB:     DB,
		Logger: logger,
	}
}

func (pro *StoragePro) Check() repo.ICheckStorage {
	return postgres.NewCheckRepo(pro.DB)
}
