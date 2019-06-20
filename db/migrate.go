package db

import (
	"log"

	"github.com/titpetric/factory"
)

const schemaJobs = `
CREATE TABLE jobs (
    name        varchar(64) NOT NULL,
    description varchar(255) NOT NULL,
    PRIMARY KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

const schemaLogs = `
CREATE TABLE logs (
    name     varchar(64) NOT NULL,
    success  bool NOT NULL,
    output   longtext NOT NULL,
    stamp    datetime NOT NULL,
    duration bigint(20) NOT NULL,
    PRIMARY KEY (name,stamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`

const schemaMigrations = `
CREATE TABLE IF NOT EXISTS migrations (
    project varchar(16) NOT NULL COMMENT 'Microservice or project name',
    filename varchar(255) NOT NULL COMMENT 'yyyymmddHHMMSS.sql',
    statement_index int(11) NOT NULL COMMENT 'Statement number from SQL file',
    status text NOT NULL COMMENT 'ok or full error message',

    PRIMARY KEY (project,filename)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`

func Migrate(db *factory.DB) error {
	// Why?!
	/*
		if err := fs.Walk(statikFS, "/", func(filename string, info os.FileInfo, err error) error {
			matched, err := filepath.Match("/*.up.sql", filename)
			if matched {
				files = append(files, filename)
			}
			return err
		}); err != nil {
			return errors.Wrap(err, "Error when listing files for migrations")
		}
	*/

	up := func() error {
		for _, query := range []string{schemaJobs, schemaLogs, schemaMigrations} {
			if _, err := db.Exec(query); err != nil {
				return err
			}
		}

		return nil
	}

	if err := db.Transaction(up); err != nil {
		log.Println("replace failed", err)
	}

	return nil
}
