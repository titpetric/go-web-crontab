package model

type Jobs struct {
	cron *Crontab

	jobs []JobItem
}

func NewJobs(cron *Crontab) (*Jobs, error) {
	jobs := &Jobs{
		cron: cron,
		jobs: []JobItem{},
	}
	return jobs, nil
}
