package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/goware/statik/fs"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//go:embed schema/mysql
var schemaMySQL embed.FS

//go:embed schema/sqlite
var schemaSqlite embed.FS

func schema(driver string) (http.FileSystem, error) {
	switch driver {
	case "sqlite":
		return http.FS(schemaSqlite), nil
	case "mysql":
		return http.FS(schemaMySQL), nil
	}
	return nil, fmt.Errorf("Unsupported driver for embedded migrations: %v", driver)
}

func statements(contents []byte, err error) ([]string, error) {
	if err != nil {
		return []string{}, err
	}
	return regexp.MustCompilePOSIX(";$").Split(string(contents), -1), nil
}

func Migrate(db *sqlx.DB, driverName string) error {
	schemaFS, err := schema(driverName)
	if err != nil {
		return err
	}

	var files []string

	if err := fs.Walk(schemaFS, "/", func(filename string, info os.FileInfo, err error) error {
		if strings.HasSuffix(filename, ".up.sql") {
			files = append(files, filename)
		}
		return err
	}); err != nil {
		return errors.Wrap(err, "Error when listing files for migrations")
	}

	sort.Strings(files)

	if len(files) == 0 {
		return errors.New("No files encoded for migration, need at least one SQL file")
	}

	migrate := func(filename string, useLog bool) error {
		status := migration{
			Project:  "go-web-crontab",
			Filename: filename,
		}
		if useLog {
			err := db.Get(&status, "select * from migrations where project=? and filename=?", status.Project, status.Filename)
			if errors.Is(err, sql.ErrNoRows) {
				err = nil
			}
			if err != nil {
				return err
			}
			if status.Status == "ok" {
				return nil
			}
		}

		up := func() error {
			stmts, err := statements(fs.ReadFile(schemaFS, filename))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Error reading migration %s", filename))
			}

			log.Println("Running migration for", filename)
			for idx, query := range stmts {
				if strings.TrimSpace(query) != "" && idx >= status.StatementIndex {
					status.StatementIndex = idx
					if _, err := db.Exec(query); err != nil {
						return err
					}
				}
			}
			log.Println("Migration done OK")
			status.Status = "ok"
			return nil
		}

		err := up()
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		if err != nil {
			status.Status = err.Error()
		}

		if useLog {
			if _, err := db.NamedExec("replace into migrations (project, filename, statement_index, status) values (:project, :filename, :statement_index, :status)", status); err != nil {
				log.Println("replace failed", err)
			}
		}
		return err
	}

	migrationFile := path.Join(path.Dir(files[0]), "migrations.sql")
	if err := migrate(migrationFile, false); err != nil {
		return err
	}

	for _, filename := range files {
		err := migrate(filename, true)
		if err != nil {
			return err
		}
	}

	return nil
}
