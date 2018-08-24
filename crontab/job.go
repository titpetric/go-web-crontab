package crontab

import (
	"os"
	"log"
	"strings"
	"time"

	"os/exec"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
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

type Job interface {
	Run(*Crontab)
	GetSchedule() string
}

func (job *JobItem) GetSchedule() string {
	return job.Schedule
}

func (job *JobItem) Run(cron *Crontab) error {
	log.Println("Running job", job.Name)
	if !job.CanRun() {
		return nil
	}
	defer job.Done()

	command := strings.Split(job.Command, " ")

	jobLog := Log{
		Name:  job.Name,
		Stamp: time.Now(),
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-job.cancel:
		return errors.Wrap(cmd.Process.Kill(), "Killing process "+job.Name)
	case err := <-done:
		if err != nil {
			return errors.Wrap(err, "Unexpected error when running "+job.Name)
		}
		jobLog.Duration = time.Since(jobLog.Stamp)
		return errors.Wrap(jobLog.save(cron.db), "Can't save "+job.Name+" run to db")
	}
}
