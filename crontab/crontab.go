package crontab

import (
	"bufio"
	"fmt"
	"log"
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
	db        *sqlx.DB
	storage   *storage.Storage
	scheduler *cron.Cron

	Jobs *Jobs
}

func NewCrontab(db *sqlx.DB) (*Crontab, error) {
	var err error

	cron := &Crontab{
		db:      db,
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

	log.Println("Starting up job runners")
	for idx, _ := range jobs {
		job := jobs[idx]
		runFunc := func() {
			if err := job.Run(cron); err != nil {
				log.Printf("error when running job: %+v", err)
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
	configs, err := filepath.Glob(configPath)
	if err != nil {
		return err
	}

	if len(configs) > 0 {
		for _, filename := range configs {
			err = cron.loadConfig(filename, scriptPath)
			if err != nil {
				return fmt.Errorf("Error loading config: %w", err)
			}
		}
	} else {
		return fmt.Errorf("No config files found: %s", configPath)
	}

	if err := os.Chdir(scriptPath); err != nil {
		return fmt.Errorf("Can't change working directory: %w", err)
	}
	return nil
}

func (cron *Crontab) loadConfig(filename, scriptPath string) error {
	log.Println("Loading config:", filename)
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

		if err := cron.storage.SaveJob(job.Name, job.Description); err != nil {
			return fmt.Errorf("Couldn't save job %s: %w", job.Name, err)
		}

		cron.Jobs.jobs = append(cron.Jobs.jobs, job)

		log.Println("Line:", lineExp)
	}

	return scanner.Err()
}
