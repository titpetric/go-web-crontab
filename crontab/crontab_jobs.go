package crontab

type CrontabJobs struct {
	cron *Crontab
	jobs []JobItem
}

func (CrontabJobs) New(cron *Crontab) (*CrontabJobs, error) {
	jobs := &CrontabJobs{
		cron: cron,
		jobs: []JobItem{},
	}
	return jobs, nil
}
