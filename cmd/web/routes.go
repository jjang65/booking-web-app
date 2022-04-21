package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jjang65/go-hello-word/pkg/config"
	"github.com/jjang65/go-hello-word/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Recoverer middleware
	mux.Use(middleware.Recoverer)

	// Custom middleware
	mux.Use(WriteToConsole)

	// NoSurf middleware for CSRF protection
	mux.Use(NoSurf)

	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
