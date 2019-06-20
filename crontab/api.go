package crontab

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/titpetric/factory/resputil"
	"github.com/titpetric/go-web-crontab/logger"
)

// API contains methods as HTTP Handlers
type API struct {
	List   http.HandlerFunc
	Get    http.HandlerFunc
	Logs   http.HandlerFunc
	Save   http.HandlerFunc
	Delete http.HandlerFunc
}

// New returns a new API
func (API) New(opts *APIOptions) (*API, error) {
	cron := opts.cron

	api := &API{
		List: func(w http.ResponseWriter, r *http.Request) {
			var err error
			response := &struct {
				Jobs []*JobItem `json:"jobs"`
			}{}

			response.Jobs, err = cron.Jobs.List()
			resputil.JSON(w, err, response)
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				resputil.JSON(w, errors.New("Missing parameter: id"), nil)
				return
			}

			var err error
			response := &struct {
				Job *JobItem `json:"job"`
			}{}

			response.Job, err = cron.Jobs.Get(id)
			resputil.JSON(w, err, response)
		},
		Logs: func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				resputil.JSON(w, errors.New("Missing parameter: id"), nil)
				return
			}

			var err error
			response := &struct {
				Job  *JobItem           `json:"job"`
				Logs []*logger.LogEntry `json:"logs"`
			}{}

			response.Job, err = cron.Jobs.Get(id)
			if err != nil {
				resputil.JSON(w, err, nil)
				return
			}

			response.Logs, err = cron.Jobs.Logs(id)
			if err != nil {
				resputil.JSON(w, err, nil)
				return
			}

			resputil.JSON(w, nil, response)
		},
		Save: func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				resputil.JSON(w, errors.New("Missing parameter: id"), nil)
				return
			}

			// @todo: parse request into a job item
			item := JobItem{}
			if err := cron.Jobs.Save(&item); err != nil {
				resputil.JSON(w, err, nil)
				return
			}

			resputil.JSON(w, nil, resputil.Success())
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			if id == "" {
				resputil.JSON(w, errors.New("Missing parameter: id"), nil)
				return
			}

			if err := cron.Jobs.Delete(id); err != nil {
				resputil.JSON(w, err, nil)
				return
			}

			resputil.JSON(w, nil, resputil.Success())
		},
	}

	return api, nil
}
