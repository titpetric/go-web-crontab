package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/go-web-crontab/model"
	"github.com/titpetric/pdo"
)

type Storage struct {
	db *pdo.PDO
}

func NewStorage(handle *sqlx.DB) *Storage {
	return &Storage{
		db: pdo.New(handle),
	}
}

func (s *Storage) SaveJob(name, description string) error {
	now := time.Now()
	job := model.Jobs{
		Name:        name,
		Description: description,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		DeletedAt:   &now,
	}

	return s.db.Replace(context.Background(), "jobs", job)
}

func (s *Storage) SaveLog(name string, stamp time.Time, duration time.Duration, output string, exitCode int) error {
	log := model.Logs{
		Name:     name,
		Stamp:    &stamp,
		Duration: int64(duration),
		Output:   output,
		ExitCode: int64(exitCode),
	}

	return s.db.Insert(context.Background(), "logs", log)
}
