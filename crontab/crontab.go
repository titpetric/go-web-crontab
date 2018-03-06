package crontab

import (
	"github.com/titpetric/factory"
)

type Crontab struct {
	db *factory.DB

	API  *API
	Jobs *CrontabJobs
}

func (Crontab) New(db *factory.DB) (*Crontab, error) {
	var err error
	cron := &Crontab{
		db: db,
	}

	opts, err := APIOptions{}.New(APIDependencies{cron})
	if err != nil {
		return nil, err
	}

	cron.Jobs, err = CrontabJobs{}.New(cron)
	if err != nil {
		return nil, err
	}

	cron.API, err = API{}.New(opts)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

func New(db *factory.DB) (*Crontab, error) {
	return Crontab{}.New(db)
}
