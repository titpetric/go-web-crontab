package storage

import (
	"context"
	"io/fs"
	"path"

	"github.com/go-bridget/mig/migrate"
	"github.com/jmoiron/sqlx"
)

// Migrate applies SQL migrations from the given filesystem to the database.
func Migrate(ctx context.Context, db *sqlx.DB, schema fs.FS) error {
	entries, err := fs.Glob(schema, "*.sql")
	if err != nil {
		return err
	}

	migrations := make(map[string][]byte, len(entries))
	for _, name := range entries {
		contents, _ := fs.ReadFile(schema, name)
		migrations[path.Base(name)] = contents
	}

	return migrate.RunWithFS(
		ctx,
		db,
		migrations,
		&migrate.Options{
			Project: "crontab",
			Apply:   true,
		},
	)
}
