package crontab

import (
	"github.com/pkg/errors"
	"github.com/titpetric/go-web-crontab/logger"
)

// CrontabJobs
type CrontabJobs struct {
	cron *Crontab
	jobs []*JobItem
}

// New creates a new cronjob
func (CrontabJobs) New(cron *Crontab) (*CrontabJobs, error) {
	jobs := &CrontabJobs{
		cron: cron,
		jobs: []*JobItem{},
	}

	return jobs, nil
}

// List ??
func (c *CrontabJobs) List() ([]*JobItem, error) {
	dbjobs := []*JobItem{}
	err := c.cron.db.Select(&dbjobs,
		"select name, description from jobs order by name asc",
	)

	// update job descriptions based on database
	for _, dbjob := range dbjobs {
		for _, job := range c.jobs {
			if job.Name == dbjob.Name {
				job.Description = dbjob.Description
				break
			}
		}
	}

	return c.jobs, err
}

// Save isn't implemented
func (c *CrontabJobs) Save(job *JobItem) error {
	// @todo: implement save job
	return errors.New("Not implemented")
}

// Delete isn't implemented
func (c *CrontabJobs) Delete(id string) error {
	// @todo: implement delete job
	return errors.New("Not implemented")
}

// Get gets a job from the name
func (c *CrontabJobs) Get(name string) (*JobItem, error) {
	jobs, err := c.List()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.Name == name {
			return job, nil
		}
	}

	return nil, errors.New("No matching job: " + name)
}

// Logs returns all logs for a cronjob
func (c *CrontabJobs) Logs(name string) ([]*logger.LogEntry, error) {
	logs := []*logger.LogEntry{}
	return logs, c.cron.db.Select(&logs,
		"SELECT name, description FROM logs WHICH name = ? ORDER BY name ASC",
		name,
	)
}
