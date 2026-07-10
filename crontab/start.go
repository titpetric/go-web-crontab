package crontab

import (
	"fmt"

	"github.com/SentimensRG/sigctx"

	"github.com/titpetric/go-web-crontab/schema"
	"github.com/titpetric/go-web-crontab/storage"
)

func Start() error {
	var ctx = sigctx.New()

	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	handle, err := storage.DB(ctx)
	if err != nil {
		return err
	}
	if err := storage.Migrate(ctx, handle, schema.Migrations()); err != nil {
		return fmt.Errorf("Error applying Crontab migrations: %w", err)
	}

	// crontab package
	cron, err := NewCrontab(handle)
	if err != nil {
		return fmt.Errorf("Error creating Crontab object: %w", err)
	}

	err = cron.Load(config.crontab.configPath, config.crontab.scriptPath)
	if err != nil {
		return fmt.Errorf("Error loading Crontab configs: %w", err)
	}

	err = cron.Start()
	if err != nil {
		return err
	}

	<-ctx.Done()

	cron.Shutdown()

	return nil
}
