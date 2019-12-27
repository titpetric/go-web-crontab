package crontab

import (
	"time"

	"github.com/SentimensRG/sigctx"
	"github.com/pkg/errors"

	"github.com/titpetric/factory"

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

	if err := migrations.Migrate(db); err != nil {
		return err
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
	<-ctx.Done()

	cron.Shutdown()

	return nil
}
