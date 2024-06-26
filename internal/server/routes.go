package server

import (
	"nasa-htmx/cmd/web"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", templ.Handler(web.HomePage()).ServeHTTP)

	r.Get("/solar", s.SolarFlareHandler)
	r.Get("/asteroid", s.AsteroidHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	r.Get("/web/solar", templ.Handler(web.SolarFlarePage()).ServeHTTP)
	r.Post("/web/solar", web.SolarFlareHandler)

	r.Get("/web/asteroid", web.AsteroidHandler)

	return r
}
