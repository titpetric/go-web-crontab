package crontab

import (
	"log"

	"github.com/robfig/cron"
)

var _jobsCron *cron.Cron

func startJobs(c *Crontab, jobs []*JobItem) {
	log.Println("Starting up job runners")
	_jobsCron = cron.New()
	for _, job := range jobs {
		runFunc := func() {
			if err := job.Run(c); err != nil {
				log.Printf("Unexpected error when running job: %+v", err)
			}
		}

		if err := _jobsCron.AddFunc(job.GetSchedule(), runFunc); err != nil {
			panic(err)
		}
	}

	_jobsCron.Start()
}

func shutdownJobs() {
	log.Println("Shutting down job runners")
	_jobsCron.Stop()
}
