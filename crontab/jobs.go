package crontab

import (
	"log"

	"github.com/robfig/cron"
)

type Jobs []Job

var _jobs_cron *cron.Cron

func (jobs Jobs) Run(c *Crontab) {
	log.Println("Starting up job runners")
	_jobs_cron = cron.New()
	for idx, _ := range jobs {
		job := jobs[idx]

		runFunc := func() {
			job.Run(c)
		}

		if err := _jobs_cron.AddFunc(job.GetSchedule(), runFunc); err != nil {
			panic(err)
		}
	}
	_jobs_cron.Start()
}

func (jobs Jobs) Shutdown() {
	log.Println("Shutting down job runners")
	_jobs_cron.Stop()
}
