package crontab

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/titpetric/go-web-crontab/lib"
	"github.com/titpetric/go-web-crontab/lib/logger"
)

type Job struct {
	lib.Semaphore
	cancel chan bool

	Name        string `db:"name"`
	Description string `db:"description"`

	Filename string
	Hostname string
	Schedule string
	Command  string
}

func (job *Job) GetSchedule() string {
	return job.Schedule
}

func (job *Job) Run(cron *Crontab) error {
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
			return fmt.Errorf("Couldn't run job %s and save to db: %w", job.Name, err)
		}

		return fmt.Errorf("Can't run command: %w", err)
	}

	var done = make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-job.cancel:
		if err := cmd.Process.Kill(); err != nil {
			if _, dberr := jobLog.Finish(cron.db, err); dberr != nil {
				return fmt.Errorf("Couldn't stop job %s and save to db: %w", job.Name, dberr)
			}

			return fmt.Errorf("Couldn't stop job %s: %w", job.Name, err)
		}
	case cmdError := <-done:
		if _, err := jobLog.Finish(cron.db, cmdError); err != nil {
			return fmt.Errorf("Couldn't finish job %s and save to db: %w", job.Name, err)
		}

		if cmdError != nil {
			return fmt.Errorf("Couldn't finish job %s: %w", job.Name, cmdError)
		}
	}

	return nil
}
