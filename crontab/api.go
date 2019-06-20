package crontab

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/titpetric/factory/resputil"
	"github.com/titpetric/go-web-crontab/logger"
)

type API struct {
	List   http.HandlerFunc
	Get    http.HandlerFunc
	Logs   http.HandlerFunc
	Save   http.HandlerFunc
	Delete http.HandlerFunc
}

func (API) New(opts *APIOptions) (*API, error) {
	cron := opts.cron

	api := &API{
		List: func(w http.ResponseWriter, r *http.Request) {
			response := &struct {
				Jobs []*JobItem `json:"jobs"`
			}{}

			request := func() (err error) {
				response.Jobs, err = cron.Jobs.List()
				return
			}

			resputil.JSON(w, request(), response)
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			response := &struct {
				Job *JobItem `json:"job"`
			}{}

			request := func() (err error) {
				id := chi.URLParam(r, "id")
				if id == "" {
					return errors.New("Missing parameter: id")
				}
				response.Job, err = cron.Jobs.Get(id)
				return
			}

			resputil.JSON(w, request(), response)
		},
		Logs: func(w http.ResponseWriter, r *http.Request) {
			response := &struct {
				Job  *JobItem           `json:"job"`
				Logs []*logger.LogEntry `json:"logs"`
			}{}

			request := func() (err error) {
				id := chi.URLParam(r, "id")
				if id == "" {
					return errors.New("Missing parameter: id")
				}
				response.Job, err = cron.Jobs.Get(id)
				if err != nil {
					return err
				}
				response.Logs, err = cron.Jobs.Logs(id)
				return
			}

			resputil.JSON(w, request(), response)
		},
		Save: func(w http.ResponseWriter, r *http.Request) {
			response := resputil.Success()

			request := func() (err error) {
				id := chi.URLParam(r, "id")
				if id == "" {
					return errors.New("Missing parameter: id")
				}
				// @todo: parse request into a job item
				item := &JobItem{}
				return cron.Jobs.Save(item)
			}

			resputil.JSON(w, request(), response)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			response := resputil.Success()

			request := func() (err error) {
				id := chi.URLParam(r, "id")
				if id == "" {
					return errors.New("Missing parameter: id")
				}
				return cron.Jobs.Delete(id)
			}

			resputil.JSON(w, request(), response)
		},
	}
	return api, nil
}
