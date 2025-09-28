//go:build integration
// +build integration

package db

import (
	"context"
	"testing"
	"time"

	"github.com/titpetric/factory"
)

func TestMigrations(t *testing.T) {
	dbOptions := &factory.DatabaseConnectionOptions{
		DSN:            "webcron:webcron@tcp(localhost:3306)/webcron?collation=utf8mb4_general_ci",
		DriverName:     "mysql",
		Logger:         "stdout",
		Retries:        100,
		RetryTimeout:   2 * time.Second,
		ConnectTimeout: 2 * time.Minute,
	}
	db, err := factory.Database.TryToConnect(context.Background(), "default", dbOptions)
	if err != nil {
		t.Fatalf("Error when connecting: %+v", err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
}
