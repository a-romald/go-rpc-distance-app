package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(staticFiles, "static")
}

func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer) // Gracefully absorb panics and prints the stack trace

	assets, _ := Assets()
	fileServer := http.FileServer(http.FS(assets))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", app.Home)
	mux.Post("/handle", app.PostHandler)
	mux.Get("/results/", app.Results)

	return mux
}
