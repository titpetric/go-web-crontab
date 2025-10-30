package crontab

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	configuration struct {
		db struct {
			dsn    string
			driver string
			logger string
		}
		crontab struct {
			configPath string
			scriptPath string
		}
	}
)

var config *configuration

func (c *configuration) Validate() error {
	if c == nil {
		return errors.New("Config is not initialized, need to call Flags()")
	}
	if c.db.dsn == "" {
		return errors.New("No DB DSN is set, can't connect to database")
	}
	if c.crontab.configPath == "" || c.crontab.scriptPath == "" {
		return errors.New("Cron config path or script path is empty")
	}
	return nil
}

func Flags(prefix ...string) {
	if config != nil {
		return
	}
	config = new(configuration)

	p := func(s string) string {
		if len(prefix) > 1 {
			return prefix[0] + "-" + s
		}
		return s
	}

	flag.StringVar(&config.crontab.configPath, p("cron-config-path"), "cron.d/*.cron", "Glob pattern for crontab configs")
	flag.StringVar(&config.crontab.scriptPath, p("cron-script-path"), "cron.scripts/", "Path to crontab scripts folder")

	flag.StringVar(&config.db.dsn, p("db-dsn"), "file:webcron.db?cache=shared", "DSN for database connection")
	flag.StringVar(&config.db.driver, p("db-driver"), "sqlite", "Driver for database connection")
	flag.StringVar(&config.db.logger, p("db-logger"), "", "Logger for DB queries (none, stdout)")
}
