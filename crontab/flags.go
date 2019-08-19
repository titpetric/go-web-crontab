package crontab

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	configuration struct {
		http struct {
			addr    string
			logging bool
			pretty  bool
			tracing bool
		}
		db struct {
			dsn    string
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
	if c.http.addr == "" {
		return errors.New("No HTTP Addr is set, can't listen for HTTP")
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

	flag.StringVar(&config.http.addr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.BoolVar(&config.http.logging, p("http-log"), true, "Enable/disable HTTP request log")
	flag.BoolVar(&config.http.pretty, p("http-pretty-json"), false, "Prettify returned JSON output")
	flag.BoolVar(&config.http.tracing, p("http-error-tracing"), false, "Return error stack frame")

	flag.StringVar(&config.db.dsn, p("db-dsn"), "webcron:webcron@tcp(db1:3306)/webcron?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.StringVar(&config.db.logger, p("db-logger"), "", "Logger for DB queries (none, stdout)")
}
