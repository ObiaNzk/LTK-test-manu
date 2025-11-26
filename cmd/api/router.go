package api

import (
	"github.com/ObiaNzk/LTK-test-manu/cmd/api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/hello", handler.HelloWorld)

	return r
}
