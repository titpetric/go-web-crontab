package storage

import (
	"context"
	"io/fs"
	"log/slog"

	"github.com/go-bridget/mig/migrate"
	"github.com/jmoiron/sqlx"
)

// Migrate applies SQL migrations from the given filesystem to the database,
// recording them under the project name "crontab". Files are selected by mig's
// default "*.up.sql" pattern. The logger reports the filename and status of
// every migration the run touched; mig writes no output of its own.
func Migrate(ctx context.Context, logger *slog.Logger, db *sqlx.DB, schema fs.FS) error {
	m, err := migrate.NewManager(db, schema, "crontab")
	if err != nil {
		return err
	}

	applied, err := m.Apply(ctx)
	for _, item := range applied {
		logger.Info("migration", "file", item.Filename, "status", item.Status)
	}
	return err
}
