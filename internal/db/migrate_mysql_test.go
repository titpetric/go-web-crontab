//go:build integration
// +build integration

package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-bridget/mig/db"
)

func TestMigrations_mysql(t *testing.T) {
	dbOptions := &db.Options{
		Credentials: db.Credentials{
			DSN:    "webcron:webcron@tcp(localhost:3306)/webcron?collation=utf8mb4_general_ci",
			Driver: "mysql",
		},
		Retries:        100,
		RetryDelay:     2 * time.Second,
		ConnectTimeout: 2 * time.Minute,
	}
	handle, err := db.ConnectWithRetry(context.Background(), dbOptions)
	if err != nil {
		t.Fatalf("Error when connecting: %+v", err)
	}
	if err := Migrate(handle, dbOptions.Credentials.Driver); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
