package model

type Jobs struct {
	cron *Crontab

	jobs []Job
}

func NewJobs(cron *Crontab) (*Jobs, error) {
	jobs := &Jobs{
		cron: cron,
		jobs: []Job{},
	}
	return jobs, nil
}
