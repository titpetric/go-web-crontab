package crontab

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/SentimensRG/sigctx"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"

	migrations "github.com/titpetric/go-web-crontab/db"
)

func Init() error {
	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	// start/configure database connection
	factory.Database.Add("default", config.db.dsn)
	db, err := factory.Database.Get()
	if err != nil {
		return err
	}

	switch config.db.profiler {
	case "stdout":
		db.Profiler = &factory.Database.ProfilerStdout
		// @todo: profiling as an external service?
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: config.http.pretty,
		Trace:  config.http.tracing,
		Logger: func(err error) {
			// @todo: error logging
		},
	})

	return migrations.Migrate(db)
}

func Start() error {
	var ctx = sigctx.New()

	log.Println("Starting http server on address " + config.http.addr)
	listener, err := net.Listen("tcp", config.http.addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't listen on addr %s", config.http.addr))
	}

	// crontab package
	cron, err := New(factory.Database.MustGet())
	if err != nil {
		return errors.Wrap(err, "Error creating Crontab object")
	}
	err = cron.Load(config.crontab.configPath, config.crontab.scriptPath)
	if err != nil {
		return errors.Wrap(err, "Error loading Crontab configs")
	}
	cron.Start()

	r := chi.NewRouter()

	// mount routes
	MountRoutes(r, config, cron)

	go http.Serve(listener, r)
	<-ctx.Done()

	cron.Shutdown()

	return nil
}
