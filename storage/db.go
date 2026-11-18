package storage

import (
	"context"
	"database/sql"

	"github.com/titpetric/platform"
)

// DB will return the sql.DB in use for the user module.
// This enables reuse of the storage package from outside
// the app without exposing implementation detail.
func DB(ctx context.Context) (*sql.DB, error) {
	return platform.Database.Connect(ctx, "crontab")
}
