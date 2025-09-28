package main

import (
	"log"
	"os"

	"github.com/titpetric/go-web-crontab/crontab"
	"github.com/titpetric/go-web-crontab/internal/services"
)

func main() {
	config := flags("crontab", crontab.Flags)

	// log to stdout not stderr
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go services.NewMonitor(config.monitorInterval)

	if err := crontab.Start(); err != nil {
		log.Fatalf("Error starting/running: %+v", err)
	}
}
