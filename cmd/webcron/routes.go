package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/titpetric/go-web-crontab/crontab"
)

// MountRoutes will register API routes
func MountRoutes(r chi.Router, opts *RouteOptions, cron *crontab.Crontab) {
	// CORS for local development...
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	if opts.enableLogging {
		r.Use(middleware.Logger)
	}

	r.Route("/api/1.0", func(r chi.Router) {
		// List all jobs
		r.Get("/jobs", cron.API.List)

		// Get job details, logs
		r.Get("/job/{id}", cron.API.Get)
		r.Get("/job/{id}/logs", cron.API.Logs)

		// Update job details
		r.Post("/job/{id}", cron.API.Save)

		// Delete job entry
		r.Delete("/job/{id}", cron.API.Delete)
	})
}
