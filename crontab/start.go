package crontab

import (
	"time"

	"github.com/SentimensRG/sigctx"
	"github.com/go-bridget/mig/db"
	"github.com/pkg/errors"

	migrations "github.com/titpetric/go-web-crontab/internal/db"
	"github.com/titpetric/go-web-crontab/internal/service"
)

func Start() error {
	var ctx = sigctx.New()

	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	options := &db.Options{
		Credentials: db.Credentials{
			DSN:    config.db.dsn,
			Driver: config.db.driver,
		},
		Retries:        100,
		RetryDelay:     2 * time.Second,
		ConnectTimeout: 2 * time.Minute,
	}

	handle, err := db.ConnectWithRetry(ctx, options)
	if err != nil {
		return err
	}

	handle.SetMaxOpenConns(1)
	handle.SetConnMaxLifetime(30 * 24 * time.Hour)

	if err := migrations.Migrate(handle, options.Credentials.Driver); err != nil {
		return err
	}

	// crontab package
	cron, err := service.NewCrontab(handle)
	if err != nil {
		return errors.Wrap(err, "Error creating Crontab object")
	}

	err = cron.Load(config.crontab.configPath, config.crontab.scriptPath)
	if err != nil {
		return errors.Wrap(err, "Error loading Crontab configs")
	}

	err = cron.Start()
	if err != nil {
		return err
	}

	<-ctx.Done()

	cron.Shutdown()

	return nil
}
