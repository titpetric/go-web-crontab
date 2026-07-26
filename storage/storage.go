package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/go-web-crontab/model"
	"github.com/titpetric/pdo"
)

type Storage struct {
	db     *pdo.PDO
	handle *sqlx.DB
}

func NewStorage(handle *sqlx.DB) *Storage {
	return &Storage{
		db:     pdo.New(handle),
		handle: handle,
	}
}

func (s *Storage) SaveJob(ctx context.Context, name, description string) error {
	now := time.Now()

	// The config sync passes an empty description, so keep any existing
	// (potentially user-edited) description instead of clobbering it on reload.
	if description == "" {
		existing, err := s.db.Get[string](ctx, "SELECT description FROM jobs WHERE name=?", name)
		if err == nil {
			description = *existing
		}
	}

	job := model.Jobs{
		Name:        name,
		Description: description,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	return s.db.Replace(ctx, "jobs", job)
}

func (s *Storage) SaveLog(ctx context.Context, name string, stamp time.Time, duration time.Duration, output string, exitCode int) error {
	log := model.Logs{
		Name:     name,
		Stamp:    &stamp,
		Duration: int64(duration),
		Output:   output,
		ExitCode: int64(exitCode),
	}

	return s.db.Insert(ctx, "logs", log)
}
