package crontab

import (
	"log"

	"github.com/robfig/cron"
)

var _jobs_cron *cron.Cron

func startJobs(c *Crontab, jobs []JobItem) {
	log.Println("Starting up job runners")
	_jobs_cron = cron.New()
	for idx, _ := range jobs {
		job := jobs[idx]
		runFunc := func() {
			job.Run(c)
		}

		log.Println(job.Name, job.GetSchedule())
		if err := _jobs_cron.AddFunc(job.GetSchedule(), runFunc); err != nil {
			panic(err)
		}
	}
	_jobs_cron.Start()
}

func shutdownJobs() {
	log.Println("Shutting down job runners")
	_jobs_cron.Stop()
}
