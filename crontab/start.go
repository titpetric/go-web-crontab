package crontab

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/titpetric/platform"

	"github.com/titpetric/go-web-crontab/frontend"
	"github.com/titpetric/go-web-crontab/schema"
	"github.com/titpetric/go-web-crontab/storage"
	"github.com/titpetric/go-web-crontab/web"
)

// Start validates configuration, launches the cron runners and web dashboard,
// and blocks until an interrupt signal is received.
func Start(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// validate configuration
	if err := config.Validate(); err != nil {
		return err
	}

	handle, err := storage.DB(ctx)
	if err != nil {
		return err
	}
	if err := storage.Migrate(ctx, logger, handle, schema.Migrations()); err != nil {
		return fmt.Errorf("Error applying Crontab migrations: %w", err)
	}

	// crontab package
	cron, err := NewCrontab(logger, handle)
	if err != nil {
		return fmt.Errorf("Error creating Crontab object: %w", err)
	}

	err = cron.Load(config.crontab.configPath, config.crontab.scriptPath)
	if err != nil {
		return fmt.Errorf("Error loading Crontab configs: %w", err)
	}

	// start web dashboard
	go startWeb(ctx, logger, config.crontab.scriptPath)

	err = cron.Start()
	if err != nil {
		return err
	}

	<-ctx.Done()

	cron.Shutdown()

	return nil
}

// startWeb launches the platform-based web server for the admin dashboard.
// It runs in a goroutine and respects the parent context for shutdown.
func startWeb(ctx context.Context, logger *slog.Logger, scriptPath string) {
	opts := platform.NewOptions()
	opts.ServerAddr = config.web.addr

	module, err := web.NewModule(frontend.Files, scriptPath)
	if err != nil {
		logger.Error("web dashboard", "error", err)
		return
	}

	svc := platform.New(opts)
	svc.Register(module)

	if err = svc.Start(ctx); err != nil {
		logger.Error("web dashboard", "error", err)
		return
	}

	logger.Info("web dashboard listening", "addr", config.web.addr)
	svc.Wait()
}
