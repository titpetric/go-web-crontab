package crontab

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/SentimensRG/sigctx"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/titpetric/factory"
	"github.com/titpetric/factory/resputil"

	migrations "github.com/titpetric/go-web-crontab/db"
)

func Start() error {
	var ctx = sigctx.New()

	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	dbOptions := &factory.DatabaseConnectionOptions{
		DSN:            config.db.dsn,
		DriverName:     "mysql",
		Logger:         config.db.logger,
		Retries:        100,
		RetryTimeout:   2 * time.Second,
		ConnectTimeout: 2 * time.Minute,
	}
	db, err := factory.Database.TryToConnect(ctx, "default", dbOptions)
	if err != nil {
		return err
	}

	// configure resputil options
	resputil.SetConfig(resputil.Options{
		Pretty: config.http.pretty,
		Trace:  config.http.tracing,
		Logger: func(err error) {
			log.Printf("Error from request: %+v", err)
		},
	})

	if err := migrations.Migrate(db); err != nil {
		return err
	}

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
