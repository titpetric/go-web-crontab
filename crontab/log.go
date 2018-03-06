package crontab

import (
	"time"

	"github.com/titpetric/factory"
)

type Log struct {
	Name     string        `db:"name"`
	Stamp    time.Time     `db:"stamp"`
	Duration time.Duration `db:"duration"`
}

type Logs []Log

func (l *Log) save(db *factory.DB) error {
	return db.Insert("logs", l)
}
