package model

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"github.com/titpetric/go-web-crontab/logger"
)

type JobItem struct {
	factory.Semaphore
	cancel chan bool

	Name        string `db:"name"`
	Description string `db:"description"`

	Filename string
	Hostname string
	Schedule string
	Command  string
}

func (job *JobItem) GetSchedule() string {
	return job.Schedule
}

func (job *JobItem) Run(cron *Crontab) error {
	if !job.CanRun() {
		return nil
	}

	defer job.Done()

	// Make a new logger. This takes in the stdout and stderr, log them into
	// both the application's std{out,err} and, when Finish() is called,
	// finalizes everything and write it to the database.
	var jobLog = logger.NewLog(job.Name)

	command := strings.Split(job.Command, " ")

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = jobLog.Stdout()
	cmd.Stderr = jobLog.Stderr()

	if err := cmd.Start(); err != nil {
		// Log when a task fails
		if _, err := jobLog.Finish(cron.db, err); err != nil {
			return errors.Wrap(err, "Couldn't run job "+job.Name+" and save to db")
		}

		return errors.Wrap(err, "Can't run command")
	}

	var done = make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-job.cancel:
		if err := cmd.Process.Kill(); err != nil {
			if _, dberr := jobLog.Finish(cron.db, err); dberr != nil {
				return errors.Wrap(dberr, "Couldn't stop job "+job.Name+" and save to db")
			}

			return errors.Wrap(err, "Couldn't stop job "+job.Name)
		}
	case cmdError := <-done:
		if _, err := jobLog.Finish(cron.db, cmdError); err != nil {
			return errors.Wrap(err, "Couldn't finish job "+job.Name+" and save to db")
		}

		return errors.Wrap(cmdError, "Couldn't finish job "+job.Name)
	}

	return nil
}
