package crontab

import (
	"github.com/titpetric/go-web-crontab/internal/model"
)

type Job interface {
	Run(*model.Crontab)
	GetSchedule() string
}
