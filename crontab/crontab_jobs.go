package crontab

import (
	"github.com/pkg/errors"
	"github.com/titpetric/go-web-crontab/logger"
)

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

func (c *CrontabJobs) List() ([]JobItem, error) {
	jobs := []JobItem{}
	err := c.cron.db.Select(&jobs, "select name, description from jobs order by name asc")
	// update jobs array from database
	for _, job := range jobs {
		for k, _ := range c.jobs {
			if c.jobs[k].Name == job.Name {
				c.jobs[k].Description = job.Description
				break
			}
		}
	}
	return c.jobs, err
}

func (c *CrontabJobs) Save(job *JobItem) error {
	// @todo: implement save job
	return errors.New("Not implemented")
}

func (c *CrontabJobs) Get(id string) (*JobItem, error) {
	jobs, err := c.List()
	if err != nil {
		return nil, err
	}
	for k, _ := range jobs {
		if jobs[k].Name == id {
			return &jobs[k], nil
		}
	}
	return nil, errors.New("No matching job: " + id)
}

func (c *CrontabJobs) Logs(id string) ([]logger.LogEntry, error) {
	logs := []logger.LogEntry{}
	err := c.cron.db.Select(&logs, "select * from logs where name=? order by stamp desc limit 0, 1", id)
	return logs, err
}

func (c *CrontabJobs) Delete(id string) error {
	// @todo: implement delete job
	return errors.New("Not implemented")
}
