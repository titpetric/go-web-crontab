package main

import (
	"log"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"

	"github.com/titpetric/go-web-crontab/crontab"
)

func main() {
	config := flags("crontab", crontab.Flags)

	// The default slog handler writes through the log package, so these two
	// calls are what decides where the service log goes and what prefixes it.
	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger := slog.Default()

	go crontab.NewMonitor(config.monitorInterval)

	if err := crontab.Start(logger); err != nil {
		logger.Error("error starting/running", "error", err)
		os.Exit(1)
	}
}
