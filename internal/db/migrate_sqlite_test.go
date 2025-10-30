package db

import (
	"context"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/go-bridget/mig/db"
)

func TestMigrations_sqlite(t *testing.T) {
	dbOptions := &db.Options{
		Credentials: db.Credentials{
			DSN:    ":memory:",
			Driver: "sqlite",
		},
		Retries:        100,
		RetryDelay:     2 * time.Second,
		ConnectTimeout: 2 * time.Minute,
	}
	db, err := db.ConnectWithRetry(context.Background(), dbOptions)
	if err != nil {
		t.Fatalf("Error when connecting: %+v", err)
	}
	if err := Migrate(db, dbOptions.Credentials.Driver); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
