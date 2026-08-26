package crontab

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"github.com/titpetric/go-web-crontab/storage"
)

type Crontab struct {
	logger     *slog.Logger
	storage    *storage.Storage
	scheduler  *cron.Cron
	scriptPath string

	Jobs *Jobs
}

func NewCrontab(logger *slog.Logger, db *sqlx.DB) (*Crontab, error) {
	var err error

	cron := &Crontab{
		logger:  logger,
		storage: storage.NewStorage(db),
		scheduler: cron.New(
			cron.WithParser(
				cron.NewParser(
					cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
				),
			),
		),
	}

	cron.Jobs, err = NewJobs(cron)
	if err != nil {
		return nil, err
	}

	return cron, nil
}

func (cron *Crontab) Start() error {
	var jobs = cron.Jobs.jobs

	cron.logger.Info("starting job runners", "jobs", len(jobs))
	for idx, _ := range jobs {
		job := jobs[idx]
		runFunc := func() {
			if err := job.Run(cron); err != nil {
				cron.logger.Error("job failed", "job", job.Name, "error", err)
			}
		}

		if _, err := cron.scheduler.AddFunc(job.GetSchedule(), runFunc); err != nil {
			return err
		}
	}
	cron.scheduler.Start()
	return nil
}

func (cron *Crontab) Shutdown() {
	<-cron.scheduler.Stop().Done()
}

func (cron *Crontab) Load(configPath, scriptPath string) error {
	cron.scriptPath = scriptPath

	configs, err := filepath.Glob(configPath)
	if err != nil {
		return err
	}

	if len(configs) > 0 {
		for _, filename := range configs {
			err = cron.loadConfig(context.TODO(), filename, scriptPath)
			if err != nil {
				return fmt.Errorf("Error loading config: %w", err)
			}
		}
	} else {
		return fmt.Errorf("No config files found: %s", configPath)
	}
	return nil
}

func (cron *Crontab) loadConfig(ctx context.Context, filename, scriptPath string) error {
	cron.logger.Info("loading config", "file", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		marker := filename + ":" + strconv.Itoa(lineCount)
		lineCount++

		// skip empty lines and comments
		if line == "" || line[0:1] == "#" {
			continue
		}

		// parse
		lineExp := regexp.MustCompile("[\t ]+").Split(line, -1)
		if len(lineExp) < 8 || len(lineExp) > 9 {
			return fmt.Errorf("Must have 8 or 9 items per line, found %d: %s", len(lineExp), marker)
		}

		command := filepath.Join(scriptPath, lineExp[len(lineExp)-1])
		if _, err := os.Stat(command); err != nil {
			return fmt.Errorf("Script %s missing, file: %s, err: %w", command, marker, err)
		}

		// prefix 0 seconds if crontab style format
		schedule := "0 " + strings.Join(lineExp[1:6], " ")
		if len(lineExp) == 9 {
			schedule = strings.Join(lineExp[1:7], " ")
		}

		job := Job{
			cancel:   make(chan bool, 1),
			Name:     lineExp[len(lineExp)-1],
			Filename: filename,
			Command:  "./" + lineExp[len(lineExp)-1],
			Hostname: lineExp[0],
			Schedule: schedule,
		}

		if err := cron.storage.SaveJob(ctx, job.Name, job.Description); err != nil {
			return fmt.Errorf("Couldn't save job %s: %w", job.Name, err)
		}

		cron.Jobs.jobs = append(cron.Jobs.jobs, job)

		cron.logger.Info("loaded job", "job", job.Name, "host", job.Hostname, "schedule", job.Schedule, "line", marker)
	}

	return scanner.Err()
}
