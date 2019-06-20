package crontab

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"path/filepath"

	"github.com/pkg/errors"
)

// Load loads the configs
func (cron *Crontab) Load(configPath, scriptPath string) error {
	configs, err := filepath.Glob(configPath)
	if err != nil {
		return err
	}

	if len(configs) > 0 {
		for _, filename := range configs {
			err = cron.loadConfig(filename, scriptPath)
			if err != nil {
				return errors.Wrap(err, "Error loading config")
			}
		}
	} else {
		return errors.New("No config files found: " + configPath)
	}

	return errors.Wrap(os.Chdir(scriptPath), "Can't change working directory")
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
			return errors.Errorf("Must have 8 or 9 items per line, found %d: %s", len(lineExp), marker)
		}

		command := filepath.Join(scriptPath, lineExp[len(lineExp)-1])
		if _, err := os.Stat(command); err != nil {
			return errors.Errorf("Script %s missing, file: %s, err: %s", command, marker, err)
		}

		// prefix 0 seconds if crontab style format
		schedule := "0 " + strings.Join(lineExp[1:6], " ")
		if len(lineExp) == 9 {
			schedule = strings.Join(lineExp[1:7], " ")
		}

		job := &JobItem{
			cancel:   make(chan bool, 1),
			Name:     lineExp[len(lineExp)-1],
			Filename: filename,
			Command:  "./" + lineExp[len(lineExp)-1],
			Hostname: lineExp[0],
			Schedule: schedule,
		}

		cron.db.Insert("jobs", job)
		cron.Jobs.jobs = append(cron.Jobs.jobs, job)

		log.Println("Line:", lineExp)
	}
	return scanner.Err()
}
