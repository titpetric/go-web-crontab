package crontab

import (
	"github.com/titpetric/factory"
)

type Crontab struct {
	db *factory.DB

	Jobs *CrontabJobs
}

func (Crontab) New(db *factory.DB) (*Crontab, error) {
	var err error
	cron := &Crontab{
		db: db,
	}

	cron.Jobs, err = CrontabJobs{}.New(cron)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

func (cron *Crontab) Start() {
	startJobs(cron, cron.Jobs.jobs)
}

func (cron *Crontab) Shutdown() {
	shutdownJobs()
}

func New(db *factory.DB) (*Crontab, error) {
	return Crontab{}.New(db)
}
